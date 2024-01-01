// Copyright 2024 Andreas Fritzler <afritzler@skiff.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	v1 "k8s.io/api/core/v1"
)

// Provider is an interface for the answers provider.
type Provider interface {
	EmitEvent(event v1.Event) error
	GetName() string
}
