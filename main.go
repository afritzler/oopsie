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

package main

import (
	"context"
	"os"

	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"

	oc "github.com/afritzler/oopsie/pkg/controller"
	op "github.com/afritzler/oopsie/pkg/provider"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func main() {

	log := log.FromContext(context.Background())
	entryLog := log.WithName("entrypoint")

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		panic("ooops failed to construct client")
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(log.Info)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: typedcorev1.New(clientset.CoreV1().RESTClient()).Events("")})
	stackRecorder := eventBroadcaster.NewRecorder(mgr.GetScheme(), corev1.EventSource{Component: "StackOverflow"})
	// TODO: register Github provider here
	providers := make([]op.Provider, 0)
	stackProvider := &op.StackOverflowProvider{
		Recorder: stackRecorder,
		Log:      log,
	}
	providers = append(providers, stackProvider)

	// Setup a new controller to reconcile
	entryLog.Info("Setting up controller")
	c, err := controller.New("oopsie-event-controller", mgr, controller.Options{
		Reconciler: &oc.ReconcileEvent{Client: mgr.GetClient(),
			Log:      log.WithName("reconciler"),
			Provider: providers,
		},
	})

	if err != nil {
		entryLog.Error(err, "unable to set up individual controller")
		os.Exit(1)
	}

	// Watch Events and enqueue Events object key
	if err := c.Watch(&source.Kind{Type: &corev1.Event{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "unable to watch events")
		os.Exit(1)
	}

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
