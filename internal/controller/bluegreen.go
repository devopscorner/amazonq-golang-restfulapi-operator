package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	restapiv1 "github.com/devopscorner/restapi-operator/api/v1"
)

type BlueGreenManager struct {
	client.Client
	Scheme *runtime.Scheme
}

func (bg *BlueGreenManager) ReconcileBlueGreen(ctx context.Context, restAPI *restapiv1.RestAPI) error {
	if restAPI.Spec.BlueGreen == nil || !restAPI.Spec.BlueGreen.Enabled {
		return nil
	}

	components := []string{"model", "view", "controller", "repository"}

	for _, component := range components {
		if err := bg.manageBlueGreenDeployment(ctx, restAPI, component); err != nil {
			return err
		}
	}

	return nil
}

func (bg *BlueGreenManager) manageBlueGreenDeployment(ctx context.Context, restAPI *restapiv1.RestAPI, component string) error {
	activeEnv := restAPI.Status.ActiveEnvironment
	if activeEnv == "" {
		activeEnv = "blue"
	}

	// Create both blue and green deployments
	for _, env := range []string{"blue", "green"} {
		deployment := bg.createBlueGreenDeployment(restAPI, component, env)
		if err := controllerutil.SetControllerReference(restAPI, deployment, bg.Scheme); err != nil {
			return err
		}

		found := &appsv1.Deployment{}
		err := bg.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			if err := bg.Create(ctx, deployment); err != nil {
				return err
			}
		}
	}

	// Manage traffic switching
	return bg.switchTraffic(ctx, restAPI, component, activeEnv)
}

func (bg *BlueGreenManager) createBlueGreenDeployment(restAPI *restapiv1.RestAPI, component, environment string) *appsv1.Deployment {
	var spec restapiv1.ComponentSpec
	switch component {
	case "model":
		spec = restAPI.Spec.Model
	case "view":
		spec = restAPI.Spec.View
	case "controller":
		spec = restAPI.Spec.Controller
	case "repository":
		spec = restAPI.Spec.Repository
	}

	replicas := int32(1)
	if restAPI.Spec.Replicas != nil {
		replicas = *restAPI.Spec.Replicas
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-%s", restAPI.Name, component, environment),
			Namespace: restAPI.Namespace,
			Labels: map[string]string{
				"app":         restAPI.Name,
				"component":   component,
				"environment": environment,
				"version":     "v1",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":         restAPI.Name,
					"component":   component,
					"environment": environment,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":         restAPI.Name,
						"component":   component,
						"environment": environment,
						"version":     "v1",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  component,
						Image: spec.Image,
						Ports: []corev1.ContainerPort{{
							ContainerPort: spec.Port,
							Protocol:      corev1.ProtocolTCP,
						}},
					}},
				},
			},
		},
	}
}

func (bg *BlueGreenManager) switchTraffic(ctx context.Context, restAPI *restapiv1.RestAPI, component, activeEnv string) error {
	serviceName := fmt.Sprintf("%s-%s-svc", restAPI.Name, component)

	service := &corev1.Service{}
	err := bg.Get(ctx, types.NamespacedName{Name: serviceName, Namespace: restAPI.Namespace}, service)
	if err != nil {
		return err
	}

	// Update service selector to point to active environment
	service.Spec.Selector = map[string]string{
		"app":         restAPI.Name,
		"component":   component,
		"environment": activeEnv,
	}

	return bg.Update(ctx, service)
}

func (bg *BlueGreenManager) PromoteDeployment(ctx context.Context, restAPI *restapiv1.RestAPI) error {
	currentActive := restAPI.Status.ActiveEnvironment
	if currentActive == "" {
		currentActive = "blue"
	}

	newActive := "green"
	if currentActive == "green" {
		newActive = "blue"
	}

	// Switch traffic to new environment
	components := []string{"model", "view", "controller", "repository"}
	for _, component := range components {
		if err := bg.switchTraffic(ctx, restAPI, component, newActive); err != nil {
			return err
		}
	}

	// Update status
	restAPI.Status.ActiveEnvironment = newActive
	now := metav1.Now()
	restAPI.Status.LastDeploymentTime = &now

	return bg.Status().Update(ctx, restAPI)
}
