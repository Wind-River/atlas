# SParts Project Atlas Service API

[TOC]

## Atlas Service Overview

The importance and dependency on software has grown exponentially across most industries over the past decade. Software solutions, whether it is an application, library, container or an entire Linux runtime are comprised of some percentage of open source software. The dependency on open source software has grown even faster over the past five year. And although the benefits of using open source are rapidly being realized, so are the complexities and costs associated with its use. Tracking which open source components were used in a solution, when and by whom across the software supply chain is vital to its continued success. 

The SParts project provides a distributed ledger to facilitate the tracking of open source used in software solutions (parts) within various different industries such as automotive, aerospace, medical devices, industrial manufacturing, consumer electronics and the Internet of Things (IoT) more generally.  Each supply chain network maintains a ledger and each ledger is replicated across by multiple ledger **nodes**. An application utilizing  the ledger can access and submit transactions to any ledger node within a given network and be assured the transaction will be correctly propagated to all other network ledger nodes. 

One challenge an application  that utilizes the ledger faces is locating a ledger node to transact with. The Atlas service (spartshub.org) serves as a ledger node lookup service. It enables different applications using a SParts ledger to located and obtain the API address of one or more of the ledger nodes belong to a specific industry supply chain network. Although there is no requirement  that each ledger node register their API address with the Atlas service, if it does, it becomes much easier to discover. For example, if the automotive industry was using a SParts ledger which was replicated across 30 different ledger nodes (hosted by the respective  automotive companies and their suppliers) then each node could register their API address with the Atlas service. Any application wanting to access the ledger can simply query the Altas service to obtain a list of available nodes to choose from.

Each supply chain network is identified by an account name consisting of up to 50 alphanumeric letters and special characters, ''.", "_" and "-". Names are case-insensitive. For example, the Zephyr project uses the name "zephyr-parts-network". All ledgers nodes that support the zephyr-parts-chain would be registered under the network name: zephyr-parts-network.

## Example Requests

[List of Networks](https://spartshub.org/atlas/api/v1/network_space)
[List of Ledger Nodes ](https://spartshub.org/atlas/api/v1/network_node_list/zephyr-parts-network)



## I)Atlas API Calls

#### Ping Request

------

Send request to see if the Atlas service is currently available.  

```
GET /atlas/api/v1/ping
```

Example of a successful response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "EmptyRecord",
	result: 	{}
}
```

Since there is not data to return the record type **EmptyRecord** is specified in the results field.  **EmptyRecord** is defined in part II of this document. If the ledger is not available then no response will be received.

#### Ledger Node Registration

------

This call is used to register a ledger node and its API address with the Atlas service.

Fill out and send a **LedgerNodeRecord** . You must include a public key which will be used to encrypt messages to verify authenticity.

| Field                | Type   | Description                                                  |
| -------------------- | ------ | ------------------------------------------------------------ |
| uuid                 | string | unique identifier -                                          |
| name                 | string | file or envelope name                                        |
| network_name         | string | The name of the network the node belongs                     |
| network_name_encrypt | string | network name encrypted (using the network's private key)     |
| alias                | string | alias for typing                                             |
| api_url              | string | api url address (e.g., 147.11.176.111:818)                   |
| public_key           | string | public key to decrypt future ledger update requests          |
| status               | string | is the server active or inactive. Ledger may update over time |
|                      |        |                                                              |

An example single artifact request:

```
{
	uuid: "7709ca8d-01f4-4de2-69ed-16b7ebae704a",
	name: "Zephypr 1.12 SPDX file",
	network_name: "zephyr-parts-network",
	alias: "zephypr_1.12",
			  label: "Zephypr 1.12 SPDX file",
			  checksum: "f855d41c49e80b9d6f2a13148e5eb838607e92f1",
			  openchain: true,
			  content_type: "spdx"
}
```

Example curl Request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K92SiHianMJRtqRiMaQ6xwzuYz7xaFRa2C8ruBQT6edSBg87Kq", "public_key" : "02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515", "artifact": {"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae705c","name": "Zephypr 1.12 SPDX file", "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "alias": "zephypr_1.12", "label": "Zephypr 1.12 SPDX file", "openchain": "true", "content_type": "spdx"} }' http://147.11.176.111:818/ledger/api/v1/artifacts
```



**Potential Errors**:

- The requesting user does not have the appropriate access credentials to perform the add.
- One or more of the required fields UUID, checksum are missing.
- The UUID is not in a valid format. 



#### Artifact URI Add*

```
POST /ledger/api/v1/artifacts/uri
```

The request must be performed by a user with Roles: Admin or Supplier.

| Field        | Type   | Description                               |
| ------------ | ------ | ----------------------------------------- |
| version      | string | name of use                               |
| checksum     | string | artifact checksum                         |
| content_type | string | type (e.g., text, binary, archive, other) |
| size         | int    | file size in bytes                        |
| uri_type     | string | (e.g., http, ipfs, ...)                   |
| location     | string | link, path                                |

An example request:

```
{  private_key: "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt",
   public_key: "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9",
   uuid: "bcb083a1-89c7-4bd2-a568-8450350e8195",
   uri:  { version: "1.0",
		   checksum: "f67d3213907a52012a4367d8ad4f093b65abc016",
		   size:	"235120"
		   content_type: ".pdf",
		   uri_type: "http",
		   location: 	  
		      "https://github.com/zephyrstorage/_content/master/f67d3213907a52012a4367d8abc016"
		}
}
```

The uri field is of type **URIRecord**.

Example curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K92SiHianMJRtqRiMaQ6xwzuYz7xaFRa2C8ruBQT6edSBg87Kq", "public_key" : "02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515", 
"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae705c", "uri": {"version": "1.0", "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "size": "235120", "content_type": ".pdf", "uri_type": "http", "location": "https://github.com/zephyrstorage/_content/master/f67d3213907a52012a4367d8abc016" } }' http://147.11.176.111:818/ledger/api/v1/artifacts/uri
```







## II) API Types

#### EmptyRecord

```
{ }
```





#### SupplierRecord

```
{
	UUID    string `json:"uuid"`               // UUID provide w/previous registration
	Name    string `json:"name"`               // Fullname
	Alias   string `json:"alias,omitempty"`    // 1-15 alphanumeric characters
	Url     string `json:"url,omitempty"`      // 2-3 sentence description
	Parts   ListOf.PartItemRecord
}
```



#### PartItemRecord 

```
{
	PartUUID string `json:"part_uuid"` // Part uuid
}
```



#### PartRecord

```
{
	Name        string `json:"name"`                  // Fullname
	Version     string `json:"version,omitempty"`     // Version if exists.
	Alias       string `json:"label,omitempty"`       // 1-15 alphanumeric characters
	Licensing   string `json:"licensing,omitempty"`   // License expression
	Description string `json:"description,omitempty"` // Part description (1-3 sentences)
	Checksum    string `json:"checksum,omitempty"`    // License expression
	UUID        string `json:"uuid"`                  // UUID provide w/previous registration
	URIList     []URIRecord `json:"uri_list,omitempty"`     //
}
```



