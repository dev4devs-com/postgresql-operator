package postgresql

import (
	v1alpha1 "github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Centralized mock objects for use in tests
var (
	dbInstance = v1alpha1.Postgresql{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: "postgresql",
		},
		Spec: v1alpha1.PostgresqlSpec{
			Image:                  "centos/postgresql-96-centos7",
			Size:                   1,
			ContainerName:          "database",
			DatabaseNameParam:      "POSTGRESQL_DATABASE",
			DatabasePasswordParam:  "POSTGRESQL_PASSWORD",
			DatabaseUserParam:      "POSTGRESQL_USER",
			DatabasePort:           5432,
			DatabaseMemoryLimit:    "512Mi",
			DatabaseMemoryRequest:  "512Mi",
			DatabaseStorageRequest: "1Gi",
			DatabaseName:           "mock-test",
			DatabasePassword:       "postgres",
			DatabaseUser:           "postgresql",
		},
	}

	dbInstanceWithoutSpec = v1alpha1.Postgresql{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: "postgresql",
		},
	}

	dbInstanceConfigMapSameKeys = v1alpha1.Postgresql{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: "postgresql",
		},
		Spec: v1alpha1.PostgresqlSpec{
			ConfigMapName: "config-samekeys",
		},
	}

	dbInstanceConfigMapOtherKeys = v1alpha1.Postgresql{
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

	configMapOtherKeyValues = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-otherkeys",
			Namespace: "postgresql",
		},
		Data: map[string]string{
			"PGDATABASE": "dbname",
			"PGPASSWORD": "root",
			"PGUSER":     "root",
		},
	}

	configMapSameKeyValues = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-samekeys",
			Namespace: "postgresql",
		},
		Data: map[string]string{
			"POSTGRESQL_DATABASE":            "dbname",
			"POSTGRESQL_PASSWORD":            "root",
			"POSTGRESQL_USERPOSTGRESQL_USER": "root",
		},
	}
)
