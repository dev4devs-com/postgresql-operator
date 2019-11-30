package backup

import (
	"fmt"
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/service"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
)

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
// - Inform the name and namespace of an Config Map as the keys which has the values which should be used (E.g. user,
// password and database name already setup for another application )
func (r *ReconcileBackup) buildDBSecretData(bkp *v1alpha1.Backup, db *v1alpha1.Database) (map[string][]byte, error) {

	dbSecret := r.newDBSecret(bkp)

	for i := 0; i < len(r.dbPod.Spec.Containers[0].Env); i++ {

		helper := r.newHelperDbSecret(i, bkp)
		var err error

		switch helper.envVarName {
		case utils.GetEnvVarKey(db.Spec.ConfigMapDatabaseNameKey, db.Spec.DatabaseNameKeyEnvVar):
			dbSecret.databaseName, err = r.getEnvVarValue(helper)
			if err != nil {
				return nil, err
			}
		case utils.GetEnvVarKey(db.Spec.ConfigMapDatabaseUserKey, db.Spec.DatabaseUserKeyEnvVar):
			dbSecret.user, err = r.getEnvVarValue(helper)
			if err != nil {
				return nil, err
			}
		case utils.GetEnvVarKey(db.Spec.ConfigMapDatabasePasswordKey, db.Spec.DatabasePasswordKeyEnvVar):
			dbSecret.pwd, err = r.getEnvVarValue(helper)
			if err != nil {
				return nil, err
			}
		}
	}

	return dbSecret.createMap(), nil
}

// getEnvVarValue will return the value that should be used for the Key informed
func (r *ReconcileBackup) getEnvVarValue(helper *HelperDbSecret) (string, error) {
	value := helper.envVarValue
	if value == "" {
		value = r.getKeyValueFromConfigMap(helper)
		if value == "" {
			return "", helper.newErrorUnableToGetKeyFrom()
		}
	}
	return value, nil
}

// newHelperDbSecret is a strtuct to keep the data in the loop in order
// to help fid the key and values which should be used
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

// newErrorUnableToGetKeyFrom returns an error when is not possible
// find the key into the configMap and namespace in order
// to create the mandatory envvar for the database
func (dt *HelperDbSecret) newErrorUnableToGetKeyFrom() error {
	return fmt.Errorf("Unable to get the key (%v) in the configMap"+
		" (%v) in the namespace (%v) to create the secret",
		dt.cfgKey, dt.cfgName, dt.cfgNamespace)
}

// getKeyValueFromConfigMap returns the value of some key defined in the ConfigMap
func (r *ReconcileBackup) getKeyValueFromConfigMap(dt *HelperDbSecret) string {
	// search for ConfigMap
	cfg, err := service.FetchConfigMap(dt.cfgName, dt.cfgNamespace, r.client)
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
		"AWS_ACCESS_KEY_ID":     []byte(bkp.Spec.AwsAccessKeyID),
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
