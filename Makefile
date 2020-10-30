# Copyright (c) 2020 CWXSTAT Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Makefile for building the Admission Controller webhook demo server + docker image.

.DEFAULT_GOAL := docker-image

IMAGE ?= quay.io/mchirico/admission-controller-eventer-demo:v1

image/webhook-server: $(shell find ./mutating-admission-controller-k8s-go -name '*.go')
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $@ ./mutating-admission-controller-k8s-go/cmd/webhook-server

.PHONY: docker-image
docker-image: image/webhook-server
	docker build -t $(IMAGE) image/
	kind load docker-image $(IMAGE)
	kubectl delete -f mutating-admission-controller-k8s-go/examples/pod-with-defaults.yaml || true
	sleep 4
	kubectl apply -f mutating-admission-controller-k8s-go/examples/pod-with-defaults.yaml

.PHONY: push-image
push-image: docker-image
	docker push $(IMAGE)
