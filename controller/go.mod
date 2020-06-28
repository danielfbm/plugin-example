module github.com/danielfbm/plugin-example/controller

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/golang/protobuf v1.3.4
	github.com/hashicorp/go-hclog v0.14.0
	github.com/hashicorp/go-plugin v1.3.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	golang.org/x/net v0.0.0-20191004110552-13f9640d40b9
	google.golang.org/grpc v1.27.1
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime v0.5.0
)

replace github.com/hashicorp/go-plugin v1.3.0 => github.com/danielfbm/go-plugin v1.3.1-0.20200526082916-ab3b911f678b
