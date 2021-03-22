# Default Namespace
The default namespace contains functions that are pertinent to the execution of the Pocket Node.
- `pocket [--datadir] [--node] [--remoteCLIURL] [--persistent_peers] [--seeds] [--madvdontneed] `
> Denotes default namespace with global flags to be used
>
> Options:
> - `--datadir`: The data directory where the configuration files for this node are specified.
> - `--node`: Takes a remote endpoint in the form <protocol>://<host>:<port>.
> - `--remoteCLIURL`: Takes a remote endpoint in the form of <protocol>://<host> (uses RPC Port).
> - `--persistent_peers`: A comma separated list of PeerURLs: '<ID>@<IP>:<PORT>,<ID2>@<IP2>:<PORT>...<IDn>@<IPn>:<PORT>'.
> - `--seeds`: A comma separated list of PeerURLs: '<ID>@<IP>:<PORT>,<ID2>@<IP2>:<PORT>...<IDn>@<IPn>:<PORT>'.
> - `--madvdontneed`: If enabled, run with GODEBUG=madvdontneed=1, --madvdontneed=(true | false).


- `pocket start  [--simulateRelay=(true | false)] [--keybase=(true | false)] [--mainnet=(true | false)] [--testnet=(true | false)] [--profileApp=(true | false)]`
> Starts the Pocket Node, picks up the config from the assigned `<datadir>`.
>
> Options:
> - `--simulateRelay`: Would you like to be able to test your relays.
> - `--keybase`: Run with keybase, if disabled allows you to stake for the current validator only. providing a keybase is still neccesary for staking for apps & sending transactions
> - `--mainnet`: Run with mainnet genesis
> - `--testnet`: Run with testnet genesis
> - `--profileApp`: bool exposes cpu & memory profiling"

- `pocket reset`
> Reset the Pocket node.
> Deletes the following files / folders:
> - .pocket/data
> - priv_val_key
> - priv_val_state
> - node_keys