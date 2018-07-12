package main

/*
	This file contains build configuration parameters for the sparts cli.
*/

/*
 * NOTICE:
 * =======
 *  Copyright (c) 2018 Wind River Systems, Inc.
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

const (
	_VERSION        = "0.8"
	_DB_Model       = "0.8" // sqlite db data model
	_HELP_FILE      = "/data/AtlasAPI.md"
	_ENCRYPT_ID_STR = "mY cRypTO String is nOT fIShy"
)

//Runtime options
const (
	_DEBUG_DISPLAY_ON  = true
	_DEBUG_REST_API_ON = true
)

// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
const (
	_CYAN_FG   = "\x1b[36;1m"
	_GREEN_FG  = "\x1b[32;1m"
	_RED_FG    = "\x1b[31;1m"
	_YELLOW_FG = "\x1b[36;1m"
	_WHITE_FG  = "\x1b[39;1m"
	_COLOR_END = "\x1b[0m"
)

// rest_api.go variables
const (
	// Atlas directory look up
	_ATLAS_PING_API = "/atlas/api/v1/ping"

	// Ledger
	_ARTIFACTS_API         = "/ledger/api/v1/artifacts"
	_LEDGER_PING_API       = "/ledger/api/v1/ping"
	_PARTS_API             = "/ledger/api/v1/parts"
	_PARTS_TO_SUPPLIER_API = "/ledger/api/v1/parts/supplier"
	_SUPPLIERS_API         = "/ledger/api/v1/suppliers"
	_ARTIFACTS_URI_API     = "/ledger/api/v1/artifacts/uri"
)
