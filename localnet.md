# V0 LocalNet <!-- omit in toc -->

<!-- https://docs.google.com/presentation/d/1mk0XogopENCI_4WXXvSYm1_DG8EhRLIpwpZQNIA5vqM/edit#slide=id.p -->

- [Prerequisites](#prerequisites)
  - [Dependencies](#dependencies)
  - [Download Repos](#download-repos)
  - [Verify Directory Structure](#verify-directory-structure)
- [Prepare Playground](#prepare-playground)
  - [Setup E2E Stack Environment](#setup-e2e-stack-environment)
  - [Pocket E2E Stack](#pocket-e2e-stack)
    - [Prepare Stack](#prepare-stack)
    - [Run Stack](#run-stack)
    - [Watch the network](#watch-the-network)
    - [Cleanup Stack](#cleanup-stack)
  - [Making configuration changes](#making-configuration-changes)
- [Tx Bot](#tx-bot)
- [TODO](#todo)

## Prerequisites

### Dependencies

This repository is intended to be deprecated within the next 12 months as of writing this document. It is only intended for experienced Pocket developers.

Newcomer dependencies are not supported in detail.

### Download Repos

```bash
mkdir v0_localnet
cd v0_localnet

git clone git@github.com:pokt-network/pocket-core.git
git clone git@github.com:pokt-network/tendermint.git
git clone git@github.com:pokt-network/tx-bot.git
git clone git@github.com:pokt-foundation/pocket-e2e-stack.git
git -C pocket-e2e-stack submodule update --init --recursive
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

## Prepare Playground

### Setup E2E Stack Environment

1. Copy the template env variables

```bash
cd pocket-e2e-stack
cp .env.template .env
cp .playground.env.example .playground.env
```

2. Update `POCKET_CORE_REPOS_PATH` in `.env` to be the full path reflecting the path where you downloaded all the repos.
3. Update `POCKET_CORE_REPOS_PATH` in `.playground.env` to be the full path reflecting the path where you downloaded all the repos
4. Update `ETH_ALTRUIST` and `POLY_ALTRUIST` appropriately.
5. Source the env variables

```bash
source .env
source .playground.env
```

### Pocket E2E Stack

```bash
cd pocket-e2e-stack
```

#### Prepare Stack

```bash
./bin/pokt-net/playground.sh scaffold
```

#### Run Stack

```bash
./bin/pokt-net/playground.sh up
```

**IMPORTANT: Inspect the logs in case something looks abnormal or an error occurred.**

#### Watch the network

```bash
watch -n 5 "curl -s -X 'POST' 'http://localhost:8084/v1/query/height'
```

The full RPC spec can be found [here](https://editor.swagger.io/?url=https://raw.githubusercontent.com/pokt-network/pocket-core/staging/doc/specs/rpc-spec.yaml).

**IMPORTANT: It might take up to 5 minutes for the first height to start incrementing while the node prepares.**

#### Cleanup Stack

```bash
./bin/pokt-net/playground.sh cleanup
```

**IMPORTANT: You will need to run `cleanup` every time you create a new setup.**

### Making configuration changes

- We are using the `homogenous` network is in `pokt-net`
- You can update configs by modifying `pocket-e2e-stack/.playground.env`
- You can further update configs by modifying `pocket-e2e-stack/playground/templates/config.template.json`
- Make sure to run `clean, scaffold & up` after changing configs

## Tx Bot

```bash
cd tx-bot
make start
# Make a selection
```

## TODO

- [ ] How to enable the monitoring stack?
- [ ] How to setup a heterogenous network?
