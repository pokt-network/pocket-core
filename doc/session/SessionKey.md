# Session Key
The `session key algorithm` is a publicly verifiable hash that is used by nodes within Pocket Network to derive the session participants.

The session key implementation is the following:
```
Hash(<Hash of Block N-1> + <Hash of Block N-2> + DeveloperID) 
N is the latest confirmed block.
DeveloperID is given by the client application.
```

DISCLAIMER: *This algorithm is far from finalized and subject to change*
