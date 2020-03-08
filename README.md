# Echoserver

[![latest release](https://badge.fury.io/gh/usvc%2Fechoserver.svg)](https://github.com/usvc/echoserver/releases)
[![pipeline status](https://gitlab.com/usvc/services/echoserver/badges/master/pipeline.svg)](https://gitlab.com/usvc/services/echoserver/-/commits/master)
[![Build Status](https://travis-ci.org/usvc/echoserver.svg?branch=master)](https://travis-ci.org/usvc/echoserver)

A simple echoserver for use with testing HTTP-based connections in production.

- - -

# Usage

## Via binary

Download the latest binaries from [https://github.com/usvc/echoserver/releases](https://github.com/usvc/echoserver/releases)

```sh
echoserver;
```

## Via `docker run` command

```sh
docker run -it -p 8888:8888 usvc/echoserver:latest;
```

## Via `docker-compose` file

```yaml
version: "3.5"
services:
  echoserver:
    build:
      context: ..
      dockerfile: ./deploy/Dockerfile
    image: usvc/echoserver:latest
    environment:
      SERVER_ADDR: 0.0.0.0
      SERVER_PORT: "8888"
    ports:
      - 8888:8888
```

- - -

# Configuration

## Environment Variables

| Key | Value |
| --- | --- |
| `SERVER_ADDR` | Hostname/IP address that server should bind to |
| `SERVER_PORT` | Port that server should listen on |

- - -

# Development Runbook

## Getting Started

1. Clone this repository
2. Run `make deps` to pull in external dependencies
3. Write some awesome stuff
4. Run `make test` to ensure unit tests are passing
5. Push

## Continuous Integration (CI) Pipeline

### On Github

Github is used to deploy binaries/libraries because of it's ease of access by other developers.

#### Releasing

Releasing of the binaries can be done via Travis CI.

1. On Github, navigate to the [tokens settings page](https://github.com/settings/tokens) (by clicking on your profile picture, selecting **Settings**, selecting **Developer settings** on the left navigation menu, then **Personal Access Tokens** again on the left navigation menu)
2. Click on **Generate new token**, give the token an appropriate name and check the checkbox on **`public_repo`** within the **repo** header
3. Copy the generated token
4. Navigate to [travis-ci.org](https://travis-ci.org) and access the cooresponding repository there. Click on the **More options** button on the top right of the repository page and select **Settings**
5. Scroll down to the section on **Environment Variables** and enter in a new **NAME** with `RELEASE_TOKEN` and the **VALUE** field cooresponding to the generated personal access token, and hit **Add**

### On Gitlab

Gitlab is used to run tests and ensure that builds run correctly.

#### Version Bumping

1. Run `make .ssh`
2. Copy the contents of the file generated at `./.ssh/id_rsa.base64` into an environment variable named **`DEPLOY_KEY`** in **Settings > CI/CD > Variables**
3. Navigate to the **Deploy Keys** section of the **Settings > Repository > Deploy Keys** and paste in the contents of the file generated at `./.ssh/id_rsa.pub` with the **Write access allowed** checkbox enabled

- **`DEPLOY_KEY`**: generate this by running `make .ssh` and copying the contents of the file generated at `./.ssh/id_rsa.base64`

#### DockerHub Publishing

1. Login to [https://hub.docker.com](https://hub.docker.com), or if you're using your own private one, log into yours
2. Navigate to [your security settings at the `/settings/security` endpoint](https://hub.docker.com/settings/security)
3. Click on **Create Access Token**, type in a name for the new token, and click on **Create**
4. Copy the generated token that will be displayed on the screen
5. Enter the following varialbes into the CI/CD Variables page at **Settings > CI/CD > Variables** in your Gitlab repository:

- **`DOCKER_REGISTRY_URL`**: The hostname of the Docker registry (defaults to `docker.io` if not specified)
- **`DOCKER_REGISTRY_USERNAME`**: The username you used to login to the Docker registry
- **`DOCKER_REGISTRY_PASSWORD`**: The generated access token

- - -

# Licensing

Code here is licensed under the [MIT license](./LICENSE) by [@zephinzer](https://gitlab.com/zephinzer).
