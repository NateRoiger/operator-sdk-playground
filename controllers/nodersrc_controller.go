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

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mygroupv1alpha1 "github.com/NateRoiger/operator-sdk-playground/api/v1alpha1"
)

// NodeRsrcReconciler reconciles a NodeRsrc object
type NodeRsrcReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=mygroup.example.com,resources=nodersrcs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mygroup.example.com,resources=nodersrcs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mygroup.example.com,resources=nodersrcs/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create

// Create a default resource on start
func (r *NodeRsrcReconciler) Start(ctx context.Context) error {
	log := r.Log.WithValues("nodersrc", "start")

	namespacedName := types.NamespacedName{
		Name:      "example-nodersrc",
		Namespace: "example-nodersrcnamespace",
	}

	log.Info("Creating node resource", "NamespacedName", namespacedName)

	ns := &corev1.Namespace{}
	if err := r.Get(ctx, types.NamespacedName{Name: namespacedName.Namespace}, ns); err != nil {
		if errors.IsNotFound(err) {
			ns = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespacedName.Namespace,
				},
			}

			if err := r.Create(ctx, ns); err != nil {
				log.Error(err, "Failed to create namespace")
				return err
			}
		} else if !errors.IsAlreadyExists(err) {
			log.Error(err, "Get Namespace Failed")
			return err
		}
	}

	nodeRsrc := &mygroupv1alpha1.NodeRsrc{}
	if err := r.Get(ctx, namespacedName, nodeRsrc); err != nil {
		if errors.IsNotFound(err) {
			nodeRsrc = &mygroupv1alpha1.NodeRsrc{
				ObjectMeta: metav1.ObjectMeta{
					Name:      namespacedName.Name,
					Namespace: namespacedName.Namespace,
				},
				Spec: mygroupv1alpha1.NodeRsrcSpec{
					Foo: "foo",
				},
			}
			if err := r.Create(ctx, nodeRsrc); err != nil {
				log.Error(err, "Failed to create node resource")
				return err
			}
		} else if !errors.IsAlreadyExists(err) {
			log.Error(err, "Get node resource failed")
			return err
		}
	}

	log.Info("Node Resource Starting", "NamespacedName", namespacedName)
	return nil
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NodeRsrc object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *NodeRsrcReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("nodersrc", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeRsrcReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.Add(r); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&mygroupv1alpha1.NodeRsrc{}).
		Complete(r)
}
