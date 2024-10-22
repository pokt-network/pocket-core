# [WIP] Chain Halt Rollback Recovery Guide <!-- omit in toc -->

:::warning

This rollback guide is a WIP.

Is is a port of this [notion doc](https://www.notion.so/Recovery-guide-641ec21aead74cae806166f5f9e61394)
by @msmania while managing a chain halt on a previous test and should only
be treated as a reference, not a definitive general purpose guide.

:::

- [Issue description](#issue-description)
- [Hotfix](#hotfix)
- [How to check if my data is correct or not](#how-to-check-if-my-data-is-correct-or-not)
- [Recovery steps](#recovery-steps)

## Issue description

[The test scenario](https://www.notion.so/RC-0-11-0-Release-Plan-a76509e73f854d0b8ca91ea62f52ca9e?pvs=21) set a new reward
delegator set on `128746` to the node `5b18a8c268ffcbf0530f61f70e3ee14f064bdf0f`, and updated it with another delegator set on `128747` as below.

```text
Address:                   5b18a8c268ffcbf0530f61f70e3ee14f064bdf0f
Public Key:                0b787e54e66b3db3a3396c2322d8314287e08990f7696c490153b078bada7e94
Jailed:                    false
Status:                    Staked
Tokens:                    18000000000
ServiceUrl:                https://poktt1698102386.c0d3r.org:443
Chains:                    [0001 0002 005A 005B 005C 005D]
Unstaking Completion Time: 0001-01-01 00:00:00 +0000 UTC
Output Address:            42846261e1798fc08e1dfd97325af7b280f815b0
Reward Delegators:         {"54751ae3431c015a6e24d711c9d1ed4e5a276479":20,"8147ed5182da6e7dea33f36d78db6327f9df6ba0":10}
```

Since the marshaling order of a map is not deterministic, the `RewardDelegators` field is considered as either `{"54751ae3431c015a6e24d711c9d1ed4e5a276479": 20, "8147ed5182da6e7dea33f36d78db6327f9df6ba0": 10}` or `{"8147ed5182da6e7dea33f36d78db6327f9df6ba0": 10, "54751ae3431c015a6e24d711c9d1ed4e5a276479": 20}`and this may fork the world state into two versions depending on the order of its fields, `54751ae3` comes first or `8147ed51`comes first.

During the upgrade, the network was stuck in round 13 on 128748. The problem is the proposer (validator1; 77E608D8AE4CD7B812F122DC82537E79DD3565CB) proposed a block on top of block 128747 where apphash was 624553DE, but the majority of the validators have a block 128747 where apphash is c407eb25.

All nodes in the network, including non-validators and seeds, can be stuck due to this bug. It explains why we had peer connection issues so often during upgrade.

## Hotfix

The fix is to sort the `RewardDelegators` field when marshaling. We decided to use [protobuf’s standard marshaler](https://pkg.go.dev/github.com/gogo/protobuf/plugin/marshalto?utm_source=godoc) to marshal a validator, which adopts the reverse alphabetical order.

https://github.com/pokt-network/pocket-core/pull/1591

## How to check if my data is correct or not

During the upgrade, the world state of Testnet was forked into two versions. With the fix, on block 128746, the delegator address `8147ed51` should have come first, which means AppHash 84cbe5d0 is correct. In the actual blockchain, however, the other incorrect one was chosen. Therefore we need to roll back Testnet.

I created a small tool https://github.com/msmania/pocket-appdb-parser.git to print the latest AppHash in application.db. Here are the steps to see the state of your node.

1.  Clone and build a tool

    ```bash
    git clone https://github.com/msmania/pocket-appdb-parser.git
    cd pocket-appdb-parser
    go build -o pocket-appdb-parser .
    ```

1.  Stop the pocket and run the following command.
    You cannot run it when pocket is running. GoLevelDB does not allow access from multiple processes.

    ```bash
    ./pocket-appdb-parser <path to application.db>
    ```

1.  If a node is on 128747, the output will be either

    1. 128747: c407eb25 (wrong)

       ```text
       128747: c407eb25c5e5192a67514132838217741011991625e410e7f16233dad7d8705c
               main:128747 ea6e1849b4bf587027401a8f105901b134f8be94d7ecab2ea170c7b5d96e4cf5
         pocketcore:128747 85c3aae78d51de66ca88a2f7d61c88a6a076013f2ca385eeee820e5d1bca2859
               auth:128747 5faa7669ef6aa9d393a584e03041c42772cf43ccf861326ca2a70544c97ca844
                pos:128747 87bc3da27011f645ae3e856e44cbec2a1691ebe2b0b3f98464d5c184a57265ac
        application:128747 2db9f5e3c2aa8fa064eb284beee5189e66344e48cb06f51b53d984af0ec2dbe7
                gov:128747
             params:128747 3c77022618e6d32441a3de1d22092f3d2e5a4221ea569cca7a7adfa22d08131c
       ```

    2. 128747: 624553de (wrong)

       ```text
       128747: 624553de014c6546f56167b47f4d92e46f72f18ae2e08e3ae254981f9914c95e
                gov:128747
        application:128747 2db9f5e3c2aa8fa064eb284beee5189e66344e48cb06f51b53d984af0ec2dbe7
             params:128747 3c77022618e6d32441a3de1d22092f3d2e5a4221ea569cca7a7adfa22d08131c
                pos:128747 56114cdc51d8217075255106c4f067e1c117d1cc3876d3da3fbb9e1a2d6689f7
               auth:128747 5faa7669ef6aa9d393a584e03041c42772cf43ccf861326ca2a70544c97ca844
         pocketcore:128747 85c3aae78d51de66ca88a2f7d61c88a6a076013f2ca385eeee820e5d1bca2859
               main:128747 ea6e1849b4bf587027401a8f105901b134f8be94d7ecab2ea170c7b5d96e4cf5
       ```

1.  If a node is on 128746, the output will be either

    1. 128746: 84cbe5d0 (correct; no need to resync)

       ```
       128746: 84cbe5d012fbd52c34775351d56762afb738888d33cfb80c96d84900ed3f3a82
         pocketcore:128746 3ecb0f4b97339e0918a57902b48b46ddbb1e5f221b905cbb09349e830ce64f21
             params:128746 3c77022618e6d32441a3de1d22092f3d2e5a4221ea569cca7a7adfa22d08131c
                pos:128746 9b6de08f1d3f4eb724fe3f9ee04dd7116103060da8755088fba8582914fc4e67
        application:128746 2db9f5e3c2aa8fa064eb284beee5189e66344e48cb06f51b53d984af0ec2dbe7
                gov:128746
               auth:128746 3c5fea92e0ec2846a21acc4faee6e78d7b6ff4ef9303c60b5152ebcee1216b3d
               main:128746 ea6e1849b4bf587027401a8f105901b134f8be94d7ecab2ea170c7b5d96e4cf5
       ```

    2. 128746: dca3d2fb (wrong)

    ```
       128746: dca3d2fb8848e6915b3745bd4db22003cfce09659436b39d289e9e5d51cabbc5
         pocketcore:128746 3ecb0f4b97339e0918a57902b48b46ddbb1e5f221b905cbb09349e830ce64f21
                gov:128746
               auth:128746 3c5fea92e0ec2846a21acc4faee6e78d7b6ff4ef9303c60b5152ebcee1216b3d
                pos:128746 a8b7de0316b42f13e0e39d872301ad7139aec20ec84269e354e9d65866784c7c
        application:128746 2db9f5e3c2aa8fa064eb284beee5189e66344e48cb06f51b53d984af0ec2dbe7
               main:128746 ea6e1849b4bf587027401a8f105901b134f8be94d7ecab2ea170c7b5d96e4cf5
             params:128746 3c77022618e6d32441a3de1d22092f3d2e5a4221ea569cca7a7adfa22d08131c
    ```

## Recovery steps

Everyone needs to run the patched version on **all nodes, not only validators but also non-validators like servicers and seeds**.

The commands vary depending on your environment.

1. Stop the pocket

   ```bash
   sudo systemctl stop pocket
   ```

2. Upgrade the binary

   ```bash
   cd <path to pocket_core repo>
   git pull origin staging
   go build -o <path to pocket> app/cmd/pocket_core/main.go
   ```

3. Apply the snapshot https://link.storjshare.io/s/jxzmjzjz4dzkalgwxlyxzzuzb6sa/pocket-snapshots/pokt-testnet@128717.tar.gz to all managed nodes
   1. Managed nodes (seeds and validators) need to be isolated..?
4. Start the pocket and pray

   ```bash
   sudo systemctl start pocket
   ```
