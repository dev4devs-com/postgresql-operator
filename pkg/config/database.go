package config

const (
	size                      = 1
	databaseName              = "solution"
	databasePassword          = "postgres"
	databaseUser              = "postgres"
	databaseNameKeyEnvVar     = "POSTGRESQL_DATABASE"
	databasePasswordKeyEnvVar = "POSTGRESQL_PASSWORD"
	databaseUserKeyEnvVar     = "POSTGRESQL_USER"
	image                     = "centos/postgresql-96-centos7"
	containerName             = "database"
	databasePort              = 5432
	databaseMemoryLimit       = "512Mi"
	databaseMemoryRequest     = "128Mi"
	databaseStorageRequest    = "1Gi"
	databaseStorageClassName  = "slow"
	databaseCpuLimit          = "60m"
	databaseCpu               = "30m"
)

type DefaultDatabaseConfig struct {
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
	DatabaseCpuLimit          string `json:"databaseCpuLimit"`
	DatabaseCpu               string `json:"databaseCpu"`
	DatabaseStorageRequest    string `json:"databaseStorageRequest"`
	DatabaseStorageClassName  string `json:"databaseStorageClassName"`
}

func NewDatabaseConfig() *DefaultDatabaseConfig {
	return &DefaultDatabaseConfig{
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
		DatabaseCpu:               databaseCpu,
		DatabaseCpuLimit:          databaseCpuLimit,
		DatabaseStorageRequest:    databaseStorageRequest,
		DatabaseStorageClassName:  databaseStorageClassName,
	}
}
