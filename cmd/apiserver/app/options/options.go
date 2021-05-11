package options

import (
	"fmt"
	"os"

	"github.com/likakuli/generic-project-template/pkg/config"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Option struct {
	ConfigFile string
	Version    bool
}

func NewOptions() *Option {
	return &Option{}
}

func (o *Option) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ConfigFile, "config", o.ConfigFile, "config file (default is $HOME/.config.toml)")
	fs.BoolVar(&o.Version, "version", o.Version, "display version info")
}

func (o *Option) Validate() error {
	if o.ConfigFile == "" {
		return fmt.Errorf("Must provide config file!")
	}

	return nil
}

func (o *Option) Complete() (*config.Config, error) {
	viper.SetConfigFile(o.ConfigFile)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	glog.V(1).Infof("unmarshal config file: %v", cfg)

	return &cfg, nil
}

// PrintFlags logs the flags in the flagset
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		glog.V(1).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}
