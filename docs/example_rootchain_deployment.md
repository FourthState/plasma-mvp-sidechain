The following example deploys the rootchain contract to the rinkeby testnet. 
It requires a full eth node and truffle to be installed.

## Deployment ##

Attach your eth node to the geth console and check that it has finished syncing

```
geth attach /opt/geth/rinkeby/chaindata/geth.ipc 
Welcome to the Geth JavaScript console!

instance: Geth/v1.8.23-stable-c9427004/linux-amd64/go1.10.4
coinbase: 0xec36ead9c897b609a4ffa5820e1b2b137d454343
at block: 4070914 (Thu, 21 Mar 2019 18:30:08 UTC)
 datadir: /opt/geth/rinkeby/chaindata
 modules: admin:1.0 clique:1.0 debug:1.0 eth:1.0 miner:1.0 net:1.0 personal:1.0 rpc:1.0 txpool:1.0 web3:1.0

> eth.syncing
false
```

Import your eth private key by adding its file to the keystore located in `chaindata/keystore` or generate a new key using the geth commands.

Ensure that your eth private key has more than 0.12 eth to deploy the contract. Use a faucet if you don't have any testnet eth.

Unlock your account using the geth console:

```
> personal.unlockAccount("0xec36ead9c897b609a4ffa5820e1b2b137d454343", "1234567890", 1000000)
true
```

Either exit geth console or using another screen navigate to `go/src/github.com/FourthState/plasma-mvp-sidechain/contracts/`.
Open truffle.js and add:
```
rinkeby: {
      host: "localhost",
      port: 8545,
      network_id: "4", // Rinkeby ID 4
      from: "0xec36ead9c897b609a4ffa5820e1b2b137d454343", // account from which to deploy
      gas: 6712390
    }

```

See our [example_truffle.js]()

Run `truffle migrate --network rinkeby`

Example Output:
```
truffle migrate --network rinkeby

Compiling ./contracts/PlasmaMVP.sol...
Writing artifacts to ./build/contracts

⚠️  Important ⚠️
If you're using an HDWalletProvider, it must be Web3 1.0 enabled or your migration will hang.


Starting migrations...
======================
> Network name:    'rinkeby'
> Network id:      4
> Block gas limit: 6996931


1_initial_migration.js
======================

   Deploying 'Migrations'
   ----------------------
   > transaction hash:    0x46d8b7e69f5b00c91c21df4003fb702ab302c712555e0013fd19c3818d0f22ff
   > Blocks: 2            Seconds: 16
   > contract address:    0x073951236805bf6332fb3F90CdE294aC2D373172
   > account:             0xEc36eaD9c897b609A4fFa5820E1B2B137D454343
   > balance:             18.586470732988776
   > gas used:            284908
   > gas price:           20 gwei
   > value sent:          0 ETH
   > total cost:          0.00569816 ETH

   > Saving artifacts
   -------------------------------------
   > Total cost:          0.00569816 ETH


2_deploy_rootchain.js
=====================

   Replacing 'PlasmaMVP'
   ---------------------
   > transaction hash:    0xe1448402b7f41fb77048e16649f32ec58faf5355ac69060312c81cd751b4ec63
   > Blocks: 0            Seconds: 12
   > contract address:    0x9c36F39E87f3EA5283e186b385232044dC2f8c30
   > account:             0xEc36eaD9c897b609A4fFa5820E1B2B137D454343
   > balance:             18.473328172988776
   > gas used:            5657128
   > gas price:           20 gwei
   > value sent:          0 ETH
   > total cost:          0.11314256 ETH

   > Saving artifacts
   -------------------------------------
   > Total cost:          0.11314256 ETH


Summary
=======
> Total deployments:   2
> Final cost:          0.11884072 ETH

```

Congratulations! You have just deployed the rootchain contract. 
In the above example, the rootchain exists at address `0x9c36F39E87f3EA5283e186b385232044dC2f8c30`
