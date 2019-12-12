package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	otypes "github.com/afritzler/oopsie/pkg/types"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"

	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// reconcileReplicaSet reconciles ReplicaSets
type ReconcileEvent struct {
	// client can be used to retrieve objects from the APIServer.
	Client   client.Client
	Log      logr.Logger
	Recorder record.EventRecorder
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &ReconcileEvent{}

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

	if event.Type == corev1.EventTypeWarning && event.Message != "" {
		log.Info(fmt.Sprintf("found event with warning: event %s with reason %s", event.Type, event.Reason))
		req, err := http.NewRequest("GET", "https://api.stackexchange.com/2.2/search", nil)
		if err != nil {
			log.Error(err, "failed to construct request")
			return reconcile.Result{}, err
		}
		q := req.URL.Query()
		q.Add("order", "desc")
		q.Add("sort", "votes")
		q.Add("intitle", event.Message)
		q.Add("site", "stackoverflow")
		req.URL.RawQuery = q.Encode()

		resp, err := http.Get(req.URL.String())
		if err != nil {
			log.Error(err, "failed to query backend")
			return reconcile.Result{}, err
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err, "failed to get response body")
			return reconcile.Result{}, err
		}

		var anyJSON otypes.Answers
		json.Unmarshal(body, &anyJSON)

		if len(anyJSON.Items) > 0 && anyJSON.Items[0].Link != "" {
			link := anyJSON.Items[0].Link
			log.Info("Fired event")
			r.Recorder.Event(&event.InvolvedObject, v1.EventTypeNormal, "Hint", link)
		}

	}
	return reconcile.Result{}, nil
}
