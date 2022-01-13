/*
Copyright 2022.

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

type EurekaApplicationPaths struct {
	// +kubebuilder:validation:MinLength=0
	// HealthCheck path to be registered in Eureka (i.e /actuator/health)
	HealthCheck string `json:"healthcheck,omitempty"`

	// +kubebuilder:validation:MinLength=0
	// Home path to be registered in Eureka (i.e /)
	Home string `json:"home,omitempty"`

	// +kubebuilder:validation:MinLength=0
	// Status path to be registered in Eureka (i.e /actuator/info)
	Status string `json:"status,omitempty"`
}

// EurekaApplicationSpec defines the desired state of EurekaApplication
type EurekaApplicationSpec struct {
	// Enable/Disable specific instance
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// Environment that should be used to register the instance
	Environment string `json:"environment,omitempty"`

	// +kubebuilder:validation:MinLength=0
	// Name of the app to be registered in Eureka
	AppName string `json:"appName,omitempty"`

	// +kubebuilder:validation:MinLength=0
	// Name of the ingress app to be registered in Eureka
	IngressName string `json:"ingressName,omitempty"`

	// Zone of the app to be registered in Eureka
	Zone string `json:"zone,omitempty"`

	// Paths to register along with the instance
	Paths EurekaApplicationPaths `json:"paths,omitempty"`
}

// EurekaApplicationStatus defines the observed state of EurekaApplication
type EurekaApplicationStatus struct {
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EurekaApplication is the Schema for the eurekaapplications API
type EurekaApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EurekaApplicationSpec   `json:"spec,omitempty"`
	Status EurekaApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EurekaApplicationList contains a list of EurekaApplication
type EurekaApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EurekaApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EurekaApplication{}, &EurekaApplicationList{})
}
