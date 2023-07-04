# [Pocket Snapshotter](https://github.com/liquify-validation/pocket-snapshotter)

## Public snapshots
If you're looking for Pocket native blockchain data snapshots, they are provided by [Liquify LTD](https://www.liquify.io/) and can be viewed via the file explorer

Explorer link: https://pocket-snapshot.liquify.com

Snapshots are updated every Monday at 00:00 UTC. The snapshots are generated on the Master (UK) and then sent over to the US and JP regions.

### Mirrors

The pocket snapshot link above is a global endpoint which is available in 3 different regions (UK, US west, Japan). The individual regions can also be accessed on the following links.

UK (Master): https://pocket-snapshot-uk.liquify.com

US: https://pocket-snapshot-us.liquify.com

JP: https://pocket-snapshot-jp.liquify.com

Note: If accessing the snapshots on Monday it may be best to use the UK (Master) endpoint since there will be a 4-12hour delay in updating the slaves in the other regions.

### Download using CLI

The snapshot repos hold the last 3 weeks of snapshots. The latest one being referenced by the file latest.txt and latest_compressed.txt.

Please fill in the location of your data_dir below

#### Uncompressed

```
wget -O latest.txt https://pocket-snapshot.liquify.com/files/latest.txt
latestFile=$(cat latest.txt)
aria2c -s6 -x6 "https://pocket-snapshot.liquify.com/files/$latestFile"
tar xvf "$latestFile" -C <pocket data_dir>
rm latest.txt
```

#### Uncompressed Inline (slower but smaller disk footprint)

The below snippet will download and extract inline this may be benifical if you have constrained disk space and cannot afford to store both the temp archive and extracted datadir

```
wget -O latest.txt https://pocket-snapshot.liquify.com/files/latest.txt
latestFile=$(cat latest.txt)
wget -c "https://pocket-snapshot.liquify.com/files/$latestFile" -O - | sudo tar -xz -C <pocket data_dir>
rm latest.txt
```

#### Compressed

```
wget -O latest.txt https://pocket-snapshot.liquify.com/files/latest_compressed.txt
latestFile=$(cat latest.txt)
aria2c -s6 -x6 "https://pocket-snapshot.liquify.com/files/$latestFile"
lz4 -c -d "$latestFile" | tar -x -C <pocket data_dir>
rm latest.txt
```

#### Compressed Inline (slower but smaller disk footprint)

```
wget -O latest.txt https://pocket-snapshot.liquify.com/files/latest_compressed.txt
latestFile=$(cat latest.txt)
wget -O - "https://pocket-snapshot.liquify.com/files/$latestFile" | lz4 -d - | tar -xv -C <pocket data_dir>
rm latest.txt
```

### Issues

For any snapshot related issues please contact Liquify via contact@liquify.io or in the [node-chat channel on discord](https://discordapp.com/channels/553741558869131266/564836328202567725).



