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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"          // BSD-3-Clause
	"github.com/russross/blackfriday" // BSD-2-Clause
)

const _META_PRIVATE_KEY = "5JrP7SpTXNXo2gLxVsjnyYzpthKAMckTFUqWC2dSccVXrsq8uXk"
const _META_PUBLIC_KEY = "036fa32c07028380269ba672b97a4da78425cac489d5e88bc061e13dba9ffd9657"

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

/**************
router.HandleFunc("/atlas/api/v1/uuid/{uuidx}", GET_UUIDTestEndPoint).Methods("GET")
func GET_UUIDTestEndPoint(http_reply http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	uuid := vars["uuidx"]
	var record UUIDRecord

	record.UUID = uuid
	record.Valid = isValidUUID(uuid)

	// reply success to indicate running.
	httpReportSuccessReply(http_reply, record)
}

********************************/

func GET_UUIDEndPoint(http_reply http.ResponseWriter, request *http.Request) {
	logEvent(request)
	if MAIN_config.Verbose_On {
		displayURLRequest(request)
	}

	var record UUIDRecord
	record.UUID = GetUUID()
	record.Valid = isValidUUID(record.UUID)

	// reply success to indicate running.
	httpReportSuccessReply(http_reply, record)
}

// Handle: POST atlas/api/v1.0/network_space/register
func POST_RegisterNetworkSpaceEndPoint(http_reply http.ResponseWriter, request *http.Request) {

	logEvent(request)

	var networkSpace NetworkSpaceRecord
	var reply PrivateKeyRecord

	err := ReadPOSTRequest(http_reply, request, &networkSpace)
	if err != nil {
		http.Error(http_reply, err.Error(), 400)
		return
	}

	// First check the system password
	signedNameAsBytes, err := hex.DecodeString(networkSpace.Password)
	if len(networkSpace.Password) == 0 || err != nil {
		// Could not Decode signed name (password) from string to bytes
		message := fmt.Sprintf("Password not valid.")
		httpReportErrorReply(http_reply, message, _EMPTY_RECORD)
		return
	}

	// check signature
	verified, err := verifySignedMessage(_META_PUBLIC_KEY, networkSpace.Name, signedNameAsBytes)
	if err != nil || verified == false {
		// signing verification is not valid
		httpReportErrorReply(http_reply, "Password not valid", _EMPTY_RECORD)
		return
	}

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
	keys, err := newKeys()
	if err != nil {
		message := fmt.Sprintf("Can't create private key: %s", err.Error())
		httpReportErrorReply(http_reply, message, _EMPTY_RECORD)
		return
	}
	networkSpace._PublicKey = keys.PublicKeyStr
	// TODO: Stop saving the private key in DB. Saved in early stage to help with testing
	// and management.
	networkSpace._PrivateKey = keys.PrivateKeyStr

	// Send back the private key.
	reply.PrivateKey = keys.PrivateKeyStr

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
		// private key was assigned to the reply above
		httpReportSuccessReply(http_reply, reply)
	}
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

	if !isValidUUID(record.UUID) {
		// error occurred - UUID not valid.
		httpReportErrorReply(http_reply, "UUID does not have a valid syntax", _EMPTY_RECORD)
		return
	}

	//TODO: Check UUID == UUID_ENCRYPT

	err = deleteLedgerNodeFromDB(record.UUID)
	if err != nil {
		// error occurred
		httpReportErrorReply(http_reply, err.Error(), _EMPTY_RECORD)
	} else {
		// Success. Simply reply we were successful
		httpReportSuccessReply(http_reply, _EMPTY_RECORD)
	}
}

