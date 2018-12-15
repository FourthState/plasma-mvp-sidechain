let PlasmaMVP = artifacts.require("PlasmaMVP");

module.exports = function(deployer, network, accounts) {
	deployer.deploy(PlasmaMVP, {from: accounts[0]});
};
