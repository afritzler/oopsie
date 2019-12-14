package provider

import (
	v1 "k8s.io/api/core/v1"
)

type Provider interface {
	EmitEvent(event v1.Event) error
}
