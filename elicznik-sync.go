package main

import (
	log "elicznik-sync/logging"
	"elicznik-sync/tauron"
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

	config Config

	rootCmd = &cobra.Command{
		Use:   "elicznik-sync",
		Short: "Measurements scraper from Tauron eLicznik",
		Long: `This program synchronizes data from Tauron with InfluxDB.
It automatically detects last sync and pulls only necessary
amount of data.`,
		Run: executeRun,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is elicznik-sync.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigFile("elicznik-sync.yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		err := viper.Unmarshal(&config)
		if err != nil {
			log.Fatal("Unable to read Viper options into configuration: %v", err)
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
			log.Fatal("Config validation errors\n%v", err)
		}
	} else {
		log.Fatal("Unable to read configuration: %v", err)
	}
}

func executeRun(ccmd *cobra.Command, args []string) {
	tdb := NewTauronDB(config.Influxdb.Url, config.Influxdb.Token, config.Influxdb.Org, config.Influxdb.Bucket)
	dateStart, err := tdb.GetLastMeasurementDate()
	if err != nil {
		log.Fatal("Error when fetching last measurement date: %v", err)
	}
	log.Debug("Last measurement date: %s", dateStart)

	// validate DateStart's timezone to not overlap
	if !dateStart.IsZero() && dateStart.Hour() != 0 {
		log.Fatal("Timezone mismatch detected (system vs. database) - ensure you set TZ correctly")
	}

	// database is empty
	if dateStart.IsZero() {
		dateStart, _ = time.ParseInLocation(tauron.ELICZNIK_DATE_FORMAT, config.Tauron.StartDate, time.Now().Location())
		log.Debug("No measurements, using start-date from the config: %s", dateStart)
	}
	yy, mm, dd := time.Now().Date()
	today := time.Date(yy, mm, dd, 0, 0, 0, 0, time.Now().Location())
	yesterday := today.AddDate(0, 0, -1)

	if dateStart.Equal(today) {
		log.Info("No new measurements to fetch")
		os.Exit(0)
	}

	log.Info("Fetching measurements from %s to %s", dateStart, today)
	elicznik := tauron.NewELicznik(config.Tauron.Login, config.Tauron.Password, config.Tauron.SmartNr)
	for month := 0; dateStart.AddDate(0, month, 0).Before(today); month++ {
		start := dateStart.AddDate(0, month, 0)
		end := start.AddDate(0, 1, -1)
		if end.After(yesterday) {
			end = yesterday
		}
		log.Info("Batch #%d: from %s to %s", month, start, end)
		measurements, incomplete, err := elicznik.GetData(start, end)
		if err != nil {
			log.Fatal("Error when fetching data from Tauron: %v", err)
		}
		tdb.WriteMeasurements(measurements)

		if incomplete {
			log.Warn("Incomplete data received, bailing out")
		}
	}
}

func main() {
	rootCmd.Execute()
}
