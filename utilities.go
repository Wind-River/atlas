package main

// Licensing: Apache-2.0 AND MIT
/*
LICENSE NOTICE:
===============
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

/************** https://github.com/nu7hatch/gouuid *************
LICENSE NOTICE:
===============

Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*****/

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	//"os"
	"strings"

	"github.com/nu7hatch/gouuid"
)

// Convers a boolean to Integer.
func BoolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func getType(object interface{}) string {

	// types have a "main." prefix. We need to remove it
	goType := strings.Replace(reflect.TypeOf(object).String(), "main.", "", 1)

	if strings.Contains(goType, "[]") {
		return strings.Replace(goType, "[]", "ListOf.", 1)
	} else {
		return goType
	}
}

// GetUUID - Unique Universal Identifier
func GetUUID() string {
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	return u4.String()
}

func ValidUUID(uuid string) bool {
	// TODO: need to implement ValidUUID
	// We check to see if the UUID is properly formatted.

	// The only check currently have is the lenght
	if len(uuid) == 36 {
		return true
	} else {
		return false
	}
}

//func getCryptoKeys() (thePrivateKey *rsa.PrivateKey, thePublicKey *rsa.PublicKey) {
func getCryptoKeys() (string, string) {
	// Generate RSA Keys
	thePrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	thePublicKey := &thePrivateKey.PublicKey

	return fmt.Sprintf("%x", thePrivateKey), fmt.Sprintf("%x", thePublicKey)
}

// Generate a shorter version of a UUID
func GetShortId(short_id string, uuid_str string) string {
	//last_6 = uuid_str[len(uuid_str)-6:]
	var str bytes.Buffer
	str.WriteString(short_id)
	str.WriteString(uuid_str[len(uuid_str)-6:])
	return str.String()
}

// Obtain IP Address of host machine.
func GetHostIPAddress() string {
	conn, err := net.Dial("udp", "example.com:80")
	if err != nil {
		log.Printf("[TOOLS] SYSADMIIIIIN : cannot use UDP")
		return "0.0.0.0"
	}
	defer conn.Close()
	torn := strings.Split(conn.LocalAddr().String(), ":")
	return torn[0]
}

// Parses a a full  file path name
// RETURNS: directory path, filename, file base name (name w/o ext), file extension
// Example: Path: "./d1/d2/d3/my_code.go"
//			Returns: "./d1/d2/d3/", "my_data.db", "my_data", ".db"
func FilenameDirectorySplit(full_file_path string) (string, string, string, string) {

	filename := filepath.Base(full_file_path)
	file_extension := filepath.Ext(full_file_path)
	base_name := filename[:len(filename)-len(file_extension)]
	dir_path := full_file_path[:(len(full_file_path) - len(filename))]

	return dir_path, filename, base_name, file_extension
}

// Pretty print (format) the json reply.
func createJSONFormat(data interface{}) (string, error) {

	// We want to pretty print the json reply. We need to wrap:
	//    json.NewEncoder(http_reply).Encode(reply)
	// with the following code:

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "   ") // tells how much to indent "  " spaces.
	err := encoder.Encode(data)

	if err != nil {
		return "", err
	} else {
		return buffer.String(), nil
	}
}

/***
// FileExist returns whether the given file or directory exists or not
func FileExists (path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
***/
