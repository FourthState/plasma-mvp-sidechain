#!/bin/sh

echo "setting node to $NODE_URL"
sed -i "/ethereum_nodeurl = /c\ethereum_nodeurl = $NODE_URL" "$HOME/.plasmad/config/plasma.toml" && \
cat $HOME/.plasmad/config/plasma.toml
echo "starting plasmad..."
plasmad start &
echo "starting plasmacli rest-server"
plasmacli rest-server &
