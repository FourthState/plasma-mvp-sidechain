#!/bin/bash
# This script is intended to install geth, plasmad, and plasmacli
# It assumes none of the dependencies have been installed
# It will format geth to be a system service

# Upgrade the system and install go, gcc, make, geth
sudo apt-get install software-properties-common
sudo add-apt-repository -y ppa:ethereum/ethereum
sudo apt update
sudo apt upgrade -y
sudo apt install gcc make ethereum -y
sudo snap install --classic go
sudo mkdir -p ~/go/bin/

# Export GO path and append to .profile file
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
echo "export GOPATH=$HOME/go" >> ~/.profile
echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.profile

source ~/.profile

# Install plasmad and plasmacli
go get github.com/FourthState/plasma-mvp-sidechain
cd ~/go/src/github.com/FourthState/plasma-mvp-sidechain/
git fetch --all
git checkout develop
make install
plasmad unsafe-reset-all

# Install npm, Truffle
apt-get install nodejs
apt-get install npm
cd ~/go/src/github.com/FourthState/plasma-mvp-sidechain/contracts/
npm install -g

# setup geth as system service
sudo useradd -m -d /opt/geth --system --shell /usr/sbin/nologin geth
sudo -u geth mkdir -p /opt/geth/rinkeby
cd /opt/geth/rinkeby/

# geth.service
echo "[Unit]
Description=Geth
After=network-online.target
[Service]
User=geth
ExecStart=/usr/bin/geth --datadir=/opt/geth/rinkeby/chaindata/ --rinkeby --rpc --rpcapi db,eth,net,web3,personal
Restart=always
RestartSec=3
LimitNOFILE=4096
[Install]
WantedBy=multi-user.target" > geth.service

sudo mv geth.service /etc/systemd/system/
sudo systemctl enable geth.service
sudo service geth start

echo ""
echo "Geth has begun syncing to rinkeby network"
echo "Run 'sudo service geth status' to check its status"
echo "Run 'sudo service geth stop' to stop the geth full node"
echo "Set configuration parameters in ~./plasmacli/plasma.toml and ~/.plasmad/config/"
echo "Copy your genesis file to ~/.plasmad/config/genesis.json or modify the existing one"
echo "Add the operators node to 'persisten peers' in ~/.plasmad/config/config.toml"
echo "Add the plasma contract address to plasma.toml in ~/.plasmacli and ~/.plasmad/config"
echo ""
