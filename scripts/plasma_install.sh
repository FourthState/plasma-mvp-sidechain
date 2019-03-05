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

# Export GO path and append to .bashrc file
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$PATH
sed -i -e "\$aexport GOPATH=\$HOME/go\nexport PATH=\$GOPATH/bin:\$PATH" ~/.bashrc

# Install dep
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Install plasmad and plasmacli
go get github.com/FourthState/plasma-mvp-sidechain
cd ~/go/src/github.com/FourthState/plasma-mvp-sidechain/
git fetch --all
git checkout develop
dep ensure
cd server/plasmad/
go install
plasmad unsafe-reset-all
cd ../../client/plasmacli/
go install

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

echo ""
echo "Run 'sudo service geth start' to begin syncing to rinkeby network"
echo "Set configuration parameters in ~./plasmacli/plasma.toml, ~/.plasmad/config/"
echo ""
