package utils

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

//BuildDatabaseNameEnvVar return the corev1.EnvVar object wth the key:value for the database name
func BuildDatabaseNameEnvVar(db *v1alpha1.Postgresql) corev1.EnvVar {
	if len(db.Spec.ConfigMapName) > 0 {
		return corev1.EnvVar{
			Name: db.Spec.DatabaseNameParam,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: db.Spec.ConfigMapName,
					},
					Key: GetEnvVarKey(db.Spec.ConfigMapDatabaseNameParam, db.Spec.DatabaseNameParam),
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  db.Spec.DatabaseNameParam,
		Value: db.Spec.DatabaseName,
	}
}

//BuildDatabaseUserEnvVar return the corev1.EnvVar object wth the key:value for the database user
func BuildDatabaseUserEnvVar(db *v1alpha1.Postgresql) corev1.EnvVar {
	if len(db.Spec.ConfigMapName) > 0 {
		return corev1.EnvVar{
			Name: db.Spec.DatabaseUserParam,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: db.Spec.ConfigMapName,
					},
					Key: GetEnvVarKey(db.Spec.ConfigMapDatabaseUserParam, db.Spec.DatabaseUserParam),
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  db.Spec.DatabaseUserParam,
		Value: db.Spec.DatabaseUser,
	}
}

//BuildDatabasePasswordEnvVar return the corev1.EnvVar object wth the key:value for the database pwd
func BuildDatabasePasswordEnvVar(db *v1alpha1.Postgresql) corev1.EnvVar {
	if len(db.Spec.ConfigMapName) > 0 {
		return corev1.EnvVar{
			Name: db.Spec.DatabasePasswordParam,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: db.Spec.ConfigMapName,
					},
					Key: GetEnvVarKey(db.Spec.ConfigMapDatabasePasswordParam, db.Spec.DatabasePasswordParam),
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  db.Spec.DatabasePasswordParam,
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
