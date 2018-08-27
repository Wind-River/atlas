package main

/*
	PURPOSE:
	   This is the main code for security related functions. For example,
	   public and private key creation.
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

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec"              // License: ISC License
	"github.com/btcsuite/btcd/chaincfg"           // License: ISC License
	"github.com/btcsuite/btcd/chaincfg/chainhash" // License: ISC License
	"github.com/btcsuite/btcutil"                 // License: ISC License
)

// When it comes to all cryptocurrency coins, there are a
// diverse set of key prefixes. These prefixes are simply a byte
// that alters how the final key looks. We use Bitcoin prefixes:
const (
	BitCoinPublicKeyPrefix  = 0x00
	BitCoinPrivateKeyPrefix = 0x80
)

type WIFKeys struct {
	PrivateKeyStr string `json:"private_key_wif"`
	PublicKeyStr  string `json:"public_key_wif"`
}

/*********************
func (key *WIFKeys) isValidPublicKey() bool {
	if len(key.PublicKeyStr) == 66 {
		return true
	} else {
		return false
	}
}
*********************/

func (key *WIFKeys) isPrivateKeyValid() bool {
	_, err := btcutil.DecodeWIF(key.PrivateKeyStr)
	if err != nil || len(key.PrivateKeyStr) != 51 {
		return false
	}
	return true
}

// Check if public key is valid.
func isValidPublicKey(publicKey string) bool {
	if len(publicKey) == 66 {
		return true
	} else {
		return false
	}
}

// newKeys generates a new set of private/public key pair
func newKeys() (WIFKeys, error) {

	keys := WIFKeys{}
	// Create private WIF key using btcutil
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return keys, err
	}
	networkParams := &chaincfg.MainNetParams
	networkParams.PubKeyHashAddrID = BitCoinPublicKeyPrefix
	networkParams.PrivateKeyID = BitCoinPrivateKeyPrefix

	////privateKey, err := createPrivateKey()
	privateKey, err := btcutil.NewWIF(secret, networkParams, false)
	if err != nil {
		return keys, err
	}
	keys.PrivateKeyStr = privateKey.String()

	publicKey, err := btcutil.NewAddressPubKey(privateKey.PrivKey.PubKey().SerializeCompressed(), networkParams)
	if err != nil {
		return keys, err
	}
	keys.PublicKeyStr = publicKey.String()

	return keys, nil
}

func encodeMessageAsString(keys WIFKeys, message string) (encodedMessage string, err error) {
	pubKeyBytes, err := hex.DecodeString(keys.PublicKeyStr)
	if err != nil {
		return "", err
	}
	BTCECPublicKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return "", err
	}
	ciphertext, err := btcec.Encrypt(BTCECPublicKey, []byte(message))
	if err != nil {
		return "", err
	}
	return string(ciphertext), nil
}

func decodeMessageAsString(keys WIFKeys, message string) (decodedMessage string, err error) {
	pkBytes, err := FromWIF(keys.PrivateKeyStr)
	if err != nil {
		//fmt.Println(err)
		return "", err
	}

	BTCECPrivateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	if err != nil {
		return "", err
	}
	ciphertext := []byte(message)
	plaintext, err := btcec.Decrypt(BTCECPrivateKey, ciphertext)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(plaintext), nil
}

// signMessage signs a message (string) with a secp256k1 private key and serializing
// the generated signature.
func signMessage(WIFPrivateKey string, message string) (encodedMessage []byte, err error) {

	// Sign a message using the private key.
	messageHash := chainhash.DoubleHashB([]byte(message))

	privateKyeAsBytes, err := FromWIF(WIFPrivateKey)
	if err != nil {
		return nil, err
	}
	BTCECPrivateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKyeAsBytes)
	if err != nil {
		return nil, err
	}

	signature, err := BTCECPrivateKey.Sign(messageHash)
	if err != nil {
		return nil, err
	}

	// Serialize and return the signature.
	return signature.Serialize(), nil
}

// verifySignedMessage verifies a secp256k1 signature against a public
// key. The signature is parsed from raw bytes.
func verifySignedMessage(publicKey, message string, encodeMessage []byte) (bool, error) {

	if !isValidPublicKey(publicKey) {
		return false, fmt.Errorf("Public key is not valid format: '%s'")
	}

	pubKeyBytes, err := hex.DecodeString(publicKey) // uncompressed pubkey
	if err != nil {
		return false, err
	}
	BTCECPublicKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return false, err
	}
	signature, err := btcec.ParseSignature(encodeMessage, btcec.S256())
	if err != nil {
		return false, err
	}
	// Verify the signature for the message using the public key.
	messageHash := chainhash.DoubleHashB([]byte(message))
	// returns true if they match.
	return signature.Verify(messageHash, BTCECPublicKey), nil
}

