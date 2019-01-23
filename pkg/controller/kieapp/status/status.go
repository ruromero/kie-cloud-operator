package status

import (
	"github.com/kiegroup/kie-cloud-operator/pkg/apis/app/v1"
	"github.com/kiegroup/kie-cloud-operator/pkg/controller/kieapp/logs"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var log = logs.GetLogger("kieapp.controller")

func SetProvisioning(cr *v1.KieApp) bool {
	log := log.With("kind", cr.Kind, "name", cr.Name, "namespace", cr.Namespace)
	conditions := len(cr.Status.Conditions)
	if conditions != 0 && cr.Status.Conditions[conditions-1].Type == v1.ProvisioningConditionType {
		log.Debug("Status: unchanged status [provisioning].")
		return false
	}
	log.Debug("Status: set provisioning")
	condition := v1.Condition{
		Type:               v1.ProvisioningConditionType,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
	}
	cr.Status.Conditions = append(cr.Status.Conditions, condition)
	return true
}

func SetDeployed(cr *v1.KieApp) bool {
	log := log.With("kind", cr.Kind, "name", cr.Name, "namespace", cr.Namespace)
	conditions := len(cr.Status.Conditions)
	if conditions != 0 && cr.Status.Conditions[conditions-1].Type == v1.DeployedConditionType {
		log.Debug("Status: unchanged status [deployed].")
		return false
	}
	log.Debug("Status: set deployed")
	condition := v1.Condition{
		Type:               v1.DeployedConditionType,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
	}
	cr.Status.Conditions = append(cr.Status.Conditions, condition)
	return true
}

func SetFailed(cr *v1.KieApp, reason v1.ReasonType, err error) {
	log := log.With("kind", cr.Kind, "name", cr.Name, "namespace", cr.Namespace)
	log.Debug("Status: set failed")
	condition := v1.Condition{
		Type:               v1.FailedConditionType,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            err.Error(),
	}
	cr.Status.Conditions = append(cr.Status.Conditions, condition)
}

func SetDeployments(cr *v1.KieApp, deployments []string) {
	log := log.With("kind", cr.Kind, "name", cr.Name, "namespace", cr.Namespace)
	log.Debug("Status: update deployments")
	cr.Status.Deployments = deployments
}
