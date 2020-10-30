#!/usr/bin/env bash

# Copyright (c) 2019 quay.io/mchirico Inc.
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

# deploy.sh
#
# Sets up the environment for the admission controller eventer demo in the active cluster.

set -euo pipefail

basedir="$(dirname "$0")/deployment"
#keydir="$(mktemp -d)"
keydir=tls-gen/basic/certs

# Generate keys into a temporary directory.
echo "Generating TLS keys ..."
"${basedir}/generate-keys.sh" "$keydir"

# Create the `eventer-demo` namespace. This cannot be part of the YAML file as we first need to create the TLS secret,
# which would fail otherwise.
echo "Creating Kubernetes objects ..."
kubectl delete namespace eventer-demo || true
kubectl create namespace eventer-demo

# Create the TLS secret for the generated keys.
kubectl -n eventer-demo create secret tls eventer-server-tls \
    --cert "${keydir}/server_certificate.pem" \
    --key "${keydir}/server_key.pem"

# Read the PEM-encoded CA certificate, base64 encode it, and replace the `${CA_PEM_B64}` placeholder in the YAML
# template with it. Then, create the Kubernetes resources.
ca_pem_b64="$(openssl base64 -A <"${keydir}/clientCom.pem")"
sed -e 's@${CA_PEM_B64}@'"$ca_pem_b64"'@g' <"${basedir}/deployment.yaml.template" \
    | kubectl create -f -

# Delete the key directory to prevent abuse (DO NOT USE THESE KEYS ANYWHERE ELSE).
#rm -rf "$keydir"

echo "The eventer server has been deployed and configured!"
