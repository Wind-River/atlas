# SParts Project Atlas Service API

[TOC]

## Atlas Service Overview

The importance and dependency on software has grown exponentially across most industries over the past decade. Software solutions, whether it is an application, library, container or an entire Linux runtime are comprised of some percentage of open source software. The dependency on open source software has grown even faster over the past five year. And although the benefits of using open source are rapidly being realized, so are the complexities and costs associated with its use. Tracking which open source components were used in a solution, when and by whom across the software supply chain is vital to its continued success. 

The SParts project provides a distributed ledger to facilitate the tracking of open source used in software solutions (parts) within various different industries such as automotive, aerospace, medical devices, industrial manufacturing, consumer electronics and the Internet of Things (IoT) more generally.  Each supply chain network maintains a ledger and each ledger is replicated across by multiple ledger **nodes**. An application utilizing  the ledger can access and submit transactions to any ledger node within a given network and be assured the transaction will be correctly propagated to all other network ledger nodes. 

One challenge an application  that utilizes the ledger faces is locating a ledger node to transact with. The Atlas service (spartshub.org) serves as a ledger node lookup service. It enables different applications using a SParts ledger to located and obtain the API address of one or more of the ledger nodes belong to a specific industry supply chain network. Although there is no requirement  that each ledger node register their API address with the Atlas service, if it does, it becomes much easier to discover. For example, if the automotive industry was using a SParts ledger which was replicated across 30 different ledger nodes (hosted by the respective  automotive companies and their suppliers) then each node could register their API address with the Atlas service. Any application wanting to access the ledger can simply query the Altas service to obtain a list of available nodes to choose from.

Each supply chain network is identified by an account name consisting of up to 50 alphanumeric letters and special characters, ''.", "_" and "-". Names are case-insensitive. For example, the Zephyr project uses the name "zephyr-parts-network". All ledgers nodes that support the zephyr-parts-chain would be registered under the network name: zephyr-parts-network.

## Example Requests

- [Ping Service](https://spartshub.org/atlas/api/v1/ping) 
- [List of Networks](https://spartshub.org/atlas/api/v1/network_space)
- [List of Ledger Nodes ](https://spartshub.org/atlas/api/v1/network_node_list/zephyr-parts-network)
- [Get UUID](https://spartshub.org/atlas/api/v1/uuid)




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

#### Get UUID

------

```
GET /atlas/api/v1/uuid
```

Will generate a new universally unique identifier returned via a UUIDRecord.

```
{
	status: "success",
	message: "Ok",
	result_type: "UUIDRecord",
	result: {
		uuid: "8447e9a6-687c-45e9-44dd-66000e2a4aa3"
	}
}
```

#### Ledger Node Registration

------

Post request to register a ledger node in a given network:

```
POST /atlas/api/v1/ledger_node/register
```

 

This call is used to register a ledger node and its API address for a existing network.

Fill out and send a **LedgerNodeRecord** . You must include a public key which will be used to encrypt messages to verify authenticity.

*Required fields 

| Field         | Type   | Description                                                  |
| ------------- | ------ | ------------------------------------------------------------ |
| *uuid         | string | - unique identifier                                          |
| name          | string | - file or envelope name                                      |
| *network_name | string | - The name of the network the node belongs                   |
| *api_url      | string | - api url address (e.g., 147.11.176.111:818)                 |
| signed_uuid   | string | - signed uuid string to verify it came from a network authorized ledger node. Each network has a designated private key to sign each uuid. |
| alias         | string | - Alias (short name)                                         |
| description   | string | - Description of node                                        |

An example single artifact request:

```
{
	*uuid: "6221ac8d-01f4-4de2-69ed-16b7ebae8127",
	name: "Wind River Test Node 1",
	*network_name: "sparts-test-network",
	*api_url: "http://35.166.246.146:818",
	signed_uuid: 		"3045022100fb4d2233ef0d155ab57e197ee2bf3233b4bc13ddc071906d2e7733298322c8b1022036ceb0d1729cff8d2cb9470375735f35608781b60f31384aa23568e9a73eba42"
	alias: "test-node-1",
	description: "The SParts project test node"
}
```



Example curl Request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"name":"Wind River Test Node 1", "UUID": "6221ac8d-01f4-4de2-69ed-16b7ebae8127", "network_name":"sparts-test-network", "api_url":"http://35.166.246.146:818", ", "signed_uuid":  "3045022100fb4d2233ef0d155ab57e197ee2bf3233b4bc13ddc071906d2e7733298322c8b1022036ceb0d1729cff8d2cb9470375735f35608781b60f31384aa23568e9a73eba42", "alias":"WR-Test-Node-1", "description":"The SParts project test node"}'  https://spartshub.org/atlas/api/v1/ledger_node/register
```



**Potential Errors**:

- The requesting user does not have the appropriate access credentials to perform the add.
- One or more of the required fields UUID, network_name, ... are missing.
- The UUID is not in a valid format.
- network does not exist. 



#### Ledger Node Deletion

------

Post request to delete a ledger node from a given network:

```
POST /atlas/api/v1/ledger_node/delete
```

 *Required fields 

| Field         | Type   | Description                                                  |
| ------------- | ------ | ------------------------------------------------------------ |
| *uuid         | string | - unique identifier                                          |
| name          | string | - file or envelope name                                      |
| *network_name | string | - The name of the network the node belongs                   |
| *api_url      | string | - api url address (e.g., 147.11.176.111:818)                 |
| signed_uuid   | string | - signed uuid string to verify it came from a network authorized ledger node. Each network has a designated private key to sign each uuid. |
| alias         | string | - Alias (short name)                                         |
| description   | string | - Description of node                                        |

Example curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d '{"name":"sparts-test-network", "name_encrypt":"xyzddkdkdkdkd" }'  https://spartshub.org/atlas/api/v1/ledger_node/delete
```



## II) API Types

#### EmptyRecord

```
{ }
```

#### LedgerNodeRecord 

```
{
	UUID        string `json:"uuid"`        // UUID
	Name        string `json:"name"`        // Fullname
	NetworkName string `json:"network_name"`// Network Space name
	Alias       string `json:"alias"`       // 1-15 alphanumberic alias
	APIURL      string `json:"api_url"`     // e.g., http://147.52.17.33:5000
	PublicKey   string `json:"public_key"`  // Public key to verify authorization
	Description string `json:"description"` // 2-3 sentence description
	Status      string `json:"status"`      // Active/Inative status
	Timestamp   string `json:"timestamp"`   // Timestamp of last update in database
}
```



#### UUIDRecord

```
{
	UUID string `json:"uuid"`
}
```

