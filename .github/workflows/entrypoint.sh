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

# Create work dir
spawn sh -c "mkdir -p /home/app/.pocket/config"
expect eof

# Pull variables from env if set
set genesis ""
catch {set genesis $env(POCKET_CORE_GENESIS)}

set chains ""
catch {set chains $env(POCKET_CORE_CHAINS)}

set config ""
catch {set config $env(POCKET_CORE_CONFIG)}

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

# If key isn't passed in, start the node
if { $env(POCKET_CORE_KEY) eq "" }  {
    log_user 0
    spawn sh -c "$command"
    send -- "$env(POCKET_CORE_PASSPHRASE)\n"
    log_user 1
} else {
# If key is passed in, load it into the local accounts
    log_user 0
    spawn pocket accounts import-raw $env(POCKET_CORE_KEY)
    sleep 1
    send -- "$env(POCKET_CORE_PASSPHRASE)\n"
    expect eof
    spawn sh -c "pocket accounts set-validator `pocket accounts list | cut -d' ' -f2- `"
    sleep 1
    send -- "$env(POCKET_CORE_PASSPHRASE)\n"
    expect eof
    log_user 1
    spawn sh -c "$command"
}

expect eof
exit
