package status

import (
	"errors"
	"testing"

	"github.com/kiegroup/kie-cloud-operator/pkg/apis/app/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSetDeployed(t *testing.T) {
	now := metav1.Now()
	cr := &v1.KieApp{}

	assert.True(t, SetDeployed(cr))

	assert.NotEmpty(t, cr.Status.Conditions)
	assert.Equal(t, v1.DeployedConditionType, cr.Status.Conditions[0].Type)
	assert.Equal(t, corev1.ConditionTrue, cr.Status.Conditions[0].Status)
	assert.True(t, now.Before(&cr.Status.Conditions[0].LastTransitionTime))
}

func TestSetDeployedSkipUpdate(t *testing.T) {
	cr := &v1.KieApp{}
	SetDeployed(cr)

	assert.NotEmpty(t, cr.Status.Conditions)
	condition := cr.Status.Conditions[0]

	assert.False(t, SetDeployed(cr))
	assert.Equal(t, 1, len(cr.Status.Conditions))
	assert.Equal(t, condition, cr.Status.Conditions[0])
}

func TestSetProvisioning(t *testing.T) {
	now := metav1.Now()
	cr := &v1.KieApp{}
	assert.True(t, SetProvisioning(cr))

	assert.NotEmpty(t, cr.Status.Conditions)
	assert.Equal(t, 2, len(cr.Status.Conditions))

	provIdx := getConditionIdx(cr, v1.ProvisioningConditionType)
	assert.NotEqual(t, -1, provIdx)
	assert.Equal(t, corev1.ConditionTrue, cr.Status.Conditions[provIdx].Status)
	assert.True(t, now.Before(&cr.Status.Conditions[provIdx].LastTransitionTime))

	deployedIdx := getConditionIdx(cr, v1.DeployedConditionType)
	assert.NotEqual(t, -1, deployedIdx)
	assert.Equal(t, corev1.ConditionFalse, cr.Status.Conditions[deployedIdx].Status)
	assert.True(t, now.Before(&cr.Status.Conditions[deployedIdx].LastTransitionTime))
}

func TestSetProvisioningSkipUpdate(t *testing.T) {
	cr := &v1.KieApp{}
	assert.True(t, SetProvisioning(cr))

	assert.NotEmpty(t, cr.Status.Conditions)
	condition := cr.Status.Conditions[0]

	assert.False(t, SetProvisioning(cr))
	assert.Equal(t, 2, len(cr.Status.Conditions))
	assert.Equal(t, condition, cr.Status.Conditions[0])

	deployedIdx := getConditionIdx(cr, v1.DeployedConditionType)
	assert.NotEqual(t, -1, deployedIdx)
	assert.Equal(t, corev1.ConditionFalse, cr.Status.Conditions[deployedIdx].Status)
}

func TestSetProvisioningAndThenDeployed(t *testing.T) {
	now := metav1.Now()
	cr := &v1.KieApp{}

	assert.True(t, SetProvisioning(cr))
	assert.True(t, SetDeployed(cr))

	assert.NotEmpty(t, cr.Status.Conditions)
	condition := cr.Status.Conditions[0]
	assert.Equal(t, 1, len(cr.Status.Conditions))
	assert.True(t, now.Before(&condition.LastTransitionTime))
	assert.Equal(t, v1.DeployedConditionType, condition.Type)
	assert.Equal(t, corev1.ConditionTrue, condition.Status)
}

func TestSetProvisioningAndThenFailed(t *testing.T) {
	now := metav1.Now()
	cr := &v1.KieApp{}

	assert.True(t, SetProvisioning(cr))
	SetFailed(cr, v1.DeploymentFailedReason, errors.New("Test"))

	assert.NotEmpty(t, cr.Status.Conditions)
	assert.Equal(t, 3, len(cr.Status.Conditions))
	provIdx := getConditionIdx(cr, v1.ProvisioningConditionType)
	assert.NotEqual(t, -1, provIdx)
	assert.Equal(t, corev1.ConditionTrue, cr.Status.Conditions[provIdx].Status)
	assert.True(t, now.Before(&cr.Status.Conditions[provIdx].LastTransitionTime))

	deployedIdx := getConditionIdx(cr, v1.DeployedConditionType)
	assert.NotEqual(t, -1, deployedIdx)
	assert.Equal(t, corev1.ConditionFalse, cr.Status.Conditions[deployedIdx].Status)
	assert.True(t, now.Before(&cr.Status.Conditions[deployedIdx].LastTransitionTime))

	failedIdx := getConditionIdx(cr, v1.FailedConditionType)
	assert.NotEqual(t, -1, deployedIdx)
	assert.Equal(t, corev1.ConditionTrue, cr.Status.Conditions[failedIdx].Status)
	assert.True(t, now.Before(&cr.Status.Conditions[failedIdx].LastTransitionTime))
	assert.Equal(t, v1.DeploymentFailedReason, cr.Status.Conditions[failedIdx].Reason)
	assert.Equal(t, "Test", cr.Status.Conditions[failedIdx].Message)
}

func TestSetFailedAndAgainFailed(t *testing.T) {
	now := metav1.Now()
	cr := &v1.KieApp{}

	SetFailed(cr, v1.DeploymentFailedReason, errors.New("Test 1"))
	SetFailed(cr, v1.DeploymentFailedReason, errors.New("Test 2"))

	assert.NotEmpty(t, cr.Status.Conditions)
	assert.Equal(t, 2, len(cr.Status.Conditions))

	deployedIdx := getConditionIdx(cr, v1.DeployedConditionType)
	assert.NotEqual(t, -1, deployedIdx)
	assert.Equal(t, corev1.ConditionFalse, cr.Status.Conditions[deployedIdx].Status)
	assert.True(t, now.Before(&cr.Status.Conditions[deployedIdx].LastTransitionTime))

	failedIdx := getConditionIdx(cr, v1.FailedConditionType)
	assert.NotEqual(t, -1, deployedIdx)
	assert.Equal(t, corev1.ConditionTrue, cr.Status.Conditions[failedIdx].Status)
	assert.True(t, now.Before(&cr.Status.Conditions[failedIdx].LastTransitionTime))
	assert.Equal(t, v1.DeploymentFailedReason, cr.Status.Conditions[failedIdx].Reason)
	assert.Equal(t, "Test 2", cr.Status.Conditions[failedIdx].Message)
}

func TestSetDeployedAndThenProvisioning(t *testing.T) {
	now := metav1.Now()
	cr := &v1.KieApp{}

	assert.True(t, SetDeployed(cr))
	assert.True(t, SetProvisioning(cr))

	assert.NotEmpty(t, cr.Status.Conditions)
	provIdx := getConditionIdx(cr, v1.ProvisioningConditionType)
	assert.NotEqual(t, -1, provIdx)
	assert.Equal(t, corev1.ConditionTrue, cr.Status.Conditions[provIdx].Status)
	assert.True(t, now.Before(&cr.Status.Conditions[provIdx].LastTransitionTime))

	deployedIdx := getConditionIdx(cr, v1.DeployedConditionType)
	assert.NotEqual(t, -1, deployedIdx)
	assert.Equal(t, corev1.ConditionFalse, cr.Status.Conditions[deployedIdx].Status)
	assert.True(t, now.Before(&cr.Status.Conditions[deployedIdx].LastTransitionTime))
}
