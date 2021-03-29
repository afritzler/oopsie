package main

import (
	"os"

	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"

	oc "github.com/afritzler/oopsie/pkg/controller"
	op "github.com/afritzler/oopsie/pkg/provider"
	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func main() {
	var log = log.Log.WithName("oopsie-controller")
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
