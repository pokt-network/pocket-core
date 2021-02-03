# Pocket Core CLI Interface Specification
## Version RC-0.6.0.0

### Overview
This document serves as a specification for the Command Line Interface of the Pocket Core application. There's no protocol verification for these commands because they map closely to protocol functions.

### Namespaces
The CLI will contain multiple namespaces listed below:

- [Default Namespace](./default.md): These functions will be called when the namespace is blank
- [Accounts](./accounts.md): Contains all the calls pertinent to accounts and their local storage.
- [Nodes](./nodes.md): Contains all the functions for Node upkeep.
- [Apps](./apps.md): Contains all the functions for app upkeep.
- [Query](./query.md): All queries to the world state are contained in this call.
- [Util](./util.md): Contains useful operations

### CLI Functions Format
Each CLI Function will be in the following format:

- Binary Name: The name of the binary for Pocket Core, for example: `pocket`
- Global Options: any number of global options, for example: `pocket --datadir /tmp/.pocket`
- Namespace: The namespace of the function, or blank for the default namespace: `accounts`
- Function Name: The name of the actual function to be called: `create`
- Function Options: Options that modify behaviour of the function `pocket query nodes --staking_status unstaking`
- (Optional): Space separated function arguments, e.g.: `pocket query nodes --staking_status unstaking <height>`
