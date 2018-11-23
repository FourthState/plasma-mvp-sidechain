let Migrations = artifacts.require("Migrations");

let PriorityQueue = artifacts.require("PriorityQueue");
let PriorityQueue_Test = artifacts.require("PriorityQueue_Test");
let RootChain = artifacts.require("RootChain");

module.exports = function(deployer, network, accounts) {
    deployer.deploy(Migrations);
    deployer.deploy(PriorityQueue, {from: accounts[0]}).then(() => {
        deployer.link(PriorityQueue, [PriorityQueue_Test, RootChain]);
    });
};
