package backup

import (
	"context"
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const statusOk = "OK"

//updateAppStatus returns error when status regards  all required resources could not be updated with OK
func (r *ReconcileBackup) updateBackupStatus(request reconcile.Request) error {
	bkp, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	// Check if all required resources were created and found
	if err := r.isAllCreated(bkp); err != nil {
		return err
	}

	// Check if BackupStatus was changed, if yes update it
	if err := r.insertUpdateBackupStatus(bkp); err != nil {
		return err
	}
	return nil
}

// Check if BackupStatus was changed, if yes update it
func (r *ReconcileBackup) insertUpdateBackupStatus(bkp *v1alpha1.Backup) error {
	if !reflect.DeepEqual(statusOk, bkp.Status.BackupStatus) {
		bkp.Status.BackupStatus = statusOk
		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// updateCronJobStatus returns error when was not possible update the CronJob status successfully
func (r *ReconcileBackup) updateCronJobStatus(request reconcile.Request) error {
	bkp, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	// Check if Cronjob Name or Status was changed, if yes update it
	cronJob, err := r.fetchCronJob(bkp)
	if err != nil {
		return err
	}

	// Check if CronJob changed, if yes update its status
	if err := r.insertUpdateCronJobStatus(cronJob, bkp); err != nil {
		return err
	}
	return nil
}

// insertUpdateCronJobStatus if CronJob name and status was changed the its status wil be updated
func (r *ReconcileBackup) insertUpdateCronJobStatus(cronJob *v1beta1.CronJob, bkp *v1alpha1.Backup) error {
	if cronJob.Name != bkp.Status.CronJobName || !reflect.DeepEqual(cronJob.Status, bkp.Status.CronJobStatus) {

		bkp.Status.CronJobStatus = cronJob.Status
		bkp.Status.CronJobName = cronJob.Name

		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// updateAWSSecretStatus returns error when was not possible update the AWS status fields in the CR successfully
func (r *ReconcileBackup) updateAWSSecretStatus(request reconcile.Request) error {
	bkp, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	aws, err := r.fetchSecret(getAwsSecretNamespace(bkp), getAWSSecretName(bkp))
	if err != nil {
		return err
	}

	// Check if the Secret with the AWS data was changed, if yes update its status
	if err := r.insertUpdateAwsSecretStatus(aws, bkp); err != nil {
		return err
	}
	return nil
}

// insertUpdateAwsSecretStatus will check and update the AWS Secret status if the Secret with the AWS data was changed
func (r *ReconcileBackup) insertUpdateAwsSecretStatus(aws *corev1.Secret, bkp *v1alpha1.Backup) error {
	data := covertDataSecretToString(aws)
	if isAwsStatusEqual(aws, bkp, data) {

		bkp.Status.AWSSecretName = aws.Name
		bkp.Status.AWSSecretData = data
		bkp.Status.AwsCredentialsSecretNamespace = aws.Namespace

		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// isAwsStatusEqual return true when something related to the aws status fields changed
func isAwsStatusEqual(aws *corev1.Secret, bkp *v1alpha1.Backup, data map[string]string) bool {
	return aws.Name != bkp.Status.AWSSecretName || !reflect.DeepEqual(data, bkp.Status.AWSSecretData) || aws.Namespace != bkp.Status.AwsCredentialsSecretNamespace
}

// updateAWSSecretStatus returns error when was not possible update the EncryptionKey status fields in the CR successfully
func (r *ReconcileBackup) updateEncSecretStatus(request reconcile.Request) error {
	bkp, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	isEncryptionKeyOptionConfig := isEncryptionKeyOptionConfig(bkp)
	if isEncryptionKeyOptionConfig {
		secret, err := r.fetchSecret(getEncSecretNamespace(bkp), getEncSecretName(bkp))
		if err != nil {
			return err
		}

		// Check if the Secret with the AWS data was changed, if yes update its status
		if err := r.insertUpdateEncKeyStatus(secret, bkp); err != nil {
			return err
		}
	}

	// Check if the config(boolean status) was changed, if yes update it
	if isEncryptionKeyOptionConfig != bkp.Status.HasEncryptionKey {

		bkp.Status.HasEncryptionKey = isEncryptionKeyOptionConfig
		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// insertUpdateEncKeyStatus will check and update the EncryptionKey Secret status if the Secret with the AWS data was changed
func (r *ReconcileBackup) insertUpdateEncKeyStatus(secret *corev1.Secret, bkp *v1alpha1.Backup) error {
	data := covertDataSecretToString(secret)
	if isEncryptKeyStatusEquals(secret, bkp, data) {

		bkp.Status.EncryptionKeySecretName = secret.Name
		bkp.Status.EncryptionKeySecretData = data
		bkp.Status.EncryptionKeySecretNamespace = secret.Namespace

		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// isEncryptKeyStatusEquals return true when something related to the aws status fields change
func isEncryptKeyStatusEquals(secret *corev1.Secret, bkp *v1alpha1.Backup, data map[string]string) bool {
	return secret.Name != bkp.Status.EncryptionKeySecretName || secret.Namespace != bkp.Status.EncryptionKeySecretNamespace || !reflect.DeepEqual(data, bkp.Status.EncryptionKeySecretData)
}

// covertDataSecretToString coverts data secret in []byte to map[string]string
func covertDataSecretToString(secret *corev1.Secret) map[string]string {
	data := make(map[string]string)
	if secret.Data != nil {
		for k, v := range secret.Data {
			value := ""
			if v != nil {
				value = string(v)
			}
			data[k] = value
		}
		for k, v := range secret.StringData {
			data[k] = v
		}
	}
	return data
}

// updateDBSecretStatus returns error when was not possible update the EncryptionKey status fields in the CR successfully
func (r *ReconcileBackup) updateDBSecretStatus(request reconcile.Request) error {
	bkp, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	dbSecret, err := r.fetchSecret(bkp.Namespace, dbSecretPrefix+bkp.Name)
	if err != nil {
		return err
	}

	// Check if the Secret with the DB Secret was changed, if yes update its status
	if err := r.insertUpdateDBSecretStatus(dbSecret, bkp); err != nil {
		return err
	}
	return nil
}

// insertUpdateDBSecretStatus will check and update the DB Secret status if the Secret with the DB data was changed
func (r *ReconcileBackup) insertUpdateDBSecretStatus(dbSecret *corev1.Secret, bkp *v1alpha1.Backup) error {
	data := covertDataSecretToString(dbSecret)
	if dbSecret.Name != bkp.Status.DBSecretName || !reflect.DeepEqual(data, bkp.Status.DBSecretData) {
		bkp.Status.DBSecretName = dbSecret.Name
		bkp.Status.DBSecretData = data

		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// updatePodDatabaseFoundStatus returns error when was not possible update the DB Pod Found status field in the CR successfully
func (r *ReconcileBackup) updatePodDatabaseFoundStatus(request reconcile.Request) error {
	bkp, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	// Check if the Pod Database Found status changed, if yes update it
	if err := r.insertUpdatePodDbFoundStatus(bkp); err != nil {
		return err
	}
	return nil
}

// insertUpdatePodDbFoundStatus will check and update the Pod Found status changed and update it
func (r *ReconcileBackup) insertUpdatePodDbFoundStatus(bkp *v1alpha1.Backup) error {
	if r.isDbPodFound() != bkp.Status.DatabasePodFound {
		bkp.Status.DatabasePodFound = r.isDbPodFound()

		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// updateServiceDbServiceFoundStatus returns error when was not possible update the DB Service Found status field in the CR successfully
func (r *ReconcileBackup) updateServiceDbServiceFoundStatus(request reconcile.Request) error {
	// Get the latest version of the CR
	bkp, err := r.fetchBackupCR(request)
	if err != nil {
		return err
	}

	// Check if the Database Found status changed, if yes update it
	if err := r.insertUpdateDbServiceFoundStatus(bkp); err != nil {
		return err
	}
	return nil
}

// insertUpdatePodDbFoundStatus will check and update the Database Found status changed and update it
func (r *ReconcileBackup) insertUpdateDbServiceFoundStatus(bkp *v1alpha1.Backup) error {
	if r.isDbServiceFound() != bkp.Status.DatabaseServiceFound {
		bkp.Status.DatabaseServiceFound = r.isDbServiceFound()
		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

//isDbServiceFound returns false when the database service which should be created by the PostgreSQL controller was not found
func (r *ReconcileBackup) isDbServiceFound() bool {
	return &r.dbService != nil && len(r.dbService.Name) > 0
}

//isDbPodFound returns false when the database pod which should be created by the PostgreSQL controller was not found
func (r *ReconcileBackup) isDbPodFound() bool {
	return &r.dbService != nil && len(r.dbService.Name) > 0
}

//isAllCreated returns error when some resource is missing
func (r *ReconcileBackup) isAllCreated(bkp *v1alpha1.Backup) error {

	// Check if was possible found the DB Pod
	if !r.isDbPodFound() {
		err := fmt.Errorf("Unable to set OK Status for Backup. The postgresql pod was not found")
		return err
	}

	// Check if was possible found the DB Service
	if !r.isDbServiceFound() {
		err := fmt.Errorf("Unable to set OK Status for Backup. The postgresql database service was not found")
		return err
	}

	// Check if DB secret was created
	dbSecretName := dbSecretPrefix + bkp.Name
	_, err := r.fetchSecret(bkp.Namespace, dbSecretName)
	if err != nil {
		err = fmt.Errorf("Unable to set OK Status for Backup. The DB Secret name %v was not found", dbSecretName)
		return err
	}

	// Check if AWS secret was created
	_, err = r.fetchSecret(getAwsSecretNamespace(bkp), getAWSSecretName(bkp))
	if err != nil {
		err := fmt.Errorf("Unable to set OK Status for Backup. The AWS Secret name %v in the namespace %v was not found", getAWSSecretName(bkp), getAwsSecretNamespace(bkp))
		return err
	}

	// Check if Enc secret was created (if was configured to be used)
	if isEncryptionKeyOptionConfig(bkp) {
		_, err := r.fetchSecret(getEncSecretNamespace(bkp), getEncSecretName(bkp))
		if err != nil {
			err := fmt.Errorf("Unable to set OK Status for Backup. The Encript Key configured was not found")
			return err
		}
	}

	//check if the cronJob was created
	_, err = r.fetchCronJob(bkp)
	if err != nil {
		err := fmt.Errorf("Unable to set OK Status for Backup. The CronJob was not found")
		return err
	}

	return nil
}
