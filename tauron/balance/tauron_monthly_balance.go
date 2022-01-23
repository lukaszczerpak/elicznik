package balance

import (
	log "elicznik/logging"
	"elicznik/util"
	"math"
	"time"
)

type TauronMonthlyBalance struct {
	Date            time.Time
	Storage         TauronEnergyStorage
	StorageTotal    int
	EnergyPurchased int
}

func NewTauronMonthlyBalance(date time.Time, storage string) *TauronMonthlyBalance {
	s := TauronEnergyStorage{}
	err := s.Unmarshal([]byte(storage))
	if err != nil {
		log.Fatalf("Error on parsing old balance record: %v", err)
	}

	return &TauronMonthlyBalance{Date: date, Storage: s}
}

func (tmb *TauronMonthlyBalance) NextBalance(fromGrid, feedIn, sf float64) {
	// method 1: compensate with storage first, then with current month feedIn then push the remainder into storage
	//balance := int(math.Round(feedIn*tmb.sf - (fromGrid - float64(s.Pull(int(fromGrid))))))
	//s.Push(util.IntMax(balance, 0))

	// method 2: compensate with current month feedIn, then the remainder with storage, afterwards push the remainder into storage
	monthBalance := int(math.Round(feedIn*sf - fromGrid))
	b := tmb.Storage.Pull(util.IntMax(-monthBalance, 0)) + monthBalance
	tmb.Storage.Push(util.IntMax(b, 0))

	tmb.StorageTotal = tmb.Storage.Total()
	tmb.EnergyPurchased = util.IntMax(-b, 0)
	tmb.Date = tmb.GetNextDate()
}

func (tmb *TauronMonthlyBalance) GetNextDate() time.Time {
	return tmb.Date.AddDate(0, 1, 0)
}
