package main

// Licensing: (Apache-2.0 AND BSD-3-Clause AND BSD-2-Clause)

/*
 * NOTICE:
 * =======
 *  Copyright (c) 2017 Wind River Systems, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software  distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
 * OR CONDITIONS OF ANY KIND, either express or implied.
 */

/*****   github.com/gorilla/mux   *****************

 NOTICE:
 =======
Copyright (c) 2012 Rodrigo Moraes. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:
	 * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
	 * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
	 * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
******/

/****** "github.com/russross/blackfriday"
NOTICE
=======
Blackfriday is distributed under the Simplified BSD License:

> Copyright © 2011 Russ Ross. All rights reserved.
>
> Redistribution and use in source and binary forms, with or without
> modification, are permitted provided that the following conditions
> are met:
> 1.  Redistributions of source code must retain the above copyright
>     notice, this list of conditions and the following disclaimer.
>
> 2.  Redistributions in binary form must reproduce the above
>     copyright notice, this list of conditions and the following
>     disclaimer in the documentation and/or other materials provided with
>     the distribution.
> THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
> "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
> LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
> FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
> COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
> INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
> BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
> LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
> CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
> LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
> ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
> POSSIBILITY OF SUCH DAMAGE.
********/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	//"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux" // BSD-3-Clause
	"github.com/russross/blackfriday"
	// BSD-2-Clause
)

// File global state variables
var theLedgerAddress string
var theLedgerPort int
var __systemReset bool

// Standard method for acknowleding success status for http requests
func httpReportSuccessReply(http_reply http.ResponseWriter, result interface{}) {
	var reply ReplyType

	reply.Status = _SUCCESS
	reply.Message = _SUCCESS_MSG

	reply.Type = getType(result)
	reply.Result = result

	httpSendReply(http_reply, reply)
}

// Standard method for acknowleding success status for http requests
func httpReportErrorReply(http_reply http.ResponseWriter, message string, result interface{}) {
	var reply ReplyType

	reply.Status = _FAILURE
	reply.Message = message
	reply.Type = getType(result)
	reply.Result = result

	httpSendReply(http_reply, reply)
}

// Print useful info
func displayURLRequest(request *http.Request) {
	fmt.Println()
	fmt.Println("-----------------------------------------------")
	fmt.Println("URL Request: ", request.URL.Path)
	fmt.Println("query params were:", request.URL.Query())
	fmt.Println("Client IP:", GetHostIPAddress())

	/*******
		// Display a copy of the request for debugging.
		requestDump, err := httputil.DumpRequest(request, true)
		if err != nil {
	  	fmt.Println(err)
		}
		fmt.Println(string(requestDump))
		*********/
}

// Display debug info about a url request
func displayURLReply(url_reply string) {
	// Display http reply content for monitoring and testing purposes
	fmt.Println("-----------------------------------------------")
	fmt.Println("URL Reply:")
	fmt.Println("---------------:")
	fmt.Println(url_reply)
}

// Pretty print (format) the json reply.
func httpSendReply(http_reply http.ResponseWriter, data interface{}) {

	// We want to pretty print the json reply. We need to wrap:
	//    json.NewEncoder(http_reply).Encode(reply)
	// with the following code:

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "   ") // tells how much to indent "  " spaces.
	err := encoder.Encode(data)

	if MAIN_config.Debug_On {
		displayURLReply(buffer.String())
	}

	if err != nil {
		msg := "error - could not encode reply"
		io.WriteString(http_reply, msg)
		logEvent(msg)
	} else {
		io.WriteString(http_reply, buffer.String())
		logEvent(buffer.String())
	}
}

// Response provided when an api call is made but not found
func httpRequestNotFound(http_reply http.ResponseWriter, request *http.Request) {

	logEvent(request)
	if MAIN_config.Debug_On {
		displayURLRequest(request)
	} // display url data

	errMsg := fmt.Sprintf("Invalid api call received: '%s'", request.URL.Path)
	httpReportErrorReply(http_reply, errMsg, _EMPTY_RECORD)
}

// ==============================================
// ====   API Handler END POINTS routines	=====
// ==============================================

