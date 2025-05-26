package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gophernet/pkg/cache"
	"gophernet/pkg/repo"
)

type Scheduler struct {
	repo   repo.IBurrowRepository
	stats  cache.IBurrowStats
	ticker *time.Ticker
	done   chan bool
}

type InitialBurrow struct {
	Name       string  `json:"name"`
	Depth      float64 `json:"depth"`
	Width      float64 `json:"width"`
	IsOccupied bool    `json:"occupied"`
	Age        int     `json:"age"`
}

func NewScheduler(repo repo.IBurrowRepository, stats cache.IBurrowStats) *Scheduler {
	return &Scheduler{
		repo:   repo,
		stats:  stats,
		ticker: time.NewTicker(1 * time.Minute),
		done:   make(chan bool),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	// Load initial burrows
	if err := s.loadInitialBurrows(ctx); err != nil {
		log.Printf("Error loading initial burrows: %v", err)
	}

	reportTicker := time.NewTicker(2 * time.Minute)

	go func() {
		defer reportTicker.Stop()
		for {
			select {
			case <-reportTicker.C:
				log.Printf("generating report")
				if err := s.generateReport(); err != nil {
					log.Printf("Error generating report: %v", err)
				}
			case <-s.ticker.C:
				if err := s.updateBurrows(ctx); err != nil {
					log.Printf("Error updating burrows: %v", err)
				}
			case <-s.done:
				s.ticker.Stop()
				return
			}
		}
	}()
	log.Println("Scheduler started")
}

func (s *Scheduler) Stop() {
	s.done <- true
	log.Println("Scheduler stopped")
}

func (s *Scheduler) generateReport() error {
	log.Printf("starting report")
	totalDepth, availableBurrows, largest, smallest := s.stats.GetStats()
	log.Printf("starting report", totalDepth, availableBurrows, largest, smallest)

	// Ensure reports directory exists
	if err := os.MkdirAll("reports", 0755); err != nil {
		log.Printf("Failed to create reports directory: %v", err)
		return fmt.Errorf("failed to create reports directory: %v", err)
	}

	filename := filepath.Join("reports", fmt.Sprintf("burrow_report_%s.txt", time.Now().Format("2006-01-02_15-04-05")))

	report := fmt.Sprintf(`Burrow System Report
Generated at: %s

Total Depth: %.2f meters
Available Burrows: %d
Largest Burrow Depth: %.2f meters
Smallest Burrow Depth: %.2f meters
`, time.Now().Format("2006-01-02 15:04:05"),
		totalDepth,
		availableBurrows,
		largest,
		smallest)

	log.Printf("Writing report to %s...", filename)
	if err := os.WriteFile(filename, []byte(report), 0644); err != nil {
		log.Printf("Failed to write report: %v", err)
		return fmt.Errorf("failed to write report: %v", err)
	}
	log.Printf("Report generated: %s", filename)
	return nil
}

func (s *Scheduler) updateBurrows(ctx context.Context) error {
	// Get all occupied burrows through repository
	burrows, err := s.repo.GetOccupiedBurrows(ctx)
	if err != nil {
		return err
	}

	log.Printf("Updating %d burrows...", len(burrows))

	// Update each burrow
	for _, b := range burrows {
		// Check if burrow is 25 days old
		if b.Age >= 90000 {
			_, _, largest, smallest := s.stats.GetStats()
			if err := s.repo.DeleteBurrow(ctx, int64(b.ID)); err != nil {
				log.Printf("Error deleting old burrow %d: %v", b.ID, err)
				continue
			}
			log.Printf("Deleted old burrow %d (age: %d)", b.ID, b.Age)
			if b.Depth == smallest {
				s.updateSmallestBurrow(ctx)
			}
			if b.Depth == largest {
				s.updateLargestBurrow(ctx)
			}
			continue
		}

		// Increment depth by a fixed amount (0.9 cm per minute)
		newDepth := b.Depth + 0.009
		if err := s.repo.UpdateBurrowDepth(ctx, int64(b.ID), newDepth); err != nil {
			log.Printf("Error updating burrow %d: %v", b.ID, err)
			continue
		}

		log.Printf("Updated burrow %d: depth %.2f -> %.2f", b.ID, b.Depth, newDepth)

		// Update stats cache
		s.stats.AddDepth(0.009)
		if !b.IsOccupied {
			s.stats.AddAvailableBurrow()
		}
		s.stats.AddDepth(newDepth)
	}

	return nil
}

func (s *Scheduler) updateSmallestBurrow(ctx context.Context) {
	remainingBurrows, err := s.repo.GetOccupiedBurrows(ctx)
	if err != nil {
		log.Printf("Error getting remaining burrows: %v", err)
		return
	}
	if len(remainingBurrows) > 0 {
		newSmallest := remainingBurrows[0].Depth
		for _, rb := range remainingBurrows[1:] {
			if rb.Depth < newSmallest {
				newSmallest = rb.Depth
			}
		}
		s.stats.SetSmallestBurrow(newSmallest)
	} else {
		s.stats.SetSmallestBurrow(0)
	}
}

func (s *Scheduler) updateLargestBurrow(ctx context.Context) {
	remainingBurrows, err := s.repo.GetOccupiedBurrows(ctx)
	if err != nil {
		log.Printf("Error getting remaining burrows: %v", err)
		return
	}
	if len(remainingBurrows) > 0 {
		newLargest := remainingBurrows[0].Depth
		for _, rb := range remainingBurrows[1:] {
			if rb.Depth > newLargest {
				newLargest = rb.Depth
			}
		}
		s.stats.SetLargestBurrow(newLargest)
	} else {
		s.stats.SetLargestBurrow(0)
	}
}

func (s *Scheduler) loadInitialBurrows(ctx context.Context) error {
	// Read the initial.json file
	data, err := os.ReadFile("data/initial.json")
	if err != nil {
		return fmt.Errorf("failed to read initial.json: %v", err)
	}

	// Parse the JSON data
	var burrows []InitialBurrow
	if err := json.Unmarshal(data, &burrows); err != nil {
		return fmt.Errorf("failed to parse initial.json: %v", err)
	}

	// Create each burrow in the database
	for _, b := range burrows {
		// Create the burrow
		if _, err := s.repo.CreateBurrow(ctx, b.Name, b.Depth, b.Width, b.IsOccupied, b.Age); err != nil {
			log.Printf("Error creating burrow: %v", err)
			continue
		}

		// Update stats
		s.stats.AddDepth(b.Depth)
		if !b.IsOccupied {
			s.stats.AddAvailableBurrow()
		}
		s.stats.SetSmallestBurrow(b.Depth)
		s.stats.SetLargestBurrow(b.Depth)
	}

	log.Printf("Loaded %d initial burrows", len(burrows))
	return nil
}
