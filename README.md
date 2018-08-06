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
$ docker build -t takitake/demo demo
```

## Demo scenarios

### Default setting

Deploy simple spring-boot application.

```
$ cat demo-manifest/1-0.default.deploy.yml
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

Deploy and watch Pod status. You can see the state transition *ContainerCreating* to *Running*.

```sh
$ kubectl apply -f demo-manifest/1-0.default.deploy.yml && kubectl get pods -w
deployment.extensions "demo" created
NAME                    READY     STATUS              RESTARTS   AGE
demo-67fb9964f4-xmrwj   0/1       ContainerCreating   0          0s
demo-67fb9964f4-xmrwj   1/1       Running   0         1s
```

Let's create the Service to access the deployed demo application.

```sh
$ kubectl apply -f demo-manifest/service.yml
service "demo" created

$ kubectl get service demo
NAME      TYPE       CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
demo      NodePort   10.111.204.139   <none>        8080:32569/TCP   23s

$ minikube status
minikube: Running
cluster: Running
kubectl: Correctly Configured: pointing to minikube-vm at 192.168.99.100
```

Now you can access the demo application with Kubernetes Service. IP should be the minikube nodeIP and Port should be the NodePort of the demo Service.

This case URL is [http://192.168.99.100:32569/](http://192.168.99.100:32569/).

Next, in order to check the rolling update, try changing the appVersion from v1.0 to v1.1

```sh
$ diff "demo-manifest/1-0.default.deploy.yml" "demo-manifest/1-1.default.deploy.yml"
20c20
<         appVersion: v1.0
---
>         appVersion: v1.1
```

```sh
# If you want to check http status also. Following command is usuful to check the status every 1 sec.
$ while do curl -s -o /dev/null -w "`date` -- %{http_code}\n" http://192.168.99.100:32569/; sleep 1s; done
```

```sh
$ kubectl apply -f demo-manifest/1-1.default.deploy.yml && kubectl get pods -w
deployment.extensions "demo" configured
NAME                    READY     STATUS        RESTARTS   AGE
demo-67fb9964f4-xmrwj   1/1       Terminating   0          1m
demo-7744c47967-c2v9m   0/1       Pending       0          0s
demo-7744c47967-c2v9m   0/1       ContainerCreating   0         0s
demo-67fb9964f4-xmrwj   0/1       Terminating   0         1m
demo-67fb9964f4-xmrwj   0/1       Terminating   0         1m
demo-7744c47967-c2v9m   1/1       Running   0         2s
demo-67fb9964f4-xmrwj   0/1       Terminating   0         1m
demo-67fb9964f4-xmrwj   0/1       Terminating   0         1m
```

After running Pod was terminated, new Pod was created..

Because `maxUnavailable` is 1 by default so only one running Pod was terminated soon.

Dividing the Pod's status into four, the status will change as follows.

![v1.1 deploy timeline](img/1-1.png)

There are two solusitions that changing the value to 0 or increasing number of replicas.
Changing the maxUnavailable to 0 this time to simplify the explanation.

### Changing strategy of rolling update

```
$ diff "demo-manifest/1-1.default.deploy.yml" "demo-manifest/1-2.strategy.deploy.yml"
12a13,16
>   strategy:
>     rollingUpdate:
>       maxSurge: 1
>       maxUnavailable: 0
20c24
<         appVersion: v1.1
---
>         appVersion: v1.2
```

```sh
$ kubectl apply -f demo-manifest/1-2.strategy.deploy.yml && kubectl get pods -w
deployment.extensions "demo" configured
NAME                    READY     STATUS        RESTARTS   AGE
demo-644fd4dcb-pthk5    0/1       Pending       0          0s
demo-7744c47967-7vbvv   1/1       Terminating   0          53s
demo-644fd4dcb-pthk5   0/1       Pending   0         0s
demo-644fd4dcb-pthk5   0/1       ContainerCreating   0         0s
demo-7744c47967-7vbvv   0/1       Terminating   0         54s
demo-644fd4dcb-pthk5   1/1       Running   0         2s
demo-7744c47967-7vbvv   0/1       Terminating   0         1m
demo-7744c47967-7vbvv   0/1       Terminating   0         1m
```

hmm, still we can lose user's request.

Firstly new Pod was created but old Pod was terminated before new Pod became Running status.

![v1.2 deploy timeline](img/1-2.png)
