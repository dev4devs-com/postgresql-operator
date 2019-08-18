package backup

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql-operator/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/utils"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Centralized mock objects for use in tests
var (

	/**
	BKP CR using mandatory specs
	*/
	bkpInstanceWithMandatorySpec = v1alpha1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-backup",
			Namespace: "postgresql",
		},
	}

	awsSecretWithMadatorySpec = corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GetAWSSecretName(&bkpInstanceWithMandatorySpec),
			Namespace: utils.GetAwsSecretNamespace(&bkpInstanceWithMandatorySpec),
		},
	}

	cronJobWithMadatorySpec = v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bkpInstanceWithMandatorySpec.Name,
			Namespace: bkpInstanceWithMandatorySpec.Namespace,
		},
	}

	dbSecretWithMadatorySpec = corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.DbSecretPrefix + bkpInstanceWithMandatorySpec.Name,
			Namespace: bkpInstanceWithMandatorySpec.Namespace,
		},
	}

	/**
	BKP CR to test when the user pass the name of the secrets
	*/

	bkpInstanceWithSecretNames = v1alpha1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-backup",
			Namespace: "postgresql",
		},
		Spec: v1alpha1.BackupSpec{
			EncryptKeySecretName:          "enc-secret-test",
			EncryptKeySecretNamespace:     "postgresql",
			AwsCredentialsSecretName:      "aws-secret-test",
			AwsCredentialsSecretNamespace: "postgresql",
		},
	}

	awsSecretWithSecretNames = corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GetAWSSecretName(&bkpInstanceWithSecretNames),
			Namespace: utils.GetAwsSecretNamespace(&bkpInstanceWithSecretNames),
		},
	}

	croJobWithSecretNames = v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bkpInstanceWithSecretNames.Name,
			Namespace: bkpInstanceWithSecretNames.Namespace,
		},
	}

	encSecretWithSecretNames = corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GetEncSecretName(&bkpInstanceWithSecretNames),
			Namespace: utils.GetEncSecretNamespace(&bkpInstanceWithSecretNames),
		},
	}

	dbSecretWithSecretNames = corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.DbSecretPrefix + bkpInstanceWithSecretNames.Name,
			Namespace: bkpInstanceWithSecretNames.Namespace,
		},
	}

	/**
	BKP CR to test when the user pass the secret data
	*/

	bkpInstanceWithEncSecretData = v1alpha1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-backup",
			Namespace: "postgresql",
		},
		Spec: v1alpha1.BackupSpec{
			GpgPublicKey:  "example-gpgPublicKey",
			GpgEmail:      "email@gmai.com",
			GpgTrustModel: "always",
		},
	}

	/**
	Mock of PostgreSQL resource
	*/

	dbInstanceWithoutSpec = v1alpha1.Postgresql{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: "postgresql",
		},
	}

	podDatabase = corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-test",
			Namespace: "postgresql",
			Labels:    utils.GetLabels(dbInstanceWithoutSpec.Name),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Image: "postgresql",
				Name:  "postgresql",
				Ports: []corev1.ContainerPort{{
					ContainerPort: 5000,
					Protocol:      "TCP",
				}},
				Env: []corev1.EnvVar{
					corev1.EnvVar{
						Name:  "PGDATABASE",
						Value: "test",
					},
					corev1.EnvVar{
						Name:  "PGUSER",
						Value: "test",
					},
					corev1.EnvVar{
						Name:  "PGPASSWORD",
						Value: "test",
					},
					{
						Name:  "PGDATA",
						Value: "/var/lib/pgsql/data/pgdata",
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "test",
						MountPath: "/var/lib/pgsql/data",
					},
				},
			}},
		},
	}

	dbInstanceWithConfigMap = v1alpha1.Postgresql{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: "postgresql",
		},
		Spec: v1alpha1.PostgresqlSpec{
			ConfigMapName:         "config-map-test",
			DatabaseNameParam:     "POSTGRESQL_DATABASE",
			DatabasePasswordParam: "POSTGRESQL_PASSWORD",
			DatabaseUserParam:     "POSTGRESQL_USER",
			DatabaseName:          "solution-database-name",
			DatabasePassword:      "postgres",
			DatabaseUser:          "postgresql",
		},
	}

	podDatabaseConfigMap = corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-test",
			Namespace: "postgresql",
			Labels:    utils.GetLabels(dbInstanceWithConfigMap.Name),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Image:           dbInstanceWithConfigMap.Spec.Image,
				Name:            dbInstanceWithConfigMap.Spec.ContainerName,
				ImagePullPolicy: dbInstanceWithConfigMap.Spec.ContainerImagePullPolicy,
				Ports: []corev1.ContainerPort{{
					ContainerPort: dbInstanceWithConfigMap.Spec.DatabasePort,
					Protocol:      "TCP",
				}},
				Env: []corev1.EnvVar{
					corev1.EnvVar{
						Name: dbInstanceWithConfigMap.Spec.DatabaseNameParam,
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: dbInstanceWithConfigMap.Spec.ConfigMapName,
								},
								Key: utils.GetEnvVarKey(dbInstanceWithConfigMap.Spec.ConfigMapDatabaseNameParam, dbInstanceWithConfigMap.Spec.DatabaseNameParam),
							},
						},
					},
					corev1.EnvVar{
						Name: dbInstanceWithConfigMap.Spec.DatabaseUserParam,
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: dbInstanceWithConfigMap.Spec.ConfigMapName,
								},
								Key: utils.GetEnvVarKey(dbInstanceWithConfigMap.Spec.ConfigMapDatabaseUserParam, dbInstanceWithConfigMap.Spec.DatabaseUserParam),
							},
						},
					},
					corev1.EnvVar{
						Name: dbInstanceWithConfigMap.Spec.DatabasePasswordParam,
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: dbInstanceWithConfigMap.Spec.ConfigMapName,
								},
								Key: utils.GetEnvVarKey(dbInstanceWithConfigMap.Spec.ConfigMapDatabasePasswordParam, dbInstanceWithConfigMap.Spec.DatabasePasswordParam),
							},
						},
					},
					{
						Name:  "PGDATA",
						Value: "/var/lib/pgsql/data/pgdata",
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      dbInstanceWithConfigMap.Name,
						MountPath: "/var/lib/pgsql/data",
					},
				},
				LivenessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						Exec: &corev1.ExecAction{
							Command: []string{
								"/usr/libexec/check-container",
								"'--live'",
							},
						},
					},
					FailureThreshold:    3,
					InitialDelaySeconds: 120,
					PeriodSeconds:       10,
					TimeoutSeconds:      10,
					SuccessThreshold:    1,
				},
				ReadinessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						Exec: &corev1.ExecAction{
							Command: []string{
								"/usr/libexec/check-container",
							},
						},
					},
					FailureThreshold:    3,
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
					TimeoutSeconds:      1,
					SuccessThreshold:    1,
				},
				TerminationMessagePath: "/dev/termination-log",
			}},
		},
	}

	serviceDatabase = corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: "postgresql",
			Labels:    utils.GetLabels(dbInstanceWithoutSpec.Name),
		},
	}

	configMapDefault = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-map-test",
			Namespace: "postgresql",
		},
		Data: map[string]string{
			"POSTGRESQL_DATABASE": "solution-database-name",
			"POSTGRESQL_PASSWORD": "postgres",
			"POSTGRESQL_USER":     "postgresql",
		},
	}

	configMapOtherKeyValues = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-otherkeys",
			Namespace: "postgresql",
		},
		Data: map[string]string{
			dbInstanceWithConfigMap.Spec.DatabaseNameParam:     "dbname",
			dbInstanceWithConfigMap.Spec.DatabasePasswordParam: "root",
			dbInstanceWithConfigMap.Spec.DatabaseUserParam:     "root",
		},
	}

	configMapInvalidDatabaseKey = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-otherkeys",
			Namespace: "postgresql",
		},
		Data: map[string]string{
			"invalid": "dbname",
			dbInstanceWithConfigMap.Spec.DatabaseUserParam:     "root",
			dbInstanceWithConfigMap.Spec.DatabasePasswordParam: "root",
		},
	}

	configMapInvalidUserKey = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-otherkeys",
			Namespace: "postgresql",
		},
		Data: map[string]string{
			dbInstanceWithConfigMap.Spec.DatabaseNameParam: "dbname",
			"invalid": "root",
			dbInstanceWithConfigMap.Spec.DatabasePasswordParam: "root",
		},
	}

	configMapInvalidPwdKey = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-otherkeys",
			Namespace: "postgresql",
		},
		Data: map[string]string{
			dbInstanceWithConfigMap.Spec.DatabaseNameParam: "dbname",
			dbInstanceWithConfigMap.Spec.DatabaseUserParam: "root",
			"invalid": "root",
		},
	}

	dbInstanceWithConfigMapAndCustomizeKeys = v1alpha1.Postgresql{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: "postgresql",
		},
		Spec: v1alpha1.PostgresqlSpec{
			ConfigMapName:                  "config-otherkeys",
			ConfigMapDatabaseNameParam:     "PGDATABASE",
			ConfigMapDatabasePasswordParam: "PGPASSWORD",
			ConfigMapDatabaseUserParam:     "PGUSER",
		},
	}
)
