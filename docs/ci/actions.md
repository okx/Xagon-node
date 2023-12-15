# GitHub Actions

This documents gives a brief overview of the existing GitHub Actions on the repo

## lint

Runs this [golang linter] over the code, this is the [linter configuration file].

### When is executed

PR opened and pushing changes to PRs.

## ok-to-test

Part of our setup for running tests on PRs from forks, you can read more about it
in the [ok-to-test document].

### When is executed

PR opened and pushing changes to PRs.

## push-docker-develop

Pushes docker images to docker hub, the images pushed are:
* `hermeznetwork/zkevm-node:develop`

### When is executed

Changes pushed to the `develop` branch.

## push-docker

Pushes docker images to docker hub, the images pushed are:
* `hermeznetwork/zkevm-node:latest`

### When is executed

Changes pushed to the `main` branch.

## test-e2e

Runs e2e tests divided in several groups executed in parallel, read more about
[CI groups].

### When is executed

PR opened and pushing changes to PRs. There are two variants, `trusted` and
`from-fork`, depending on the procedence of the PR, more about it in the
[ok-to-test document].

## test-full-non-e2e

Runs all the non-e2e tests.

### When is executed

PR opened and pushing changes to PRs. There are two variants, `trusted` and
`from-fork`, depending on the procedence of the PR, more about it in the
[ok-to-test document].

## updatedeps

The `x1-node` repo requires some external resources for working. We call
these resources custom dependencies (as opposed to the golang packages required
by the code).

The goal of the `updatedeps` action is to keep these custom dependencies up to
date. It checks the external resources with content required by this repo and
proposes a PR in case it finds any changes. The code executed can be found in
the [dependencies package].

Currently we are checking [three types of custom dependencies]:
* External docker images used in the [docker compose file]. For each image the
code compares the digest existing in the docker compose file with the digest
returned by docker hub API, if they differ it includes the new one in the docker
compose file.
* Protocol buffer files from [comms protocol repo]: after checking the files
for changes the client/server golang code is generated from them.
* Test vectors from the [test vectors repo].

With all the potential changes we create a new PR and the tests are run on it,
so that we can review and eventually approve the changes to be included in the
`x1-node` repo.

### When is executed

It runs as an scheduled action, every 3 hours.

[golang linter]: https://golangci-lint.run/
[linter configuration file]: ../../.golangci.yml
[ok-to-test document]: ./ok-to-test.md
[CI groups]: ./groups.md
[dependencies package]: ../../scripts/cmd/dependencies
[three types of custom dependencies]: ../../scripts/cmd/dependencies.go
[docker compose file]: ../../docker-compose.yml
[comms protocol repo]: https://github.com/0xPolygonHermez/zkevm-comms-protocol/
[test vectors repo]: https://github.com/0xPolygonHermez/zkevm-testvectors
