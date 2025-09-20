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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RestAPISpec defines the desired state of RestAPI.
type RestAPISpec struct {
	// Application configuration
	Image    string            `json:"image"`
	Replicas *int32            `json:"replicas,omitempty"`
	EnvVars  map[string]string `json:"envVars,omitempty"`

	// MVC+R Pattern Components
	Model      ComponentSpec `json:"model"`
	View       ComponentSpec `json:"view"`
	Controller ComponentSpec `json:"controller"`
	Repository ComponentSpec `json:"repository"`

	// Auto-scaling configuration
	AutoScaling *AutoScalingSpec `json:"autoScaling,omitempty"`

	// Health monitoring
	HealthCheck *HealthCheckSpec `json:"healthCheck,omitempty"`

	// Blue-green deployment
	BlueGreen *BlueGreenSpec `json:"blueGreen,omitempty"`
}

type ComponentSpec struct {
	Enabled bool              `json:"enabled"`
	Image   string            `json:"image,omitempty"`
	Port    int32             `json:"port,omitempty"`
	EnvVars map[string]string `json:"envVars,omitempty"`
}

type AutoScalingSpec struct {
	Enabled                 bool   `json:"enabled"`
	MinReplicas             *int32 `json:"minReplicas,omitempty"`
	MaxReplicas             int32  `json:"maxReplicas"`
	TargetCPUUtilization    *int32 `json:"targetCPUUtilization,omitempty"`
	TargetMemoryUtilization *int32 `json:"targetMemoryUtilization,omitempty"`
}

type HealthCheckSpec struct {
	Enabled             bool   `json:"enabled"`
	Path                string `json:"path,omitempty"`
	InitialDelaySeconds *int32 `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       *int32 `json:"periodSeconds,omitempty"`
	TimeoutSeconds      *int32 `json:"timeoutSeconds,omitempty"`
	FailureThreshold    *int32 `json:"failureThreshold,omitempty"`
}

type BlueGreenSpec struct {
	Enabled          bool   `json:"enabled"`
	Strategy         string `json:"strategy,omitempty"`
	PromotionTimeout *int32 `json:"promotionTimeout,omitempty"`
}

// RestAPIStatus defines the observed state of RestAPI.
type RestAPIStatus struct {
	Phase              string             `json:"phase,omitempty"`
	Replicas           int32              `json:"replicas,omitempty"`
	ReadyReplicas      int32              `json:"readyReplicas,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
	ActiveEnvironment  string             `json:"activeEnvironment,omitempty"`
	LastDeploymentTime *metav1.Time       `json:"lastDeploymentTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// RestAPI is the Schema for the restapis API.
type RestAPI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RestAPISpec   `json:"spec,omitempty"`
	Status RestAPIStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RestAPIList contains a list of RestAPI.
type RestAPIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RestAPI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RestAPI{}, &RestAPIList{})
}
