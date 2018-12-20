<<<<<<< HEAD
let RootChain = artifacts.require("RootChain");

module.exports = function(deployer, network, accounts) {
	deployer.deploy(RootChain, {from: accounts[0]});
=======
let PlasmaMVP = artifacts.require("PlasmaMVP");

module.exports = function(deployer, network, accounts) {
	deployer.deploy(PlasmaMVP, {from: accounts[0]});
>>>>>>> b3167013cb609ec55bd2a944e44a4d169ed332c9
};
