let Migrations = artifacts.require("Migrations");

<<<<<<< HEAD
let PriorityQueue = artifacts.require("PriorityQueue");
let PriorityQueue_Test = artifacts.require("PriorityQueue_Test");
let RootChain = artifacts.require("RootChain");

module.exports = function(deployer, network, accounts) {
    deployer.deploy(Migrations);
    deployer.deploy(PriorityQueue, {from: accounts[0]}).then(() => {
        deployer.link(PriorityQueue, [PriorityQueue_Test, RootChain]);
    });
=======
module.exports = function(deployer, network, accounts) {
    deployer.deploy(Migrations);
>>>>>>> b3167013cb609ec55bd2a944e44a4d169ed332c9
};
