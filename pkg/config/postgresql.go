package config

const (
	size                      = 1
	databaseName              = "solution-database-name"
	databasePassword          = "postgres"
	databaseUser              = "postgresql"
	databaseNameKeyEnvVar     = "POSTGRESQL_DATABASE"
	databasePasswordKeyEnvVar = "POSTGRESQL_PASSWORD"
	databaseUserKeyEnvVar     = "POSTGRESQL_USER"
	image                     = "centos/postgresql-96-centos7"
	containerName             = "database"
	databasePort              = 5432
	databaseMemoryLimit       = "512Mi"
	databaseMemoryRequest     = "512Mi"
	databaseStorageRequest    = "1Gi"
)

type DefaultPostgreSQLConfig struct {
	Size                      int32  `json:"size"`
	Image                     string `json:"image"`
	DatabaseName              string `json:"databaseName"`
	DatabasePassword          string `json:"databasePassword"`
	DatabaseUser              string `json:"databaseUser"`
	DatabaseNameKeyEnvVar     string `json:"databaseNameKeyEnvVar"`
	DatabasePasswordKeyEnvVar string `json:"databasePasswordKeyEnvVar"`
	DatabaseUserKeyEnvVar     string `json:"databaseUserKeyEnvVar"`
	ContainerName             string `json:"containerName"`
	DatabasePort              int32  `json:"databasePort"`
	DatabaseMemoryLimit       string `json:"databaseMemoryLimit"`
	DatabaseMemoryRequest     string `json:"databaseMemoryRequest"`
	DatabaseStorageRequest    string `json:"databaseStorageRequest"`
}

func NewPostgreSQLConfig() *DefaultPostgreSQLConfig {
	return &DefaultPostgreSQLConfig{
		Size:                      size,
		Image:                     image,
		DatabaseName:              databaseName,
		DatabasePassword:          databasePassword,
		DatabaseUser:              databaseUser,
		DatabaseNameKeyEnvVar:     databaseNameKeyEnvVar,
		DatabasePasswordKeyEnvVar: databasePasswordKeyEnvVar,
		DatabaseUserKeyEnvVar:     databaseUserKeyEnvVar,
		ContainerName:             containerName,
		DatabasePort:              databasePort,
		DatabaseMemoryLimit:       databaseMemoryLimit,
		DatabaseMemoryRequest:     databaseMemoryRequest,
		DatabaseStorageRequest:    databaseStorageRequest,
	}
}
