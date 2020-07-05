package utils

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresql/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/config"
)

var defaulDatabaseConfig = config.NewDatabaseConfig()

// AddDatabaseMandatorySpecs will add the specs which are mandatory for Database CR in the case them
// not be applied
func AddDatabaseMandatorySpecs(db *v1alpha1.Database) {

	/*
	   CR DB Resource
	   ---------------------
	*/

	if db.Spec.Size == 0 {
		db.Spec.Size = defaulDatabaseConfig.Size
	}

	/*
		Environment Variables
		---------------------
		The following values are used to create the ConfigMap and the Environment Variables which will use these values
	*/

	if db.Spec.DatabaseName == "" {
		db.Spec.DatabaseName = defaulDatabaseConfig.DatabaseName
	}

	if db.Spec.DatabasePassword == "" {
		db.Spec.DatabasePassword = defaulDatabaseConfig.DatabasePassword
	}

	if db.Spec.DatabaseUser == "" {
		db.Spec.DatabaseUser = defaulDatabaseConfig.DatabaseUser
	}

	/*
	   Database Container
	   ---------------------------------
	*/

	//Following are the values which will be used as the key label for the environment variable of the database image.
	if db.Spec.DatabaseNameKeyEnvVar == "" {
		db.Spec.DatabaseNameKeyEnvVar = defaulDatabaseConfig.DatabaseNameKeyEnvVar
	}

	if db.Spec.DatabasePasswordKeyEnvVar == "" {
		db.Spec.DatabasePasswordKeyEnvVar = defaulDatabaseConfig.DatabasePasswordKeyEnvVar
	}

	if db.Spec.DatabaseUserKeyEnvVar == "" {
		db.Spec.DatabaseUserKeyEnvVar = defaulDatabaseConfig.DatabaseUserKeyEnvVar
	}

	if db.Spec.Image == "" {
		db.Spec.Image = defaulDatabaseConfig.Image
	}

	if db.Spec.ContainerName == "" {
		db.Spec.ContainerName = defaulDatabaseConfig.ContainerName
	}

	if db.Spec.DatabaseMemoryLimit == "" {
		db.Spec.DatabaseMemoryLimit = defaulDatabaseConfig.DatabaseMemoryLimit
	}

	if db.Spec.DatabaseMemoryRequest == "" {
		db.Spec.DatabaseMemoryRequest = defaulDatabaseConfig.DatabaseMemoryRequest
	}

	if db.Spec.DatabaseStorageRequest == "" {
		db.Spec.DatabaseStorageRequest = defaulDatabaseConfig.DatabaseStorageRequest
	}

	if db.Spec.DatabaseCpu == "" {
		db.Spec.DatabaseCpu = defaulDatabaseConfig.DatabaseCpu
	}

	if db.Spec.DatabaseCpuLimit == "" {
		db.Spec.DatabaseCpuLimit = defaulDatabaseConfig.DatabaseCpuLimit
	}

	if db.Spec.DatabasePort == 0 {
		db.Spec.DatabasePort = defaulDatabaseConfig.DatabasePort
	}

	if len(db.Spec.DatabaseStorageClassName) < 1 {
		db.Spec.DatabaseStorageClassName = defaulDatabaseConfig.DatabaseStorageClassName
	}
}
