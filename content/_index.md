---
title: wp-matrix
---

[![Build Status](https://ci.thegeeklab.de/api/badges/thegeeklab/wp-matrix/status.svg)](https://ci.thegeeklab.de/repos/thegeeklab/wp-matrix)
[![Docker Hub](https://img.shields.io/badge/dockerhub-latest-blue.svg?logo=docker&logoColor=white)](https://hub.docker.com/r/thegeeklab/wp-matrix)
[![Quay.io](https://img.shields.io/badge/quay-latest-blue.svg?logo=docker&logoColor=white)](https://quay.io/repository/thegeeklab/wp-matrix)
[![Go Report Card](https://goreportcard.com/badge/github.com/thegeeklab/wp-matrix)](https://goreportcard.com/report/github.com/thegeeklab/wp-matrix)
[![GitHub contributors](https://img.shields.io/github/contributors/thegeeklab/wp-matrix)](https://github.com/thegeeklab/wp-matrix/graphs/contributors)
[![Source: GitHub](https://img.shields.io/badge/source-github-blue.svg?logo=github&logoColor=white)](https://github.com/thegeeklab/wp-matrix)
[![License: Apache-2.0](https://img.shields.io/github/license/thegeeklab/wp-matrix)](https://github.com/thegeeklab/wp-matrix/blob/main/LICENSE)

Woodpecker CI plugin to send messages to a Matrix room.

<!-- prettier-ignore-start -->
<!-- spellchecker-disable -->
{{< toc >}}
<!-- spellchecker-enable -->
<!-- prettier-ignore-end -->

## Usage

```YAML
kind: pipeline
name: default

steps:
- name: notify
  image: quay.io/thegeeklab/matrix
  settings:
    homeserver: https://matrix.org
    roomid: abcdefghijklmnopqrstuvwxyz:matrix.org
    username: octocat
    password: secret
```

### Parameters

<!-- prettier-ignore-start -->
<!-- spellchecker-disable -->
{{< propertylist name=wp-matrix.data sort=name >}}
<!-- spellchecker-enable -->
<!-- prettier-ignore-end -->

## Build

Build the binary with the following command:

```Shell
make build
```

Build the container image with the following command:

```Shell
docker build --file Containerfile.multiarch --tag thegeeklab/wp-matrix .
```

## Test

```Shell
docker run --rm \
  -e PLUGIN_ROOMID=0123456789abcdef:matrix.org \
  -e PLUGIN_USERNAME=yourbot \
  -e PLUGIN_PASSWORD=p455w0rd \
  -v $(pwd):/build:z \
  -w /build \
  thegeeklab/wp-matrix
```
