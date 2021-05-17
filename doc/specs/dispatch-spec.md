# Dispatch

```text
Dispatch Request Protocol
-------------------------
Node selection begins at the beginning of the session.

Nodes are cross referenced with the current world state 
info to see if they are eligible (not jailed).

Node payment/challenge verification looks at the node selection 
at the beginning of the session and state information from the 
last block of the session to see who is eligible for payment.

Contract
------------------------
Unjail time >= session time
Jailed for missed blocks < blocksPerSession

Notes
------------------------
You are only eligible for sessions if you are not jailed

If you are jailed during a session you do not get paid and will 
not get selected for the next session

Replacement nodes, provide refreshed service to the clients
```

