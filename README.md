# CMYK Services

> Felix Felicis - VS Code workspace [cmyk-svc.code-workspace](cmyk-svc.code-workspace)

> [CMYK](https://github.com/chrisntb/cmyk) - Engineering and architecture documentation

Go services using the [Fiber](https://gofiber.io/) framework and a Postgres database.

## Documentation

See [docs/README.md](docs/README.md).

## Prerequisites

See [docs/Tools.md](docs/Tools.md) for the required tools.

You need to setup the infrastructure, see [cmyk-system](https://github.com/chrisntb/cmyk-system).

If using a kubernetes cluster, once setup, copy its controller `Kubeconfig` file from `/home/ubuntu/.kube/config` to this file on your local machine `~/.kube/config`.
After that, when you run the service, it should be able to access the K8s cluster:

```shell
make run

Running the app locally
go run cmd/api/main.go

2026/02/01 11:11:33 created env client: MockModeEnv=NOT_SET, IsMockMode=false
2026/02/01 11:11:33 created k8s client
2026/02/01 11:11:33 created mock client

 ┌───────────────────────────────────────────────────┐
 │                  Fiber v2.52.11                   │
 │               http://127.0.0.1:4000               │
 │       (bound on host 0.0.0.0 and port 4000)       │
 │                                                   │
 │ Handlers ............. 15 Processes ........... 1 │
 │ Prefork ....... Disabled  PID ............. 62128 │
 └───────────────────────────────────────────────────┘
```

## Tasks

Usage:

```shell
make help
```

### Run

```shell
make run
# OR using a SOCKS5 proxy
SOCKS5_PROXY="socks5://127.0.0.1:8123" make run
# OR using mock data
make run_

curl -v http://127.0.0.1:4000/api/health
```

### Build

```shell
make build

# Run the build
./tmp/cmyk
```

### Dependencies

Re-initialize:

```shell
make reinit
```

### Checks

Quality:

```shell
make check
```

Run unit tests:

```shell
make test
```

Run API tests:

```shell
make run

cd tests/api
bru run --env [local|dev]
```

### API Documentation

Generate OpenAPI specification:

```shell
make openapi

# For usage
swag init --help
```
