#!/bin/bash

rm -f nohup.out

if ! [ -z "$NODE_URL" ]; then
  echo "attempting to override \$NODE_URL to $NODE_URL"
  sed "/ethereum_nodeurl = /c\ethereum_nodeurl = $NODE_URL" "$HOME/.plasmad/config/plasma.toml" > "$HOME/.plasmad/config/plasma_with_address.toml"
  cat "$HOME/.plasmad/config/plasma_with_address.toml" > "$HOME/.plasmad/config/plasma.toml"
fi

echo "plasma.toml is: "
cat "$HOME/.plasmad/config/plasma.toml"

echo "starting plasmad..."
nohup plasmad start >> nohup.out &

echo "starting plasmacli rest-server"
nohup plasmacli rest-server >> nohup.out &

if ! [ -z "$KUBERNETES_SERVICE_HOST" ] && ! [ -z "$KUBERNETES_SERVICE_PORT" ]; then
  echo "looks like we're running in kube, will tail nohup.out while waiting for any subprocess to finish"
  cat nohup.out &
  wait -n
else
  echo "looks like we're not running in kube, exiting immediately"
fi
