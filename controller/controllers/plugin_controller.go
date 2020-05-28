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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pluginsv1alpha1 "github.com/danielfbm/plugin-example/controller/api/v1alpha1"
	ext "github.com/danielfbm/plugin-example/controller/extension"
)

// PluginReconciler reconciles a Plugin object
type PluginReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	PluginManager ext.Manager
}

// +kubebuilder:rbac:groups=plugins.danielfbm.github.com,resources=plugins,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=plugins.danielfbm.github.com,resources=plugins/status,verbs=get;update;patch

func (r *PluginReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("plugin", req.NamespacedName)

	// Get plugin
	plgin := &pluginsv1alpha1.Plugin{}
	if err := r.Client.Get(ctx, req.NamespacedName, plgin); err != nil {
		err = client.IgnoreNotFound(err)
		return ctrl.Result{}, err
	}

	// Verifying if it is a local plugin and tried to delete
	// we should recover the local plugin, update and return
	if plgin.DeletionTimestamp != nil && plgin.Spec.Local != nil && len(plgin.Finalizers) > 0 {
		plgin.DeletionTimestamp = nil
		err := r.Client.Update(ctx, plgin)
		return ctrl.Result{}, err
	}

	// Check if plugin with the same name is loaded
	// and if is the same plugin
	pluginClient, err := r.PluginManager.Get(plgin.Name)
	if err != nil {
		r.PluginManager.
	}

	return ctrl.Result{}, nil
}

func (r *PluginReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pluginsv1alpha1.Plugin{}).
		Complete(r)
}
