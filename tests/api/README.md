# API Tests

API tests are implemented using [Bruno](https://www.usebruno.com/).

Run all API tests:

```shell
mkdir -p tmp
bru run --env [local|dev] --output tmp/all.json
cat tmp/all.json
```

## Smoke Tests

Sometimes it is useful to run a few quick smoke tests using `curl`.

Health:

```shell
curl -v -X GET http://localhost:4000/api/health \
  -H "Accept: application/json"
```

Read Nodes:

```shell
curl -v -X GET http://localhost:4000/api/v1/nodes \
  -H "Accept: application/json"
```

Read Pods:

```shell
curl -v -X GET http://localhost:4000/api/v1/pods \
  -H "Accept: application/json"
```

## Health

```shell
mkdir -p tmp
bru run --env [local|dev] 01_health --output tmp/01_health.json
cat tmp/01_health.json
```

## Nodes

```shell
mkdir -p tmp
bru run --env [local|dev] 02_nodes --output tmp/02_nodes.json
cat tmp/02_nodes.json
```

## Pods

```shell
mkdir -p tmp
bru run --env [local|dev] 03_pods --output tmp/03_pods.json
cat tmp/03_pods.json
```

## Kueue

```shell
mkdir -p tmp
bru run --env [local|dev] 04_kueue --output tmp/04_kueue.json
cat tmp/04_kueue.json
```

## KAI Scheduler

```shell
mkdir -p tmp
bru run --env [local|dev] 05_kai-scheduler --output
cat tmp/05_kai-scheduler.json
```
