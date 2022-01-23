package balance

import (
	"elicznik/util"
	"encoding/json"
)

//type EnergyStorage interface {
//	Unmarshall(data string) error
//	Marshall() (string, error)
//	Push(wh int)
//	Pull(wh int) int
//}

type TauronEnergyStorage struct {
	storage []int
	len     int
}

func NewTauronEnergyStorage() *TauronEnergyStorage {
	storage := TauronEnergyStorage{len: 12, storage: []int{}}
	return &storage
}

func (t *TauronEnergyStorage) Unmarshal(data []byte) error {
	var arr []int
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	t.storage = arr
	t.len = 12
	return nil
}

func (t *TauronEnergyStorage) Marshal() ([]byte, error) {
	return json.Marshal(t.storage)
}

func (t *TauronEnergyStorage) Push(wh int) {
	t.storage = append(t.storage, wh)
	t.cleanup()
}

func (t *TauronEnergyStorage) Pull(wh int) int {

	if len(t.storage) == 0 || wh == 0 {
		return 0
	}

	var splitIndex int
	var sum = 0

	for i, v := range t.storage {
		splitIndex = i
		sum += v

		if sum >= wh {
			break
		}
	}

	t.storage = t.storage[splitIndex:]
	t.storage[0] = util.IntMax(0, sum-wh)

	t.cleanup()

	return util.IntMin(wh, sum)
}

func (t *TauronEnergyStorage) cleanup() {
	if len(t.storage) > t.len {
		t.storage = t.storage[len(t.storage)-t.len:]
	}

	var splitIndex = 0
	for i, v := range t.storage {
		if v != 0 {
			break
		}
		splitIndex = i + 1
	}
	t.storage = t.storage[splitIndex:]
}

func (t *TauronEnergyStorage) Total() int {
	total := 0
	for _, v := range t.storage {
		total += v
	}

	return total
}
