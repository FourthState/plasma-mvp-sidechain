# Contributing

Thank you for considering making contributions to Fourth State's Plasma MVP implementation! We welcome contributions from anyone! See the [open issues](https://github.com/FourthState/plasma-mvp-sidechain/issues) for things we need help with!

Contribute to design discussions and conversation by joining our [Discord Server](https://discord.gg/YTB5A4P)

## Prerequisites
* [Golang](https://golang.org/doc/install)

## How to get started:

Add this repository into your $GOPATH/src/github.com/FourthState directory:

`go get "github.com/FourthState/plasma-mvp-sidechain"`

Testing:

`make test`

### Forking

Using `go get` on a forked repository will result in all the import paths being wrong. So instead we will add a new remote for the original repo.

To create a fork and remote branch to work on:

- Create a fork on github

- Go to the original repo locally (ie. `$GOPATH/src/github.com/FourthState/plasma-mvp-sidechain`)

- `git remote rename origin upstream`

- `git remote add origin git@github.com:user/plasma-mvp-sidechain.git`

`origin` : refers to fork

`upstream` : refers to original repo


Now, create a branch that is named off the feature you are trying to implement. See these [guidelines](https://nvie.com/posts/a-successful-git-branching-model/)

Make your changes. Add tests and comment those changes. 

If your tests pass, push to your fork and [submit a pull request](https://github.com/FourthState/plasma-mvp-rootchain/pulls) to the develop branch.


## Proposals:

If you would like to propose a protocol change, open up an issue. If the reviewers decide the proposed change is in line with the project's aim, then a writeup should also be added to the [research repository](https://github.com/FourthState/plasma-research). It is also advisable to publish the proposed change to [Eth Research](https://ethresear.ch/), so other plasma implementations can benefit from the proposed change. 

