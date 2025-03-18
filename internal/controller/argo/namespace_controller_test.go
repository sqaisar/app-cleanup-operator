package controller

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Namespace Controller", func() {
	Context("When reconciling an Argo Application resource", func() {
		const resourceName = "test-application"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // Adjust if necessary
		}

		app := &unstructured.Unstructured{}
		app.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "argoproj.io",
			Version: "v1alpha1",
			Kind:    "Application",
		})

		BeforeEach(func() {
			By("creating the Argo Application resource")

			err := k8sClient.Get(ctx, typeNamespacedName, app)
			if err != nil && errors.IsNotFound(err) {
				app = &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "argoproj.io/v1alpha1",
						"kind":       "Application",
						"metadata": map[string]interface{}{
							"name":       resourceName,
							"namespace":  "default",
							"finalizers": []interface{}{namespaceFinalizer},
						},
						"spec": map[string]interface{}{
							"destination": map[string]interface{}{
								"namespace": "target-namespace",
							},
							"project": "target-project",
						},
					},
				}
				Expect(k8sClient.Create(ctx, app)).To(Succeed())
			}
		})

		AfterEach(func() {
			By("Cleaning up the Argo Application resource")
			err := k8sClient.Get(ctx, typeNamespacedName, app)
			if err == nil {
				Expect(k8sClient.Delete(ctx, app)).To(Succeed())
			}
		})

		It("should successfully reconcile the Application custom resource", func() {
			By("Reconciling the Argo Application")

			controllerReconciler := &NamespaceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			// First, ensure the resource is there
			Expect(k8sClient.Get(ctx, typeNamespacedName, app)).To(Succeed())

			// Delete the Application to simulate deletion
			Expect(k8sClient.Delete(ctx, app)).To(Succeed())

			// Fetch the updated resource to ensure deletionTimestamp is set by Kubernetes
			deletedApp := &unstructured.Unstructured{}
			deletedApp.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   "argoproj.io",
				Version: "v1alpha1",
				Kind:    "Application",
			})
			Eventually(func() bool {
				err := k8sClient.Get(ctx, typeNamespacedName, deletedApp)
				if err != nil {
					return false
				}
				return !deletedApp.GetDeletionTimestamp().IsZero()
			}, "5s", "100ms").Should(BeTrue(), "Deletion timestamp should be set after deletion")

			// Explicitly call the reconcile function
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("verifying that the namespace has been cleaned up and finalizer removed")

			updatedApp := &unstructured.Unstructured{}
			updatedApp.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   "argoproj.io",
				Version: "v1alpha1",
				Kind:    "Application",
			})

			// Validate finalizer has been successfully removed.
			Eventually(func() []string {
				if err := k8sClient.Get(ctx, typeNamespacedName, updatedApp); err != nil {
					return nil // could return empty if deleted already
				}
				return updatedApp.GetFinalizers()
			}, "10s", "250ms").ShouldNot(ContainElement(namespaceFinalizer), "Finalizer should be removed after successful reconciliation")
		})

	})
})
