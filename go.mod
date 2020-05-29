module github.com/afritzler/oopsie

go 1.14

require (
	github.com/go-logr/logr v0.1.0
	github.com/prometheus/common v0.10.0
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v11.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.0
)

replace k8s.io/client-go => k8s.io/client-go v0.18.3
