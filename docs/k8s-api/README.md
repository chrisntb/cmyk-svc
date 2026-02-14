# K8s API Specifications

To generate K8s API specifications for a cluster, log into the controller and start the `proxy`:

```shell
kubectl proxy --port=8080
```

Then get the `index`:

```shell
curl http://127.0.0.1:8080/openapi/v3 \
  > k8s-openapi-v3_index.json
```

Inspect the `index` for the core API path:

```shell
curl http://127.0.0.1:8080/openapi/v3/api/v1 \
  > k8s-openapi-v3_core-v1.json
```

Inspect the `index` for other API paths, e.g. for Kueue:

```shell
curl http://127.0.0.1:8080/openapi/v3/apis/kueue.x-k8s.io/v1beta2 \
  > k8s-openapi-v3_kueue.x-k8s.io_v1beta2.json
```
