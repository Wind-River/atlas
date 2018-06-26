package main

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

type RequestType struct {
	requestType string      `json:"request_type"`
	args        interface{} `json:"request"`
}

type AddArtifactRecord struct {
	PublicKey  string         `json:"public_key"`
	PrivateKey string         `json:"private_key"`
	Artifact   ArtifactRecord `json:"artifact"`
}

type AddPartRecord struct {
	PublicKey  string     `json:"public_key"`
	PrivateKey string     `json:"private_key"`
	Part       PartRecord `json:"part"`
}

type AddURIRecord struct {
	PublicKey  string    `json:"public_key"`
	PrivateKey string    `json:"private_key"`
	UUID       string    `json:"uuid"`
	URI        URIRecord `json:"uri"`
}

type ArtifactOfEnvelopeRecord struct {
	PublicKey  string               `json:"public_key"`
	PrivateKey string               `json:"private_key"`
	Relation   ArtifactEnvelopePair `json:"relation"`
}

type ArtifactEnvelopePair struct {
	ArtifactUUID string `json:"artifact_uuid"`
	EnvelopeUUID string `json:"envelope_uuid"`
}

type ArtifactOfPart struct {
	PublicKey  string           `json:"public_key"`
	PrivateKey string           `json:"private_key"`
	Relation   ArtifactPartPair `json:"relation"`
}

type ArtifactPartPair struct {
	ArtifactUUID string `json:"artifact_uuid"`
	PartUUID     string `json:"part_uuid"`
}

type PartOfSupplierRecord struct {
	PublicKey  string           `json:"public_key"`
	PrivateKey string           `json:"private_key"`
	Relation   PartSupplierPair `json:"relation"`
}

type PartSupplierPair struct {
	PartUUID     string `json:"part_uuid"`
	SupplierUUID string `json:"supplier_uuid"`
}

type UserRecord struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Authorized string `json:"authorized"`
	PublicKey  string `json:"public_key"`
}

type ArtifactRecord struct {
	UUID         string         `json:"uuid"`
	Name         string         `json:"name"`
	Alias        string         `json:"short_id,omitempty"`
	Label        string         `json:"label,omitempty"` // Display name
	Checksum     string         `json:"checksum"`
	OpenChain    string         `json:"openchain,omitempty"`
	ContentType  string         `json:"content_type,omitempty"`
	Timestamp    string         `json:"timestamp,omitempty"`
	ArtifactList []ArtifactItem `json:"artifact_list,omitempty"`
	URIList      []URIRecord    `json:"uri_list, omitempty"`
}

type ArtifactItem struct {
	UUID string `json:"uuid"` // Artifact Universal Unique IDentifier
	Path string `json:"path"` // Path of artifact within the envelope
}

type URIRecord struct {
	Version     string `json:"version"`
	Checksum    string `json:"checksum"`
	ContentType string `json:"content_type"`   // text, envelope, binary, archive
	Size        string `json:"size,omitempty"` // size in bytes
	URIType     string `json:"uri_type"`       // e.g., http, ipfs
	Location    string `json:"location"`       // actual link
}

type PartRecord struct {
	Name        string `json:"name"`                  // Fullname
	Version     string `json:"version,omitempty"`     // Version if exists.
	Label       string `json:"label,omitempty"`       // 1-5 alphanumeric characters (unique)
	Licensing   string `json:"licensing,omitempty"`   // License expression
	Description string `json:"description,omitempty"` // Part description (1-3 sentences)
	Checksum    string `json:"checksum,omitempty"`    // License expression
	UUID        string `json:"uuid"`                  // UUID provide w/previous registration
	URI         string `json:"src_uri,omitempty"`     //
}

type PartItemRecord struct {
	PartUUID string `json:"part_id"` // Part uuid
}

type SupplierRecord struct {
	UUID  string `json:"uuid"`            // UUID provide w/previous registration
	Name  string `json:"name"`            // Fullname
	Alias string `json:"alias,omitempty"` // 1-15 alphanumeric characters
	Url   string `json:"url,omitempty"`   // 2-3 sentence description
	Parts []PartItemRecord
}
