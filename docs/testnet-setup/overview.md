# Running your own testnet

This is a short guide to running your own testnet with multiple nodes.

Follow these steps to create the testnet:

#### 1. Find node_id for each node

```
$ plasmad tendermint show_node_id
de83c1a52f9c2aba189a2d01d7fa3e2e7c1736dd
```
    
#### 2. Open ~/.plasmad/config/config.toml and add the other nodes on the network to the persistent_peers field like so, node_id@ip_address:26656

```
# .plasmad/config/config.toml
# Comma separated list of nodes to keep persistent connections to
persistent_peers = "c00ce0b868bd5d5576d23f0ad1090f3f478b7961@35.193.21.171:26656, d35da1c7365c5a3345a23d01d1081f4d87634abc@192.168.1.12:26656"
```

#### 3. Set genesis file with the genesis UTXO's you want and the single validator's public key.

```json
 "app_state": {
    "Validator": {
          "type": "tendermint/PubKeyEd25519",
          "value": "nchF//ddMszMxK+2bQ4xQjTdeyHCRg08NuPDGZCWHw="
    },
    "UTXOs": [
      {
        "Address": "0xdc8820baA512f5827a8A5b45a07b0045Da407700",
        "Denom": "100",
        "Position": [
          "0",
          "0",
          "0",
          "1"
        ]
      }
    ]
  }
```

#### 4. Run plasmad start for each node

There are example config and genesis files in this folder that can help make setting up your own testnet easier.
