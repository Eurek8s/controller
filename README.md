# Eurek8s Controller

![image](https://github.com/Eurek8s/controller/blob/main/docs/static/logo.png?raw=true)

Eurek8s is a controller to help using Netflix's Eureka with Kubernetes or just to help during the migration of Eureka to
another solution.

## Installing

### Helm

There is a [Helm Chart](https://github.com/Eurek8s/helm-charts/tree/main/charts/eurek8s-controller) to help you install
Eurek8s.

## Configuring

Eurek8s has only one configuration: the CONFIG environment variable.
This setting expects a JSON containging a map of Eureka clusters.
For example:

```
CONFIG='{"qa":["http://qa1.example.com","http://qa2.example.com"],"staging":["http://staging1.example.com"]}'
```

## Developing

### Running and deploying the controller

This project is created using Kubebuilder. You just need to follow the instructions to run Kubebuilder projects

https://book.kubebuilder.io/cronjob-tutorial/running.html

### Releasing new versions

Release Drafter will take care of the changelogs but the package itself (Docker image) will be tagged equals Github's
release name.

For example: To release a `v1.2.3` you should create a GitHub Release using `v.1.2.3` as the name

