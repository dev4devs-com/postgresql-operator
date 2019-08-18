package service

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func FetchPostgreSQL(name, namespace string, client client.Client) (*v1alpha1.Postgresql, error) {
	db := &v1alpha1.Postgresql{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, db)
	return db, err
}

func FetchBackupCR(name, namespace string, client client.Client) (*v1alpha1.Backup, error) {
	bkp := &v1alpha1.Backup{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, bkp)
	return bkp, err
}
