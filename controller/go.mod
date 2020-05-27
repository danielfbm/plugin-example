module github.com/danielfbm/plugin-example/controller

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/hashicorp/go-hclog v0.0.0-20180709165350-ff2cf002a8dd
	github.com/hashicorp/go-plugin v1.3.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime v0.5.0
)

replace github.com/hashicorp/go-plugin v1.3.0 => github.com/danielfbm/go-plugin v1.3.1-0.20200526082916-ab3b911f678b
