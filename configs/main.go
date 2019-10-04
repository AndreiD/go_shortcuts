package configs

import (
	"flag"
	"log"
	"os"
	"os/user"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Configuration .
type Configuration interface {
	Get(string) (string, error)
	Init(set *flag.FlagSet)
}

// ViperConfiguration .
type ViperConfiguration struct {
}

func (vc *ViperConfiguration) setDefaults() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("setDefaults: %+v\n", err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("setDefaults: %+v\n", err)
	}
	viper.SetDefault("debug", false)
	viper.SetDefault("user", usr.Name+"@"+hostname)
}

// Init .
func (vc *ViperConfiguration) Init() {

	// toggle flag here. [overwrites the config file if used!]
	flag.Bool("debug", true, "Debug mode: true or false")

	vc.setDefaults()

	// config paths
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	viper.SetConfigName("config_markets")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("Config file not found; ignore error if desired")
		} else {
			log.Fatalf("Config file error %s", err.Error())
		}
	}

	// add another config
	viper.SetConfigName("config_app")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("Config file error %s", err.Error())
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("an error occured while running viper.BindPFlags(): %+v\n", err)
	}
}

// Get .
func (vc *ViperConfiguration) Get(param string) string {
	return viper.GetString(param)
}

// GetBool .
func (vc *ViperConfiguration) GetBool(param string) bool {
	return viper.GetBool(param)
}

// NewConfiguration .
func NewConfiguration() (cfg *ViperConfiguration) {
	cfg = &ViperConfiguration{}
	return

}
