package main

// Licensing: Apache-2.0
/*
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

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	config_file = "./data/atlas_config.json"
)

type Configuration struct {
	ConfigReloadAllowed bool   `json:"config_reload_allowed"`
	DatabaseFile        string `json:"database_file"`
	Debug_On            bool   `json:"debug_on"`
	Debug_DB_On         bool   `json:"debug_db_on"`
	HelpFile            string `json:"help_file"`
	HttpPort            int    `json:"http_port"`
	HttpsFullchainPEM   string `json:"https_fullchain_pem"`
	HttpsPrivatePEM     string `json:"https_privkey_pem"`
	Logging_On          bool   `json:"logging_on"`
	Verbose_On          bool   `json:"verbose_on"`
}

func GetConfigurationInfo(configuration *Configuration, first_time bool) (bool, error) {

	// When this func is called the first time we want to load the config file
	// regardless of whether the value of config_reload_allowed is true.
	// If this config variable is false on furture invocations we do not
	// want to allow the config file to be reloaded. We created a temp struct
	// to load and check config_reload_allowed first. If this varible is true
	// then we can proceed to load the current values.

	var temp_config Configuration

	file, err := os.Open(config_file)
	defer file.Close()
	if err != nil {
		if strings.Contains(err.Error(), "system cannot find the file specified") {
			return false, fmt.Errorf("atlas system error - can't access config file")
		} else {
			return false, err
		}
	}
	decoder := json.NewDecoder(file)
	temp_config = Configuration{}
	//*configuration = Configuration{}
	err = decoder.Decode(&temp_config)
	if err != nil {
		fmt.Println("error:", err)
		return false, err
	}
	if first_time || temp_config.ConfigReloadAllowed {
		*configuration = temp_config
		if MAIN_config.Verbose_On {
			fmt.Println("Configuration:")
			fmt.Println("-----------------------------------------------")
			fmt.Println("debug on	  = ", configuration.Debug_On)
			fmt.Println("debug db on	  = ", configuration.Debug_DB_On)
			fmt.Println("db file	          = ", configuration.DatabaseFile)
			fmt.Println("help file 	  = ", configuration.HelpFile)
			fmt.Println("http port	  = ", configuration.HttpPort)
			////fmt.Println("ledger_api_passwd = ", configuration.Ledger_API_Password)
			fmt.Println("logging_on 	  = ", configuration.Logging_On)
			fmt.Println("verbose on	  = ", configuration.Debug_DB_On)
			fmt.Println("config  reload	  = ", configuration.ConfigReloadAllowed)
		}
		return true, nil
	} else {
		return false, nil
	}
}

// Gloabl request counter
var requestCount = 1
var host_pid = os.Getpid() // process id
var MAIN_config Configuration
var http_ip_address = GetHostIPAddress()

func main() {
	fmt.Println()
	fmt.Println()

	// Read configuration file to set a number of global settings
	// This should be called before initializing the DB
	GetConfigurationInfo(&MAIN_config, true)

	// Initialize DB and RESTful API
	InitializeDB()
	InitializeRestAPI()

	fmt.Println()
	fmt.Println("-----------------------------------------------")
	fmt.Println("Starting Conductor ...")
	fmt.Println("Host IP:	=", http_ip_address)
	fmt.Println("Host Port:	=", MAIN_config.HttpPort)
	fmt.Println("Host PID:	=", host_pid)

	// Initialize DB and RESTful API
	//InitializeDB()
	//InitializeRestAPI()

	// Listen and responsed to requests.
	RunWaitAndRespond(MAIN_config.HttpPort)
}
