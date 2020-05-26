module github.com/danielfbm/plugin-example/demo

go 1.14

require (
	github.com/hashicorp/go-hclog v0.13.0
	github.com/hashicorp/go-plugin v1.3.0
)

replace github.com/hashicorp/go-plugin v1.3.0 => github.com/danielfbm/go-plugin v1.3.1-0.20200526082916-ab3b911f678b
