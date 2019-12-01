package utils

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

//BuildDatabaseNameEnvVar return the corev1.EnvVar object wth the key:value for the database name
func BuildDatabaseNameEnvVar(db *v1alpha1.Database) corev1.EnvVar {
	if len(db.Spec.ConfigMapName) > 0 {
		return corev1.EnvVar{
			Name: db.Spec.DatabaseNameKeyEnvVar,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: db.Spec.ConfigMapName,
					},
					Key: GetEnvVarKey(db.Spec.ConfigMapDatabaseNameKey, db.Spec.DatabaseNameKeyEnvVar),
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  db.Spec.DatabaseNameKeyEnvVar,
		Value: db.Spec.DatabaseName,
	}
}

//BuildDatabaseUserEnvVar return the corev1.EnvVar object wth the key:value for the database user
func BuildDatabaseUserEnvVar(db *v1alpha1.Database) corev1.EnvVar {
	if len(db.Spec.ConfigMapName) > 0 {
		return corev1.EnvVar{
			Name: db.Spec.DatabaseUserKeyEnvVar,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: db.Spec.ConfigMapName,
					},
					Key: GetEnvVarKey(db.Spec.ConfigMapDatabaseUserKey, db.Spec.DatabaseUserKeyEnvVar),
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  db.Spec.DatabaseUserKeyEnvVar,
		Value: db.Spec.DatabaseUser,
	}
}

//BuildDatabasePasswordEnvVar return the corev1.EnvVar object wth the key:value for the database pwd
func BuildDatabasePasswordEnvVar(db *v1alpha1.Database) corev1.EnvVar {
	if len(db.Spec.ConfigMapName) > 0 {
		return corev1.EnvVar{
			Name: db.Spec.DatabasePasswordKeyEnvVar,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: db.Spec.ConfigMapName,
					},
					Key: GetEnvVarKey(db.Spec.ConfigMapDatabasePasswordKey, db.Spec.DatabasePasswordKeyEnvVar),
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  db.Spec.DatabasePasswordKeyEnvVar,
		Value: db.Spec.DatabasePassword,
	}
}

//GetEnvVarKey check if the customized key is in place for the configMap and returned the valid key
func GetEnvVarKey(cgfKey, defaultKey string) string {
	if len(cgfKey) > 0 {
		return cgfKey
	}
	return defaultKey
}
