# Journey to changing a node's output address

A staked node's output address can be changed only when a transaction message
`MsgStake` is signed by the owner of the node's current output address.  The
following charts illustrate how a `MsgStake` is handled by the key functions
and the transaction is accepted or rejected by the blockchain.

## `baseApp/runTx`

The first validation occurs in `x/auth/ValidateTransaction`, where we make sure
an incoming `MsgStake` is signed by any of:

- The new output address specified in the transaction itself
- Node's operator address
- Node's current output address

```mermaid
flowchart TD
  A(baseApp/runTx) --> B(x/auth/ValidateTransaction)
	B --> C{Both NCUST and OEDIT have been enabled?}
	C -- Yes --> D{Signed by any of:\nthe new output address\nthe operator address\nthe current output address}
  D -- Yes --> F([See baseApp/runMsg])
  D -- No --> G([sdk:4])
	C -- No --> E{Signed by any of:\nthe new output address\nthe operator address}
  E -- Yes --> F
  E -- No --> G
```

## `baseApp/runMsg`

After `x/auth/ValidateTransaction`, the transaction is validated in
`x/nodes/keeper/ValidateValidatorStaking`, which validates all stake
transactions.  If the target operator is already staked, this function basically
calls `x/nodes/keeper/ValidateValidatorMsgSigner` twice, one for the new
parameters of the node, and one for the current parameters of the node.  If
the transaction is to edit the output address, we skip the first call.
If any calls to `x/nodes/keeper/ValidateValidatorMsgSigner` fails, the
transaction is rejected with the error code `pos:125`.

```mermaid
flowchart TD
  A(baseApp/runMsg) --> B(x/nodes/handleStake)
  B --> C(x/nodes/keeper/ValidateValidatorStaking)
  C --> D{Both NCUST and OEDIT have been enabled?\nChanging the output address?\nSigned by the current output address?}
  D -- All yes --> E(x/nodes/keeper/ValidateValidatorMsgSigner\nwith the current state)
  E --> F{Signed by any of:\nthe operator\nthe current output address}
  F -- Yes --> G([See x/nodes/keeper/ValidateEditStake])
  F -- No --> H([pos:125])
  D -- Any no --> I(x/nodes/keeper/ValidateValidatorMsgSigner\nwith the new state)
  I --> J{Signed by any of:\nthe operator\nthe new output address}
  J -- Yes --> E
  J -- No --> H
```

## `x/nodes/keeper/ValidateEditStake`

After `x/nodes/keeper/ValidateValidatorMsgSigner`,
`x/nodes/keeper/ValidateValidatorStaking` calls the final validation function
`x/nodes/keeper/ValidateEditStake`, which validates all editstake transactions.
The chart below focuses on the case of changing the output address.

```mermaid
flowchart TD
	A(x/nodes/keeper/ValidateEditStake) -- We're changing the output address --> B{Both NCUST and OEDIT have been enabled?}
	B -- Yes --> C{Signed by the current output address?}
	C -- Yes --> D([See x/nodes/keeper/StakeValidator])
	C -- No --> E([pos:127])
	B -- No --> F([pos:124])
```

## `x/nodes/keeper/StakeValidator`

After all validations above pass, we just go to
`x/nodes/keeper/EditStakeValidator` where we change the world state.

```mermaid
flowchart TD
  A(x/nodes/keeper/StakeValidator) --> B(x/nodes/keeper/EditStakeValidator)
  B --> C([Happy ending!])
```
