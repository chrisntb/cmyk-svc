# Context

## Tools

- Go service
  - Use alphabetical ordering for all structs
  - When creating files, do not use spaces in the file name
  - JSON APIs implemented using Fiber `https://gofiber.io/`
  - API specification using OpenAPI
  - API tests using Bruno
    - `tests/api`
- Postgres database

## Architecture

- Uses the Kubernetes (K8s) client library to visualize and manage the state of the K8s cluster
  - K8s API version `https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/`
- Uses Kueue for advanced job admission and placement logic, integrating with the native scheduler `https://kueue.sigs.k8s.io/docs/overview/`
  - Kueue API version `https://kueue.sigs.k8s.io/docs/reference/kueue.v1beta2/`

## Tasks

- Build
  - `make build`
- Static analysis
  - `make check`
- Test
  - `make test`
