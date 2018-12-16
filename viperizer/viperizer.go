package viperizer

import (
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	proxyFile = "proxy.yaml"
	serverFile = "grpcserver.yaml"
)

type Viperizer struct {
	Service string
	Config 	*viper.Viper
}

func NewGrpcServerViperizer(service string) (*Viperizer, error) {
	viper.SetConfigName("grpc_server") // name of config file (without extension)
	viper.AddConfigPath(os.Getenv("$HOME")) // name of config file (without extension)
	viper.AddConfigPath(".")                // optionally look for config in the working directory
	viper.AutomaticEnv()
	viper.SetEnvPrefix(service+"_SERVER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// If a config file is found, read it in."A generator for gRPC based Applications"
	if err := viper.ReadInConfig(); err != nil {
		log.Info("failed to read config file, writing defaults to -->"+serverFile)
		if err := viper.WriteConfigAs(serverFile); err != nil {
			return nil, err
		}

	} else {
		log.Info("Using config file-->", viper.ConfigFileUsed())
		if err := viper.WriteConfig(); err != nil {
			return nil, err
		}
	}
	vi := &Viperizer{
		Service: service,
		Config:  nil,
	}
	vi.SetGrpcServerDefaults()
	return vi, nil
}

func NewProxyViperizer(service string) (*Viperizer, error) {
	viper.SetConfigName("proxy") // name of config file (without extension)
	viper.AddConfigPath(os.Getenv("$HOME")) // name of config file (without extension)
	viper.AddConfigPath(".")                // optionally look for config in the working directory
	viper.AutomaticEnv()
	viper.SetEnvPrefix(service+"_PROXY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// If a config file is found, read it in."A generator for gRPC based Applications"
	if err := viper.ReadInConfig(); err != nil {
		log.Info("failed to read config file, writing defaults to -->"+proxyFile)
		if err := viper.WriteConfigAs(proxyFile); err != nil {
			return nil, err
		}

	} else {
		log.Info("Using config file-->", viper.ConfigFileUsed())
		if err := viper.WriteConfig(); err != nil {
			return nil, err
		}
	}
	vi := &Viperizer{
		Service: service,
	}
	vi.SetProxyDefaults()
	return vi, nil
}

func (v *Viperizer) GetService() string {
	return v.Service
}

func (v *Viperizer) GetViper() *viper.Viper {
	return v.Config
}

func (v *Viperizer) SetGrpcServerDefaults() error {
	v.Config.SetDefault("server.tls", false)
	v.Config.SetDefault("server.tracing", true)
	v.Config.SetDefault("server.metrics_endpoint", true)
	v.Config.SetDefault("server.live_endpoint", true)
	v.Config.SetDefault("server.ready_endpoint", true)
	v.Config.SetDefault("server.pprof_endpoint", true)
	v.Config.SetDefault("server.db_host", "localhost")
	v.Config.SetDefault("server.db_port", ":5432")
	v.Config.SetDefault("server.db_name", "postgresdb")
	v.Config.SetDefault("server.db_user", "admin")
	v.Config.SetDefault("server.port", ":8443")
	v.Config.SetDefault("server.routine_threshold", 300)
	v.Config.SetDefault("server.jaeger_metrics", true)
	log.Info("updating grpc server defaults-->"+serverFile)
	if err := v.WriteConfig(serverFile); err != nil {
		return err
	}
	return nil
}

func (v *Viperizer) SetProxyDefaults() error {
	v.Config.SetDefault("proxy.prefix", true)
	v.Config.SetDefault("proxy.tracing", true)
	v.Config.SetDefault("proxy.tls", false)
	v.Config.SetDefault("proxy.metrics_endpoint", true)
	v.Config.SetDefault("proxy.live_endpoint", true)
	v.Config.SetDefault("proxy.ready_endpoint", true)
	v.Config.SetDefault("proxy.pprof_endpoint", true)
	v.Config.SetDefault("proxy.db_host", "localhost")
	v.Config.SetDefault("proxy.db_port", ":5432")
	v.Config.SetDefault("proxy.db_name", "postgresdb")
	v.Config.SetDefault("proxy.db_user", "admin")
	v.Config.SetDefault("proxy.grpc_port", ":8443")
	v.Config.SetDefault("proxy.routine_threshold", 300)
	v.Config.SetDefault("proxy.jaeger_metrics", true)
	v.Config.SetDefault("proxy.backend", true)
	v.Config.SetDefault("proxy.log_level", true)
	v.Config.SetDefault("proxy.swagger_file", true)
	v.Config.SetDefault("proxy.allow_origin", true)
	v.Config.SetDefault("proxy.allow_creds", true)
	v.Config.SetDefault("proxy.allow_methods", true)
	v.Config.SetDefault("proxy.allow_headers", true)

	log.Info("updating proxy defaults-->"+proxyFile)
	if err := v.WriteConfig(proxyFile); err != nil {
		return err
	}
	return nil
}

func (v *Viperizer) ReadConfig(file string) error {
	// If a config file is found, read it in."A generator for gRPC based Applications"
	if err := v.Config.ReadInConfig(); err != nil {
		log.Info("failed to read config file, writing defaults to -->"+file)
		if err = v.WriteConfig(file); err != nil {
			return err
		}
		return err
	} else {
		log.Info("Using config file-->", v.Config.ConfigFileUsed())
		if err = v.WriteConfig(file); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (v *Viperizer) WriteConfig(file string) error {
	if err := v.Config.WriteConfigAs(file); err != nil {
		return err
	} else {
		log.Info("grpc config file created-->", v.Config.ConfigFileUsed())
		if err := v.Config.WriteConfig(); err != nil {
			return err
		}
	}
	return nil
}