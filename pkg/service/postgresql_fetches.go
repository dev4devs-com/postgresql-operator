package service

import (
	goctx "context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FetchDatabasePod search in the cluster for 1 Pod managed by the Database Controller
func FetchDatabasePod(bkp *v1alpha1.Backup, db *v1alpha1.Database, client client.Client) (*corev1.Pod, error) {
	listOps := buildDatabaseCriteria(db)
	dbPodList := &corev1.PodList{}
	err := client.List(goctx.TODO(), dbPodList, listOps)
	if err != nil {
		return nil, err
	}

	if len(dbPodList.Items) == 0 {
		return nil, err
	}

	pod := dbPodList.Items[0]
	return &pod, nil
}

//FetchDatabaseService search in the cluster for 1 Service managed by the Database Controller
func FetchDatabaseService(bkp *v1alpha1.Backup, db *v1alpha1.Database, client client.Client) (*corev1.Service, error) {
	listOps := buildDatabaseCriteria(db)
	dbServiceList := &corev1.ServiceList{}
	err := client.List(goctx.TODO(), dbServiceList, listOps)
	if err != nil {
		return nil, err
	}

	if len(dbServiceList.Items) == 0 {
		return nil, err
	}

	srv := dbServiceList.Items[0]
	return &srv, nil
}

//buildDatabaseCreteria returns client.ListOptions required to fetch the secondary resource created by
func buildDatabaseCriteria(db *v1alpha1.Database) *client.ListOptions {
	labelSelector := labels.SelectorFromSet(utils.GetLabels(db.Name))
	listOps := &client.ListOptions{Namespace: db.Namespace, LabelSelector: labelSelector}
	return listOps
}