// Handle: GET /api/atlas/help
func GetHelpEndPoint(http_reply http.ResponseWriter, request *http.Request) {

	logEvent(request)
	if MAIN_config.Verbose_On {
		displayURLRequest(request)
	}
	// reply success to indicate running.
	b, err := ioutil.ReadFile(MAIN_config.HelpFile) // just pass the file name obtain from config file
	if err != nil {
		fmt.Print(err)
	}
	// Convert markdown to html and send as html content
	output := blackfriday.MarkdownCommon(b)
	http_reply.Header().Set("Content-Type", "text/html")
	io.WriteString(http_reply, string(output))

	// Alternatively - we could redirect to an online web page.
	// http.Redirect(http_reply, request, "https://sparts.readthedocs.io/en/latest/ledger/api.html", 301)
}

// Handle:  GET /api/sparts/ping
func GET_Ping_EndPoint(http_reply http.ResponseWriter, request *http.Request) {

	logEvent(request)
	if MAIN_config.Verbose_On {
		displayURLRequest(request)
	}
	// reply success to indicate running.
	httpReportSuccessReply(http_reply, _EMPTY_RECORD)

}

// Handle: POST atlas/api/v1.0/network_space/register
func POST_RegisterNetworkSpaceEndPoint(http_reply http.ResponseWriter, request *http.Request) {

	logEvent(request)

	var networkSpace NetworkSpaceRecord
	var reply PublicKeyRecord

	err := ReadPOSTRequest(http_reply, request, &networkSpace)
	if err != nil {
		http.Error(http_reply, err.Error(), 400)
		return
	}

	/******
		if MAIN_config.Verbose_On {
			displayURLRequest(request)
		} // display url data

		if request.Body == nil {
			http.Error(http_reply, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(request.Body).Decode(&networkSpace)
		if err != nil {
			http.Error(http_reply, err.Error(), 400)
			return
		}
	****/

	fmt.Println("Network Name is: ", networkSpace.Name)
	fmt.Println("Network Address is: ", networkSpace.Description)
	fmt.Println("Network Address is: ", networkSpace.Status)

	// Check that the Name syntax is correct
	r, err := regexp.Compile("^[A-Za-z0-9][A-Za-z0-9._-]*$")
	if err != nil ||
		len(networkSpace.Name) == 0 ||
		len(networkSpace.Name) > 80 ||
		!r.MatchString(networkSpace.Name) {
		// Invalid syntax
		message := fmt.Sprintf("Network Name '%s' is not properly formated", networkSpace.Name)
		httpReportErrorReply(http_reply, message, _EMPTY_RECORD)
		return
	}

	// Name is properly formatted. Proceed by obtaining public key
	networkSpace.PublicKey = GetUUID()

	err = AddNetworkSpaceToDB(networkSpace)
	if err != nil {
		var message string
		// error occurred. See if it was database related.
		if strings.Contains(err.Error(), "no such table") {
			message = "atlas system error: database not responding"
		} else {
			message = err.Error()
		}
		httpReportErrorReply(http_reply, message, _EMPTY_RECORD)
	} else {
		// Success.
		reply.PublicKey = networkSpace.PublicKey
		httpReportSuccessReply(http_reply, reply)
	}
}

func ReadPOSTRequest(http_reply http.ResponseWriter, request *http.Request, record interface{}) error {

	if MAIN_config.Verbose_On {
		displayURLRequest(request)
	} // display url data

	logEvent(request)

	if request.Body == nil {
		/////http.Error(http_reply, "Please send a request body", 400)
		return fmt.Errorf("Missing request body content")
	}
	err := json.NewDecoder(request.Body).Decode(&record)
	if err != nil {
		////http.Error(http_reply, err.Error(), 400)
		return err
	}

	return nil

}

func POST_DeleteLedgerNodeEndPoint(http_reply http.ResponseWriter, request *http.Request) {
	logEvent(request)
	/****
		TODO:
		Add private key for network to add a node
		Add private key for each node - such that
	   		node sends network name encrypt w/network pri key to first register node (node get's from network)
	   		node uuid encrypt w/node pri key to update
	   		node uuid encrypt w/node pri key for node to delete itself
			node uuid encrypt w/network pri key for network to delete one of its nodes
		******/

	var record LedgerNodeDeleteReq

	err := ReadPOSTRequest(http_reply, request, &record)
	if err != nil {
		http.Error(http_reply, err.Error(), 400)
		return
	}

	if !ValidUUID(record.UUID) {
		// error occurred - UUID not valid.
		httpReportErrorReply(http_reply, "UUID does not have a valid syntax", _EMPTY_RECORD)
		return
	}

	//TODO: Check UUID == UUID_ENCRYPT

	err = DeleteLedgerNodeToDB(record.UUID)
	if err != nil {
		// error occurred
		httpReportErrorReply(http_reply, err.Error(), _EMPTY_RECORD)
	} else {
		// Success. Simply reply we were successful
		httpReportSuccessReply(http_reply, _EMPTY_RECORD)
	}
}

