#!/usr/bin/expect

# Command to run
set command $argv
set timeout -1

# Send `pocket stop` when interrupted to prevent corruption
proc graceful_exit {} {
    send_user "Gracefully exiting Pocket...\n"
    spawn sh -c "pocket stop"
}

proc graceful_mesh_exit {pid} {
    send_user "Gracefully exiting Pocket...\n"
    exec kill -SIGTERM $pid
}


# Pull variables from env if set
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
    set genesis_file [open /home/app/.pocket/config/genesis.json w]
    puts $genesis_file $genesis
    close $genesis_file
    send_user "GENESIS loaded from env\n"
}
if {$chains != ""} {
    set chains_file [open /home/app/.pocket/config/chains.json w]
    puts $chains_file $chains
    close $chains_file
    send_user "CHAINS loaded from env\n"
}
if {$config != ""} {
    set config_file [open /home/app/.pocket/config/config.json w]
    puts $config_file $config
    close $config_file
    send_user "CONFIG loaded from env\n"
}

# if not --keybase=false
# e.g. "pocket start --keybase=false --mainnet --datadir=/home/app/.pocket/"
if {[regexp -nocase "keybase=false" $command]} {
	spawn sh -c "$command"
} elseif { $core_key eq "" }  {
	  # If key isn't passed in, start the node
    log_user 0
    spawn sh -c "$command"
    send -- "$core_passphrase\n"
    log_user 1
} else {
    # If key is passed in, load it into the local accounts
    log_user 0
    spawn pocket accounts import-raw $core_key
    sleep 1
    send -- "$core_passphrase\n"
    expect eof
    spawn sh -c "pocket accounts set-validator `pocket accounts list | cut -d' ' -f2- `"
    sleep 1
    send -- "$core_passphrase\n"
    expect eof
    log_user 1
    spawn sh -c "$command"
}

set pid [exp_pid]
if {![regexp -nocase "start-mesh" $command]} {
  trap graceful_exit {SIGINT SIGTERM}
} else {
  trap "graceful_mesh_exit $pid" {SIGINT SIGTERM}
}
expect eof
exit
