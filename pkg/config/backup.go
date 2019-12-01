package config

const (
	schedule        = "0 0 * * *"
	bakupImage      = "quay.io/integreatly/backup-container:1.0.8"
	databaseVersion = "9.6"
	databaseCRName  = "database"
)

type DefaultBackupConfig struct {
	Schedule        string `json:"schedule"`
	Image           string `json:"image"`
	DatabaseVersion string `json:"databaseVersion"`
	DatabaseCRName  string `json:"databaseCRName"`
}

func NewDefaultBackupConfig() *DefaultBackupConfig {
	return &DefaultBackupConfig{
		Schedule:        schedule,
		Image:           bakupImage,
		DatabaseVersion: databaseVersion,
		DatabaseCRName:  databaseCRName,
	}
}
