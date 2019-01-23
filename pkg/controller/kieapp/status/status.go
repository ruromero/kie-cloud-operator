package status

import (
	"github.com/kiegroup/kie-cloud-operator/pkg/apis/app/v1"
	"github.com/kiegroup/kie-cloud-operator/pkg/controller/kieapp/logs"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var log = logs.GetLogger("kieapp.controller")

// SetProvisioning - Sets the condition type to Provisioning and status True if not yet set.
// Existing previous errors won't be removed until SetDeployed is invoked.
func SetProvisioning(cr *v1.KieApp) bool {
	log := log.With("kind", cr.Kind, "name", cr.Name, "namespace", cr.Namespace)
	condIdx := getConditionIdx(cr, v1.ProvisioningConditionType)
	if condIdx != -1 {
		log.Debug("Status: unchanged status [provisioning].")
		return false
	}
	setDeployed(cr, false)
	log.Debug("Status: set provisioning")
	condition := v1.Condition{
		Type:               v1.ProvisioningConditionType,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
	}
	if condIdx == -1 {
		cr.Status.Conditions = append(cr.Status.Conditions, condition)
	}
	return true
}

// SetDeployed - Updates (if changed) the condition with the DeployedCondition and True status
func SetDeployed(cr *v1.KieApp) bool {
	return setDeployed(cr, true)
}

// SetFailed - Sets the failed condition to the status
func SetFailed(cr *v1.KieApp, reason v1.ReasonType, err error) {
	log := log.With("kind", cr.Kind, "name", cr.Name, "namespace", cr.Namespace)
	log.Debug("Status: set failed")
	setDeployed(cr, false)
	condition := v1.Condition{
		Type:               v1.FailedConditionType,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            err.Error(),
	}
	condIdx := getConditionIdx(cr, v1.FailedConditionType)
	if condIdx == -1 {
		cr.Status.Conditions = append(cr.Status.Conditions, condition)
	} else {
		cr.Status.Conditions[condIdx] = condition
	}
}

// SetDeployments - sets the deployment names to the status
func SetDeployments(cr *v1.KieApp, deployments []string) {
	log := log.With("kind", cr.Kind, "name", cr.Name, "namespace", cr.Namespace)
	log.Debug("Status: update deployments")
	cr.Status.Deployments = deployments
}

// setDeployed Adds the DeployedCondition if it doesn't exist or replaces it by
// the previous one only if the status is different
func setDeployed(cr *v1.KieApp, isDeployed bool) bool {
	log := log.With("kind", cr.Kind, "name", cr.Name, "namespace", cr.Namespace)
	status := corev1.ConditionFalse
	if isDeployed {
		status = corev1.ConditionTrue
	}
	condIdx := getConditionIdx(cr, v1.DeployedConditionType)
	if condIdx != -1 && cr.Status.Conditions[condIdx].Status == status {
		log.Debugf("Status: unchanged status [deployed:%v].", status)
		return false
	}
	log.Debugf("Status: changed status [deployed:%v].", status)

	condition := v1.Condition{
		Type:               v1.DeployedConditionType,
		Status:             status,
		LastTransitionTime: metav1.Now(),
	}
	if isDeployed {
		cr.Status.Conditions = cr.Status.Conditions[:0]
	}
	if isDeployed || condIdx == -1 {
		cr.Status.Conditions = append(cr.Status.Conditions, condition)
	} else {
		cr.Status.Conditions[condIdx] = condition
	}
	return true
}

func getConditionIdx(cr *v1.KieApp, condType v1.ConditionType) int {
	for i, condition := range cr.Status.Conditions {
		if condition.Type == condType {
			return i
		}
	}
	return -1
}
