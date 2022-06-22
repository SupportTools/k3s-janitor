# k3s-janitor
This tool is designed to cleanup k3s/rke2 nodes by removing exited containers, unused images, and dangling volumes.

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/supporttools)](https://artifacthub.io/packages/search?repo=supporttools)


## Installation

### YAML manifest
```bash
curl -sSL https://raw.githubusercontent.com/k3s-janitor/k3s-janitor/main/deploy.yaml | kubectl apply -f -
```

#### Helm chart
```bash
helm repo add support-tools https://charts.support.tools
helm repo update
helm install k3s-janitor support-tools/k3s-janitor --namespace k3s-janitor
```

## Verify the installation
```bash
kubectl get rolloutstatus -w -n k3s-janitor
```
