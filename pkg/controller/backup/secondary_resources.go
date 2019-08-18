package backup

import (
	"context"
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
)

// Set in the ReconcileBackup the Pod database created by PostgreSQL
// NOTE: This data is required in order to create the secrets which will access the database container to do the backup
func (r *ReconcileBackup) setDatabasePod(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql) error {
	dbPod, err := r.fetchPostgreSQLPod(bkp, db)
	if err != nil || dbPod == nil {
		r.dbPod = nil
		err := fmt.Errorf("Unable to find the PostgreSQL Pod")
		return err
	}
	r.dbPod = dbPod
	return nil
}

// Set in the ReconcileBackup the service database created by PostgreSQL
// NOTE: This data is required in order to create the secrets which will access the database container to do the backup
func (r *ReconcileBackup) setDatabaseService(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql) error {
	dbService, err := r.fetchPostgreSQLService(bkp, db)
	if err != nil || dbService == nil {
		r.dbService = nil
		err := fmt.Errorf("Unable to find the PostgreSQL Service")
		return err
	}
	r.dbService = dbService
	return nil
}

// Check if the cronJob is created, if not create one
func (r *ReconcileBackup) createCronJob(bkp *v1alpha1.Backup) error {
	if _, err := r.fetchCronJob(bkp); err != nil {
		if err := r.client.Create(context.TODO(), buildCronJob(bkp, r.scheme)); err != nil {
			return err
		}
	}
	return nil
}

// Check if the encryptionKey is created, if not create one
// NOTE: The user can config in the CR to use a pre-existing one by informing the name
func (r *ReconcileBackup) createEncryptionKey(bkp *v1alpha1.Backup) error {
	if isEncryptionKeyOptionConfig(bkp) {
		if _, err := r.fetchSecret(getEncSecretNamespace(bkp), getEncSecretName(bkp)); err != nil {
			// The user can just inform the name of the Secret which is already applied in the cluster
			if isEncKeySetupByName(bkp) {
				return err
			} else {
				secretData, secretStringData := createEncDataMaps(bkp)
				encSecret := buildSecret(bkp, encSecretPrefix, secretData, secretStringData, r.scheme)
				if err := r.client.Create(context.TODO(), encSecret); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// createAwsSecret checks if the secret with the aws data is created, if not create one
// NOTE: The user can config in the CR to use a pre-existing one by informing the name
func (r *ReconcileBackup) createAwsSecret(bkp *v1alpha1.Backup) error {
	if _, err := r.fetchSecret(getAwsSecretNamespace(bkp), getAWSSecretName(bkp)); err != nil {
		// The user can just inform the name of the Secret which is already applied in the cluster
		if !isAwsKeySetupByName(bkp) {
			secretData := createAwsDataByteMap(bkp)
			awsSecret := buildSecret(bkp, awsSecretPrefix, secretData, nil, r.scheme)
			if err := r.client.Create(context.TODO(), awsSecret); err != nil {
				return err
			}
		}
	}
	return nil
}

// createDatabaseSecret checks if the secret with the database is created, if not create one
func (r *ReconcileBackup) createDatabaseSecret(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql) error {
	dbSecretName := dbSecretPrefix + bkp.Name
	if _, err := r.fetchSecret(bkp.Namespace, dbSecretName); err != nil {
		secretData, err := r.buildDBSecretData(bkp, db)
		if err != nil {
			return err
		}
		dbSecret := buildSecret(bkp, dbSecretPrefix, secretData, nil, r.scheme)
		if err := r.client.Create(context.TODO(), dbSecret); err != nil {
			return err
		}
	}
	return nil
}
