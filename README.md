# App Cleanup Operator
[![Release](https://github.com/sqaisar/app-cleanup-operator/actions/workflows/publish-image.yml/badge.svg)](https://github.com/sqaisar/app-cleanup-operator/actions/workflows/publish-image.yml)

This Kubernetes operator performs cleanup tasks automatically, helping simplify maintenance and operational tasks within your cluster.

---

## 🚀 Installation Guide

Follow these simple steps to install **App Cleanup Operator** into your Kubernetes cluster.

### 🔧 Install with Kubectl

Define the desired Controller version and execute the deployment manifest via `kubectl`:

```bash
# Define the operator version you want to install (replace this with your desired version)
CONTROLLER_VERSION=v1.0.5

# Install the App Cleanup Operator into your Kubernetes cluster
kubectl apply -f https://github.com/sqaisar/app-cleanup-operator/releases/download/${CONTROLLER_VERSION}/operator.yaml
```

**Explanation:**
- `CONTROLLER_VERSION` sets the specific version of the controller you desire.
- The `kubectl apply` downloads and applies the operator's YAML manifest directly from official releases.

---

## ✅ Verifying the Installation

Check clearly if the operator pods are running successfully after installation:

```bash
kubectl get pods -n <namespace>
```

All pods related to the operator should display status `Running` and ready.

---

## 📚 Updating the Installed Version

Check the available versions at [GitHub Releases](https://github.com/sqaisar/app-cleanup-operator/releases).

When updating clearly define your new `CONTROLLER_VERSION` and re-apply manifests as shown:

```bash
# Define new version
CONTROLLER_VERSION=<new-version>

# Upgrade the operator
kubectl apply -f https://github.com/sqaisar/app-cleanup-operator/releases/download/${CONTROLLER_VERSION}/operator.yaml
```

---

## 🛠️ Customizing and Configuring

(List clearly any additional options or configuration tweaks users might prefer or require here.)

---

## 🙋 Getting Help

If you encounter issues or need guidance:
- Open an Issue clearly at the GitHub repository: [Issues Page](https://github.com/sqaisar/app-cleanup-operator/issues).
- Consult clearly documented troubleshooting steps in this README.

---

## Contributing guidelines

### Prerequisites
- go version v1.23.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/app-cleanup-operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/app-cleanup-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/app-cleanup-operator:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/app-cleanup-operator/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

