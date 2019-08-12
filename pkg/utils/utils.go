package utils

//GetConfigMapEnvVarKey check if the customized key is in place for the configMap and returned the valid key
func GetConfigMapEnvVarKey(cgfKey, defaultKey string) string {
	if len(cgfKey) > 0 {
		return cgfKey
	}
	return defaultKey
}