func POST_DeleteNetworkSpaceEndPoint(http_reply http.ResponseWriter, request *http.Request) {
	logEvent(request)

	var record NetworkSpaceDeleteReq

	err := ReadPOSTRequest(http_reply, request, &record)
	if err != nil {
		http.Error(http_reply, err.Error(), 400)
		return
	}

	fmt.Println("The Record:", record)

	// TODO: remove this condition once the ledgers nodes implement signing the private key
	// ==============================================================================
	//if len(record.SignedName) > 0 {

	// convert signed uuid from string to bytes
	signedNameAsBytes, err := hex.DecodeString(record.SignedName)
	if err != nil {
		// Could not Decode signed uuid from string to bytes
		message := fmt.Sprintf("Could not decode signed uuid from string: %s", err.Error())
		httpReportErrorReply(http_reply, message, _EMPTY_RECORD)
		return
	}

	publicKeyStr, err := getNetworkPublicKeyFromDB(record.NetworkName)
	if err != nil {
		// could not obtain network's public key
		message := fmt.Sprintf("Could not decode signed uuid from string: %s", err.Error())
		httpReportErrorReply(http_reply, message, _EMPTY_RECORD)
		return
	}
	// If public key does not exist we assume network does not exist in db
	if len(publicKeyStr) == 0 {
		message := fmt.Sprintf("Network space '%s' does not exist", record.NetworkName)
		httpReportErrorReply(http_reply, message, _EMPTY_RECORD)
		return
	}

	// check signature
	verified, err := verifySignedMessage(publicKeyStr, record.NetworkName, signedNameAsBytes)
	if err != nil || verified == false {
		// signing verification is not valid
		if MAIN_config.Verbose_On {
			fmt.Println("Signed UUID Verification FAILED")
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		message := fmt.Sprintf("Signed Name verification FAILED.")
		if len(record.SignedName) == 0 {
			message = message + " Signed Name parameter missing."
		} else if err != nil {
			message = message + " " + err.Error()
		} else if verified == false {
			message = message + " Signature cannot be veritfied signed by private key."
		}
		httpReportErrorReply(http_reply, message, _EMPTY_RECORD)
		return
	}
	// If we get here we have successful verified the signature.
	if MAIN_config.Verbose_On {
		fmt.Println("Signed Name Verification Succeeded...")
	}

	//} // if len(record.SignedUUID) > 0 {
	// We will remove this condition once ledger nodes successfully implement the signature
	// ==============================================================================

	// We have a properly signed network name. Proceed next to remove the
	// Ledger noded registered with the network before removing the network itself
	ledgerList, _ := GetLedgerNodeListDB(record.NetworkName)
	// If it returns in err or empty we assume there are not ledger nodes.
	for i := range ledgerList {
		err = deleteLedgerNodeFromDB(ledgerList[i].UUID)
		if err != nil {
			// error occurred. Log info and skip
			logEvent(fmt.Sprintf("Error - deleting node '%s': %s", ledgerList[i].UUID, err))
			//httpReportErrorReply(http_reply, err.Error(), _EMPTY_RECORD)
		}
	}

	// Now delete the Netwpork space
	err = deleteNetworkSpaceFromDB(record.NetworkName)
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

	err := ReadPOSTRequest(http_reply, request, &record)
	if err != nil {
		http.Error(http_reply, err.Error(), 400)
		return
	}

	logEvent(request)

	if MAIN_config.Verbose_On {
		displayURLRequest(request) // display url data
		fmt.Println("UUID is: ", record.UUID)
		fmt.Println("Name is: ", record.Name)
		fmt.Println("Networks is: ", record.NetworkName)
		fmt.Printf("Address is: '%s'\n", record.APIURL)
		fmt.Printf("Signed UUID is: '%s'\n", record.SignedUUID)
	}

	if !isValidUUID(record.UUID) {
		// error occurred - UUID not valid.
		httpReportErrorReply(http_reply, "UUID does not have a valid syntax", _EMPTY_RECORD)
		return
	}

	// TODO: remove this condition once the ledgers nodes implement signing the private key
	// ==============================================================================
	if len(record.SignedUUID) > 0 {

		// convert signed uuid from string to bytes
		uuidSignedAsBytes, err := hex.DecodeString(record.SignedUUID)
		if err != nil {
			// Could not Decode signed uuid from string to bytes
			message := fmt.Sprintf("Could not decode signed uuid from string: %s", err.Error())
			httpReportErrorReply(http_reply, message, _EMPTY_RECORD)
			return
		}

		publicKeyStr, err := getNetworkPublicKeyFromDB(record.NetworkName)
		if err != nil {
			// could not obtain network's public key
			message := fmt.Sprintf("Could not decode signed uuid from string: %s", err.Error())
			httpReportErrorReply(http_reply, message, _EMPTY_RECORD)
			return
		}
		if len(publicKeyStr) == 0 {
			message := fmt.Sprintf("Network space '%s' does not exist", record.NetworkName)
			httpReportErrorReply(http_reply, message, _EMPTY_RECORD)
			return
		}

		verified, err := verifySignedMessage(publicKeyStr, record.UUID, uuidSignedAsBytes)
		if err != nil || verified == false {
			if MAIN_config.Verbose_On {
				fmt.Println("Signed UUID Verification FAILED")
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			message := fmt.Sprintf("Signed UUID verification FAILED.")
			if len(record.SignedUUID) == 0 {
				message = message + " Signed UUID parameter missing."
			} else if err != nil {
				message = message + " " + err.Error()
			} else if verified == false {
				message = message + " Signature does not match public key."
			}
			httpReportErrorReply(http_reply, message, _EMPTY_RECORD)
			return
		}
		if MAIN_config.Verbose_On {
			fmt.Println("Signed UUID Verification Succeed...")
		}

	} // if len(record.SignedUUID) > 0 {
	// ==============================================================================

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
		//fmt.Println("Well .....")
		httpReportSuccessReply(http_reply, make([]NetworkSpaceRecord, 0))
	} else {
		httpReportSuccessReply(http_reply, networkList)
	}
}

func GET_LedgerListEndPoint(http_reply http.ResponseWriter, request *http.Request) {

	//var reply ReplyType

	logEvent(request)
	vars := mux.Vars(request)
	networkName := vars["network_name"]
	ledgerList, err := GetLedgerNodeListDB(networkName)
	if err != nil {
		// error occurred
		httpReportErrorReply(http_reply, err.Error(), _EMPTY_RECORD)
	} else {
		// Success.
		for i := range ledgerList {
			ledgerList[i].NetworkName = networkName
			//fmt.Print("Yes", networkName)
		}
		if ledgerList == nil {
			httpReportSuccessReply(http_reply, make([]LedgerNodeRecord, 0))
		} else {
			httpReportSuccessReply(http_reply, ledgerList)
		}

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

	router.HandleFunc("/atlas/api/v1/uuid", GET_UUIDEndPoint).Methods("GET")

	/*
	   curl -i -H "Content-Type: application/json" -X POST -d '{"name":"sparts-test-network", "status":"Public/Active",  "public_key":"",  "description":"The Sparts Test network"}'  http://localhost:811/atlas/api/v1/network_space/register
	*/
	router.HandleFunc("/atlas/api/v1/network_space/register", POST_RegisterNetworkSpaceEndPoint).Methods("POST")
	router.HandleFunc("/atlas/api/v1/network_space", GET_NetworkSpacesEndPoint).Methods("GET")
	router.HandleFunc("/atlas/api/v1/network_node_list/{network_name}", GET_LedgerListEndPoint).Methods("GET")
	/*
		curl -i -H "Content-Type: application/json" -X POST -d '{"name":"sparts-test-network", "name_encrypt":"" }'  https://spartshub.org/atlas/api/v1/network_space/delete
	*/
	router.HandleFunc("/atlas/api/v1/network_space/delete", POST_DeleteNetworkSpaceEndPoint).Methods("POST")

	/*
		curl -i -H "Content-Type: application/json" -X POST -d  '{"name":"Wind River Test Node 1", "UUID": "6221ac8d-01f4-4de2-69ed-16b7ebae8127", "network_name":"sparts-test-network", "api_url":"http://35.166.246.146:818", "alias":"WR-Test-Node-1", "description":"A SParts test network node #1"}'  https://spartshub.org/atlas/api/v1/ledger_node/register
		curl -i -H "Content-Type: application/json" -X POST -d  '{"name":"Wind River Test Node 1", "UUID": "6221ac8d-01f4-4de2-69ed-16b7ebae8127", "network_name":"sparts-test-network-3", "api_url":"http://35.166.246.007:818", "alias":"WR-Test-Node-1", "signed_uuid":  "30440220115e2098957ec635cce636973d50bb29747ee5e97cd7bf7a1553451e5047814702200da130d72a27d1774c7b40ecfa4027deba7f32956258cec539bb6759e5e9e05f", "description":"A SParts test network node #1"}'  http://localhost:811/atlas/api/v1/ledger_node/register
	*/
	router.HandleFunc("/atlas/api/v1/ledger_node/register", POST_RegisterLedgerNodeEndPoint).Methods("POST")

	/*
	  curl -i -H "Content-Type: application/json" -X POST -d '{"uuid":"4122ac8d-01f4-4de2-69ed-16b7ebae812c", "uuid_encrypt":""}'  https://spartshub.org/atlas/api/v1/ledger_node/delete
	*/
	router.HandleFunc("/atlas/api/v1/ledger_node/delete", POST_DeleteLedgerNodeEndPoint).Methods("POST")

	// General requests
	router.HandleFunc("/atlas/api/v1/ping", GET_Ping_EndPoint).Methods("GET")
	router.HandleFunc("/atlas/api/v1/config/reload", GET_ConfigReloadEndPoint).Methods("GET")
	router.HandleFunc("/favicon.ico", GET_favicon_ico_EndPoint).Methods("GET")

	router.HandleFunc("/atlas/api/v1/sinfo", _getSecureInfo).Methods("POST")
	router.HandleFunc("/atlas/api/v1/test", PostTestEndPoint).Methods("GET")

	// API call not supported
	router.NotFoundHandler = http.HandlerFunc(httpRequestNotFound)
}

type BeforeAfter struct {
	Before string `json:"before"`
	Signed string `json:"signed"`
	After  string
	Keys   WIFKeys
}

func PostTestEndPoint(http_reply http.ResponseWriter, request *http.Request) {
	type ByteStrReply struct {
		//StringAnswer string `json:"the_string"`
		ByteAnswer string `json:"the_bytes"`
		Verify     bool   `json:"verify"`
		HexAnswer  string `json:"the_hex_str"`
	}

	var reply ByteStrReply

	keys, err := newKeys()
	if err != nil {
		fmt.Println(err)
		return
	}

	keys.PrivateKeyStr = "5JdktTxpvVnW3fwVgVPJpbZuyf4uTRSXgShymXPpCn8HzzAS44r"
	keys.PublicKeyStr = "022e7f91218d910537a7137b5dc6a5c96103682b6a92ac857074ab6eda725d4fb9"
	// signedUUIDAsHexString = "30440220115e2098957ec635cce636973d50bb29747ee5e97cd7bf7a1553451e5047814702200da130d72a27d1774c7b40ecfa4027deba7f32956258cec539bb6759e5e9e05f"
	fmt.Printf("Private: %s\nPublic: %s\n", keys.PrivateKeyStr, keys.PublicKeyStr)

	uuid := "6221ac8d-01f4-4de2-69ed-16b7ebae8127"
	signedUUID, err := signMessage(keys.PrivateKeyStr, uuid)
	if err != nil {
		fmt.Println(err)
		return
	}

	//dst := make([]byte, hex.EncodedLen(len(uuidSigned)))
	//hex.Encode(dst, uuidSigned)
	reply.HexAnswer = fmt.Sprintf("%x", signedUUID)

	//fmt.Printf("signed UUID as byte1: %x\n", uuidSigned)
	//fmt.Printf("signed UUID HEX string: %s\n", reply.HexAnswer)

	uuid2Signed, _ := hex.DecodeString(reply.HexAnswer)
	//uuid2Signed, _ = hex.DecodeString("30440220115e2098957ec635cce636973d50bb29747ee5e97cd7bf7a1553451e5047814702200da130d72a27d1774c7b40ecfa4027deba7f32956258cec539bb6759e5e9e05f")
	//fmt.Printf("signed UUID byte2 as hex: %x\n", uuid2Signed)

	verify, err := verifySignedMessage(keys.PublicKeyStr, uuid, uuid2Signed)
	if err != nil {
		fmt.Println("Test failed")
		reply.Verify = false
		//return false
	}
	if verify {
		fmt.Println("Test PASSED.")
		//return true
		reply.Verify = true
	} else {
		fmt.Println("Test FAILED.")
		reply.Verify = false
		//return false
	}

	httpReportSuccessReply(http_reply, reply)
}

func GetTestEndPoint(http_reply http.ResponseWriter, request *http.Request) {

	if MAIN_config.Verbose_On {
		displayURLRequest(request)
	}

	// reply success to indicate running.
	httpReportSuccessReply(http_reply, _EMPTY_RECORD)

}

type SecureinfoRecord struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	SecreteMsg string `json:"secrete_msg"`
	SignedMsg  string `json:"signed_msg"`
}

func _getSecureInfo(http_reply http.ResponseWriter, request *http.Request) {
	var record SecureinfoRecord

	err := ReadPOSTRequest(http_reply, request, &record)
	if err != nil {
		http.Error(http_reply, err.Error(), 400)
		return
	}
	var keys WIFKeys
	keys.PrivateKeyStr = record.PrivateKey
	keys.PublicKeyStr = record.PublicKey
	signedUUID, err := signMessage(keys.PrivateKeyStr, record.SecreteMsg)
	if err != nil {
		httpReportErrorReply(http_reply, err.Error(), _EMPTY_RECORD)
		return
	}
	record.SignedMsg = convertBytesToHexString(signedUUID)
	httpReportSuccessReply(http_reply, record)
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
