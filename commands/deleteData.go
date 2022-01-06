package commands

import (
	"elicznik/common"
	log "elicznik/logging"
	"elicznik/taurondb"
	"elicznik/util"
	"time"

	"github.com/spf13/cobra"
)

var deleteDataCmd = &cobra.Command{
	Use:   "delete START_DATE STOP_DATE",
	Short: "Deletes measurements from database for the specified time window.",
	Long:  `Dates must be provided in a format YYYY-MM-DD, ie.: 2021-08-13`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		startTime, err := time.ParseInLocation(util.DATE_FORMAT, args[0], config.General.Location)
		if err != nil {
			log.Fatalf("Invalid start time: %v", err)
		}

		stopTime, err := time.ParseInLocation(util.DATE_FORMAT, args[1], config.General.Location)
		if err != nil {
			log.Fatalf("Invalid stop time: %v", err)
		}

		deleteData(&config, startTime, stopTime)
	},
}

func init() {
	rootCmd.AddCommand(deleteDataCmd)
}

func deleteData(cfg *common.AppConfig, startTime, stopTime time.Time) {
	// <START_DATE, STOP_DATE> range is inclusive thus STOP_DATE's time must be 23:59:59
	stopTime = time.Date(stopTime.Year(), stopTime.Month(), stopTime.Day(), 23, 59, 59, 0, stopTime.Location())

	db := taurondb.New(cfg.General.Location, cfg.Influxdb.Url, cfg.Influxdb.Token, cfg.Influxdb.Org, cfg.Influxdb.Bucket, cfg.Influxdb.Measurement)
	defer db.Close()

	log.Infof("Period from %v to %v => deleting data",
		startTime.Format(util.DATE_FORMAT), stopTime.Format(util.DATE_FORMAT))

	db.DeleteMeasurements(startTime, stopTime)
}
