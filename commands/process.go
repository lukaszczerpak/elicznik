package commands

import (
	"elicznik/common"
	log "elicznik/logging"
	"elicznik/tauron/balance"
	"elicznik/tauron/db"
	"elicznik/util"
	"github.com/spf13/cobra"
	"time"
)

var processDataCmd = &cobra.Command{
	Use:   "process",
	Short: "Processes elicznik data and calculates monthly balances",
	Long: `Command automatically detects period for which the calculation needs to be done.
It also marks balance records with "complete=false" tag to enforce recalculating them
next time the command is executed.`,
	Run: func(cmd *cobra.Command, args []string) {
		processData(&config)
	},
}

func init() {
	rootCmd.AddCommand(processDataCmd)
}

func processData(cfg *common.AppConfig) {
	db := db.New(cfg.General.Location, cfg.Influxdb.Url, cfg.Influxdb.Token, cfg.Influxdb.Org, cfg.Influxdb.Bucket, cfg.Influxdb.Measurement)
	defer db.Close()

	log.Infof("Processing data")

	tmb, err := db.GetLastCompleteBalance()
	if err != nil {
		log.Fatalf("Error when getting last complete balance: %v", err)
	}
	if tmb == nil {
		t, _ := time.ParseInLocation(util.DATE_FORMAT, cfg.Tauron.StartDate, cfg.General.Location)
		t = t.AddDate(0, -1, 0)
		tmb = balance.NewTauronMonthlyBalance(time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, cfg.General.Location), "[]")
	}

	for tmb.GetNextDate().Before(time.Now()) {
		fromGrid, feedIn, err := db.GetMonthStats(tmb.GetNextDate())
		if err != nil {
			log.Fatalf("Error when processing data: %v", err)
		}

		tmb.NextBalance(fromGrid, feedIn, cfg.Tauron.StorageFactor)

		complete, err := db.IsMonthComplete(tmb.Date)
		if err != nil {
			log.Fatalf("Error when processing data: %v", err)
		}

		log.Infof("Month %s => complete data: %v", tmb.Date.Format(util.YEAR_MONTH_FORMAT), complete)
		log.Debugf("%s: FromGridSum=%v, FeedInSum=%v", tmb.Date.Format(util.YEAR_MONTH_FORMAT), fromGrid, feedIn)
		log.Debugf("%s: EnergyPurchase=%v, StorageTotal=%v, Storage=%v, Complete=%v", tmb.Date.Format(util.YEAR_MONTH_FORMAT),
			tmb.EnergyPurchased, tmb.StorageTotal, tmb.Storage, complete)

		db.WriteBalance(tmb, complete)
	}
}
