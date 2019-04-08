#!/bin/sh

echo "setting contract address to $PLASMA_ADDRESS"
sed -i "/ethereum_plasma_contract_address = /c\ethereum_plasma_contract_address = \"$PLASMA_ADDRESS\"" "$HOME/.plasmad/config/plasma.toml"
echo "setting node to $NODE_URL"
sed -i "/ethereum_nodeurl = /c\ethereum_nodeurl = $NODE_URL" "$HOME/.plasmad/config/plasma.toml" && \
cat $HOME/.plasmad/config/plasma.toml
echo "starting plasmad..."
plasmad start &
echo "starting plasmacli rest-server"
plasmacli rest-server &
