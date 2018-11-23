let RootChain = artifacts.require("RootChain");

module.exports = function(deployer, network, accounts) {
	deployer.deploy(RootChain, {from: accounts[0]});
};
