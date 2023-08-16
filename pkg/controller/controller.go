// Copyright 2022 Andreas Fritzler
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"context"
	"fmt"

	"github.com/afritzler/oopsie/pkg/provider"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileEvent reconciles Events
type ReconcileEvent struct {
	Client   client.Client
	Scheme   *runtime.Scheme
	Provider []provider.Provider
}

// Reconcile reconciles all events.
func (r *ReconcileEvent) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := log.FromContext(ctx)

	if request == (reconcile.Request{}) {
		return reconcile.Result{}, nil
	}
	event := &corev1.Event{}
	if err := r.Client.Get(ctx, request.NamespacedName, event); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, fmt.Errorf("could not find events for %s: %w", request.NamespacedName, err)
		}
		return reconcile.Result{}, fmt.Errorf("could not get events for %s: %w", request.NamespacedName, err)
	}

	log.V(2).Info("found event", "Event", event)

	for _, p := range r.Provider {
		if err := p.EmitEvent(*event); err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to emit event for %s: %s", request.NamespacedName, err)
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileEvent) SetupWithManager(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Event{}).
		Complete(r)
}
