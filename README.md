# golang plugin example

## Intro

1. Create a simple Plugin for a go application
2. Analyse [Terraform's plugin implementation](https://www.terraform.io/docs/plugins/basics.html)
3. Enabling a rich plugin ecosystem on Kubernetes

## Simple plugin for a go application

Uses [Hashicorp's go-plugin](https://github.com/hashicorp/go-plugin.git) to create a `hello-world` kind of plugin

More [here](demo/README.md)

## Sample controller that uses plugins with `kubebuilder`

- User [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) to create a controller
- Creates a local and network plugin abstraction
- Checks and executes plugins on a controller

More [here](controller/README.md)
