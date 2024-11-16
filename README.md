# fake-dcgm-exporter ![CI Status](https://github.com/crazytaxii/fake-dcgm-exporter/actions/workflows/publish-image.yml/badge.svg)

A simulated [DCGM Exporter](https://github.com/NVIDIA/dcgm-exporter) provides fake metrics when development environment doesn't have any NVIDIA GPU sometimes.

## Quickstart on Kubernetes

Modify the configuration file according to your requirements. The fake-dcgm-exporter will simulate **8** *NVIDIA A100-SXM4-40GB* GPUs with driver version *535.104.12* by default.

Deploy with:

```bash
$ kubectl create -f manifests/all.yaml
```

Test with:

```bash
$ SERVICE_IP=$(kubectl get svc -n monitoring dcgm-exporter -o jsonpath='{.spec.clusterIP}')
$ curl http://${SERVICE_IP}:9400/metrics

# or enable gzip compressed
$ curl --compressed http://${SERVICE_IP}:9400/metrics
```

## Quickstart on Docker

Deploy with:

```bash
$ mkdir -p /etc/fake-dcgm-exporter
$ curl -sL -o /etc/fake-dcgm-exporter/config.yaml https://raw.githubusercontent.com/crazytaxii/fake-dcgm-exporter/refs/heads/main/config.yaml
# Modify the configuration file according to your requirements.

$ docker run --name=fake-dcgm-exporter --rm -d --network=host -v /etc/fake-dcgm-exporter:/etc/fake-dcgm-exporter crazytaxii/fake-dcgm-exporter:latest
```

Test with:

```bash
$ curl http://127.0.0.1:9400/metrics

# or enable gzip compressed
$ curl --compressed http://127.0.0.1:9400/metrics
```
