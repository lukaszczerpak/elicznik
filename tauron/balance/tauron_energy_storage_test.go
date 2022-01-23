package balance

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewTauronEnergyStorage(t *testing.T) {
	storage := NewTauronEnergyStorage()
	assert.NotNil(t, storage)
	assert.Len(t, storage.storage, 0)
	assert.Equal(t, storage.len, 12)
}

var serializationTests = []struct {
	data string
	obj  TauronEnergyStorage
}{
	{"[]", TauronEnergyStorage{storage: []int{}, len: 12}},
	{"[20]", TauronEnergyStorage{storage: []int{20}, len: 12}},
	{"[20,0,0]", TauronEnergyStorage{storage: []int{20, 0, 0}, len: 12}},
	{"[20,0,0,10]", TauronEnergyStorage{storage: []int{20, 0, 0, 10}, len: 12}},
}

func TestTauronEnergyStorage_Unmarshall(t *testing.T) {
	for _, tt := range serializationTests {
		t.Run(tt.data, func(t *testing.T) {
			var got = TauronEnergyStorage{}
			if got.Unmarshal([]byte(tt.data)); !reflect.DeepEqual(got, tt.obj) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.obj)
			}
		})
	}
}

func TestTauronEnergyStorage_Marshall(t *testing.T) {
	for _, tt := range serializationTests {
		t.Run(tt.data, func(t *testing.T) {
			if got, err := tt.obj.Marshal(); err != nil || string(got) != tt.data {
				t.Errorf("Marshal() = %v, want %v", string(got), tt.obj)
			}
		})
	}
}

func TestTauronEnergyStorage_Push(t *testing.T) {
	tes := TauronEnergyStorage{storage: []int{}, len: 12}
	tes.Push(20)
	tes.Push(0)
	tes.Push(0)
	tes.Push(10)
	assert.EqualValues(t, []int{20, 0, 0, 10}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{20, 0, 0, 10}, len: 12}
	tes.Push(200)
	assert.EqualValues(t, []int{20, 0, 0, 10, 200}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{20, 0, 0, 10, 200}, len: 12}
	tes.Push(0)
	assert.EqualValues(t, []int{20, 0, 0, 10, 200, 0}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{20, 0, 0, 10, 200, 0}, len: 12}
	tes.Push(10)
	tes.Push(20)
	tes.Push(30)
	tes.Push(40)
	tes.Push(50)
	tes.Push(60)
	assert.EqualValues(t, []int{20, 0, 0, 10, 200, 0, 10, 20, 30, 40, 50, 60}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{20, 0, 0, 10, 200, 0, 10, 20, 30, 40, 50, 60}, len: 12}
	tes.Push(70)
	assert.EqualValues(t, []int{10, 200, 0, 10, 20, 30, 40, 50, 60, 70}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{10, 200, 0, 10, 20, 30, 40, 50, 60, 70}, len: 12}
	tes.Push(80)
	tes.Push(90)
	assert.EqualValues(t, []int{10, 200, 0, 10, 20, 30, 40, 50, 60, 70, 80, 90}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{10, 200, 0, 10, 20, 30, 40, 50, 60, 70, 80, 90}, len: 12}
	tes.Push(99)
	assert.EqualValues(t, []int{200, 0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 99}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{}, len: 12}
	tes.Push(0)
	assert.EqualValues(t, []int{}, tes.storage)
}

func TestTauronEnergyStorage_Pull(t *testing.T) {
	tes := TauronEnergyStorage{storage: []int{}, len: 12}
	r := tes.Pull(100)
	assert.EqualValues(t, []int{}, tes.storage)
	assert.Equal(t, 0, r)

	tes = TauronEnergyStorage{storage: []int{20, 0, 0, 50}, len: 12}
	r = tes.Pull(10)
	assert.EqualValues(t, []int{10, 0, 0, 50}, tes.storage)
	assert.Equal(t, 10, r)

	tes = TauronEnergyStorage{storage: []int{20, 0, 0, 50}, len: 12}
	r = tes.Pull(20)
	assert.EqualValues(t, []int{50}, tes.storage)
	assert.Equal(t, 20, r)

	tes = TauronEnergyStorage{storage: []int{20, 0, 0, 50}, len: 12}
	r = tes.Pull(30)
	assert.EqualValues(t, []int{40}, tes.storage)
	assert.Equal(t, 30, r)

	tes = TauronEnergyStorage{storage: []int{20, 0, 0, 50}, len: 12}
	r = tes.Pull(70)
	assert.EqualValues(t, []int{}, tes.storage)
	assert.Equal(t, 70, r)

	tes = TauronEnergyStorage{storage: []int{20, 0, 0, 50}, len: 12}
	r = tes.Pull(80)
	assert.EqualValues(t, []int{}, tes.storage)
	assert.Equal(t, 70, r)
}

func TestTauronEnergyStorage_cleanup(t *testing.T) {
	tes := TauronEnergyStorage{storage: []int{5, 0, 0, 0, 10}, len: 12}
	tes.cleanup()
	assert.EqualValues(t, []int{5, 0, 0, 0, 10}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{5, 0, 0, 0}, len: 12}
	tes.cleanup()
	assert.EqualValues(t, []int{5, 0, 0, 0}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{0, 0, 0, 10}, len: 12}
	tes.cleanup()
	assert.EqualValues(t, []int{10}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{0, 0, 0, 0, 0}, len: 12}
	tes.cleanup()
	assert.EqualValues(t, []int{}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{}, len: 12}
	tes.cleanup()
	assert.EqualValues(t, []int{}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, len: 12}
	tes.cleanup()
	assert.EqualValues(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, len: 12}
	tes.cleanup()
	assert.EqualValues(t, []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, tes.storage)

	tes = TauronEnergyStorage{storage: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}, len: 12}
	tes.cleanup()
	assert.EqualValues(t, []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}, tes.storage)
}
