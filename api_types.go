package main

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

type ReplyType struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Type    string      `json:"result_type"`
	Result  interface{} `json:"result,omitempty"`
}

type EmptyRecord struct {
}

var _EMPTY_RECORD EmptyRecord

const (
	_SUCCESS     = "success"
	_FAILURE     = "failed"
	_SUCCESS_MSG = "Ok"
)

type NetworkSpaceRecord struct {
	Name        string `json:"name"`                  // Fullname
	Description string `json:"description,omitempty"` // 2-3 sentence description
	Status      string `json:"status,omitempty"`      // Network Status - e.g., Public/Active
	PublicKey   string `json:"public_key,omitempty"`  // Public key to vertify authorization
	Timestamp   string `json:"timestamp,omitempty"`
}

type LedgerNodeRecord struct {
	UUID        string `json:"uuid"`                  // UUID
	Name        string `json:"name"`                  // Fullname
	NetworkName string `json:"network_name"`          // Network Space name
	Alias       string `json:"alias,omitempty"`       // 1-15 alphanumberic alias
	APIURL      string `json:"api_url"`               // e.g., http://147.52.17.33:5000
	PublicKey   string `json:"public_key,omitempty"`  // Public key to verify authorization
	Description string `json:"description,omitempty"` // 2-3 sentence description
	Status      string `json:"status,omitempty"`      // Active/Inative status
	Timestamp   string `json:"timestamp,omitempty"`   // Timestamp of last update in database
}

type LedgerNodeDeleteReq struct {
	UUID        string `json:"uuid"`         // UUID
	UUIDEncrypt string `json:"uuid_encrypt"` // UUID encypted with private key
}

type PublicKeyRecord struct {
	PublicKey string `json:"public_key"`
}
