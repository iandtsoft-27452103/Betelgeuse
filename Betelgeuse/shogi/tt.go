package shogi

import (
	"sync"
)

var mu sync.RWMutex

func IniTT(tt *TT) {
	tt.value = make(map[uint64]int16)
	tt.color = make(map[uint64]Color)
	tt.is_check = make(map[uint64]bool)
	tt.move = make(map[uint64]uint32)
	tt.ply = make(map[uint64]uint16)
}

func Store(tt *TT) {
	mu.Lock()
	// not implemented
	defer mu.Unlock()
}
