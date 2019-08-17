package backup

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Centralized mock objects for use in tests
var (
	bkpInstance = v1alpha1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-backup",
			Namespace: "postgresql",
		},
		Spec: v1alpha1.BackupSpec{
			Image:              "quay.io/integreatly/backup-container:latest",
			Schedule:           "0 0 * * *",
			AwsS3BucketName:    "example-awsS3BucketName",
			AwsAccessKeyId:     "example-awsAccessKeyId",
			AwsSecretAccessKey: "example-awsSecretAccessKey",
		},
	}

	bkpInstanceWithSecretNames = v1alpha1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-backup",
			Namespace: "postgresql",
		},
		Spec: v1alpha1.BackupSpec{
			EncryptionKeySecretName:       "enc-secret-test",
			EncryptionKeySecretNamespace:  "postgresql",
			AwsCredentialsSecretName:      "aws-secret-test",
			AwsCredentialsSecretNamespace: "postgresql",
		},
	}

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

	bkpInstanceNonDefaultNamespace = v1alpha1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-db-backup",
			Namespace: "postgresql-namespace",
		},
		Spec: v1alpha1.BackupSpec{
			Image:              "quay.io/integreatly/backup-container:latest",
			Schedule:           "0 0 * * *",
			AwsS3BucketName:    "example-awsS3BucketName",
			AwsAccessKeyId:     "example-awsAccessKeyId",
			AwsSecretAccessKey: "example-awsSecretAccessKey",
		},
	}

	bkpInstanceWithoutSpec = v1alpha1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-backup",
			Namespace: "postgresql",
		},
	}

	dbInstanceWithoutSpec = v1alpha1.Postgresql{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: "postgresql",
		},
	}

	lsDB = map[string]string{"app": "postgresql", "postgresql_cr": dbInstanceWithoutSpec.Name}

	podDatabase = corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-test",
			Namespace: "postgresql",
			Labels:    lsDB,
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

	podDatabaseConfigMap = corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-test",
			Namespace: "postgresql",
			Labels:    lsDB,
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
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: configMapOtherKeyValues.Name,
								},
							},
						},
						Name: dbInstanceWithConfigMap.Spec.DatabaseNameParam,
					},
					corev1.EnvVar{
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: configMapOtherKeyValues.Name,
								},
							},
						},
						Name: dbInstanceWithConfigMap.Spec.DatabaseUserParam,
					},
					corev1.EnvVar{
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: configMapOtherKeyValues.Name,
								},
							},
						},
						Name: dbInstanceWithConfigMap.Spec.DatabasePasswordParam,
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

	serviceDatabase = corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: "postgresql",
			Labels:    lsDB,
		},
	}

	awsSecretMock = corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-secret-test",
			Namespace: "postgresql",
			Labels:    lsDB,
		},
	}

	encSecretMock = corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "enc-secret-test",
			Namespace: "postgresql",
			Labels:    lsDB,
		},
	}

	dbInstanceWithConfigMap = v1alpha1.Postgresql{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: "postgresql",
		},
		Spec: v1alpha1.PostgresqlSpec{
			ConfigMapName: "config-otherkeys",
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

	configMapOtherKeyValuesInvalidKeys = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-otherkeys",
			Namespace: "postgresql",
		},
		Data: map[string]string{
			"PGDATABASE": "dbname",
			"DBPASSWORD": "root",
			"DBUSER":     "root",
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
