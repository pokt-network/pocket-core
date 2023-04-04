# V0 LocalNet <!-- omit in toc -->

<!-- https://docs.google.com/presentation/d/1mk0XogopENCI_4WXXvSYm1_DG8EhRLIpwpZQNIA5vqM/edit#slide=id.p -->

- [Prerequisites](#prerequisites)
  - [Dependencies](#dependencies)
  - [Download Repos](#download-repos)
- [Playground](#playground)
  - [Setup E2E Stack Environment](#setup-e2e-stack-environment)
  - [Pocket E2E Stack](#pocket-e2e-stack)
    - [1. Prepare Stack](#1-prepare-stack)
    - [2. Run Stack](#2-run-stack)
    - [3. Watch the network](#3-watch-the-network)
    - [4. Cleanup Stack](#4-cleanup-stack)
    - [5. Cleaning containers](#5-cleaning-containers)
  - [Tx Bot](#tx-bot)
  - [Configuration changes](#configuration-changes)
- [TODO](#todo)

## Prerequisites

### Dependencies

This repository is intended to be deprecated within the next 12 months as of writing this document. It is only intended for experienced Pocket developers.

It is intentionally non-exhaustive so newcomers may find it difficult to follow as it is not the intended audience.

### Download Repos

```bash
mkdir v0_localnet
cd v0_localnet

git clone git@github.com:pokt-network/pocket-core.git
git clone git@github.com:pokt-network/tendermint.git
git clone git@github.com:pokt-network/tx-bot.git
git clone --recurse-submodules git@github.com:pokt-network/pocket-e2e-stack.git
```

Run `pwd` and identify the current path as it will be referenced to as `POCKET_CORE_REPOS_PATH` below

## Playground

### Setup E2E Stack Environment

1. Copy the template env variables

```bash
cd pocket-e2e-stack
cp .env.template .env
cp .playground.env.example .playground.env
```

2. Update `POCKET_CORE_REPOS_PATH` in the appropriate `.env` variables:

```bash
# Update `POCKET_CORE_REPOS_PATH` in `.env` to be the full path reflecting where you downloaded all the repos.
sed -i 's|^POCKET_CORE_REPOS_PATH=.*$|POCKET_CORE_REPOS_PATH='$PWD'|' pocket-e2e-stack/.env
# Update `POCKET_CORE_REPOS_PATH` in `.playground.env` to be the full path reflecting the path where you downloaded all the
sed -i 's|^POCKET_CORE_REPOS_PATH=.*$|POCKET_CORE_REPOS_PATH='$PWD'|' pocket-e2e-stack/.playground.env
```

3. Update the Ethereum & Polygon altruist nodes

```bash
# Update `ETH_ALTRUIST` and `POLY_ALTRUIST` appropriately.
sed -i 's|^# ETH_ALTRUIST=.*$|ETH_ALTRUIST=<Your or public Ethereum endpoint>|' pocket-e2e-stack/.env
sed -i 's|^# POLY_ALTRUIST=.*$|POLY_ALTRUIST=<Your or public Polygon endpoint>|' pocket-e2e-stack/.env
```

```bash
source .env
source .playground.env
```

### Pocket E2E Stack

```bash
cd pocket-e2e-stack
```

#### 1. Prepare Stack

```bash
./bin/pokt-net/playground.sh scaffold
```

#### 2. Run Stack

```bash
./bin/pokt-net/playground.sh up
```

**What am I running?**

- Four containerized `Validator` nodes
- A local version of the `Tendermint` library
- Hot reloading code in both `Tendermint` and `Pocket-Core`
- A bot that can be used to send transactions to the network
- Pocket RPC endpoints are exposed at 8082, 8083, 8084, and 8085/tcp on each node respectively
- Tendermint RPC endpoints are exposed at 26658, 26659, 26660, 26661/tcp respectively

https://user-images.githubusercontent.com/1892194/225497350-3817b262-f5d5-483b-bb8f-47fb0b614321.mov

**IMPORTANT: Inspect the logs in case something looks abnormal or an error occurred.**

#### 3. Watch the network

```bash
watch -n 5 "curl -s -X 'POST' 'http://localhost:8084/v1/query/height'"
```

The full RPC spec can be found [here](https://editor.swagger.io/?url=https://raw.githubusercontent.com/pokt-network/pocket-core/staging/doc/specs/rpc-spec.yaml).

**IMPORTANT: It might take up to 5 minutes for the first height to start incrementing while the node prepares.**

#### 4. Cleanup Stack

```bash
./bin/pokt-net/playground.sh cleanup
```

**IMPORTANT: You will need to run `cleanup` every time you create a new setup.**

#### 5. Cleaning containers

If you make any changes to the scripts or need a cleaner start, you'll need to remove the containers and docker image:

```bash
docker rm node1.dev.pokt node2.dev.pokt node3.dev.pokt node4.dev.pokt && docker rmi playground/pocket-core
```

### Tx Bot

```bash
cd tx-bot
make start
# Make a selection
```

### Configuration changes

- We are using the `homogenous` network is in `pokt-net`
- You can update configs by modifying `pocket-e2e-stack/.playground.env`
- You can further update configs by modifying `pocket-e2e-stack/playground/templates/config.template.json`
- Make sure to run `clean, scaffold & up` after changing configs

## TODO

- [ ] How to enable the monitoring stack?
- [ ] How to setup a heterogenous network?
