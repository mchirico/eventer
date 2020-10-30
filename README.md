

![CI/CD](https://github.com/mchirico/eventer/workflows/CI/CD/badge.svg)
[![codecov](https://codecov.io/gh/mchirico/client-go/branch/master/graph/badge.svg)](https://codecov.io/gh/mchirico/client-go)



# eventer

This project uses both a Mutating Admission Controller and
client-go to control the K8s environment.

## Steps

```bash



cd kind
make

# Stop. At this point, make sure everything is running.

cd ./..
./gencerts.sh
./deploy.sh






```bash
go get k8s.io/kubernetes || true

```
