package main

import (
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	"os"

	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("example-controller")

func main() {
	logf.SetLogger(zap.Logger(false))
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
	recorder := eventBroadcaster.NewRecorder(mgr.GetScheme(), corev1.EventSource{Component: "StackOverflow :-)"})

	// Setup a new controller to reconcile ReplicaSets
	entryLog.Info("Setting up controller")
	c, err := controller.New("foo-controller", mgr, controller.Options{
		Reconciler: &reconcileReplicaSet{client: mgr.GetClient(), recorder: recorder, log: log.WithName("reconciler")},
	})
	if err != nil {
		entryLog.Error(err, "unable to set up individual controller")
		os.Exit(1)
	}

	//informer := cache.NewSharedIndexInformer(&cache.ListWatch{
	//	ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
	//		return clientset.CoreV1().Events("").List(options)
	//	},
	//	WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
	//		return clientset.CoreV1().Events("").Watch(options)
	//	},
	//},
	//&corev1.Event{},
	//0,
	//cache.Indexers{})

	// Watch Events and enqueue Events object key
	if err := c.Watch(&source.Kind{Type: &corev1.Event{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "unable to watch events")
		os.Exit(1)
	}

	//// Watch ReplicaSets and enqueue ReplicaSet object key
	//if err := c.Watch(&source.Kind{Type: &appsv1.ReplicaSet{}}, &handler.EnqueueRequestForObject{}); err != nil {
	//	entryLog.Error(err, "unable to watch ReplicaSets")
	//	os.Exit(1)
	//}
	//
	//// Watch Pods and enqueue owning ReplicaSet key
	//if err := c.Watch(&source.Kind{Type: &corev1.Pod{}},
	//	&handler.EnqueueRequestForOwner{OwnerType: &appsv1.ReplicaSet{}, IsController: true}); err != nil {
	//	entryLog.Error(err, "unable to watch Pods")
	//	os.Exit(1)
	//}

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}