// testSignatures test the signing functionality by signing and verfiying the signed result
func testSignatures() bool {
	keys, err := newKeys()
	if err != nil {
		fmt.Println("Test failed")
		fmt.Println(err)
		return false
	}

	fmt.Printf("Private: %s\nPublic: %s\n", keys.PrivateKeyStr, keys.PublicKeyStr)
	try := "Just try It!"

	trySigned, err := signMessage(keys.PrivateKeyStr, try)
	if err != nil {
		fmt.Println("Test failed")
		fmt.Println(err)
		return false
	}

	verify, err := verifySignedMessage(keys.PublicKeyStr, try, trySigned)
	if err != nil {
		fmt.Println("Test failed")
		return false
	}
	if verify {
		fmt.Println("Test PASSED.")
		return true
	} else {
		fmt.Println("Test FAILED.")
		return false
	}
}

// FromWIF converts a Wallet Import Format string to a Bitcoin private key //and derives the corresponding Bitcoin public key.
func FromWIF(wif string) ([]byte, error) {
	/* See https://en.bitcoin.it/wiki/Wallet_import_format */

	/* Base58 Check Decode the WIF string */
	ver, priv_bytes, err := b58checkdecode(wif)
	if err != nil {
		return nil, err
	}

	/* Check that the version byte is 0x80 */
	if ver != 0x80 {
		return nil, fmt.Errorf("Invalid WIF version 0x%02x, expected 0x80.", ver)
	}

	/* If the private key bytes length is 33, check that suffix byte is 0x01 (for compression) and strip it off */
	if len(priv_bytes) == 33 {
		if priv_bytes[len(priv_bytes)-1] != 0x01 {
			return nil, fmt.Errorf("Invalid private key, unknown suffix byte 0x%02x.", priv_bytes[len(priv_bytes)-1])
		}
		priv_bytes = priv_bytes[0:32]
	}

	return priv_bytes, nil
}

// CheckWIF checks that string wif (private key) is a valid Wallet Import Format or Wallet Import Format Compressed string. If it is not, err is populated with the reason.
func checkPrivateWIF(wif string) (valid bool, err error) {
	// See https://en.bitcoin.it/wiki/Wallet_import_format

	// Base58 Check Decode the WIF string
	ver, priv_bytes, err := b58checkdecode(wif)
	if err != nil {
		return false, err
	}

	// Check that the version byte is 0x80
	if ver != 0x80 {
		return false, fmt.Errorf("Invalid WIF version 0x%02x, expected 0x80.", ver)
	}

	// Check that private key bytes length is 32 or 33
	if len(priv_bytes) != 32 && len(priv_bytes) != 33 {
		return false, fmt.Errorf("Invalid private key bytes length %d, expected 32 or 33.", len(priv_bytes))
	}

	// If the private key bytes length is 33, check that suffix byte is 0x01 (for compression) //
	if len(priv_bytes) == 33 && priv_bytes[len(priv_bytes)-1] != 0x01 {
		return false, fmt.Errorf("Invalid private key bytes, unknown suffix byte 0x%02x.", priv_bytes[len(priv_bytes)-1])
	}

	return true, nil
}

// b58decode decodes a base-58 encoded string into a byte slice b.
func b58decode(s string) (b []byte, err error) {
	// See https://en.bitcoin.it/wiki/Base58Check_encoding */

	const BITCOIN_BASE58_TABLE = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	// Initialize
	x := big.NewInt(0)
	m := big.NewInt(58)

	// Convert string to big int
	for i := 0; i < len(s); i++ {
		b58index := strings.IndexByte(BITCOIN_BASE58_TABLE, s[i])
		if b58index == -1 {
			return nil, fmt.Errorf("Invalid base-58 character encountered: '%c', index %d.", s[i], i)
		}
		b58value := big.NewInt(int64(b58index))
		x.Mul(x, m)
		x.Add(x, b58value)
	}

	/* Convert big int to big endian bytes */
	b = x.Bytes()

	return b, nil
}

// b58checkdecode decodes base-58 check encoded string s into a version ver and byte slice b.
func b58checkdecode(s string) (ver uint8, b []byte, err error) {
	/* Decode base58 string */
	b, err = b58decode(s)
	if err != nil {
		return 0, nil, err
	}

	/* Add leading zero bytes */
	for i := 0; i < len(s); i++ {
		if s[i] != '1' {
			break
		}
		b = append([]byte{0x00}, b...)
	}

	/* Verify checksum */
	if len(b) < 5 {
		return 0, nil, fmt.Errorf("Invalid base-58 check string: missing checksum.")
	}

	// Create a new SHA256 context
	sha256_h := sha256.New()

	// SHA256 Hash #1
	sha256_h.Reset()
	sha256_h.Write(b[:len(b)-4])
	hash1 := sha256_h.Sum(nil)

	// SHA256 Hash #2
	sha256_h.Reset()
	sha256_h.Write(hash1)
	hash2 := sha256_h.Sum(nil)

	/* Compare checksum */
	if bytes.Compare(hash2[0:4], b[len(b)-4:]) != 0 {
		return 0, nil, fmt.Errorf("Invalid base-58 check string: invalid checksum.")
	}

	// Strip checksum bytes
	b = b[:len(b)-4]

	// Extract and strip version
	ver = b[0]
	b = b[1:]

	return ver, b, nil
}
