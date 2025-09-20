/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	restapiv1 "github.com/aws/restapi-operator/api/v1"
)

// RestAPIReconciler reconciles a RestAPI object
type RestAPIReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.aws.com,resources=restapis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.aws.com,resources=restapis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.aws.com,resources=restapis/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete

func (r *RestAPIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Fetch RestAPI instance
	restAPI := &restapiv1.RestAPI{}
	if err := r.Get(ctx, req.NamespacedName, restAPI); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Handle MVC+R components
	if err := r.reconcileMVCRComponents(ctx, restAPI); err != nil {
		log.Error(err, "Failed to reconcile MVC+R components")
		return ctrl.Result{RequeueAfter: time.Minute}, err
	}

	// Handle auto-scaling
	if restAPI.Spec.AutoScaling != nil && restAPI.Spec.AutoScaling.Enabled {
		if err := r.reconcileAutoScaling(ctx, restAPI); err != nil {
			log.Error(err, "Failed to reconcile auto-scaling")
			return ctrl.Result{RequeueAfter: time.Minute}, err
		}
	}

	// Update status
	if err := r.updateStatus(ctx, restAPI); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{RequeueAfter: time.Minute}, err
	}

	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

func (r *RestAPIReconciler) reconcileMVCRComponents(ctx context.Context, restAPI *restapiv1.RestAPI) error {
	components := map[string]restapiv1.ComponentSpec{
		"model":      restAPI.Spec.Model,
		"view":       restAPI.Spec.View,
		"controller": restAPI.Spec.Controller,
		"repository": restAPI.Spec.Repository,
	}

	for name, spec := range components {
		if !spec.Enabled {
			continue
		}

		deployment := r.createDeployment(restAPI, name, spec)
		if err := controllerutil.SetControllerReference(restAPI, deployment, r.Scheme); err != nil {
			return err
		}

		found := &appsv1.Deployment{}
		err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			if err := r.Create(ctx, deployment); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			if err := r.Update(ctx, deployment); err != nil {
				return err
			}
		}

		// Create service
		service := r.createService(restAPI, name, spec)
		if err := controllerutil.SetControllerReference(restAPI, service, r.Scheme); err != nil {
			return err
		}

		foundSvc := &corev1.Service{}
		err = r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundSvc)
		if err != nil && errors.IsNotFound(err) {
			if err := r.Create(ctx, service); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (r *RestAPIReconciler) createDeployment(restAPI *restapiv1.RestAPI, component string, spec restapiv1.ComponentSpec) *appsv1.Deployment {
	replicas := int32(1)
	if restAPI.Spec.Replicas != nil {
		replicas = *restAPI.Spec.Replicas
	}

	image := spec.Image
	if image == "" {
		image = restAPI.Spec.Image
	}

	envVars := []corev1.EnvVar{}
	for k, v := range restAPI.Spec.EnvVars {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}
	for k, v := range spec.EnvVars {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}

	container := corev1.Container{
		Name:  component,
		Image: image,
		Ports: []corev1.ContainerPort{{
			ContainerPort: spec.Port,
			Protocol:      corev1.ProtocolTCP,
		}},
		Env: envVars,
	}

	// Add health checks if configured
	if restAPI.Spec.HealthCheck != nil && restAPI.Spec.HealthCheck.Enabled {
		probe := &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: restAPI.Spec.HealthCheck.Path,
					Port: intstr.FromInt32(spec.Port),
				},
			},
			InitialDelaySeconds: *restAPI.Spec.HealthCheck.InitialDelaySeconds,
			PeriodSeconds:       *restAPI.Spec.HealthCheck.PeriodSeconds,
			TimeoutSeconds:      *restAPI.Spec.HealthCheck.TimeoutSeconds,
			FailureThreshold:    *restAPI.Spec.HealthCheck.FailureThreshold,
		}
		container.LivenessProbe = probe
		container.ReadinessProbe = probe
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", restAPI.Name, component),
			Namespace: restAPI.Namespace,
			Labels: map[string]string{
				"app":       restAPI.Name,
				"component": component,
				"version":   "v1",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":       restAPI.Name,
					"component": component,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       restAPI.Name,
						"component": component,
						"version":   "v1",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{container},
				},
			},
		},
	}
}

