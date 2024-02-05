# Pocket Snapshotter <!-- omit in toc -->

_tl;dr Access Liquify's Pocket Snapshotter [here](https://github.com/liquify-validation/pocket-snapshotter)_

- [Public snapshots](#public-snapshots)
  - [Update Frequency](#update-frequency)
- [Mirrors](#mirrors)
  - [Download using CLI](#download-using-cli)
    - [Uncompressed](#uncompressed)
    - [Uncompressed Inline (slower but smaller disk footprint)](#uncompressed-inline-slower-but-smaller-disk-footprint)
    - [Compressed](#compressed)
    - [Compressed Inline (slower but smaller disk footprint)](#compressed-inline-slower-but-smaller-disk-footprint)
  - [Issues](#issues)

## Public snapshots

If you're looking for Pocket native blockchain data snapshots, they are provided by [Liquify LTD](https://www.liquify.io/) and can be viewed via the Explorer link: [File Explorer here](https://pocket-snapshot.liquify.com).

[https://pocket-snapshot.liquify.com](pocket-snapshot.liquify.com)

### Update Frequency

Snapshots are updated every **Monday at 00:00 UTC**. The snapshots are generated on the Master (UK) and then sent over to the US and JP regions.

## Mirrors

The pocket snapshot link above is a global endpoint which is available in 3 different regions (UK, US west, Japan). The individual regions can also be accessed on the following links.

- UK (Master):[pocket-snapshot-uk.liquify.com](https://pocket-snapshot-uk.liquify.com)
- US: [pocket-snapshot-us.liquify.com](https://pocket-snapshot-us.liquify.com)
- JP: [pocket-snapshot-jp.liquify.com](https://pocket-snapshot-jp.liquify.com)

_Note: If accessing the snapshots on Monday it may be best to use the UK (Master) endpoint since there will be a 4-12 hour delay in updating the slaves in the other regions._

### Download using CLI

The snapshot repos hold the last **3 weeks of snapshots**. The latest one being referenced by the file `latest.txt` and `latest_compressed.txt`.

To copy-paste the commands below, please update `POCKET_DATA_DIR` appropriately.

```bash
export POCKET_DATA_DIR=<absolute path to your data dir>
```

#### Uncompressed

```bash
wget -O latest.txt https://pocket-snapshot.liquify.com/files/latest.txt
latestFile=$(cat latest.txt)
aria2c -s6 -x6 "https://pocket-snapshot.liquify.com/files/$latestFile"
tar xvf "$latestFile" -C ${POCKET_DATA_DIR}
rm latest.txt
```

#### Uncompressed Inline (slower but smaller disk footprint)

The below snippet will download and extract the snapshot inline. This may be beneficial if you have constrained disk space and cannot afford to store both the temp archive and extracted datadir.

```bash
wget -O latest.txt https://pocket-snapshot.liquify.com/files/latest.txt
latestFile=$(cat latest.txt)
wget -c "https://pocket-snapshot.liquify.com/files/$latestFile" -O - | sudo tar -xv -C {POCKET_DATA_DIR}
rm latest.txt
```

#### Compressed

```bash
wget -O latest.txt https://pocket-snapshot.liquify.com/files/latest_compressed.txt
latestFile=$(cat latest.txt)
aria2c -s6 -x6 "https://pocket-snapshot.liquify.com/files/$latestFile"
lz4 -c -d "$latestFile" | tar -x -C ${POCKET_DATA_DIR}
rm latest.txt
```

#### Compressed Inline (slower but smaller disk footprint)

```bash
wget -O latest.txt https://pocket-snapshot.liquify.com/files/latest_compressed.txt
latestFile=$(cat latest.txt)
wget -O - "https://pocket-snapshot.liquify.com/files/$latestFile" | lz4 -d - | tar -xv -C {POCKET_DATA_DIR}
rm latest.txt
```

### Issues

For any snapshot related issues, please [email Liquify](mailto:contact@liquify.io) or in the [node-chat channel on discord](https://discordapp.com/channels/553741558869131266/564836328202567725).

![Screenshot](https://github.com/pokt-network/pocket-core/assets/1892194/079b8dc5-4536-46b9-be69-7ae6b162c883)
