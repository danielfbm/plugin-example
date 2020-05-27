/*


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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PluginSpec defines the desired state of Plugin
type PluginSpec struct {
	// Local local plugin
	Local *LocalPluginSpec `json:"local,omitempty"`
	// Network network plugin
	Network *NetworkPluginSpec `json:"network,omitempty"`
}

type LocalPluginSpec struct {
	Path string `json:"path"`
}

type NetworkPluginSpec struct {
	Address string `json:"address"`
}

// PluginStatus defines the observed state of Plugin
type PluginStatus struct {
	Conditions []Condition `json:"conditions"`
}

// Condition generic conditions
type Condition struct {
	Type    string `json:"type"`
	Ready   string `json:"status"`
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

// +kubebuilder:object:root=true

// Plugin is the Schema for the plugins API
type Plugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PluginSpec   `json:"spec,omitempty"`
	Status PluginStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PluginList contains a list of Plugin
type PluginList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Plugin `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Plugin{}, &PluginList{})
}
