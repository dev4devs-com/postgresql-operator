package backup

import (
	"context"
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/resource"
	"github.com/dev4devs-com/postgresql-operator/pkg/service"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
)

// Set in the ReconcileBackup the Pod database created by Database
// NOTE: This data is required in order to create the secrets which will access the database container to do the backup
func (r *ReconcileBackup) getDatabasePod(bkp *v1alpha1.Backup, db *v1alpha1.Database) error {
	dbPod, err := service.FetchDatabasePod(bkp, db, r.client)
	if err != nil || dbPod == nil {
		r.dbPod = nil
		err := fmt.Errorf("Unable to find the Database Pod")
		return err
	}
	r.dbPod = dbPod
	return nil
}

// Set in the ReconcileBackup the service database created by Database
// NOTE: This data is required in order to create the secrets which will access the database container to do the backup
func (r *ReconcileBackup) getDatabaseService(bkp *v1alpha1.Backup, db *v1alpha1.Database) error {
	dbService, err := service.FetchDatabaseService(bkp, db, r.client)
	if err != nil || dbService == nil {
		r.dbService = nil
		err := fmt.Errorf("Unable to find the Database Service")
		return err
	}
	r.dbService = dbService
	return nil
}

// Check if the cronJob is created, if not create one
func (r *ReconcileBackup) createCronJob(bkp *v1alpha1.Backup) error {
	if _, err := service.FetchCronJob(bkp.Name, bkp.Namespace, r.client); err != nil {
		cron, err := resource.NewBackupCronJob(bkp, r.scheme)
		if err != nil {
			return err
		}
		if err := r.client.Create(context.TODO(), cron); err != nil {
			return err
		}
	}
	return nil
}

// Check if the encryptionKey is created, if not create one
// NOTE: The user can config in the CR to use a pre-existing one by informing the name
func (r *ReconcileBackup) createEncryptionKey(bkp *v1alpha1.Backup) error {
	if utils.IsEncryptionKeyOptionConfig(bkp) {
		if _, err := service.FetchSecret(utils.GetEncSecretNamespace(bkp), utils.GetEncSecretName(bkp), r.client); err != nil {
			// The user can just inform the name of the Secret which is already applied in the cluster
			if utils.IsEncKeySetupByName(bkp) {
				return err
			}
			secretData, secretStringData := createEncDataMaps(bkp)
			encSecret, err := resource.NewBackupSecret(bkp, utils.EncSecretPrefix, secretData, secretStringData, r.scheme)
			if err != nil {
				return err
			}
			if err := r.client.Create(context.TODO(), encSecret); err != nil {
				return err
			}
		}
	}
	return nil
}

// createAwsSecret checks if the secret with the aws data is created, if not create one
// NOTE: The user can config in the CR to use a pre-existing one by informing the name
func (r *ReconcileBackup) createAwsSecret(bkp *v1alpha1.Backup) error {
	if _, err := service.FetchSecret(utils.GetAwsSecretNamespace(bkp), utils.GetAWSSecretName(bkp), r.client); err != nil {
		// The user can just inform the name of the Secret which is already applied in the cluster
		if !utils.IsAwsKeySetupByName(bkp) {
			secretData := createAwsDataByteMap(bkp)
			awsSecret, err := resource.NewBackupSecret(bkp, utils.AwsSecretPrefix, secretData, nil, r.scheme)
			if err != nil {
				return err
			}
			if err := r.client.Create(context.TODO(), awsSecret); err != nil {
				return err
			}
		}
	}
	return nil
}

// createDatabaseSecret checks if the secret with the database is created, if not create one
func (r *ReconcileBackup) createDatabaseSecret(bkp *v1alpha1.Backup, db *v1alpha1.Database) error {
	dbSecretName := utils.DbSecretPrefix + bkp.Name
	if _, err := service.FetchSecret(bkp.Namespace, dbSecretName, r.client); err != nil {
		secretData, err := r.buildDBSecretData(bkp, db)
		if err != nil {
			return err
		}
		dbSecret, err := resource.NewBackupSecret(bkp, utils.DbSecretPrefix, secretData, nil, r.scheme)
		if err != nil {
			return err
		}
		if err := r.client.Create(context.TODO(), dbSecret); err != nil {
			return err
		}
	}
	return nil
}
