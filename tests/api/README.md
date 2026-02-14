# API Tests

API tests are implemented using [Bruno](https://www.usebruno.com/).

Run all API tests:

```shell
mkdir -p tmp
bru run --env [local|dev] --output tmp/results.json
```

Run one set of API test:

```shell
mkdir -p tmp
bru run --env [local|dev] 01_health --output tmp/01_health.json
```

Run one API test:

```shell
mkdir -p tmp
bru run --env [local|dev] 01_health/01_ReadHealth.bru --output tmp/01_ReadHealth.bru.json
```

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
