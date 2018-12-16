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
	GrpcConfig 	*viper.Viper
	ProxyConfig *viper.Viper
}

func NewGrpcServerViperizer(vi *viper.Viper, service string) error {
	vi.SetConfigName("grpcserver") // name of config file (without extension)
	vi.AddConfigPath(os.Getenv("$HOME")) // name of config file (without extension)
	vi.AddConfigPath(".")                // optionally look for config in the working directory
	vi.AutomaticEnv()
	vi.SetEnvPrefix(service+"_SERVER")
	vi.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vi.SetDefault("proxy.service", service)
	vi.SetDefault("server.tls", false)
	vi.SetDefault("server.tracing", true)
	vi.SetDefault("server.metrics_endpoint", true)
	vi.SetDefault("server.live_endpoint", true)
	vi.SetDefault("server.ready_endpoint", true)
	vi.SetDefault("server.pprof_endpoint", true)
	vi.SetDefault("server.db_host", "localhost")
	vi.SetDefault("server.db_port", ":5432")
	vi.SetDefault("server.db_name", "postgresdb")
	vi.SetDefault("server.db_user", "admin")
	vi.SetDefault("server.port", ":8443")
	vi.SetDefault("server.routine_threshold", 300)
	vi.SetDefault("server.jaeger_metrics", true)

	// If a config file is found, read it in."A generator for gRPC based Applications"
	if err := vi.ReadInConfig(); err != nil {
		log.Info("failed to read config file, writing defaults to -->"+serverFile)
		if err := vi.WriteConfigAs(serverFile); err != nil {
			return err
		}

	} else {
		log.Info("Using config file-->", vi.ConfigFileUsed())
		if err := vi.WriteConfig(); err != nil {
			return err
		}
	}

	return nil
}

func NewProxyViperizer(vi *viper.Viper, service string) error {

	vi.SetConfigName("proxy") // name of config file (without extension)
	vi.AddConfigPath(os.Getenv("$HOME")) // name of config file (without extension)
	vi.AddConfigPath(".")                // optionally look for config in the working directory
	vi.AutomaticEnv()
	vi.SetEnvPrefix(service+"_PROXY")
	vi.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vi.SetDefault("proxy.service", service)
	vi.SetDefault("proxy.prefix", true)
	vi.SetDefault("proxy.tracing", true)
	vi.SetDefault("proxy.tls", false)
	vi.SetDefault("proxy.metrics_endpoint", true)
	vi.SetDefault("proxy.live_endpoint", true)
	vi.SetDefault("proxy.ready_endpoint", true)
	vi.SetDefault("proxy.pprof_endpoint", true)
	vi.SetDefault("proxy.db_host", "localhost")
	vi.SetDefault("proxy.db_port", ":5432")
	vi.SetDefault("proxy.db_name", "postgresdb")
	vi.SetDefault("proxy.db_user", "admin")
	vi.SetDefault("proxy.grpc_port", ":8443")
	vi.SetDefault("proxy.routine_threshold", 300)
	vi.SetDefault("proxy.jaeger_metrics", true)
	vi.SetDefault("proxy.backend", true)
	vi.SetDefault("proxy.log_level", true)
	vi.SetDefault("proxy.swagger_file", true)
	vi.SetDefault("proxy.allow_origin", true)
	vi.SetDefault("proxy.allow_creds", true)
	vi.SetDefault("proxy.allow_methods", true)
	vi.SetDefault("proxy.allow_headers", true)

	// If a config file is found, read it in."A generator for gRPC based Applications"
	if err := vi.ReadInConfig(); err != nil {
		log.Info("failed to read config file, writing defaults to -->"+proxyFile)
		if err := vi.WriteConfigAs(proxyFile); err != nil {
			return err
		}

	} else {
		log.Info("Using config file-->", vi.ConfigFileUsed())
		if err := vi.WriteConfig(); err != nil {
			return err
		}
	}
	return nil
}

func (v *Viperizer) GetService() string {
	return v.Service
}

func (v *Viperizer) GetProxyConfig() *viper.Viper {
	return v.ProxyConfig
}

func (v *Viperizer) GetGrpcConfig() *viper.Viper {
	return v.GrpcConfig
}

func (v *Viperizer) ReadGrpcConfig() error {
	// If a config file is found, read it in."A generator for gRPC based Applications"
	if err := v.GrpcConfig.ReadInConfig(); err != nil {
		log.Info("failed to read config file, writing defaults to -->"+serverFile)
		if err = v.WriteGrpcConfig(); err != nil {
			return err
		}
		return err
	} else {
		log.Info("Using config file-->", v.GrpcConfig.ConfigFileUsed())
		if err = v.WriteGrpcConfig(); err != nil {
			return err
		}
		return nil
	}
	return nil
}


func (v *Viperizer) ReadProxyConfig() error {
	// If a config file is found, read it in."A generator for gRPC based Applications"
	if err := v.ProxyConfig.ReadInConfig(); err != nil {
		log.Info("failed to read config file, writing defaults to -->"+proxyFile)
		if err = v.WriteProxyConfig(); err != nil {
			return err
		}
		return err
	} else {
		log.Info("Using config file-->", v.ProxyConfig.ConfigFileUsed())
		if err = v.WriteProxyConfig(); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (v *Viperizer) WriteGrpcConfig() error {
	if err := v.GrpcConfig.WriteConfigAs(serverFile); err != nil {
		return err
	} else {
		log.Info("grpc config file created-->", v.GrpcConfig.ConfigFileUsed())
		if err := v.GrpcConfig.WriteConfig(); err != nil {
			return err
		}
	}
	return nil
}

func (v *Viperizer) WriteProxyConfig() error {
	if err := v.ProxyConfig.WriteConfigAs(proxyFile); err != nil {
		return err
	} else {
		log.Info("proxy config file created-->", v.ProxyConfig.ConfigFileUsed())
		if err := v.ProxyConfig.WriteConfig(); err != nil {
			return err
		}
	}
	return nil
}