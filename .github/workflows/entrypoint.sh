#!/usr/bin/expect

# Send `pocket stop` when interrupted to prevent corruption
proc graceful_exit {} {
    send_user "Gracefully exiting Pocket...\n"
    spawn sh -c "pocket stop"
}

trap graceful_exit {SIGINT SIGTERM}

# Command to run
set command $argv
set timeout -1

# Pull variables from env if set
set defaultDatadir "/home/app/.pocket"
set datadir $defaultDatadir
catch {set datadir $env(POCKET_CORE_DATADIR)}
set datadirParam "--datadir=${datadir}"

set genesis ""
catch {set genesis $env(POCKET_CORE_GENESIS)}

set chains ""
catch {set chains $env(POCKET_CORE_CHAINS)}

set config ""
catch {set config $env(POCKET_CORE_CONFIG)}

set core_key ""
catch {set core_key $env(POCKET_CORE_KEY)}

set core_passphrase ""
catch {set core_passphrase $env(POCKET_CORE_PASSPHRASE)}

# Create dynamic config files
if {$genesis != ""} {
    set genesis_file [open "{$datadir}/config/genesis.json" w]
    puts $genesis_file $genesis
    close $genesis_file
    send_user "GENESIS loaded from env\n"
}
if {$chains != ""} {
    set chains_file [open "${datadir}/config/chains.json" w]
    puts $chains_file $chains
    close $chains_file
    send_user "CHAINS loaded from env\n"
}
if {$config != ""} {
    set config_file [open "${datadir}/config/config.json" w]
    puts $config_file $config
    close $config_file
    send_user "CONFIG loaded from env\n"
}

if {[regexp -nocase "datadir=" $command] && ![regexp -nocase "datadir=${datadir}" $command]} {
  send_user "WARNING: --datadir provided with a different path than the one defined in the Dockerfile. This could lead to errors using pocket CLI commands.\n"
} elseif {![regexp -nocase "datadir=" $command]} {
  send_user "INFO: param --datadir was not provided; attaching ${datadirParam} on every command\n"
  set command "${command} ${datadirParam}"
}

if {![regexp -nocase "datadir=${defaultDatadir}" $command]} {
	send_user "WARNING: --datadir is not the default one
Please review:
1. Mount your config folder to the same path you specify in --datadir
2. Review your config.json to ensure the following configs match the value of --datadir
  2.1. tendermint_config.RootDir
  2.2. RPC.RootDir
  2.3. P2P.RootDir
  2.4. Mempool.RootDir
  2.5. Consensus.RootDir
  2.6. pocket_config.RootDir
"
}

proc check_passphrase {str} {
  send_user "checking passphrase: ${str}\n"
  if {$str == ""} {
    send_user "missing POCKET_CORE_PASSPHRASE environment variable"
    exit 1
  }
}

# if not --keybase=false
# e.g. "pocket start --keybase=false --mainnet --datadir=/home/app/.pocket/"
if {[regexp -nocase "keybase=false" $command]} {
	spawn sh -c "$command"
} elseif { $core_key eq "" }  {
	  # If key isn't passed in, start the node
	  check_passphrase $core_passphrase
    log_user 0
    spawn sh -c "$command"
    send -- "$core_passphrase\n"
    log_user 1
} else {
    check_passphrase $core_passphrase
    # If key is passed in, load it into the local accounts
    log_user 0
    spawn pocket accounts import-raw $datadirParam $core_key
    sleep 1
    send -- "$core_passphrase\n"
    expect eof
    spawn sh -c "pocket accounts set-validator ${datadirParam} `pocket accounts list ${datadirParam} | cut -d' ' -f2- `"
    sleep 1
    send -- "$core_passphrase\n"
    expect eof
    log_user 1
    spawn sh -c "$command"
}

expect eof
exit
