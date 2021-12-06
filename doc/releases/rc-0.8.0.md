# RC-0.8.0

### Important Release Notes

1) **SessionDB may be deleted after this release**

* `rm -rf <datadir>/session.db`

2) **Force Unstake no longer burns tokens**
   (If Validator/Servicer falls below minimum stake)
- Send the Validator/Servicer to jail
- Validator may not unjail unless they have more than the minimum stake tokens
- Jailed Validators may unstake/be-unstaked (will wait the unstaking period)
- Coins are returned from unstaked validators to the output account

3) **BlockTxs & AccountTxs returns Page_Total & Total_Txs**
- total_count is deprecated in rpc
- page_total returns total count in the page
- total_txs returns total possible in query
