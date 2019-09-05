package backup

import (
	"context"
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/service"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const statusOk = "OK"

//updateAppStatus returns error when status regards  all required resource could not be updated with OK
func (r *ReconcileBackup) updateBackupStatus(request reconcile.Request) error {
	bkp, err := service.FetchBackupCR(request.Name, request.Namespace, r.client)
	if err != nil {
		return err
	}

	statusMsgUpdate := statusOk
	// Check if all required resource were created and found
	if err := r.isAllCreated(bkp); err != nil {
		statusMsgUpdate = err.Error()
	}

	// Check if BackupStatus was changed, if yes update it
	if err := r.insertUpdateBackupStatus(bkp, statusMsgUpdate); err != nil {
		return err
	}
	return nil
}

// Check if BackupStatus was changed, if yes update it
func (r *ReconcileBackup) insertUpdateBackupStatus(bkp *v1alpha1.Backup, statusMsgUpdate string) error {
	if !reflect.DeepEqual(statusMsgUpdate, bkp.Status.BackupStatus) {
		bkp.Status.BackupStatus = statusOk
		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// updateCronJobStatus returns error when was not possible update the CronJob status successfully
func (r *ReconcileBackup) updateCronJobStatus(request reconcile.Request) error {
	bkp, err := service.FetchBackupCR(request.Name, request.Namespace, r.client)
	if err != nil {
		return err
	}

	// Check if Cronjob Name or Status was changed, if yes update it
	cronJob, err := service.FetchCronJob(bkp.Name, bkp.Namespace, r.client)
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
	bkp, err := service.FetchBackupCR(request.Name, request.Namespace, r.client)
	if err != nil {
		return err
	}

	aws, err := service.FetchSecret(utils.GetAwsSecretNamespace(bkp), utils.GetAWSSecretName(bkp), r.client)
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
	if isAwsStatusEqual(aws, bkp) {

		bkp.Status.AWSSecretName = aws.Name
		bkp.Status.AwsCredentialsSecretNamespace = aws.Namespace

		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// isAwsStatusEqual return true when something related to the aws status fields changed
func isAwsStatusEqual(aws *corev1.Secret, bkp *v1alpha1.Backup) bool {
	return aws.Name != bkp.Status.AWSSecretName  || aws.Namespace != bkp.Status.AwsCredentialsSecretNamespace
}

// updateAWSSecretStatus returns error when was not possible update the EncryptionKey status fields in the CR successfully
func (r *ReconcileBackup) updateEncSecretStatus(request reconcile.Request) error {
	bkp, err := service.FetchBackupCR(request.Name, request.Namespace, r.client)
	if err != nil {
		return err
	}

	isEncryptionKeyOptionConfig := utils.IsEncryptionKeyOptionConfig(bkp)
	if isEncryptionKeyOptionConfig {
		secret, err := service.FetchSecret(utils.GetEncSecretNamespace(bkp), utils.GetEncSecretName(bkp), r.client)
		if err != nil {
			return err
		}

		// Check if the Secret with the AWS data was changed, if yes update its status
		if err := r.insertUpdateEncKeyStatus(secret, bkp); err != nil {
			return err
		}
	}

	// Check if the config(boolean status) was changed, if yes update it
	if isEncryptionKeyOptionConfig != bkp.Status.HasEncryptKey {

		bkp.Status.HasEncryptKey = isEncryptionKeyOptionConfig
		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// insertUpdateEncKeyStatus will check and update the EncryptionKey Secret status if the Secret with the AWS data was changed
func (r *ReconcileBackup) insertUpdateEncKeyStatus(secret *corev1.Secret, bkp *v1alpha1.Backup) error {
	if isEncryptKeyStatusEquals(secret, bkp) {

		bkp.Status.EncryptKeySecretName = secret.Name
		bkp.Status.EncryptKeySecretNamespace = secret.Namespace

		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// isEncryptKeyStatusEquals return true when something related to the aws status fields change
func isEncryptKeyStatusEquals(secret *corev1.Secret, bkp *v1alpha1.Backup) bool {
	return secret.Name != bkp.Status.EncryptKeySecretName || secret.Namespace != bkp.Status.EncryptKeySecretNamespace
}

// updateDBSecretStatus returns error when was not possible update the EncryptionKey status fields in the CR successfully
func (r *ReconcileBackup) updateDBSecretStatus(request reconcile.Request) error {
	bkp, err := service.FetchBackupCR(request.Name, request.Namespace, r.client)
	if err != nil {
		return err
	}

	dbSecret, err := service.FetchSecret(bkp.Namespace, utils.DbSecretPrefix+bkp.Name, r.client)
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
	if dbSecret.Name != bkp.Status.DBSecretName {
		bkp.Status.DBSecretName = dbSecret.Name
		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// updatePodDatabaseFoundStatus returns error when was not possible update the DB Pod Found status field in the CR successfully
func (r *ReconcileBackup) updatePodDatabaseFoundStatus(request reconcile.Request) error {
	bkp, err := service.FetchBackupCR(request.Name, request.Namespace, r.client)
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
	if r.isDbPodFound() != bkp.Status.IsDatabasePodFound {
		bkp.Status.IsDatabasePodFound = r.isDbPodFound()

		if err := r.client.Status().Update(context.TODO(), bkp); err != nil {
			return err
		}
	}
	return nil
}

// updateDbServiceFoundStatus returns error when was not possible update the DB Service Found status field in the CR successfully
func (r *ReconcileBackup) updateDbServiceFoundStatus(request reconcile.Request) error {
	bkp, err := service.FetchBackupCR(request.Name, request.Namespace, r.client)
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
	if r.isDbServiceFound() != bkp.Status.IsDatabaseServiceFound {
		bkp.Status.IsDatabaseServiceFound = r.isDbServiceFound()
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
		err := fmt.Errorf("Error: PostgreSQL Pod is missing")
		return err
	}

	// Check if was possible found the DB Service
	if !r.isDbServiceFound() {
		err := fmt.Errorf("Error: PostgreSQL Service is missing")
		return err
	}

	// Check if DB secret was created
	dbSecretName := utils.DbSecretPrefix + bkp.Name
	_, err := service.FetchSecret(bkp.Namespace, dbSecretName, r.client)
	if err != nil {
		err = fmt.Errorf("Error: DB Secret is missing. (%v)", dbSecretName)
		return err
	}

	// Check if AWS secret was created
	awsSecretName := utils.GetAwsSecretNamespace(bkp)
	awsSecretNamespace := utils.GetAWSSecretName(bkp)
	_, err = service.FetchSecret(awsSecretNamespace, awsSecretName, r.client)
	if err != nil {
		err := fmt.Errorf("Error: AWS Secret is missing. (name:%v,namespace:%v)", awsSecretName, awsSecretNamespace)
		return err
	}

	// Check if Enc secret was created (if was configured to be used)
	if utils.IsEncryptionKeyOptionConfig(bkp) {
		encSecretName := utils.GetEncSecretName(bkp)
		encSecretNamespace := utils.GetEncSecretNamespace(bkp)
		_, err := service.FetchSecret(encSecretNamespace, encSecretName, r.client)
		if err != nil {
			err := fmt.Errorf("Error: Encript Key Secret is missing. (name:%v,namespace:%v)", encSecretName, encSecretNamespace)
			return err
		}
	}

	//check if the cronJob was created
	_, err = service.FetchCronJob(bkp.Name, bkp.Namespace, r.client)
	if err != nil {
		err := fmt.Errorf("Error: CronJob is missing")
		return err
	}

	return nil
}
