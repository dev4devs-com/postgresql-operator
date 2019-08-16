package postgresql

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func (r *ReconcilePostgresql) fetchPostgreSQLCR( request reconcile.Request) (*v1alpha1.Postgresql, error) {
	db := &v1alpha1.Postgresql{}
	err := r.client.Get(context.TODO(), request.NamespacedName, db)
	return db, err
}

//fetchDBService returns the service resource created for this instance
func (r *ReconcilePostgresql) fetchDBService( db *v1alpha1.Postgresql) (*corev1.Service, error) {
	service := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: db.Name, Namespace: db.Namespace}, service)
	return service, err
}

//fetchDBDeployment returns the deployment resource created for this instance
func (r *ReconcilePostgresql) fetchDBDeployment( db *v1alpha1.Postgresql) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: db.Name, Namespace: db.Namespace}, deployment)
	return deployment, err
}

//fetchDBPersistentVolumeClaim returns the PersistentVolumeClaim resource created for this instance
func (r *ReconcilePostgresql) fetchDBPersistentVolumeClaim( db *v1alpha1.Postgresql) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: db.Name, Namespace: db.Namespace}, pvc)
	return pvc, err
}
