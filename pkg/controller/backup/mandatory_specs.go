package backup

import (
	"github.com/dev4devs-com/postgresql-operator/pkg/apis/postgresqloperator/v1alpha1"
	"github.com/dev4devs-com/postgresql-operator/pkg/config"
)

var defaultBackupConfig = config.NewDefaultBackupConfig()

// addMandatorySpecsDefinitions will add the specs which are mandatory for Backup CR in the case them
// not be applied
func addMandatorySpecsDefinitions(bkp *v1alpha1.Backup) {

	/*
		 Backup Container
		---------------------
		See https://github.com/integr8ly/backup-container-image
	*/

	if bkp.Spec.Schedule == "" {
		bkp.Spec.Schedule = defaultBackupConfig.Schedule
	}

	if bkp.Spec.Image == "" {
		bkp.Spec.Image = defaultBackupConfig.Image
	}

	if bkp.Spec.DatabaseVersion == "" {
		bkp.Spec.DatabaseVersion = defaultBackupConfig.DatabaseVersion
	}
}