// Handle: POST /api/atlas/ledger_node/register
func POST_RegisterLedgerNodeEndPoint(http_reply http.ResponseWriter, request *http.Request) {

	var record LedgerNodeRecord

	logEvent(request)

	if MAIN_config.Verbose_On {
		displayURLRequest(request) // display url data
		fmt.Println("UUID is: ", record.UUID)
		fmt.Println("Name is: ", record.Name)
		fmt.Println("Networks is: ", record.NetworkName)
		fmt.Printf("Address is: '%s'\n", record.APIURL)
	}

	if request.Body == nil {
		http.Error(http_reply, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&record)
	if err != nil {
		http.Error(http_reply, err.Error(), 400)
		return
	}

	if !ValidUUID(record.UUID) {
		// error occurred - UUID not valid.
		httpReportErrorReply(http_reply, "UUID does not have a valid syntax", _EMPTY_RECORD)
		return
	}

	err = AddLedgerNodeToDB(record)
	if err != nil {
		// error occurred
		httpReportErrorReply(http_reply, err.Error(), _EMPTY_RECORD)
	} else {
		// Success. Simply reply we were successful
		httpReportSuccessReply(http_reply, _EMPTY_RECORD)
	}
}

// Handle: GET /api/atlas/network_space
// Returns:
//
func GET_NetworkSpacesEndPoint(http_reply http.ResponseWriter, request *http.Request) {

	logEvent(request)

	var networkList []NetworkSpaceRecord
	networkList = GetNetworkSpaceListDB()

	if networkList == nil {
		// We have an empty list. Create an empty list
		httpSendReply(http_reply, make([]NetworkSpaceRecord, 0))
	} else {
		httpSendReply(http_reply, networkList)
	}
}

func GET_LedgerListEndPoint(http_reply http.ResponseWriter, request *http.Request) {

	logEvent(request)

	vars := mux.Vars(request)
	networkName := vars["network_name"]
	ledgerList, err := GetLedgerNodeListDB(networkName)
	if err != nil {
		// error occurred
		httpReportErrorReply(http_reply, err.Error(), _EMPTY_RECORD)
	} else {
		// Success. Simply reply we were successful
		httpReportSuccessReply(http_reply, ledgerList)
	}
}

// Handle: GET /api/atlas/config/reload
// Reloads server config file to allow config updates to used
func GET_ConfigReloadEndPoint(http_reply http.ResponseWriter, request *http.Request) {

	logEvent(request)
	//if MAIN_config.ConfigReloadAllowed {AIN_config.ConfigReloadAllowed {
	configFileLoaded, err := GetConfigurationInfo(&MAIN_config, false)
	if configFileLoaded {
		log.Println("Config file reloaded.")
		httpReportSuccessReply(http_reply, _EMPTY_RECORD)
	} else {
		if err != nil {
			httpReportErrorReply(http_reply, err.Error(), _EMPTY_RECORD)
		} else {
			httpReportErrorReply(http_reply, "The config file reload option is not enabled", _EMPTY_RECORD)
		}
	}
}

/***
func GetHelpEndPoint(http_reply http.ResponseWriter, request *http.Request) {
	http.Redirect(http_reply, request, "http://www.google.com", 301)
}
***/

// Handle: GET /favicon.ico
func GET_favicon_ico_EndPoint(writer http.ResponseWriter, request *http.Request) {
	// I get the following  additional url request browser icon) on some servers after each
	// normal url request:
	// 	/favicon.ico
	//So we will ignore it by doing nothing.
}

//func logRequest(request *http.Request) {
func logEvent(request interface{}) {
	if !MAIN_config.Logging_On {
		return
	}

	switch getType(request) {
	case "*http.Request":
		httpRequest := request.(*http.Request)
		log.Printf("Request: %s\n", httpRequest.URL.Path)
	case "string":
		msg := request.(string)
		log.Printf("%s\n", msg)
	default:
		log.Printf("Logging for Type '%s' not specified\n", getType(request))
	}

}

// ==============================================
// ====    Set up API routering calls		=====
// ==============================================

var router = mux.NewRouter()

// Initialie the RESTful API calls
func InitializeRestAPI() {
	fmt.Println("Initializing REST API ...")

	/*******
	// Initialize Ledger API address based on values from last run
	GetLedgerAPIAddress(&theLedgerAddress, &theLedgerPort)
	// ToDo - ping to see if the Ledger API is still up.
	***********/
	__systemReset = true

	// Get help
	router.HandleFunc("/atlas/api/v1/help", GetHelpEndPoint).Methods("GET")
	router.HandleFunc("/", GetHelpEndPoint).Methods("GET")

	/*
	   curl -i -H "Content-Type: application/json" -X POST -d '{"name":"zephyr", "api_address":"http://147.52.17.33:5000", "description":"The Zephyr supply chain network"}'  http://localhost:3075/atlas/api/v1/network_space/register
	*/
	router.HandleFunc("/atlas/api/v1/network_space/register", POST_RegisterNetworkSpaceEndPoint).Methods("POST")
	/*
		curl -i -H "Content-Type: application/json" -X POST -d  '{"name":"WR Node", "UUID": "7709ca8d-01f4-4de2-69ed-16b7ebae704a", "network_name":"zephyr","api_url":"http://147.52.17.33:5000", "alias":"WR", "description":"The zzephyr supply chain network"}'  http://localhost:3075/atlas/api/v1/ledger_node/register
	*/
	router.HandleFunc("/atlas/api/v1/ledger_node/register", POST_RegisterLedgerNodeEndPoint).Methods("POST")
	/***
	  curl -i -H "Content-Type: application/json" -X POST -d  '{"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae704a", "uuid_encrypt":"xyzddkdkdkdkd"}'  http://localhost:3075/api/atlas/ledger_node/delete
	  ***/
	router.HandleFunc("/atlas/api/v1/ledger_node/delete", POST_DeleteLedgerNodeEndPoint).Methods("POST")

	router.HandleFunc("/atlas/api/v1/network_space", GET_NetworkSpacesEndPoint).Methods("GET")
	router.HandleFunc("/atlas/api/v1/ledgerlist/{network_name}", GET_LedgerListEndPoint).Methods("GET")

	// General requests
	router.HandleFunc("/atlas/api/v1/ping", GET_Ping_EndPoint).Methods("GET")
	router.HandleFunc("/atlas/api/v1/config/reload", GET_ConfigReloadEndPoint).Methods("GET")
	router.HandleFunc("/favicon.ico", GET_favicon_ico_EndPoint).Methods("GET")

	// API call not supported
	router.NotFoundHandler = http.HandlerFunc(httpRequestNotFound)
}

func PostTestEndPoint(http_reply http.ResponseWriter, request *http.Request) {

}

func GetTestEndPoint(http_reply http.ResponseWriter, request *http.Request) {

	if MAIN_config.Verbose_On {
		displayURLRequest(request)
	}
	// reply success to indicate running.
	httpReportSuccessReply(http_reply, _EMPTY_RECORD)

}

// Main wait, listen and response to API requests.
func RunWaitAndRespond(http_port int) {

	// Create port string, e.g., for port 8080 we create ":8080" needed for ListenAndServe ()
	port_str := ":" + strconv.Itoa(http_port)

	// Check if https port is being used and lauch web server.
	fmt.Println("Listening on port", port_str, "...")
	if port_str == ":443" {
		if _, err := os.Stat(MAIN_config.HttpsFullchainPEM); err != nil {
			fmt.Printf("https cerficate file '%s' does not exist\n", MAIN_config.HttpsFullchainPEM)
			log.Printf("https cerficate file '%s' does not exist\n", MAIN_config.HttpsFullchainPEM)
			// terminate program.
			os.Exit(0)
		}

		if _, err := os.Stat(MAIN_config.HttpsPrivatePEM); err != nil {
			fmt.Printf("https cerficate file '%s' does not exist\n", MAIN_config.HttpsPrivatePEM)
			log.Printf("https cerficate file '%s' does not exist\n", MAIN_config.HttpsPrivatePEM)
			// terminate program.
			os.Exit(0)
		}
		err := http.ListenAndServeTLS(port_str, MAIN_config.HttpsFullchainPEM, MAIN_config.HttpsPrivatePEM, router)
		if err != nil {
			fmt.Printf("Error - encountered problem starting web server on port %s\n", port_str)
			log.Printf("Error - encountered problem starting web server on port %s\n", port_str)
			os.Exit(0)
		}
	} else { // for all other ports
		err := http.ListenAndServe(port_str, router)
		if err != nil {
			fmt.Printf("Error - encountered problem starting web server on port %s\n", port_str)
			log.Printf("Error - encountered problem starting web server on port %s\n", port_str)
			os.Exit(0)
		}
	}
}
