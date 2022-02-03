package commands

import (
	"elicznik/common"
	log "elicznik/logging"
	"elicznik/tauron/api"
	"elicznik/tauron/db"
	"elicznik/util"
	"time"

	"github.com/spf13/cobra"
)

var fetchDataCmd = &cobra.Command{
	Use:   "fetch START_DATE STOP_DATE",
	Short: "Fetches measurements from Tauron API and stores in database.",
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

		fetchData(&config, startTime, stopTime)
	},
}

func init() {
	fetchDataCmd.Flags().BoolVar(&config.General.DeleteBeforeWrite, "delete", false, "Delete old data before writing new data")
	rootCmd.AddCommand(fetchDataCmd)
}

func fetchData(cfg *common.AppConfig, startTime, stopTime time.Time) {
	yy, mm, dd := time.Now().Date()
	today := time.Date(yy, mm, dd, 0, 0, 0, 0, cfg.General.Location)

	// check start < stop
	if startTime.After(stopTime) {
		log.Fatalf("START_TIME must not be after STOP_TIME")
	}

	// check start and stop are before today midnight
	if !startTime.Before(today) || !stopTime.Before(today) {
		log.Fatalf("START_TIME and STOP_TIME must be before today midnight")
	}

	log.Infof("Fetching measurements from %s to %s", startTime.Format(util.DATE_FORMAT), stopTime.Format(util.DATE_FORMAT))
	elicznik := api.New(cfg.Tauron.Login, cfg.Tauron.Password, cfg.Tauron.SmartNr)
	tdb := db.New(cfg.General.Location, cfg.Influxdb.Url, cfg.Influxdb.Token, cfg.Influxdb.Org, cfg.Influxdb.Bucket, cfg.Influxdb.Measurement)

	for month := 0; !startTime.AddDate(0, month, 0).After(stopTime); month++ {
		rangeStartTime := startTime.AddDate(0, month, 0)
		rangeStopTime := rangeStartTime.AddDate(0, 1, -1)
		if rangeStopTime.After(stopTime) {
			rangeStopTime = stopTime
		}

		log.Infof("Batch #%d: Period from %v to %v", month,
			rangeStartTime.Format(util.DATE_FORMAT), rangeStopTime.Format(util.DATE_FORMAT))

		if !cfg.General.DeleteBeforeWrite {
			exists, err := tdb.CheckIfDataExists(rangeStartTime,
				time.Date(rangeStopTime.Year(), rangeStopTime.Month(), rangeStopTime.Day(), 23, 0, 0, 0, rangeStopTime.Location()),
				cfg.Influxdb.Measurement)
			if err != nil {
				log.Errorf("Checking DB failed: %v", err)
				continue
			}
			if exists {
				log.Infof("Period from %v to %v => data exists in db, skipping",
					rangeStartTime.Format(util.DATE_FORMAT), rangeStopTime.Format(util.DATE_FORMAT))
				continue
			}
		}

		measurements, completeData := elicznik.GetMeasurements(rangeStartTime, rangeStopTime)

		if cfg.General.DeleteBeforeWrite {
			if !completeData {
				log.Warnf("Incomplete data received, skipping this period")
				break
			}

			tdb.DeleteMeasurements(rangeStartTime, rangeStopTime)
		}

		tdb.WriteMeasurements(measurements)

		if !completeData {
			if len(measurements) == 0 {
				log.Warnf("Received no data, bailing out")
			} else {
				log.Warnf("Incomplete data received - storing period from %v to %v and bailing out",
					measurements[0].Time.Format(util.DATE_FORMAT),
					measurements[len(measurements)-1].Time.Format(util.DATE_FORMAT))
			}
			break
		}
	}
}
