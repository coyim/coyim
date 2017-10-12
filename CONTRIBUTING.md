# How to contribute
Here's the brief:

* We welcome contributions of all kinds, including but not limited to features, bug fixes, quality assurance, documentation, security review or asking questions
* If you do not know how to code yet you can help us translating CoyIM through [Zanata](https://translate.zanata.org/zanata/project/view/coyim)
* Pull requests are based off, integrated into & rebased against master
* Write automated tests, ideally using TDD. CI needs to be green in order to merge.
* Contact us for questions & suggestions:
  * IRC: #coyim @ irc.oftc.net ([join via webchat](https://webchat.oftc.net))
  * Email: [coyim@olabini.se](mailto:coyim@olabini.se)
  * Twitter: [@coyproject](https://twitter.com/coyproject)

  This document outlines our way of working, gives hints and outlines the steps to make your contribution to CoyIM as smooth as possible. You're not required to read this before getting started. We're explaining the way we work to make sure you're having a good experience and can make best use of the time you're contributing to our project.

## Getting started

Coy is written in [Golang](https://golang.org/) and uses
[GTK+3](http://www.gtk.org/) as its UI toolkit.

### Requirements

- `git`
- `golang` with `cgo` support.
- `gtk`
- `make`
- `ruby`

Installing these requirements differs on each system.

### Instructions

1. Make sure you have created your `GOPATH` directory (default to `$HOME/go`). You can also set a different path using:
```sh
export GOPATH=$HOME/<your-go-directory>
```

2. Also, make sure you have appended `$GOPATH/bin` to your path, like this:
```sh
export PATH=$PATH:$GOPATH/bin
```

3. Export an environment variable with your GTK version:
```sh
export GTK_VERSION=$(pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2)
```

4. Install the code inside go workspace:
```sh
go get -u -tags "gtk_${GTK_VERSION}" github.com/coyim/coyim
```
This will clone the repo inside the `$GOPATH/src/github.com/coyim/coyim` directory. For further instructions, check https://golang.org/doc/code.html#Workspaces.

5. Download the project dependencies:
```
make deps-dev
```

6. Build and run the tests:
```
make
```

7. Build user interface:
```
make build-gui
```

## Contributions steps

This is the lifecycle of a contribution. We follow a simplified fork + pull request workflow:

* To start, fork this repository and create a branch that's based off the latest commit in the `master` branch
* Implement the change
* Send a pull request against the master branch. Please make sure the automated tests are passing, as indicated by GitHub on the pull requests.
* Please keep your feature branch updated. Rebase your branch against upstream changes on the master branch, resolve any conflicts and make sure the tests are staying green.
* Your pull request will reviewed and merged

### What to work on

Generally, all issues that have no user assigned are awaiting work and free to play. If you want to make sure, or you think it will take more than a couple of days to complete your work, please reach out to us using the contact info above.

### Guidelines

When implementing your change, please follow this advice:

* Your change should be described in an issue, or latest in the pull request description.
* For bugs, please describe, in an issue or pull request:
  1. Steps to reproduce the behavior
  2. Expected behavior
  3. Actual behavior. Please also include as much meta-information as reasonable, e.g. time & date, software version etc.
* Pull requests need not to be finished work only; you can also submit changes in consecutive Pull Requests as long as CI stays green. Also, please send a PR with the intention of discussion & feedback. Please mark those Pull Requests appropriately.
* We review your pull request. This review is prioritised and done as part of our prioritisation. During this time, we ask you to keep it up to date by frequently rebasing against master.

### Review Criteria

When reviewing your contribution, we apply the following criteria:

* Test must be green. This usually includes an automatic check of the style guide using e.g. gofmt. All tests should be executed locally before you push, as well as on CI. If you struggle to reproduce a failure on CI locally, please notify us on IRC so we can resolve the issue.
* We won't tolerate abusive, exploitative or harassing behavior in every context of our project and refuse collaboration with any individual who exposes such behavior.

TODO: Should we include a Code of Conduct, like
https://github.com/discourse/discourse/blob/master/docs/code-of-conduct.md ?

## ThoughtWork's role

ThoughtWorks started the development of CoyIM on October 2015. The company invested its own resources and provided a team of software delivery experts to lay the foundation for the project.

ThoughtWork's goal wasn't to make money from CoyIM, but to combine their capability to deliver software and passion for defending a free Internet. Since then, we wanted to build software to counter the widespread of mass surveillance, and ensure digital privacy and anonymity for every person on the planet.
