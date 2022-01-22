# Eurek8s Controller

![image](https://github.com/Eurek8s/controller/blob/main/docs/static/logo.png?raw=true)

Eurek8s is a controller to help using Netflix's Eureka with Kubernetes or just to help during the migration of Eureka to
another solution.

## Developing

### Running and deploying the controller

This project is created using Kubebuilder. You just need to follow the instructions to run Kubebuilder projects

https://book.kubebuilder.io/cronjob-tutorial/running.html

### Releasing new versions

Release Drafter will take care of the changelogs but the package itself (Docker image) will be tagged equals Github's
release name.

For example: To release a `v1.2.3` you should create a GitHub Release using `v.1.2.3` as the name

