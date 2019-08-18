package utils

//GetEnvVarKey check if the customized key is in place for the configMap and returned the valid key
func GetEnvVarKey(cgfKey, defaultKey string) string {
	if len(cgfKey) > 0 {
		return cgfKey
	}
	return defaultKey
}
