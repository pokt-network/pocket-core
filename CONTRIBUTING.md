## Pocket Core Contribution Guide

This Contribution Guide aims to guide contributors through the process of proposing, developing, testing and upstreaming changes to the Pocket Core client.

### Proposing changes

#### Communicating proposals

The first step towards contribution is to effectively propose your change by opening up a new issue using the **Contribution Proposal** issue template. Communicating your proposal effectively will be of the upmost importance throughout the lifecycle of your contribution, so make sure your description is clear and concise. Feel free to also use the `#core-research` channel in our [Official Discord server](https://bit.ly/POKTARCADEdscrd) to ask any questions regarding the proposal you want to make to the repository.

#### Consensus-breaking changes

A consensus-breaking change means a change that would require 66% of the Validator Power in the network to be adopted. Furthermore these changes need to be voted in and approved by the DAO. To propose a consensus-breaking change please follow the Pocket Improvement Proposal documentation found [here](https://docs.pokt.network/home/paths/governor/submit-a-proposal/pip-pocket-improvement-proposal).

#### Quality Assurance

Proposals must be accompanied by a Quality Assurance plan outlining if implemented, how they will be tested end-to-end. More on this on the Quality Assurance section of this guide.

### Request for proposals

A request for proposals (RFP) is used when the author wants to specify a desired behaviour of the software, but they don't have enough context or information to make a proposal for the change, inviting others to propose solutions and surfacing the need for the desired behaviour.

#### Browsing RFPs

Issues in the repository marked with the `rfp` tag will be considered requests for proposals. More than one proposal could be linked to a RFP which should be tracked via issue linking in the RFP issue.

#### Submitting a RFP

To submit a RFP, please use the Request for Proposals issue template when creating a new issue. This template will automatically include the tag and any necessary information to be able to move on with the RFP.

### Development Guide

#### Forking Pocket Core

The first step towards modifying Pocket Core is to fork the repository. All changes will be upstreamed from forks. To fork a repository on Github please follow this [guide](https://docs.github.com/en/get-started/quickstart/fork-a-repo).

Within your fork you are free to work however you want, but keep in mind that in the end you need to Pull Request into the official Pocket Core repository so your changes are included within an official release.

#### Setting up the Go Environment

Please follow the [Official Installation Guide](https://go.dev/doc/install) to complete this step. Pocket Core uses `go 1.18` so make sure to install the appropiate version before beginning development.

#### Installing dependencies

Pocket Core uses Go Modules, listed in the `go.mod` file in the root of the project. To add a new dependency to the `go.mod` follow the official guide [here](https://pkg.go.dev/cmd/go#hdr-Add_dependencies_to_current_module_and_install_them).

##### Forked dependencies

Pocket Core uses 2 forked dependencies:

- `github.com/pokt-network/tendermint` which is a fork of `github.com/tendermint/tendermint`.
- `github.com/pokt-network/tm-db` which is a fork of `github.com/tendermint/tm-db`.

If you need to make changes to these dependencies, you will need to fork them and follow all the steps in this guide to this point before proceeding further.

#### Building Pocket Core Binary

From the root of the project run:

`go build -o pocket app/cmd/pocket_core/main.go`

And you will build a `pocket` binary in the root folder which can now be run.

#### Coding style

- Code must adhere to the official Go formatting guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt)).
- Code must be documented adhering to the official Go commentary guidelines.

### Quality Assurance

When proposing new changes, you also need to propose how these changes will be tested and added to the project's Unit Testing and Regression Testing suites.

#### Unit Testing

Pocket Core contains unit tests where applicated and it is encouraged that if you are making changes to the codebase, that you include the appropiate unit tests which will be run by the repository configured CI/CD (Continous Integration/Continous Deployment) pipeline.

#### Regression Testing

The Regression Testing suite of Pocket Core is a series of scenarios that need to be manually run and submitted in a Regression Testing Report with each PR. The test suite can be found [here](./doc/qa/regression). Please add any scenarios regarding your changes as needed.

### The Pull Request Process

#### Submitting a PR to a Release Window

Every PR will target an upcoming Release Window for the Pocket Core repository, these windows will be listed as Github Projects in the repository and will be listed in the repository [Projects Page](https://github.com/pokt-network/pocket-core/projects?type=classic). Every project in this page will represent a Release Window which contains the following:

- A PR Cut-out date by which PR's must be approved to be included in the release.
- An integration branch to submit PR's to for this release.

Open a PR against the **integration branch** in the official repository: `https://github.com/pokt-network/pocket-core`. To learn how to open a PR in Github please follow this [guide](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork).

#### The PR review process

Every Pull Request will require at least 2 reviewers from the Core team. The reviews will be done in a 2 phase approach:

1. A functional PR review where the code will be extensively reviewed by a Core Team member and feedback will be provided until the code quality meets the functional requirements of the intended Proposal.
2. An Integration approval, every friday the Core Team will be evaluating PR's ready to be included, and selecting those into QA builds that will be tested as part of the Regression testing performed before each released is launched.

#### The PR Success Criteria

The Core Team will be requiring the following before merging a Pull Request:

- It has a linked Contribution Proposal issue specifying the functionality implemented in the PR.
- There is an associated pull request to the [Pocket Core Functional Tests repository](https://github.com/pokt-network/pocket-core-func-tests) with the appropiate QA scenarios associated with the functionality or fix contained in the PR.
- Enough evidence (either automated or manual) of testing of these new scenarios included in the PR. This evidence must be presented in the format template contained [here](https://github.com/pokt-network/pocket-core/tree/staging/doc/qa/regression).

Every proposal is different in scope and complexity, so the following points will increase the likelihood of you submitting a successful Pull Request:

- Your proposal issue is clear, concise and informative.
- Your PR is within the scope of the proposal.
- Your code follows the code style outlined in this guide.
- Your Quality Assurance additions are clear and well documented.
- The CI/CD pipeline automated testing is **all green**.
- You provide any necessary documentation for your implementation and reasoning in implementation decisions.

### Releases

#### Beta releases

Beta releases in Pocket Core are the first step for changes to be deployed. Marked with the `BETA-` prefix in the release tag, these releases are the first step to test changes in testing environments such as `Pocket Testnet` and other local environments. A Beta release should not be used in production, and if used, it should be at that user's own risk.

#### Release Candidate releases

RC (or release candidates) are releases meant for production deployment. Release candidates have been tested extensively by the developer and other third parties, and are of stable behaviour and considered safe. Your proposal can be considered released once it is included in a Release Candidate and has been deployed to `Pocket Mainnet`.

### Security Bugs Disclosure

#### Starting a disclosure

To start the disclosure process, please send an email containing all evidence, documents and/or links for the vulnerability you are disclosing to `security@pokt.network`. You will receive a response within a 24 hour window outlining next steps.

#### Public disclosure

Once the vulnerability has been patched and deployed to the appropiate environments, the team will create a public disclosure announcement, acknowledging the vulnerability and giving credit to the original discloser or disclosers in case more than one person identifies and discloses the vulnerability.

