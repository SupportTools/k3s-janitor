# k3s-janitor
This tool is designed to cleanup k3s/rke2 nodes by removing exited containers, unused images, and dangling volumes.

## Stats

### Drone
[![Build Status](https://drone.support.tools/api/badges/SupportTools/k3s-janitor/status.svg?ref=refs/heads/main)](https://drone.support.tools/SupportTools/k3s-janitor)

### HELM Chart
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/supporttools)](https://artifacthub.io/packages/search?repo=supporttools)

### GitHub
![GitHub stars](https://img.shields.io/github/stars/supporttools/k3s-janitor?style=social)
![GitHub followers](https://img.shields.io/github/followers/mattmattox?style=social)
![GitHub Org's stars](https://img.shields.io/github/stars/supporttools?style=social)
![GitHub last commit](https://img.shields.io/github/last-commit/supporttools/k3s-janitor)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/supporttools/k3s-janitor)

### Docker Hub
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/supporttools/k3s-janitor)
![Docker Pulls](https://img.shields.io/docker/pulls/supporttools/k3s-janitor)
![Docker Stars](https://img.shields.io/docker/stars/supporttools/k3s-janitor)

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
