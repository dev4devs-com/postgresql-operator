package service

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FetchPostgreSQLPod search in the cluster for 1 Pod managed by the Postgresql Controller
func FetchPostgreSQLPod(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql, client client.Client) (*corev1.Pod, error) {
	listOps := buildPostgreSQLCriteria(bkp, db)
	dbPodList := &corev1.PodList{}
	err := client.List(context.TODO(), listOps, dbPodList)
	if err != nil {
		return nil, err
	}

	if len(dbPodList.Items) == 0 {
		return nil, err
	}

	pod := dbPodList.Items[0]
	return &pod, nil
}

//FetchPostgreSQLService search in the cluster for 1 Service managed by the Postgresql Controller
func FetchPostgreSQLService(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql, client client.Client) (*corev1.Service, error) {
	listOps := buildPostgreSQLCriteria(bkp, db)
	dbServiceList := &corev1.ServiceList{}
	err := client.List(context.TODO(), listOps, dbServiceList)
	if err != nil {
		return nil, err
	}

	if len(dbServiceList.Items) == 0 {
		return nil, err
	}

	srv := dbServiceList.Items[0]
	return &srv, nil
}

//buildPostgreSQLCreteria returns client.ListOptions required to fetch the secondary resource created by
func buildPostgreSQLCriteria(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql) *client.ListOptions {
	labelSelector := labels.SelectorFromSet(utils.GetLabels(db.Name))
	listOps := &client.ListOptions{Namespace: db.Namespace, LabelSelector: labelSelector}
	return listOps
}
