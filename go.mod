module github.com/signalfx/splunk-otel-collector-operator

go 1.16

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/go-logr/logr v0.4.0
	github.com/golangci/golangci-lint v1.44.0
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.20.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.4
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.4
	k8s.io/kubectl v0.22.4
	sigs.k8s.io/controller-runtime v0.10.3
	sigs.k8s.io/controller-tools v0.7.0
	sigs.k8s.io/kustomize/kustomize/v4 v4.4.1
)
