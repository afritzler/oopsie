package main

import (
	"context"

	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// reconcileReplicaSet reconciles ReplicaSets
type reconcileReplicaSet struct {
	// client can be used to retrieve objects from the APIServer.
	client client.Client
	log    logr.Logger
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &reconcileReplicaSet{}

func (r *reconcileReplicaSet) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// set up a convenient log object so we don't have to type request over and over again
	log := r.log.WithValues("request", request)

	// Fetch the ReplicaSet from the cache
	rs := &appsv1.ReplicaSet{}
	err := r.client.Get(context.TODO(), request.NamespacedName, rs)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find ReplicaSet")
		return reconcile.Result{}, nil
	}

	if err != nil {
		log.Error(err, "Could not fetch ReplicaSet")
		return reconcile.Result{}, err
	}

	// Print the ReplicaSet
	log.Info("Reconciling ReplicaSet", "container name", rs.Spec.Template.Spec.Containers[0].Name)

	// Set the label if it is missing
	if rs.Labels == nil {
		rs.Labels = map[string]string{}
	}
	if rs.Labels["hello"] == "world" {
		return reconcile.Result{}, nil
	}

	// Update the ReplicaSet
	rs.Labels["hello"] = "world"
	err = r.client.Update(context.TODO(), rs)
	if err != nil {
		log.Error(err, "Could not write ReplicaSet")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}