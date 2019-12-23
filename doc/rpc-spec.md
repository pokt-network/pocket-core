# Pocket Core RPC Specification
## Version 0.0.1

### Overview
This document serves as a specification for the RPC of the Pocket Core application. There's no protocol verification for these commands, however because they map closely to protocol functions.

### Namespaces
The RPC will contain multiple namespaces listed below:

- Client: Contains all the calls pertinent to pocket core clients (relay and dispatch)
- Query: All queries to the world state are contained in this call.

### RPC Functions Format
Each RPC Function will be in the following format:

- Version: The version for Pocket Core rpc, for example: `/v1`
- Namespace: The namespace of the function: `query`
- Function Name: The name of the actual function to be called: `tx`

ALl RPC calls in Pocket Core will be HTTP `POST`

### Client Namespace
The default namespace contains functions that are pertinent to the execution of the Pocket Node.

// todo

### Query Namespace
The `query` namespace handles all queries to the current world state built on the Pocket node.

// todo
