# Contributing

Thank you for considering making contributions to Fourth State's Plasma MVP implementation! We welcome contributions from anyone! See the [open issues](https://github.com/FourthState/plasma-mvp-rootchain/issues) for things we need help with!

Contribute to design discussions and conversation by joining our [Discord Server](https://discord.gg/YTB5A4P).

## How to get started:

Fork, then clone the repo:

If you have ssh keys:
``git clone git@github.com:your-username/plasma-mvp-rootchain``

Otherwise:
``git clone https://github.com/your-username/plasma-mvp-rootchain``

Install dependencies with:
``npm install``

**Note**: requires Solidity 0.4.24 and Truffle 4.1.14

Make sure the tests pass:
1. Start ganache-cli: ``ganache-cli -m=plasma_mvp``
2. Run tests: ``truffle test``

Create a branch that is named off the feature you are trying to implement. See these [guidelines](https://nvie.com/posts/a-successful-git-branching-model/)

Make your changes. Add tests and comment those changes. 

If your tests pass, push to your fork and [submit a pull request](https://github.com/FourthState/plasma-mvp-rootchain/pulls) to the master branch. 

## Proposals:

If you would like to propose a protocol change, open up an issue. If the reviewers decide the proposed change is in line with the project's aim, then a writeup should also be added to the [research repository](https://github.com/FourthState/plasma-research). It is also advisable to publish the proposed change to [Eth Research](https://ethresear.ch/), so other plasma implementations can benefit from the proposed change. 

