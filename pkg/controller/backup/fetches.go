package backup

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func (r *ReconcileBackup) fetchBackupCR(request reconcile.Request) (*v1alpha1.Backup, error) {
	bkp := &v1alpha1.Backup{}
	err := r.client.Get(context.TODO(), request.NamespacedName, bkp)
	return bkp, err
}

//fetchCronJob search in the cluster for the CronJob managed by the Backup Controller
func (r *ReconcileBackup) fetchCronJob(bkp *v1alpha1.Backup) (*v1beta1.CronJob, error) {
	cronJob := &v1beta1.CronJob{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: bkp.Name, Namespace: bkp.Namespace}, cronJob)
	return cronJob, err
}

//fetchSecret search in the cluster for the Secret managed by the Backup Controller
func (r *ReconcileBackup) fetchSecret(namespace, name string) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, secret)
	return secret, err
}

//fetchConfigMap search in the cluster for the ConfigMap managed by the Backup Controller
func (r *ReconcileBackup) fetchConfigMap(name, namespace string) (*corev1.ConfigMap, error) {
	cfg := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, cfg)
	return cfg, err
}

//fetchPostgreSQLPod search in the cluster for 1 Pod managed by the Postgresql Controller
func (r *ReconcileBackup) fetchPostgreSQLPod(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql) (*corev1.Pod, error) {
	listOps := buildPostgreSQLCriteria(bkp, db)
	dbPodList := &corev1.PodList{}
	err := r.client.List(context.TODO(), listOps, dbPodList)
	if err != nil {
		return nil, err
	}

	if len(dbPodList.Items) == 0 {
		return nil, err
	}

	pod := dbPodList.Items[0]
	return &pod, nil
}

//fetchPostgreSQLService search in the cluster for 1 Service managed by the Postgresql Controller
func (r *ReconcileBackup) fetchPostgreSQLService(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql) (*corev1.Service, error) {
	listOps := buildPostgreSQLCriteria(bkp, db)
	dbServiceList := &corev1.ServiceList{}
	err := r.client.List(context.TODO(), listOps, dbServiceList)
	if err != nil {
		return nil, err
	}

	if len(dbServiceList.Items) == 0 {
		return nil, err
	}

	srv := dbServiceList.Items[0]
	return &srv, nil
}

//fetchPostgreSQLInstance search in the cluster for 1 pod managed by the Postgresql Controller
func (r *ReconcileBackup) fetchPostgreSQLInstance(bkp *v1alpha1.Backup) (*v1alpha1.Postgresql, error) {
	db := &v1alpha1.Postgresql{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "postgresql", Namespace: bkp.Namespace}, db)
	return db, err
}

//buildPostgreSQLCreteria returns client.ListOptions required to fetch the secondary resources created by
func buildPostgreSQLCriteria(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql) *client.ListOptions {
	ls := map[string]string{"app": "postgresql", "postgresql_cr": db.Name}
	labelSelector := labels.SelectorFromSet(ls)
	listOps := &client.ListOptions{Namespace: bkp.Namespace, LabelSelector: labelSelector}
	return listOps
}
