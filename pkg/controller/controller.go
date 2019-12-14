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

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &ReconcileEvent{}

// Reconcile reconciles all events.
func (r *ReconcileEvent) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// set up a convenient log object so we don't have to type request over and over again
	log := r.Log.WithValues("request", request)

	if request == (reconcile.Request{}) {
		return reconcile.Result{}, nil
	}
	event := &corev1.Event{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, event)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find Events")
		return reconcile.Result{}, nil
	}

	if err != nil {
		log.Error(err, "Could not fetch Events")
		return reconcile.Result{}, err
	}

	// Output found events
	log.V(2).Info(fmt.Sprintf("found events: event %s with reason %s", event.Type, event.Reason))

	// Loop over provider
	for _, p := range r.Provider {
		p.EmitEvent(*event)
	}

	return reconcile.Result{}, nil
}
