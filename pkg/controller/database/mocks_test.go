package database

import (
	v1alpha1 "github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Centralized mock objects for use in tests
var (
	dbInstanceWithoutSpec = v1alpha1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "database",
			Namespace: "postgresql-operator",
		},
	}

	dbInstanceConfigMapSameKeys = v1alpha1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "database",
			Namespace: "postgresql-operator",
		},
		Spec: v1alpha1.DatabaseSpec{
			ConfigMapName: "config-samekeys",
		},
	}

	dbInstanceConfigMapOtherKeys = v1alpha1.Database{
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

	configMapOtherKeyValues = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-otherkeys",
			Namespace: "postgresql-operator",
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
			Namespace: "postgresql-operator",
		},
		Data: map[string]string{
			"POSTGRESQL_DATABASE":            "dbname",
			"POSTGRESQL_PASSWORD":            "root",
			"POSTGRESQL_USERPOSTGRESQL_USER": "root",
		},
	}
)
