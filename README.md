Kubernetes rolling update safely
================================

This file is generated by *README.md.tmpl*.
Run `go run readme.go -o README.md` to generate README.md.

## Demo environment

```
minikube version: v0.28.0
```

## Preparations

```shell
# start minikube
$ minikube start

# change docker host to minikube
$ eval $(minikube docker-env)

# build demo application
$ cd demo
$ docker build -t takitake/demo demo
$ cd -
```

## Demo scenarios

### Default setting

Deploy simple spring-boot application.

This is deployment manifest.
```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: demo
  name: demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: demo
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: demo
      annotations:
        appVersion: v1.0
    spec:
      containers:
      - image: takitake/demo
        imagePullPolicy: IfNotPresent
        name: demo
        resources: {}
status: {}
```

Deploy and watch Pod status.
```shell
$ kubectl apply -f 1-1.default.deploy.yml && kubectl get pods -w
```

```
20c20
<         appVersion: v1.0
---
>         appVersion: v1.1
```
