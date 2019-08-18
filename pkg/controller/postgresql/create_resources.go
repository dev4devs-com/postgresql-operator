package postgresql

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/resource"
	"github.com/dev4devs-com/postgresql-operator/pkg/service"
)

// Check if PersistentVolumeClaim for the app exist, if not create one
func (r *ReconcilePostgresql) createPvc(db *v1alpha1.Postgresql) error {
	if _, err := service.FetchPersistentVolumeClaim(db.Name, db.Namespace, r.client); err != nil {
		if err := r.client.Create(context.TODO(), resource.NewPostgresqlPvc(db, r.scheme)); err != nil {
			return err
		}
	}
	return nil
}

// Check if Service for the app exist, if not create one
func (r *ReconcilePostgresql) createService(db *v1alpha1.Postgresql) error {
	if _, err := service.FetchService(db.Name, db.Namespace, r.client); err != nil {
		if err := r.client.Create(context.TODO(), resource.NewPostgresqlService(db, r.scheme)); err != nil {
			return err
		}
	}
	return nil
}

// Check if Deployment for the app exist, if not create one
func (r *ReconcilePostgresql) createDeployment(db *v1alpha1.Postgresql) error {
	_, err := service.FetchDeployment(db.Name, db.Namespace, r.client)
	if err != nil {
		if err := r.client.Create(context.TODO(), resource.NewPostgresqlDeployment(db, r.scheme)); err != nil {
			return err
		}
	}
	return nil
}
