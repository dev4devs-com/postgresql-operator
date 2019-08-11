package postgresql

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func getDBLabels(name string) map[string]string {
	return map[string]string{"app": "postgresql", "postgresql_cr": name}
}

//buildDatabaseNameEnvVar return the corev1.EnvVar object wth the key:value for the database name
func buildDatabaseNameEnvVar(db *v1alpha1.Postgresql) corev1.EnvVar {
	if len(db.Spec.ConfigMapName) > 0 {
		return corev1.EnvVar{
			Name: db.Spec.DatabaseNameParam,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: db.Spec.ConfigMapName,
					},
					Key: getConfigMapEnvVarKey(db.Spec.ConfigMapDatabaseNameParam, db.Spec.DatabaseNameParam),
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  db.Spec.DatabaseNameParam,
		Value: db.Spec.DatabaseName,
	}
}

func getConfigMapEnvVarKey(cgfKey, defaultKey string) string {
	if len(cgfKey) > 0 {
		return cgfKey
	}
	return defaultKey
}

//buildDatabaseUserEnvVar return the corev1.EnvVar object wth the key:value for the database user
func buildDatabaseUserEnvVar(db *v1alpha1.Postgresql) corev1.EnvVar {
	if len(db.Spec.ConfigMapName) > 0 {
		return corev1.EnvVar{
			Name: db.Spec.DatabaseUserParam,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: db.Spec.ConfigMapName,
					},
					Key: getConfigMapEnvVarKey(db.Spec.ConfigMapDatabaseUserParam, db.Spec.DatabaseUserParam),
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  db.Spec.DatabaseUserParam,
		Value: db.Spec.DatabaseUser,
	}
}

//buildDatabasePasswordEnvVar return the corev1.EnvVar object wth the key:value for the database pwd
func buildDatabasePasswordEnvVar(db *v1alpha1.Postgresql) corev1.EnvVar {
	if len(db.Spec.ConfigMapName) > 0 {
		return corev1.EnvVar{
			Name: db.Spec.DatabasePasswordParam,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: db.Spec.ConfigMapName,
					},
					Key: getConfigMapEnvVarKey(db.Spec.ConfigMapDatabasePasswordParam, db.Spec.DatabasePasswordParam),
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  db.Spec.DatabasePasswordParam,
		Value: db.Spec.DatabasePassword,
	}
}
