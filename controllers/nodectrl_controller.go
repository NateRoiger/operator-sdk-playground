/*
Copyright 2021.

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
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mygroupv1alpha1 "github.com/NateRoiger/operator-sdk-playground/api/v1alpha1"
)

// NodeCtrlReconciler reconciles a NodeCtrl object
type NodeCtrlReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=mygroup.example.com,namespace=operator-sdk-playground-system,resources=nodectrls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mygroup.example.com,namespace=operator-sdk-playground-system,resources=nodectrls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mygroup.example.com,namespace=operator-sdk-playground-system,resources=nodectrls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NodeCtrl object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *NodeCtrlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logr := r.Log.WithValues("nodectrl", req.NamespacedName)

	nodeCtrl := &mygroupv1alpha1.NodeCtrl{}
	if err := r.Get(ctx, req.NamespacedName, nodeCtrl); err != nil {
		return ctrl.Result{}, err
	}

	nodeRsrcNamespacedName := types.NamespacedName{Name: nodeCtrl.Spec.Name, Namespace: nodeCtrl.Spec.Namespace}
	nodeRsrc := &mygroupv1alpha1.NodeRsrc{}
	if err := r.Get(ctx, nodeRsrcNamespacedName, nodeRsrc); err != nil {
		logr.Error(err, "Failed to retrieve node resource", "Resource", nodeRsrcNamespacedName)
		return ctrl.Result{}, nil // No Retry
	}

	logr.Info(fmt.Sprintf("Got Resource %+v", *nodeRsrc), "Resource", nodeRsrcNamespacedName)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeCtrlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mygroupv1alpha1.NodeCtrl{}).
		Complete(r)
}
