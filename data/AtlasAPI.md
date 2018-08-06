# SParts Project Atlas Service API

[TOC]

## Atlas Service Overview

The importance and dependency on software has grown exponentially across most industries over the past decade. Software solutions, whether it is an application, library, container or an entire Linux runtime are comprised of some percentage of open source software. The dependency on open source software has grown even faster over the past five year. And although the benefits of using open source are rapidly being realized, so are the complexities and costs associated with its use. Tracking which open source components were used in a solution, when and by whom across the software supply chain is vital to its continued success. 

The SParts project provides a distributed ledger to facilitate the tracking of open source used in software solutions (parts) within various different industries such as automotive, aerospace, medical devices, industrial manufacturing, consumer electronics and the Internet of Things (IoT) more generally.  Each supply chain network maintains a ledger and each ledger is replicated across by multiple ledger **nodes**. An application utilizing  the ledger can access and submit transactions to any ledger node within a given network and be assured the transaction will be correctly propagated to all other network ledger nodes. 

One challenge an application  that utilizes the ledger faces is locating a ledger node to transact with. The Atlas service (spartshub.org) serves as a ledger node lookup service. It enables different applications using a SParts ledger to located and obtain the API address of one or more of the ledger nodes belong to a specific industry supply chain network. Although there is no requirement  that each ledger node register their API address with the Atlas service, if it does, it becomes much easier to discover. For example, if the automotive industry was using a SParts ledger which was replicated across 30 different ledger nodes (hosted by the respective  automotive companies and their suppliers) then each node could register their API address with the Atlas service. Any application wanting to access the ledger can simply query the Altas service to obtain a list of available nodes to choose from.

Each supply chain network is identified by an account name consisting of up to 50 alphanumeric letters and special characters, ''.", "_" and "-". Names are case-insensitive. For example, the Zephyr project uses the name "zephyr-parts-network". All ledgers nodes that support the zephyr-parts-chain would be registered under the network name: zephyr-parts-network.

## Example Requests

- [List of Networks](https://spartshub.org/atlas/api/v1/network_space)

- [List of Ledger Nodes ](https://spartshub.org/atlas/api/v1/network_node_list/zephyr-parts-network)



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

This call is used to register a ledger node and its API address for a existing network.

Fill out and send a **LedgerNodeRecord** . You must include a public key which will be used to encrypt messages to verify authenticity.

| Field        | Type   | Description                                |
| ------------ | ------ | ------------------------------------------ |
| uuid         | string | unique identifier -                        |
| name         | string | file or envelope name                      |
| network_name | string | The name of the network the node belongs   |
| api_url      | string | api url address (e.g., 147.11.176.111:818) |
| alias        | string | Alias (short name)                         |
| description  | string | Description of node                        |

An example single artifact request:

```
{
	uuid: "4122ac8d-01f4-4de2-69ed-16b7ebae812c",
	name: "Wind River Test Node 1",
	network_name: "sparts-test-network",
	api_url: "http://35.166.246.146:818",
	alias: "test-node-1",
	description: "The SParts project test node"
}
```

Example curl Request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"name":"Wind River Test Node 1", "UUID": "4122ac8d-01f4-4de2-69ed-16b7ebae812c", "network_name":"sparts-test-network", "api_url":"http://35.166.246.146:818", "alias":"WR-Test-Node-1", "description":"The SParts project test node"}'  https://spartshub.org/atlas/api/v1/ledger_node/register
```



**Potential Errors**:

- The requesting user does not have the appropriate access credentials to perform the add.
- One or more of the required fields UUID, network_name, ... are missing.
- The UUID is not in a valid format.
- network does not exist. 



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


