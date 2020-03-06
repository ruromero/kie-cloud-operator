package controller

import (
	"github.com/RHsyseng/operator-utils/pkg/utils/kubernetes"
	api "github.com/kiegroup/kie-cloud-operator/pkg/apis/app/v2"
	"github.com/kiegroup/kie-cloud-operator/pkg/controller/kieapp"
	"github.com/kiegroup/kie-cloud-operator/pkg/controller/kieapp/logs"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var log = logs.GetLogger("kie-cloud-operator.controller")

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	addManager := func(mgr manager.Manager) error {
		k8sService := GetInstance(mgr)
		reconciler := kieapp.Reconciler{Service: &k8sService}
		extReconciler := kubernetes.NewExtendedReconciler(&k8sService, &reconciler, &api.KieApp{})
		err := extReconciler.RegisterFinalizer(&kieapp.ConsoleLinkFinalizer{})
		if err != nil {
			log.Errorf("Unable to register finalizer. ", err)
		}
		return kieapp.Add(mgr, &reconciler)
	}
	AddToManagerFuncs = []func(manager.Manager) error{addManager}
}
