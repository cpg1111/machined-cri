package nspawnlet

import (
	"k8s.io/kubernetes/pkg/kubelet"
)

// Deps are the dependencies for nspawnlet,
// inherits from kubelet's dependencies
type Deps struct {
	kubelet.KubeletDeps
}
