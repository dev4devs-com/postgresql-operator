package backup

import (
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
)

const (
	awsSecretPrefix = "aws-"
	dbSecretPrefix  = "db-"
	encSecretPrefix = "encryption-"
)

func getBkpLabels(name string) map[string]string {
	return map[string]string{"app": "postgresql", "backup_cr": name}
}

// DbSecret keep the data which will be used in the DB secret
type DbSecret struct {
	databaseName string
	user         string
	pwd          string
	host         string
	superuser    string
	dbVersion    string
}

// HelperDbSecret just help build the Map Data for the DB Secret
type HelperDbSecret struct {
	envVarName   string
	envVarValue  string
	cfgName      string
	cfgKey       string
	cfgNamespace string
}

// buildDBSecretData will returns the data required to create the database secret according to the configuration
// NOTE: The user can:
// - Customize the environment variables keys as values that should be used with
// - Inform the name and namespace of an Config Map as the keys which has the values which should be used (E.g. user, password and database name already setup for another application )
func (r *ReconcileBackup) buildDBSecretData(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql) (map[string][]byte, error) {

	dbSecret := r.newDBSecret(bkp)

	for i := 0; i < len(r.dbPod.Spec.Containers[0].Env); i++ {

		helper := r.newHelperDbSecret(i, bkp)
		var err error

		switch helper.envVarName {
		case utils.GetEnvVarKey(db.Spec.ConfigMapDatabaseNameParam, db.Spec.DatabaseNameParam):
			dbSecret.databaseName, err = r.getEnvVarValue(dbSecret.databaseName, dbSecret, helper)
			if err != nil {
				return nil, err
			}
		case utils.GetEnvVarKey(db.Spec.ConfigMapDatabaseUserParam, db.Spec.DatabaseUserParam):
			dbSecret.user, err = r.getEnvVarValue(dbSecret.user, dbSecret, helper)
			if err != nil {
				return nil, err
			}
		case utils.GetEnvVarKey(db.Spec.ConfigMapDatabasePasswordParam, db.Spec.DatabasePasswordParam):
			dbSecret.pwd, err = r.getEnvVarValue(dbSecret.pwd, dbSecret, helper)
			if err != nil {
				return nil, err
			}
		}
	}
	return dbSecret.createMap(), nil
}

// getEnvVarValue will return the value that should be used for the Key informed
func (r *ReconcileBackup) getEnvVarValue(value string, dbSecret *DbSecret, helper *HelperDbSecret) (string, error) {
	value = helper.envVarValue
	if value == "" {
		value = r.getKeyValueFromConfigMap(helper)
		if value == "" {
			return "", helper.newErrorUnableToGetKeyFrom()
		}
	}
	return value, nil
}

// newHelperDbSecret is a strtuct to keep the data in the loop in order to help fid the key and values which should be used
func (r *ReconcileBackup) newHelperDbSecret(i int, bkp *v1alpha1.Backup) *HelperDbSecret {
	dt := new(HelperDbSecret)
	dt.envVarName = r.dbPod.Spec.Containers[0].Env[i].Name
	dt.envVarValue = r.dbPod.Spec.Containers[0].Env[i].Value
	dt.cfgNamespace = bkp.Namespace
	if r.dbPod.Spec.Containers[0].Env[i].ValueFrom != nil {
		dt.cfgName = r.dbPod.Spec.Containers[0].Env[i].ValueFrom.ConfigMapKeyRef.Name
		dt.cfgKey = r.dbPod.Spec.Containers[0].Env[i].ValueFrom.ConfigMapKeyRef.Key
	}
	return dt
}

// newDBSecret will create the DbSecret with the data which is required to add in its secret
func (r *ReconcileBackup) newDBSecret(bkp *v1alpha1.Backup) *DbSecret {
	db := new(DbSecret)
	db.host = r.dbService.Name + "." + bkp.Namespace + ".svc"
	db.superuser = "false"
	db.dbVersion = bkp.Spec.DatabaseVersion
	return db
}

// newErrorUnableToGetKeyFrom returns an error when is not possible find the key into the configMap and namespace in order
// to create the mandatory envvar for the database
func (dt *HelperDbSecret) newErrorUnableToGetKeyFrom() error {
	return fmt.Errorf("Unable to get the key (%v) in the configMap (%v) in the namespace (%v) to create the secret",
		dt.cfgKey, dt.cfgName, dt.cfgNamespace)
}

