package database

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/resource"
	"github.com/dev4devs-com/postgresql-operator/pkg/service"
)

// Check if PersistentVolumeClaim for the app exist, if not create one
func (r *ReconcileDatabase) createPvc(db *v1alpha1.Database) error {
	if _, err := service.FetchPersistentVolumeClaim(db.Name, db.Namespace, r.client); err != nil {
		if err := r.client.Create(context.TODO(), resource.NewDatabasePvc(db, r.scheme)); err != nil {
			return err
		}
	}
	return nil
}

// Check if Service for the app exist, if not create one
func (r *ReconcileDatabase) createService(db *v1alpha1.Database) error {
	if _, err := service.FetchService(db.Name, db.Namespace, r.client); err != nil {
		if err := r.client.Create(context.TODO(), resource.NewDatabaseService(db, r.scheme)); err != nil {
			return err
		}
	}
	return nil
}

// Check if Deployment for the app exist, if not create one
func (r *ReconcileDatabase) createDeployment(db *v1alpha1.Database) error {
	_, err := service.FetchDeployment(db.Name, db.Namespace, r.client)
	if err != nil {
		if err := r.client.Create(context.TODO(), resource.NewDatabaseDeployment(db, r.scheme)); err != nil {
			return err
		}
	}
	return nil
}
