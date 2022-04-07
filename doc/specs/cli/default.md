# Default Namespace

## Global Flags

```text
pocket [--datadir] [--node] [--remoteCLIURL] [--persistent_peers] [--seeds] [--madvdontneed]
```

Denotes default namespace with global flags to be used.

Options:

* `--datadir`: The data directory where the configuration files for this node are specified.
* `--node`: Takes a remote endpoint in the form ://:.
* `--remoteCLIURL`: Takes a remote endpoint in the form of :// \(uses RPC Port\).
* `--persistent_peers`: A comma separated list of PeerURLs: '&lt;ID&gt;@:,&lt;ID2&gt;@:...&lt;IDn&gt;@:'.
* `--seeds`: A comma separated list of PeerURLs: '&lt;ID&gt;@:,&lt;ID2&gt;@:...&lt;IDn&gt;@:'.
* `--madvdontneed`: If enabled, run with GODEBUG=madvdontneed=1, --madvdontneed=\(true \| false\).

## Start Pocket Core

```text
pocket start [--simulateRelay=(true | false)] [--keybase=(true | false)] [--mainnet=(true | false)] [--testnet=(true | false)] [--profileApp=(true | false)]
```

Starts the Pocket Node, picks up the config from the assigned `<datadir>`.

Options:

* `--simulateRelay`: Would you like to be able to test your relays.
* `--keybase`: Run with keybase, if disabled allows you to stake for the current validator only. providing a keybase is
  still neccesary for staking for apps & sending transactions
* `--mainnet`: Run with mainnet genesis
* `--testnet`: Run with testnet genesis
* `--profileApp`: bool exposes cpu & memory profiling
* `--useCache`: If added, runs with a cache for the IAVL store, which trades increases RAM usage and reduces CPU usage
  in consensus operations.

## Stop Pocket Core

```text
pocket stop
```

Stops the Pocket Node.

## Reset Pocket Core

```text
pocket reset
```

Reset the Pocket node. Deletes the following files/folders:

* .pocket/data
* priv\_val\_key
* priv\_val\_state
* node\_keys

## Show CLI Help

```text
pocket help
```

