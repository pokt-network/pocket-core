# Pocket Core Metrics

Pocket Core provides 2 layers of metrics
- [Tendermint](https://docs.tendermint.com/master/nodes/metrics.html)
- Pocket

Both layers of metrics are exposed through prometheus on individual ports:
- Tendermint Default Port: `26656`
- Pocket Prometheus Default Port: `8083`

For Tendermint Prometheus info please refer to [this documentation](https://docs.tendermint.com/master/nodes/metrics.html)


Pocket Metrics work expose service metrics per hosted chain for the validator.
By default Pocket metrics are enabled.

| Name        | Type        | Tags        | Description |
| ----------- | ----------- | ----------- | ----------- |
| relay_count_for_ | Counter |  | The number of relays executed against a hosted blockchain  |
| challenge_count_for_ | Counter |  | The number of challenges executed against a hosted blockchain |
| err_count_for_ | Counter |  | The number of errors executed against a hosted blockchain  |
| avg_relay_time_for_ | Histogram |  | The average relay time in ms exevcuted agianst a hosted blockchain  |
| sessions_count_for | Counter |  | The number of unique sessions generated for a hosted blockchain |
| tokens_earned_for_ | Counter |  | The number of tokens earned in uPOKT for a hosted blockchain |