func (r *RestAPIReconciler) createService(restAPI *restapiv1.RestAPI, component string, spec restapiv1.ComponentSpec) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-svc", restAPI.Name, component),
			Namespace: restAPI.Namespace,
			Labels: map[string]string{
				"app":       restAPI.Name,
				"component": component,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":       restAPI.Name,
				"component": component,
			},
			Ports: []corev1.ServicePort{{
				Port:       spec.Port,
				TargetPort: intstr.FromInt32(spec.Port),
				Protocol:   corev1.ProtocolTCP,
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

func (r *RestAPIReconciler) reconcileAutoScaling(ctx context.Context, restAPI *restapiv1.RestAPI) error {
	components := []string{"model", "view", "controller", "repository"}

	for _, component := range components {
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

		if !spec.Enabled {
			continue
		}

		hpa := r.createHPA(restAPI, component)
		if err := controllerutil.SetControllerReference(restAPI, hpa, r.Scheme); err != nil {
			return err
		}

		found := &autoscalingv2.HorizontalPodAutoscaler{}
		err := r.Get(ctx, types.NamespacedName{Name: hpa.Name, Namespace: hpa.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			if err := r.Create(ctx, hpa); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (r *RestAPIReconciler) createHPA(restAPI *restapiv1.RestAPI, component string) *autoscalingv2.HorizontalPodAutoscaler {
	minReplicas := int32(1)
	if restAPI.Spec.AutoScaling.MinReplicas != nil {
		minReplicas = *restAPI.Spec.AutoScaling.MinReplicas
	}

	targetCPU := int32(80)
	if restAPI.Spec.AutoScaling.TargetCPUUtilization != nil {
		targetCPU = *restAPI.Spec.AutoScaling.TargetCPUUtilization
	}

	metrics := []autoscalingv2.MetricSpec{{
		Type: autoscalingv2.ResourceMetricSourceType,
		Resource: &autoscalingv2.ResourceMetricSource{
			Name: corev1.ResourceCPU,
			Target: autoscalingv2.MetricTarget{
				Type:               autoscalingv2.UtilizationMetricType,
				AverageUtilization: &targetCPU,
			},
		},
	}}

	if restAPI.Spec.AutoScaling.TargetMemoryUtilization != nil {
		targetMemory := *restAPI.Spec.AutoScaling.TargetMemoryUtilization
		metrics = append(metrics, autoscalingv2.MetricSpec{
			Type: autoscalingv2.ResourceMetricSourceType,
			Resource: &autoscalingv2.ResourceMetricSource{
				Name: corev1.ResourceMemory,
				Target: autoscalingv2.MetricTarget{
					Type:               autoscalingv2.UtilizationMetricType,
					AverageUtilization: &targetMemory,
				},
			},
		})
	}

	return &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-hpa", restAPI.Name, component),
			Namespace: restAPI.Namespace,
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       fmt.Sprintf("%s-%s", restAPI.Name, component),
			},
			MinReplicas: &minReplicas,
			MaxReplicas: restAPI.Spec.AutoScaling.MaxReplicas,
			Metrics:     metrics,
		},
	}
}

func (r *RestAPIReconciler) updateStatus(ctx context.Context, restAPI *restapiv1.RestAPI) error {
	restAPI.Status.Phase = "Running"
	restAPI.Status.ActiveEnvironment = "blue"
	if restAPI.Status.LastDeploymentTime == nil {
		now := metav1.Now()
		restAPI.Status.LastDeploymentTime = &now
	}

	return r.Status().Update(ctx, restAPI)
}

// SetupWithManager sets up the controller with the Manager.
func (r *RestAPIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&restapiv1.RestAPI{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&autoscalingv2.HorizontalPodAutoscaler{}).
		Named("restapi").
		Complete(r)
}
