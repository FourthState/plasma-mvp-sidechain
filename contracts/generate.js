#!/usr/bin/env node
var shell = require('shelljs');
var fs = require("fs");

console.log('Generating Go Wrappers...');
console.log('Cleaning previous build...');
shell.rm('-rf', 'build wrappers abi');
shell.exec('truffle compile');

// files to generate
generate('./build/contracts/PlasmaMVP.json');

function generate(path) {
    shell.mkdir('-p', ['wrappers', 'abi']);

    let contract = JSON.parse(fs.readFileSync(path, {encoding: 'utf8'}));

    //const filename = contract.contractName;
    let snakeCasedFilename = (
        contract.contractName.replace(/([a-z])([A-Z])/g, '$1_$2').replace(/([A-Z])([A-Z][a-z])/g, '$1_$2')
    ).toLowerCase();

    
    fs.writeFileSync(`abi/${contract.contractName}.abi`, JSON.stringify(contract.abi))
    shell.exec(`abigen --abi abi/${contract.contractName}.abi --pkg wrappers --type ${contract.contractName} --out wrappers/${snakeCasedFilename}.go`);
}
