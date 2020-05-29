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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"github.com/hashicorp/go-hclog"

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
		if errors.IsNotFound(err) {
			netAddr, localPath := "", ""
			if plgin.Spec.Local != nil && plgin.Spec.Local.Path != "" {
				localPath = plgin.Spec.Local.Path
			} else if plgin.Spec.Network != nil && plgin.Spec.Network.Address != "" {
				netAddr = plgin.Spec.Network.Address
			}
			err = r.PluginManager.Load(plgin.Name, ext.PluginLoadOptions{
				Config: &plugin.ClientConfig{Logger: hclog.Default()}
				LocalPluginPath: localPath,
				NetworkPluginAddress: netAddr,
			})
			if err != nil {
				setCondition(plgin, pluginsv1alpha1.Condition{
					Type:    "Load",
					Ready:   strconv.FormatBool(false),
					Message: err.Error(),
					Reason:  "LoadError",
				})
				err = r.Client.Update(ctx, plgin)
			}
		}
		return ctrl.Result{}, err
	}

	// Validate each plugin and return values
	setCondition(plgin, fooCheck(pluginClient))
	setCondition(plgin, barCheck(pluginClient))
	err = r.Client.Status().Update(ctx, plgin)

	return ctrl.Result{}, err
}

// func (r *PluginReconciler) Update(ctx context.Context, plgin *pluginsv1alpha1.Plugin) {

// }

func (r *PluginReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pluginsv1alpha1.Plugin{}).
		Complete(r)
}


func setCondition(plgin *pluginsv1alpha1.Plugin, cond pluginsv1alpha1.Condition) {
	if plgin.Status.Conditions == nil {
		plgin.Status.Conditions = make([]pluginsv1alpha1.Condition,0, 2)
	}
	found := false
	for i, current := range plgin.Status.Conditions {
		if current.Type == cond.Type {
			found = true
			plgin.Status.Conditions[i] = cond
			break
		}
	}
	if !found {
		plgin.Status.Conditions = append(plgin.Status.Conditions, cond)
	}
}

func fooCheck(pluginClient plugin.Client) (cond pluginsv1alpha1.Condition) {
	cond = pluginsv1alpha1.Condition{Type: "FooPlugin"}
	var raw interface{}
	var err error
	defer func() {
		if err != nil {
			cond.Ready = "false"
			cond.Message = err.Error()
		}
	}()
	
	if raw, err = pluginClient.Dispense("foo"); err != nil {
		cond.Reason = "DispenseFailed"
		return
	}
	fooClient, ok := raw.(ext.Foo)
	if !ok {
		err = fmt.Errorf("Not a Foo interface")
		cond.Reason = "WrongInterface"
		return
	}
	if cond.Message, err = fooClient.Foos(); err != nil {
		cond.Reason = "FooResult"
		return
	}
	cond.Ready = "true"
	return
}


func barCheck(pluginClient plugin.Client) (cond pluginsv1alpha1.Condition) {
	cond = pluginsv1alpha1.Condition{Type: "BarPlugin"}
	var raw interface{}
	var err error
	defer func() {
		if err != nil {
			cond.Ready = "false"
			cond.Message = err.Error()
		}
	}()
	
	if raw, err = pluginClient.Dispense("bar"); err != nil {
		cond.Reason = "DispenseFailed"
		return
	}
	barClient, ok := raw.(ext.Bar)
	if !ok {
		err = fmt.Errorf("Not a Bar interface")
		cond.Reason = "WrongInterface"
		return
	}
	cond.Message = strings.Join(barClient.Bars(), ",")
	cond.Ready = "true"
	return
}