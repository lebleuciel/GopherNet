package cache

import (
	"sync"
)

type IBurrowStats interface {
	AddDepth(depth float64)
	AddAvailableBurrow()
	GetStats() (float64, int, float64, float64)
	SetSmallestBurrow(size float64)
	SetLargestBurrow(size float64)
}

type BurrowStats struct {
	mu               sync.RWMutex
	totalDepth       float64
	availableBurrows int
	largestBurrow    float64
	smallestBurrow   float64
}

func NewBurrowStats() *BurrowStats {
	return &BurrowStats{
		smallestBurrow: -1,
	}
}

func (bs *BurrowStats) AddDepth(depth float64) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.totalDepth += depth
}

func (bs *BurrowStats) AddAvailableBurrow() {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.availableBurrows++
}

func (bs *BurrowStats) GetStats() (float64, int, float64, float64) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.totalDepth, bs.availableBurrows, bs.largestBurrow, bs.smallestBurrow
}

func (bs *BurrowStats) SetSmallestBurrow(size float64) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.smallestBurrow = size
}

func (bs *BurrowStats) SetLargestBurrow(size float64) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.largestBurrow = size
}
