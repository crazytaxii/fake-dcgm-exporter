# fake-dcgm-exporter ![CI Status](https://github.com/crazytaxii/fake-dcgm-exporter/actions/workflows/publish-image.yml/badge.svg)

A simulated [DCGM Exporter](https://github.com/NVIDIA/dcgm-exporter) provides fake metrics when development environment doesn't have any NVIDIA GPU sometimes.

The fake-dcgm-exporter simulates **8(by default) NVIDIA** GPUs with builtin driver version and hardware data.

The fake-dcgm-exporter supports 3 types of **NVIDIA** GPU:

- A100-SXM4-40GB(**default**): `A100` in configuration
- GeForce RTX 4090: `RTX4090` in configuration
- A800-SXM4-80GB: `A800` in configuration

## Quickstart on Kubernetes

Modify the configuration file according to your requirements.

Deploy with:

```bash
$ kubectl create -f manifests/all.yaml
```

Test with:

```bash
$ SERVICE_IP=$(kubectl get svc -n monitoring dcgm-exporter -o jsonpath='{.spec.clusterIP}')
$ curl http://${SERVICE_IP}:9400/metrics
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
```
