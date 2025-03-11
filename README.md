![Baton Logo](./baton-logo.png)

# `baton-airbyte` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-airbyte.svg)](https://pkg.go.dev/github.com/conductorone/baton-airbyte) ![main ci](https://github.com/conductorone/baton-airbyte/actions/workflows/main.yaml/badge.svg)

`baton-airbyte` is a connector for Airbyte built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It allows you to sync user information from your Airbyte instance into ConductorOne.

Check out [Baton](https://github.com/conductorone/baton) to learn more about the project in general.

# Getting Started

## Prerequisites

To use this connector, you will need:
- An Airbyte instance
- Your Airbyte Client ID and Client Secret
  - These can be obtained from your Airbyte instance under Settings > Applications
  - Applications allow you to generate tokens to access the Airbyte API with the same permissions as your user

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-airbyte
baton-airbyte
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_DOMAIN_URL=domain_url -e BATON_API_KEY=apiKey -e BATON_USERNAME=username ghcr.io/conductorone/baton-airbyte:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-airbyte/cmd/baton-airbyte@main

baton-airbyte

baton resources
```

# Data Model

`baton-airbyte` will pull down information about the following resources:
- Users
- Workspaces
- Organizations

# Authentication

The connector requires the following credentials to authenticate with Airbyte:
- Airbyte Client ID (`BATON_AIRBYTE_CLIENT_ID`)
- Airbyte Client Secret (`BATON_AIRBYTE_CLIENT_SECRET`)
- Airbyte Domain URL (`BATON_DOMAIN_URL`)

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-airbyte` Command Line Usage

```
baton-airbyte

Usage:
  baton-airbyte [flags]
  baton-airbyte [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
   --domain-url string                 The domain URL of your Airbyte instance ($BATON_DOMAIN_URL)
   --airbyte-client-id string         The Airbyte client ID used to authenticate with Airbyte ($BATON_AIRBYTE_CLIENT_ID)
   --airbyte-client-secret string     The Airbyte client secret used to authenticate with Airbyte ($BATON_AIRBYTE_CLIENT_SECRET)
   --client-id string                 The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
   --client-secret string             The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                      The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                             help for baton-airbyte
      --log-format string                The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string                 The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -p, --provisioning                     If this connector supports provisioning, this must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --ticketing                          This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                            version for baton-airbyte

Use "baton-airbyte [command] --help" for more information about a command.
```
