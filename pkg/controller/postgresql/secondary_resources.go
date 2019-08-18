package postgresql

import (
	"context"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
)

// Check if PersistentVolumeClaim for the app exist, if not create one
func (r *ReconcilePostgresql) createPvc(db *v1alpha1.Postgresql) error {
	if _, err := r.fetchDBPvc(db); err != nil {
		if err := r.client.Create(context.TODO(), buildPVCForDB(db, r.scheme)); err != nil {
			return err
		}
	}
	return nil
}

// Check if Service for the app exist, if not create one
func (r *ReconcilePostgresql) createService(db *v1alpha1.Postgresql) error {
	if _, err := r.fetchDBService(db); err != nil {
		if err := r.client.Create(context.TODO(), buildDBService(db, r.scheme)); err != nil {
			return err
		}
	}
	return nil
}

// Check if Deployment for the app exist, if not create one
func (r *ReconcilePostgresql) createDeployment(db *v1alpha1.Postgresql) error {
	_, err := r.fetchDBDeployment(db)
	if err != nil {
		if err := r.client.Create(context.TODO(), buildDBDeployment(db, r.scheme)); err != nil {
			return err
		}
	}
	return nil
}
