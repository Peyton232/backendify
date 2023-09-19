# Backendify

At Backendify, our goal is to make our customer's lives easier. No one has to deal with the complexity of having multiple providers for the same kind of data!
We have several services that return company data when given country iso code and a company id.
You will build a proxy service that provides the same API for all of them.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
- [Testing](#testing)

## Getting Started

This section explains how to get started with the project, including prerequisites, installation, and configuration.

### Prerequisites

- Docker: Make sure you have Docker installed on your system.

### Installation

You can choose one of the following methods to run the project locally:

#### Method 1: Docker

Build and run the Docker container using the following commands:

```bash
docker build -t project-name .
docker run -p 9000:9000 project-name arg1 arg2
```

#### Method 2: Using Go

- Download project dependencies:

```bash
go mod download
```

- Run the project:

```bash
go run main.go arg1 arg2
```

#### Configuration

The configuration file is config.yaml. For local development, it is recommended to set mockFlag to true to mock responses from external APIs.

#### Testing

To run tests, execute the following command:

```bash
go test -v ./...
```
# backendify
