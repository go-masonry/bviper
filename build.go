package bviper

import (
	"container/list"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/spf13/viper"
	"strings"
)

const (
	EnvKeyDelimiterFrom = "_"
	EnvKeyDelimiterTo   = "."
)

type viperConfig struct {
	envDelimiterFrom, envDelimiterTo string
	mainConfigFilePath               string
	extraConfigFilePaths             []string
}

type defaultViperBuilder struct {
	ll *list.List
}

// Builder creates a simple viper instance with a given config file path
// configFilePath should point to a file that have one of these extensions:
//
// "json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "dotenv", "env", "ini"
func Builder() cfg.Builder {
	return &defaultViperBuilder{
		ll: list.New(),
	}
}

func (vb *defaultViperBuilder) SetConfigFile(path string) cfg.Builder {
	vb.ll.PushBack(func(cfg *viperConfig) {
		cfg.mainConfigFilePath = path
	})
	return vb
}

func (vb *defaultViperBuilder) AddExtraConfigFile(path string) cfg.Builder {
	vb.ll.PushBack(func(cfg *viperConfig) {
		cfg.extraConfigFilePaths = append(cfg.extraConfigFilePaths, path)
	})
	return vb
}

func (vb *defaultViperBuilder) SetEnvDelimiterReplacer(from, to string) cfg.Builder {
	vb.ll.PushBack(func(cfg *viperConfig) {
		if len(from) > 0 && len(to) > 0 {
			cfg.envDelimiterFrom = from
			cfg.envDelimiterTo = to
		}
	})
	return vb
}

func (vb *defaultViperBuilder) Build() (cfg.Config, error) {
	// apply options
	viperCfg := &viperConfig{
		envDelimiterFrom: EnvKeyDelimiterFrom,
		envDelimiterTo:   EnvKeyDelimiterTo,
	}
	for e := vb.ll.Front(); e != nil; e = e.Next() {
		f := e.Value.(func(cfg *viperConfig))
		f(viperCfg)
	}
	// build instance with multiple files
	viperWithOptions := viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(viperCfg.envDelimiterTo, viperCfg.envDelimiterFrom)), // yeah not the right order...
	)
	viperWithOptions.AutomaticEnv()
	// TODO consider having an ENV variable with File path
	if len(viperCfg.mainConfigFilePath) > 0 {
		files := append([]string{viperCfg.mainConfigFilePath}, viperCfg.extraConfigFilePaths...)
		for _, file := range files {
			viperWithOptions.SetConfigFile(file)
			if err := viperWithOptions.MergeInConfig(); err != nil {
				return nil, err
			}
		}
	}
	return &viperWrapper{
		instance: viperWithOptions,
	}, nil
}

// CustomBuilder allows one to use an already initialized and configure Viper instance
//
// Note:
//	It's important to receive a viper instance that is already initialized using viper.ReadInConfig() or viper.ReadConfig(in io.Reader) or viper.ReadRemoteConfig()
// Also note that all other Builder function will do nothing since we assume this instance is configured properly
func CustomBuilder(custom *viper.Viper) cfg.Builder {
	return &customViperBuilder{
		viperInstance: custom,
	}
}

type customViperBuilder struct {
	viperInstance *viper.Viper
}

func (cvb *customViperBuilder) SetConfigFile(path string) cfg.Builder {
	return cvb
}

func (cvb *customViperBuilder) AddExtraConfigFile(path string) cfg.Builder {
	return cvb
}

func (cvb *customViperBuilder) SetEnvDelimiterReplacer(from, to string) cfg.Builder {
	return cvb
}

func (cvb *customViperBuilder) Build() (cfg.Config, error) {
	return &viperWrapper{
		instance: cvb.viperInstance,
	}, nil
}

var _ cfg.Builder = (*customViperBuilder)(nil)
var _ cfg.Builder = (*defaultViperBuilder)(nil)
