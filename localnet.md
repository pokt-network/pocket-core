# V0 LocalNet <!-- omit in toc -->

<!-- https://docs.google.com/presentation/d/1mk0XogopENCI_4WXXvSYm1_DG8EhRLIpwpZQNIA5vqM/edit#slide=id.p -->

- [Prerequisites](#prerequisites)
  - [Dependencies](#dependencies)
  - [Download Repos](#download-repos)
  - [Verify Directory Structure](#verify-directory-structure)
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
git clone git@github.com:pokt-network/pocket-e2e-stack.git
git -C pocket-e2e-stack submodule update --init --recursive # This pulls in https://github.com/pokt-network/local_playground
```

### Verify Directory Structure

Run `tree`:

```bash
tree -L 1
```

You should have all the repositories in a single directory like so:

```bash
.
├── pocket-core
├── pocket-e2e-stack
├── tendermint
└── tx-bot
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

1. Update `POCKET_CORE_REPOS_PATH` in `.env` to be the full path reflecting where you downloaded all the repos.
2. Update `POCKET_CORE_REPOS_PATH` in `.playground.env` to be the full path reflecting the path where you downloaded all the repos
3. Update `ETH_ALTRUIST` and `POLY_ALTRUIST` appropriately.
4. Source the env variables

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

**IMPORTANT: Inspect the logs in case something looks abnormal or an error occurred.**

#### 3. Watch the network

```bash
watch -n 5 "curl -s -X 'POST' 'http://localhost:8084/v1/query/height'
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
