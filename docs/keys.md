In order to spend utxos on the sidechain, we will need keys corresponding to the addresses that own those utxos.
We can use the keys command to generate and manage these keys.
We use a mapping between a name for our key and its address to simplify usage of these addresses.

The keystore's default location is `~/.plasmacli/keys/`

## Generating Keys ##

For example, the following generates a new key with the associated name "mykey":

```
plasmacli keys add mykey
Enter new passphrase for your key:
Repeat passphrase:

**Important** do not lose your passphrase.
It is the only way to recover your account
You should export this account and store it in a secure location
NAME: mykey ADDRESS: 0x3b00b1deee88ac18a2a53c988ff30ba1f561d3a5
```

#### Import unencrypted private keys ####

Ganache-cli display

```
 Ganache CLI v6.1.8 (ganache-core: 2.2.1)

Available Accounts
==================
(0) 0xea6ed4bb7cba09c391c11a15d5472e806caa3986 (~100 ETH)
(1) 0x6a06b2fd816021568b5b8abad00eab4679ab3450 (~100 ETH)

Private Keys
==================
(0) 0x46b7f5573b110d21ca809d1e17d97c582e6e402b1d7daa7d09264c13015e73ae
(1) 0x2a21b4bcbc0d812adf41ff098f336a86325e3c7df31d5c576ac7b760ed2d56f0

```

Import by passing in the unecrypted key via command line:

```
plasmacli keys import imported_key 46b7f5573b110d21ca809d1e17d97c582e6e402b1d7daa7d09264c13015e73ae
Enter new passphrase for your key:
Repeat passphrase:
Successfully imported.
NAME: imported_key      ADDRESS: 0xea6ed4bb7cba09c391c11a15d5472e806caa3986
```

Import using a file which stores the unencrypted key:

Put the private key you wish to import in a file

keyfile:
```
46b7f5573b110d21ca809d1e17d97c582e6e402b1d7daa7d09264c13015e73ae
```

Use plasmacli
```
plasmacli keys import imported_keyfile --file ~/keyfile
Enter new passphrase for your key:
Repeat passphrase:
Successfully imported.
NAME: imported_keyfile      ADDRESS: 0x6a06b2fd816021568b5b8abad00eab4679ab3450
```

#### List Keys ####

```
plasmacli keys list
NAME:               ADDRESS:
acc1                0xec36ead9c897b609a4ffa5820e1b2b137d454343
acc2                0x5475b99e01ac3bb08b24fd754e2868dbb829bc3a
imported_key        0xea6ed4bb7cba09c391c11a15d5472e806caa3986
imported_keyfile    0x6a06b2fd816021568b5b8abad00eab4679ab3450
mykey               0x3b00b1deee88ac18a2a53c988ff30ba1f561d3a5
```

#### Delete Keys ####
You must know a key password to delete a key

```
plasmacli keys delete imported_keyfile
Enter passphrase:
Account deleted.
```

#### Update Keys ####
Update the password:

```
plasmacli keys update imported_key
Enter passphrase:
Enter new passphrase for your key:
Repeat passphrase:
Account passphrase has been updated.
```

Update the name of a key:

```
plasmacli keys update imported_key --name new_key_name
Account name has been updated.
```


