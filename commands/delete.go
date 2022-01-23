package commands

import (
	"elicznik/common"
	log "elicznik/logging"
	"elicznik/tauron/db"
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

		deleteType, _ := cmd.Flags().GetString("type")
		deleteData(&config, startTime, stopTime, deleteType)
	},
}

func init() {
	deleteType := newEnum([]string{"measurement", "balance"}, "measurement")
	deleteDataCmd.Flags().VarP(deleteType, "type", "t", "type of data to delete ['measurement' or 'balance']")
	deleteDataCmd.MarkFlagRequired("type")
	rootCmd.AddCommand(deleteDataCmd)
}

func deleteData(cfg *common.AppConfig, startTime, stopTime time.Time, deleteType string) {
	db := db.New(cfg.General.Location, cfg.Influxdb.Url, cfg.Influxdb.Token, cfg.Influxdb.Org, cfg.Influxdb.Bucket, cfg.Influxdb.Measurement)
	defer db.Close()

	switch deleteType {
	case "measurement":
		// <START_DATE, STOP_DATE> range is inclusive thus STOP_DATE's time must be 23:59:59
		stopTime = time.Date(stopTime.Year(), stopTime.Month(), stopTime.Day(), 23, 59, 59, 0, stopTime.Location())
		log.Infof("Period from %v to %v => deleting measurement data",
			startTime.Format(util.DATE_FORMAT), stopTime.Format(util.DATE_FORMAT))
		//db.DeleteMeasurements(startTime, stopTime)
	case "balance":
		startTime = time.Date(startTime.Year(), startTime.Month(), 1, 0, 0, 0, 0, stopTime.Location())
		stopTime = startTime.AddDate(0, 1, 0).Add(-1 * time.Second)
		log.Infof("Period from %v to %v => deleting balance data",
			startTime.Format(util.DATE_FORMAT), stopTime.Format(util.DATE_FORMAT))
		db.DeleteBalance(startTime, stopTime)
	}
}
