![Awesome grid logo](https://raw.githubusercontent.com/qa-kit/awesome-grid/master/doc/assets/logo.png)

# Awesome grid [![Build Status](https://travis-ci.com/qa-kit/awesome-grid.svg?branch=master)](https://travis-ci.com/qa-kit/awesome-grid) [![Go Report Card](https://goreportcard.com/badge/github.com/qa-kit/awesome-grid)](https://goreportcard.com/report/github.com/qa-kit/awesome-grid) [![codecov](https://codecov.io/gh/qa-kit/awesome-grid/branch/master/graph/badge.svg)](https://codecov.io/gh/qa-kit/awesome-grid)

Project provides simple tool for launching your UI-tests with selenium-base docker images in kubernetes cluster.

For every test tool will create unique pod and after test finished pod will be removed.

## Quick start with minikube
* Install minikube https://kubernetes.io/docs/tasks/tools/install-minikube/
* Start minikube

  `minikube start`
* Clone source

  `git clone git@github.com:qa-kit/awesome-grid.git`
* Change directory to repo

  `cd awesome-grid`
* Apply default permissions of cluster

  `kube apply -f build/kube`
* Build docker image

  ```
  eval $(minikube docker-env)
  docker build . -t awesome-grid
  ```
* Start deployment

  ```
  kubectl run awesome-grid --image=awesome-grid --labels="app=awesome-grid" --image-pull-policy=Never --port 4444
  ```
* Get url of pod

  `minikube service awesome-grid --url`
* Start tests with this url.

## Building
Building a docker image and
```
docker build . -t registry.example.com/awesome-grid
```
After that you could push image to your private registry and use it in k8s cluster
```
docker push registry.example.com/awesome-grid
```

## Features
* Flexible JSON-config for grid
* JSON template for deployment template for k8s
* On-demand grid scaling
* Complete end-to-end tests provide stability of current release
* Opportunity of tuning selenium images for best performance.

## Contributing
If you'd like to contribute, please fork the repository and use a feature branch. Pull requests are warmly welcome.

We try to follow and hope you will too:
* Codestyle by `gofmt`
* Effective go https://golang.org/doc/effective_go.html
* Best practice of project layout https://github.com/golang-standards/project-layout
* Use codecov https://codecov.io/gh/qa-kit/awesome-grid
* Travis-ci tests and builds https://travis-ci.com/qa-kit/awesome-grid/

## Links

- Repository: https://github.com/qa-kit/awesome-grid/
- Issue tracker: https://github.com/qa-kit/awesome-grid/issues

## Licensing

The code in this project is licensed under MIT license.
