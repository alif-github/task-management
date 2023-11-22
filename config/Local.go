package config

type LocalConfig struct {
	Configuration
	Server struct {
		Protocol     string `json:"protocol"`
		Host         string `json:"host"`
		Port         string `json:"port"`
		Version      string `json:"version"`
		PrefixPath   string `json:"prefix_path"`
		Application  string `json:"application"`
		SignatureKey string `json:"signature_key" envconfig:"SIGNATURE_KEY"`
	} `json:"server"`
	Postgresql struct {
		Driver            string `json:"driver"`
		Address           string `json:"address" envconfig:"POSTGRESQL_ADDRESS"`
		DefaultSchema     string `json:"default_schema"`
		MaxOpenConnection int    `json:"max_open_connection"`
		MaxIdleConnection int    `json:"max_idle_connection"`
	} `json:"postgresql"`
	LogFile []string `json:"log_file"`
}

func (input LocalConfig) GetServer() Server {
	return Server{
		Protocol:     input.Server.Protocol,
		Host:         input.Server.Host,
		Port:         input.Server.Port,
		Version:      input.Server.Version,
		PrefixPath:   input.Server.PrefixPath,
		Application:  input.Server.Application,
		SignatureKey: input.Server.SignatureKey,
	}
}

func (input LocalConfig) GetPostgresql() Postgresql {
	return Postgresql{
		Driver:            input.Postgresql.Driver,
		Address:           input.Postgresql.Address,
		DefaultSchema:     input.Postgresql.DefaultSchema,
		MaxOpenConnection: input.Postgresql.MaxOpenConnection,
		MaxIdleConnection: input.Postgresql.MaxIdleConnection,
	}
}

func (input LocalConfig) GetLogFile() []string {
	return input.LogFile
}
