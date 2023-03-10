# V0 LocalNet <!-- omit in toc -->

- [Download Repos](#download-repos)
- [Prepare Playground](#prepare-playground)
  - [Setup E2E Stack Environment](#setup-e2e-stack-environment)
  - [Pocket E2E Stack](#pocket-e2e-stack)
    - [Prepare Stack](#prepare-stack)
    - [Run Stack](#run-stack)
    - [Cleanup Stack](#cleanup-stack)

## Download Repos

```bash
git clone git@github.com:pokt-network/pocket-core.git
git clone git@github.com:pokt-network/tendermint.git
git clone git@github.com:pokt-network/tx-bot.git
git clone git@github.com:pokt-foundation/pocket-e2e-stack.git
git -C pocket-e2e-stack submodule update --init --recursive
```

You should have all the repositories in a single directory like so:

```bash
$ tree -L 1
.
├── doc.md
├── pocket-core
├── pocket-e2e-stack
├── tendermint
└── tx-bot
```

Identify the full path directory where these repositories are stored as `POCKET_CORE_REPOS_PATH`.

## Prepare Playground

### Setup E2E Stack Environment

1. Copy the template env variables

```bash
cd pocket-e2e-stack
cp .env.template .env
cp .playground.env.example .playground.env

source .env
source .playground.env
```

2. Update `POCKET_CORE_REPOS_PATH` in `.env` to be the full path reflecting the path where you downloaded all the repos.
3. Update `POCKET_CORE_REPOS_PATH` in `.playground.env` to be the full path reflecting the path where you downloaded all the repos

### Pocket E2E Stack

```bash
cd cd pocket-e2e-stack
```

#### Prepare Stack

```bash
./bin/pokt-net/playground.sh scaffold
```

#### Run Stack

```bash
./bin/pokt-net/playground.sh up
```

#### Cleanup Stack

```bash
./bin/pokt-net/playground.sh cleanup
```

<!-- ```bash
sudo -- sh -c "echo 127.0.0.1 monitoring.dev.pokt >> /etc/hosts"
``` -->
