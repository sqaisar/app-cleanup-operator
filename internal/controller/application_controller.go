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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ApplicationReconciler reconciles Application objects
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=argoproj.io,resources=applications,verbs=get;list;watch
// +kubebuilder:rbac:groups=argoproj.io,resources=applications/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;delete

const (
	namespaceFinalizer = "namespace.app-cleanup.io"
	argoAppGVK         = "argoproj.io/v1alpha1, Kind=Application"
)

func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Starting reconciliation", "request", req)

	// Create unstructured Application object
	app := &unstructured.Unstructured{}
	app.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "argoproj.io",
		Version: "v1alpha1",
		Kind:    "Application",
	})

	// Fetch the Application
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Application not found - ignoring deletion")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to fetch Application: %v", err)
	}

	// Check deletion and finalizer presence
	if app.GetDeletionTimestamp().IsZero() || !hasFinalizer(app.GetFinalizers(), namespaceFinalizer) {
		logger.Info("Skipping - no finalizer present or not in deletion")
		return ctrl.Result{}, nil
	}

	logger.Info("Processing namespace cleanup finalizer")

	// Get target namespace
	namespace, found, err := unstructured.NestedString(app.Object, "spec", "destination", "namespace")
	if err != nil || !found || namespace == "" {
		return ctrl.Result{}, fmt.Errorf("failed to get namespace from spec.destination.namespace: %v", err)
	}

	// Delete namespace
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	if err := r.Delete(ctx, ns); err != nil && !errors.IsNotFound(err) {
		return ctrl.Result{}, fmt.Errorf("failed to delete namespace %s: %v", namespace, err)
	}

	// Remove finalizer
	finalizers := removeString(app.GetFinalizers(), namespaceFinalizer)
	app.SetFinalizers(finalizers)
	if err := r.Update(ctx, app); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %v", err)
	}

	logger.Info("Successfully processed finalizer", "namespace", namespace)
	return ctrl.Result{}, nil
}

// Helper functions
func hasFinalizer(finalizers []string, target string) bool {
	for _, f := range finalizers {
		if f == target {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) []string {
	result := []string{}
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}

func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "argoproj.io/v1alpha1",
				"kind":       "Application",
			},
		}).
		Named("application-cleanup-controller").
		Complete(r)
}