// getKeyValueFromConfigMap returns the value of some key defined in the ConfigMap
func (r *ReconcileBackup) getKeyValueFromConfigMap(dt *HelperDbSecret) string {
	// search for ConfigMap
	cfg, err := r.fetchConfigMap(dt.cfgName, dt.cfgNamespace)
	if err != nil {
		return ""
	}
	// Get ENV value
	return cfg.Data[dt.cfgKey]
}

// createMap returns the a map with the data in the []byte format required to create the database secret
func (data *DbSecret) createMap() map[string][]byte {
	return map[string][]byte{
		"POSTGRES_USERNAME":  []byte(data.user),
		"POSTGRES_PASSWORD":  []byte(data.pwd),
		"POSTGRES_DATABASE":  []byte(data.databaseName),
		"POSTGRES_HOST":      []byte(data.host),
		"POSTGRES_SUPERUSER": []byte(data.superuser),
		"VERSION":            []byte(data.dbVersion),
	}
}

func createAwsDataByteMap(bkp *v1alpha1.Backup) map[string][]byte {
	dataByte := map[string][]byte{
		"AWS_S3_BUCKET_NAME":    []byte(bkp.Spec.AwsS3BucketName),
		"AWS_ACCESS_KEY_ID":     []byte(bkp.Spec.AwsAccessKeyId),
		"AWS_SECRET_ACCESS_KEY": []byte(bkp.Spec.AwsSecretAccessKey),
	}
	return dataByte
}

func createEncDataMaps(bkp *v1alpha1.Backup) (map[string][]byte, map[string]string) {
	dataByte := map[string][]byte{
		"GPG_PUBLIC_KEY": []byte(bkp.Spec.GpgPublicKey),
	}

	dataString := map[string]string{
		"GPG_RECIPIENT":   bkp.Spec.GpgEmail,
		"GPG_TRUST_MODEL": bkp.Spec.GpgTrustModel,
	}
	return dataByte, dataString
}

// getAWSSecretName returns the name of the secret
// NOTE: The user can just inform the name and namespace of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one in the same namespace where the backup is applied
func getAWSSecretName(bkp *v1alpha1.Backup) string {
	if isAwsKeySetupByName(bkp) {
		return bkp.Spec.AwsCredentialsSecretName
	}
	return awsSecretPrefix + bkp.Name
}

// getAwsSecretNamespace returns the namespace where the secret is applied already
// NOTE: The user can just inform the name and namespace of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one in the same namespace where the backup is applied
func getAwsSecretNamespace(bkp *v1alpha1.Backup) string {
	if isAwsKeySetupByName(bkp) && bkp.Spec.AwsCredentialsSecretNamespace != "" {
		return bkp.Spec.AwsCredentialsSecretNamespace
	}
	return bkp.Namespace
}

// getEncSecretNamespace returns the namespace where the secret is applied already
// NOTE: The user can just inform the name and namespace of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one in the same namespace where the backup is applied
func getEncSecretNamespace(bkp *v1alpha1.Backup) string {
	if isEncKeySetupByNameAndNamaspace(bkp) {
		return bkp.Spec.EncryptKeySecretNamespace
	}
	return bkp.Namespace
}

// getEncSecretName returns the name of the secret
// NOTE: The user can just inform the name and namespace of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one in the same namespace where the backup is applied
func getEncSecretName(bkp *v1alpha1.Backup) string {
	if isEncKeySetupByName(bkp) {
		return bkp.Spec.EncryptKeySecretName
	}
	return encSecretPrefix + bkp.Name
}

// isEncryptionKeyOptionConfig returns true when the CR has the configuration to allow it be used
func isEncryptionKeyOptionConfig(bkp *v1alpha1.Backup) bool {
	return bkp.Spec.AwsCredentialsSecretName != "" ||
		(bkp.Spec.GpgTrustModel != "" && bkp.Spec.GpgEmail != "" && bkp.Spec.GpgPublicKey != "")
}

// isEncKeySetupByName returns true when it is setup to get an pre-existing secret applied in the cluster.
// NOTE: The user can just inform the name of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one
func isEncKeySetupByName(bkp *v1alpha1.Backup) bool {
	return bkp.Spec.EncryptKeySecretName != ""
}

// isAwsKeySetupByName returns true when it is setup to get an pre-existing secret applied in the cluster.
// NOTE: The user can just inform the name of the Secret which is already applied in the cluster OR
// the data required for the operator be able to create one
func isAwsKeySetupByName(bkp *v1alpha1.Backup) bool {
	return bkp.Spec.AwsCredentialsSecretName != ""
}

func isEncKeySetupByNameAndNamaspace(bkp *v1alpha1.Backup) bool {
	return isEncKeySetupByName(bkp) && bkp.Spec.EncryptKeySecretNamespace != ""
}
