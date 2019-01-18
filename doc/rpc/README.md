# JSON RPC/REST
//TODO
## Endpoints
Default API Endpoints:

## Check In-Client Reference
Example: `````localhost:8546/v1/dispatch/serve`````

Returns: ```{ "endpoint": "/v1/dispatch/serve",
              	"method": "POST",
              	"params": [
              		"devid"
              	],
              	"returns": "DATA - Session ID",
              	"example": "curl --data {devid:1234}' http://localhost:8546/v1/dispatch/serve"
              }```

| Type  |URL      |
| :-----| :----|
| Client | http://localhost:8545 |
| Relay  | http://localhost:8546 |

## Client
You can start the client HTTP APIs with the `````--clientrpc````` flag

Change the default port (8545) and listing address (localhost) with the````--clientrpcport```` flag

## Relay
You can start the relay HTTP APIs with the `````--relayrpc````` flag

Change the default port (8545) and listing address (localhost) with the````--relayrpcport```` flag

## JSON RPC/REST API Structure
Client Index:  
````http://localhost:8545/v1/````  

Relay Index:  
````http://localhost:8546/v1/````  

Account:  
````http://localhost:8545/v1/account/````    
````http://localhost:8545/v1/account/active/````    
````http://localhost:8545/v1/account/balance/````    
````http://localhost:8545/v1/account/joined/````    
````http://localhost:8545/v1/account/karma/````  
````http://localhost:8545/v1/account/last_active/````  
````http://localhost:8545/v1/account/transaction_count/````  
````http://localhost:8545/v1/account/session_count/````  
````http://localhost:8545/v1/account/status/````  

Client:  
````http://localhost:8545/v1/client/````  
````http://localhost:8545/v1/client/id/````  
````http://localhost:8545/v1/client/version/````  
````http://localhost:8545/v1/client/syncing/````  

Networking:  
````http://localhost:8545/v1/network/````  
````http://localhost:8545/v1/network/id/````  
````http://localhost:8545/v1/network/peer_count/````  
````http://localhost:8545/v1/network/peer_list/````  
````http://localhost:8545/v1/network/peers/````  

Personal: (Requires Passphrase)  
````http://localhost:8545/v1/personal/````  
````http://localhost:8545/v1/personal/active/````  
````http://localhost:8545/v1/personal/list_accounts/````  
````http://localhost:8545/v1/personal/network/````  
````http://localhost:8545/v1/personal/network/enter/````  
````http://localhost:8545/v1/personal/network/exit/````  
````http://localhost:8545/v1/personal/primary_address/````  
````http://localhost:8545/v1/personal/send/````  
````http://localhost:8545/v1/personal/send/raw/````  
````http://localhost:8545/v1/personal/sign/````  
````http://localhost:8545/v1/personal/status/````  
````http://localhost:8545/v1/personal/stake/````  
````http://localhost:8545/v1/personal/stake/add/````  
````http://localhost:8545/v1/personal/stake/remove/````  

Pocket:  
````http://localhost:8545/v1/pocket/````  
````http://localhost:8545/v1/pocket/block/````  
````http://localhost:8545/v1/pocket/block/hash/````  
````http://localhost:8545/v1/pocket/block/hash/transaction_count/````  
````http://localhost:8545/v1/pocket/block/hash/transaction_receipt/````  
````http://localhost:8545/v1/pocket/block/number/````  
````http://localhost:8545/v1/pocket/block/number/transaction_count/````  
````http://localhost:8545/v1/pocket/block/number/transaction_receipt/````  
````http://localhost:8545/v1/pocket/version/````  

Relay:  
````http://localhost:8546/v1/relay/````  
````http://localhost:8546/v1/relay/read/````  
````http://localhost:8546/v1/relay/write/````  

Transaction:  
````http://localhost:8545/v1/transaction/````  
````http://localhost:8545/v1/transaction/hash/````  

# JSON RPC/REST Methods

/TODO
