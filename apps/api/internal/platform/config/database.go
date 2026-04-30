package config

func HasDatabase() bool {
	return Current.DatabaseURL != ""
}
