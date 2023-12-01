# Releases Playground

This is a playground for using the [CircleCI Releases feature](https://circleci.com/docs/release/releases-overview/).

The demo app is the snippetbox from the Let's Go book, migrated to kubernetes from docker-compose.

## Prerequisites

1) [minikube](https://minikube.sigs.k8s.io/docs/start/)
2) [kubectl](https://kubernetes.io/docs/tasks/tools/)
3) [helm](https://helm.sh/docs/intro/install/)

## How To Setup

1) Start minikube

```bash
minikube start
```

2) Apply k8s deployment

```bash
kubectl apply -f snippetbox_kubernetes.yaml
```

3) Confirm pods started

```bash
kubectl get pods
```

4) Install Releases chart, if needed

```bash
helm repo add release-agent https://circleci-public.github.io/cci-k8s-release-agent
```

5) Update helm chart

```bash
helm repo update
```

6) Generate release token
7) Install the release agent

```bash
helm upgrade --install circleci-release-agent-system release-agent/circleci-release-agent \
--set tokenSecret.token=[YOUR_CCI_INTEGRATION_TOKEN] --create-namespace \
--namespace circleci-release-agent-system
```

8) Walk through required and option labels in k8s deployment
9) Confirm components appear in Releases tab
10) Change component version and confirm new version appears in Releases tab