package backup

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"
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
			Name:      "backup",
			Namespace: "postgresql-operator",
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
			Name:      "backup",
			Namespace: "postgresql-operator",
		},
		Spec: v1alpha1.BackupSpec{
			EncryptKeySecretName:      "enc-secret-test",
			EncryptKeySecretNamespace: "postgresql-operator",
			AwsSecretName:             "aws-secret-test",
			AwsSecretNamespace:        "postgresql-operator",
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
			Name:      "backup",
			Namespace: "postgresql-operator",
		},
		Spec: v1alpha1.BackupSpec{
			GpgPublicKey:  "example-gpgPublicKey",
			GpgEmail:      "email@gmai.com",
			GpgTrustModel: "always",
		},
	}

	/**
	Mock of Database resource
	*/

	dbInstanceWithoutSpec = v1alpha1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "database",
			Namespace: "postgresql-operator",
		},
	}

	podDatabase = corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "database-test",
			Namespace: "postgresql-operator",
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
						Value: "/var/lib/pgsql/data",
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

	dbInstanceWithConfigMap = v1alpha1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "database",
			Namespace: "postgresql-operator",
		},
		Spec: v1alpha1.DatabaseSpec{
			ConfigMapName:             "config-map-test",
			DatabaseNameKeyEnvVar:     "POSTGRESQL_DATABASE",
			DatabasePasswordKeyEnvVar: "POSTGRESQL_PASSWORD",
			DatabaseUserKeyEnvVar:     "POSTGRESQL_USER",
			DatabaseName:              "solution-database-name",
			DatabasePassword:          "postgres",
			DatabaseUser:              "postgresql",
		},
	}

	podDatabaseConfigMap = corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "database-test",
			Namespace: "postgresql-operator",
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
						Name: dbInstanceWithConfigMap.Spec.DatabaseNameKeyEnvVar,
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: dbInstanceWithConfigMap.Spec.ConfigMapName,
								},
								Key: utils.GetEnvVarKey(
									dbInstanceWithConfigMap.Spec.ConfigMapDatabaseNameKey,
									dbInstanceWithConfigMap.Spec.DatabaseNameKeyEnvVar),
							},
						},
					},
					corev1.EnvVar{
						Name: dbInstanceWithConfigMap.Spec.DatabaseUserKeyEnvVar,
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: dbInstanceWithConfigMap.Spec.ConfigMapName,
								},
								Key: utils.GetEnvVarKey(dbInstanceWithConfigMap.Spec.ConfigMapDatabaseUserKey,
									dbInstanceWithConfigMap.Spec.DatabaseUserKeyEnvVar),
							},
						},
					},
					corev1.EnvVar{
						Name: dbInstanceWithConfigMap.Spec.DatabasePasswordKeyEnvVar,
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: dbInstanceWithConfigMap.Spec.ConfigMapName,
								},
								Key: utils.GetEnvVarKey(dbInstanceWithConfigMap.Spec.ConfigMapDatabasePasswordKey,
									dbInstanceWithConfigMap.Spec.DatabasePasswordKeyEnvVar),
							},
						},
					},
					{
						Name:  "PGDATA",
						Value: "/var/lib/pgsql/data",
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
			Name:      "database",
			Namespace: "postgresql-operator",
			Labels:    utils.GetLabels(dbInstanceWithoutSpec.Name),
		},
	}

	configMapDefault = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-map-test",
			Namespace: "postgresql-operator",
		},
		Data: map[string]string{
			"POSTGRESQL_DATABASE": "solution-database-name",
			"POSTGRESQL_PASSWORD": "postgres",
			"POSTGRESQL_USER":     "postgresql",
		},
	}

	configMapInvalidDatabaseKey = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-otherkeys",
			Namespace: "postgresql-operator",
		},
		Data: map[string]string{
			"invalid": "dbname",
			dbInstanceWithConfigMap.Spec.DatabaseUserKeyEnvVar:     "root",
			dbInstanceWithConfigMap.Spec.DatabasePasswordKeyEnvVar: "root",
		},
	}

	configMapInvalidUserKey = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-otherkeys",
			Namespace: "postgresql-operator",
		},
		Data: map[string]string{
			dbInstanceWithConfigMap.Spec.DatabaseNameKeyEnvVar: "dbname",
			"invalid": "root",
			dbInstanceWithConfigMap.Spec.DatabasePasswordKeyEnvVar: "root",
		},
	}

	configMapInvalidPwdKey = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-otherkeys",
			Namespace: "postgresql-operator",
		},
		Data: map[string]string{
			dbInstanceWithConfigMap.Spec.DatabaseNameKeyEnvVar: "dbname",
			dbInstanceWithConfigMap.Spec.DatabaseUserKeyEnvVar: "root",
			"invalid": "root",
		},
	}

	dbInstanceWithConfigMapAndCustomizeKeys = v1alpha1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "database",
			Namespace: "postgresql-operator",
		},
		Spec: v1alpha1.DatabaseSpec{
			ConfigMapName:                "config-otherkeys",
			ConfigMapDatabaseNameKey:     "PGDATABASE",
			ConfigMapDatabasePasswordKey: "PGPASSWORD",
			ConfigMapDatabaseUserKey:     "PGUSER",
		},
	}
)
