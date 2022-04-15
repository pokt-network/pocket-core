# RC-0.8.0

### Important Release Notes

1) **Non-Custodial and the 'Output' Address**

```
// Rules
    The Operator Address is the only valid signer for blocks & relays
    The Output Address is where reward and staked funds are directed
    The Operator and the Output Address must be set when Staking
    Neither the Output nor the Operator Address may be edited once set
    Both the Operator and the Output Address are valid signers for 'Node Transactions'
      - (Stake, EditStake, Unstake, Unjail)
    Only the Operator may sign Claim Or Proof Transactions

// Legacy Migration
    Output is empty for legacy nodes and the current Address is the Operator
    If Output is empty -> the Operator may set it (once)
    If Output is empty -> the Operator is treated as the Output
```

2) **SessionDB may be deleted after this release**

* `rm -rf <datadir>/session.db`

3) **Force Unstake no longer burns tokens**
   (If Validator/Servicer falls below minimum stake)

- Send the Validator/Servicer to jail
- Validator may not unjail unless they have more than the minimum stake tokens
- Jailed Validators may unstake/be-unstaked (will wait the unstaking period)
- Coins are returned from unstaked validators to the output account

4) **BlockTxs & AccountTxs returns Page_Total & Total_Txs**

- total_count is deprecated in rpc
- page_count returns total count in the page
- total_txs returns total possible in query

5) **Chains.json is now automatically refreshed**

- Every minute based on chains.json file

6) **Tendermint Evidence is double-checked before proposal block creation and evidence receive**

- Changes ported from Tendermint PR [5574](https://github.com/tendermint/tendermint/pull/5574)
