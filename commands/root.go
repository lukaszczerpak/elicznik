package commands

import (
	"elicznik/common"
	log "elicznik/logging"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	cfgFile string

	config common.AppConfig

	rootCmd = &cobra.Command{
		Use:   "elicznik",
		Short: "Measurements scraper from Tauron eLicznik",
		Long:  `This program synchronizes data from Tauron with InfluxDB.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is elicznik.yaml)")
	rootCmd.PersistentFlags().BoolVar(&config.General.Debug, "debug", false, "verbose logging")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigFile("elicznik.yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		err := viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("Unable to read Viper options into configuration: %v", err)
		}

		validate := validator.New()
		// register function to get tag name from json tags.
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			if name := strings.SplitN(fld.Tag.Get("mapstructure"), ",", 2)[0]; name != "-" {
				return name
			}
			return ""
		})

		if err := validate.Struct(&config); err != nil {
			log.Fatalf("Config validation errors\n%v", err)
		}

		location, err := time.LoadLocation(config.General.Timezone)
		if err != nil {
			log.Fatalf("Unknown timezone: %v", config.General.Timezone)
		}
		config.General.Location = location

		if config.General.Debug {
			log.EnableDebug()
		}
	} else {
		log.Fatalf("Unable to read configuration: %v", err)
	}
}
