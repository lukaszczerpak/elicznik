package commands

import (
	"elicznik/common"
	log "elicznik/logging"
	"elicznik/tauron/db"
	"elicznik/util"
	"time"

	"github.com/spf13/cobra"
)

var syncDataCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetches missing measurements from Tauron API and stores in database.",
	Long:  `If first checks last measurement in db and fetches missing data between the last measurement and today midnight`,
	Run: func(cmd *cobra.Command, args []string) {
		config.General.DeleteBeforeWrite = false
		processDataFlag, _ := cmd.Flags().GetBool("process")
		syncData(&config, processDataFlag)
	},
}

func init() {
	syncDataCmd.Flags().BoolP("process", "p", false, "Process balance data after the sync.")
	rootCmd.AddCommand(syncDataCmd)
}

func syncData(cfg *common.AppConfig, processDataFlag bool) {
	tdb := db.New(cfg.General.Location, cfg.Influxdb.Url, cfg.Influxdb.Token, cfg.Influxdb.Org, cfg.Influxdb.Bucket, cfg.Influxdb.Measurement)

	lastMeasurementDate, err := tdb.GetLastMeasurementDate()
	if err != nil {
		log.Fatalf("Error when fetching last measurement date: %v", err)
	}
	log.Debugf("Last measurement date: %s", lastMeasurementDate.Format(util.DATE_FORMAT))

	var startTime time.Time

	// database is empty
	if lastMeasurementDate.IsZero() {
		startTime, _ = time.ParseInLocation(util.DATE_FORMAT, cfg.Tauron.StartDate, cfg.General.Location)
		log.Debugf("No measurements, using start-date from the config: %s", startTime)
	} else {
		startTime = lastMeasurementDate.AddDate(0, 0, 1)
	}

	yy, mm, dd := time.Now().Date()
	today := time.Date(yy, mm, dd, 0, 0, 0, 0, cfg.General.Location)
	yesterday := today.AddDate(0, 0, -1)

	if startTime.Equal(today) {
		log.Infof("No new measurements to fetch")
		return
	}

	fetchData(cfg, startTime, yesterday)

	if processDataFlag {
		processData(cfg)
	}
}
