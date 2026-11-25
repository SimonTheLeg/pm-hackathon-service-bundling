package controller

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	mcbuilder "sigs.k8s.io/multicluster-runtime/pkg/builder"
	mcmanager "sigs.k8s.io/multicluster-runtime/pkg/manager"
	mcreconcile "sigs.k8s.io/multicluster-runtime/pkg/reconcile"

	juiceshopv1alpha1 "github.com/platform-mesh/example-app-provider/api/v1alpha1"
)

// AppReconciler reconciles App objects across all engaged clusters.
type AppReconciler struct {
	Manager mcmanager.Manager
}

// Reconcile moves the current state of an App closer to the desired state.
func (r *AppReconciler) Reconcile(ctx context.Context, req mcreconcile.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("cluster", req.ClusterName, "namespace", req.Namespace)

	cl, err := r.Manager.GetCluster(ctx, req.ClusterName)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get cluster: %w", err)
	}
	client := cl.GetClient()

	app := &juiceshopv1alpha1.App{}
	if err := client.Get(ctx, req.NamespacedName, app); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info(fmt.Sprintf("app %s got deleted", req.Name))
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, fmt.Errorf("failed to get app: %w", err)
	}

	log.Info(fmt.Sprintf("processed app %s", app.Name))

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the multicluster Manager.
func (r *AppReconciler) SetupWithManager(mgr mcmanager.Manager) error {
	return mcbuilder.ControllerManagedBy(mgr).
		Named("app-controller").
		For(&juiceshopv1alpha1.App{}).
		Complete(r)
}
