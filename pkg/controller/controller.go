package controller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/afritzler/oopsie/pkg/provider"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileEvent reconciles Events
type ReconcileEvent struct {
	// client can be used to retrieve objects from the APIServer.
	Client   client.Client
	Log      logr.Logger
	Provider []provider.Provider
}

// Reconcile reconciles all events.
func (r *ReconcileEvent) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// set up a convenient log object so we don't have to type request over and over again
	log := r.Log.WithValues("request", request)

	if request == (reconcile.Request{}) {
		return reconcile.Result{}, nil
	}
	event := &corev1.Event{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, event)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, fmt.Errorf("could not find events for %s: %s", request.NamespacedName, err)
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not get events for %s: %s", request.NamespacedName, err)
	}

	// Output found events
	log.V(2).Info(fmt.Sprintf("found events: event %s with reason %s", event.Type, event.Reason))

	// Loop over provider
	for _, p := range r.Provider {
		if err := p.EmitEvent(*event); err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to emit event for %s: %s", request.NamespacedName, err)
		}
	}

	return reconcile.Result{}, nil
}
