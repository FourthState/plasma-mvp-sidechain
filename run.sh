#!/bin/sh

echo "wants to set node to $NODE_URL"
sed "/ethereum_nodeurl = /c\ethereum_nodeurl = $NODE_URL" "$HOME/.plasmad/config/plasma.toml" > "$HOME/.plasmad/config/plasma_with_address.toml"
cat "$HOME/.plasmad/config/plasma_with_address.toml" > "$HOME/.plasmad/config/plasma.toml"
cat $HOME/.plasmad/config/plasma.toml
echo "starting plasmad..."
plasmad start &
echo "starting plasmacli rest-server"
plasmacli rest-server &
