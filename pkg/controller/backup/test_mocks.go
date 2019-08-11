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
			Image:                    "quay.io/integreatly/backup-container:latest",
			Schedule:                 "0 0 * * *",
			EncryptionKeySecretName:  "enc-secret-test",
			AwsCredentialsSecretName: "aws-secret-test",
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

	lsDB = map[string]string{"app": "postgresql", "postgresql_cr": "postgresql"}

	podDatabase = corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-db",
			Namespace: "postgresql",
			Labels:    lsDB,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Image: "postgresql-db",
				Name:  "postgresql-db",
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

	serviceDatabase = corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-db",
			Namespace: "postgresql",
			Labels:    lsDB,
		},
	}
)
