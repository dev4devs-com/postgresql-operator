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

type DatabaseSecretData struct {
	databaseName string
	user         string
	pwd          string
	host         string
	superuser    string
	dbVersion    string
}

// buildDBSecretData will returns the data required to create the database secret according to the configuration
// NOTE: The user can:
// - Customize the environment variables keys as values that should be used with
// - Inform the name and namespace of an Config Map as the keys which has the values which should be used (E.g. user, password and database name already setup for another application )
func (r *ReconcileBackup) buildDBSecretData(bkp *v1alpha1.Backup, db *v1alpha1.Postgresql) (map[string][]byte, error) {

	dbData := &DatabaseSecretData{
		host:      r.dbService.Name + "." + bkp.Namespace + ".svc",
		superuser: "false",
		dbVersion: bkp.Spec.DatabaseVersion,
	}

	for i := 0; i < len(r.dbPod.Spec.Containers[0].Env); i++ {

		envVarName := r.dbPod.Spec.Containers[0].Env[i].Name
		envVarValue := r.dbPod.Spec.Containers[0].Env[i].Value

		var cfgName, cfgKey string
		if r.dbPod.Spec.Containers[0].Env[i].ValueFrom != nil {
			cfgName = r.dbPod.Spec.Containers[0].Env[i].ValueFrom.ConfigMapKeyRef.Name
			cfgKey = r.dbPod.Spec.Containers[0].Env[i].ValueFrom.ConfigMapKeyRef.Key
		}

		switch envVarName {
		case utils.GetConfigMapEnvVarKey(db.Spec.ConfigMapDatabaseNameParam, db.Spec.DatabaseNameParam):
			dbData.databaseName = envVarValue
			if dbData.databaseName == "" {
				database, err := r.getKeyValueFromConfigMap(cfgName, bkp.Namespace, cfgKey)
				if database == "" || err != nil {
					err := fmt.Errorf("Unable to get the user the database env var %v to create the secret",
						utils.GetConfigMapEnvVarKey(db.Spec.ConfigMapDatabaseNameParam, db.Spec.DatabaseNameParam))
					return nil, err
				}
			}
		case utils.GetConfigMapEnvVarKey(db.Spec.ConfigMapDatabaseUserParam, db.Spec.DatabaseUserParam):
			dbData.user = envVarValue
			if dbData.user == "" {
				user, err := r.getKeyValueFromConfigMap(cfgName, bkp.Namespace, cfgKey)
				if user == "" || err != nil {
					err := fmt.Errorf("Unable to get the user  the database env var %v to create the secret",
						utils.GetConfigMapEnvVarKey(db.Spec.ConfigMapDatabaseUserParam, db.Spec.DatabaseUserParam))
					return nil, err
				}
			}
		case utils.GetConfigMapEnvVarKey(db.Spec.ConfigMapDatabasePasswordParam, db.Spec.DatabasePasswordParam):
			dbData.pwd = envVarValue
			if dbData.pwd == "" {
				pwd, err := r.getKeyValueFromConfigMap(cfgName, bkp.Namespace, cfgKey)
				if pwd == "" || err != nil {
					err := fmt.Errorf("Unable to get the pwd for the database env var %v to create the secret",
						utils.GetConfigMapEnvVarKey(db.Spec.ConfigMapDatabasePasswordParam, db.Spec.DatabasePasswordParam))
					return nil, err
				}
			}
		}
	}
	return createDbDataByteMap(dbData), nil
}

// getKeyValueFromConfigMap returns the value of some key defined in the ConfigMap
func (r *ReconcileBackup) getKeyValueFromConfigMap(configMapName, configMapNamespace, configMapKey string) (string, error) {
	// search for ConfigMap
	cfg, err := r.fetchConfigMap(configMapName, configMapNamespace)
	if err != nil {
		return "", err
	}
	// Get ENV value
	return cfg.Data[configMapKey], nil
}

// createDbDataByteMap returns the a map with the data in the []byte format required to create the database secret
func createDbDataByteMap(data *DatabaseSecretData) map[string][]byte {
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
