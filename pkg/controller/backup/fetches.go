package backup

import (
	"context"
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func (r *ReconcileBackup) fetchBkpInstance(reqLogger logr.Logger, request reconcile.Request) (*v1alpha1.Backup, error) {
	reqLogger.Info("Checking if the Backup already exists")
	bkp := &v1alpha1.Backup{}
	//Fetch the PostgreSQL Backup db
	err := r.client.Get(context.TODO(), request.NamespacedName, bkp)
	return bkp, err
}

// fetchCronJob return the cronJob created pod created by Backup
func (r *ReconcileBackup) fetchCronJob(reqLogger logr.Logger, bkp *v1alpha1.Backup) (*v1beta1.CronJob, error) {
	reqLogger.Info("Checking if the cronJob already exists")
	cronJob := &v1beta1.CronJob{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: bkp.Name, Namespace: bkp.Namespace}, cronJob)
	return cronJob, err
}

// fetchCronJob return the cronJob created pod created by Backup
func (r *ReconcileBackup) fetchSecret(reqLogger logr.Logger, secretNamespace, secretName string) (*corev1.Secret, error) {
	reqLogger.Info("Checking if the secret already exists", "secret.name", secretName, "secret.Namespace", secretNamespace)
	secret := &corev1.Secret{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: secretNamespace}, secret)
	return secret, err
}

func (r *ReconcileBackup) fetchConfigMap(bkp *v1alpha1.Backup, cfgName string) (*corev1.ConfigMap, error) {
	log.Info("Looking for ConfigMap to get database data", "configMapName", cfgName)
	cfg := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: cfgName, Namespace: bkp.Namespace}, cfg)
	return cfg, err
}

func (r *ReconcileBackup) fetchBDPod(bkp *v1alpha1.Backup, reqLogger logr.Logger, request reconcile.Request) (*corev1.Pod, error) {
	listOps, err := r.getListOpsToSearchDBObject(bkp, reqLogger)
	if err != nil {
		return nil, err
	}

	// Search DB pods
	reqLogger.Info("Searching for DB pod ...")
	dbPodList := &corev1.PodList{}
	err = r.client.List(context.TODO(), listOps, dbPodList)
	if err != nil {
		return nil, err
	}

	if len(dbPodList.Items) == 0 {
		err = fmt.Errorf("Unable to find database pod. Maybe, it was not create yet")
		return nil, err
	}

	// Getting the pod ( it has just one )
	pod := dbPodList.Items[0]
	reqLogger.Info("DB Pod was found", "pod.Name", pod.Name)
	return &pod, nil
}

func (r *ReconcileBackup) fetchServiceDB(bkp *v1alpha1.Backup, reqLogger logr.Logger, request reconcile.Request) (*corev1.Service, error) {
	listOps, err := r.getListOpsToSearchDBObject(bkp, reqLogger)
	if err != nil {
		return nil, err
	}

	// Search DB pods
	reqLogger.Info("Searching for Service pod ...")
	dbServiceList := &corev1.ServiceList{}
	err = r.client.List(context.TODO(), listOps, dbServiceList)
	if err != nil {
		return nil, err
	}

	if len(dbServiceList.Items) == 0 {
		err = fmt.Errorf("Unable to find database service. Maybe, it was not create yet")
		return nil, err
	}

	// Getting the pod ( it has just one )
	srv := dbServiceList.Items[0]
	reqLogger.Info("DB Service was found", "srv.Name", srv.Name)
	return &srv, nil
}

func (r *ReconcileBackup) getListOpsToSearchDBObject(bkp *v1alpha1.Backup, reqLogger logr.Logger) (*client.ListOptions, error) {
	reqLogger.Info("Checking if the Database Service exists ...")
	reqLogger.Info("Checking operator namespace ...")

	// Fetch PostgreSQL Database
	reqLogger.Info("Checking PostgreSQL exists ...")
	db := &v1alpha1.Postgresql{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "postgresql", Namespace: bkp.Namespace}, db)
	if err != nil {
		return nil, err
	}
	// Create criteria
	reqLogger.Info("Creating criteria to looking for Service ...")
	ls := map[string]string{"app": "postgresql", "postgresql_cr": db.Name, "name": "postgresql"}
	labelSelector := labels.SelectorFromSet(ls)
	listOps := &client.ListOptions{Namespace: bkp.Namespace, LabelSelector: labelSelector}
	return listOps, nil
}
