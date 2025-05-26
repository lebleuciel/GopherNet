package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"gophernet/pkg/db/ent"
	"gophernet/pkg/dto"
	"gophernet/pkg/repo"
	"gophernet/pkg/utils"
)

type Scheduler struct {
	repo   repo.IBurrowRepository
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

func NewScheduler(repo repo.IBurrowRepository) *Scheduler {
	return &Scheduler{
		repo:   repo,
		ticker: time.NewTicker(1 * time.Minute),
		done:   make(chan bool),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	//TODO: Remove later,added for testing purposes
	// Clear database
	if err := s.repo.DeleteAllBurrows(ctx); err != nil {
		log.Printf("Error clearing database: %v", err)
	}

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
				log.Printf("Generating report")
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

// BurrowStats represents the statistics of a burrow system
type BurrowStats struct {
	TotalDepth     float64
	LargestVolume  float64
	SmallestVolume float64
	LargestBurrow  *ent.Burrow
	SmallestBurrow *ent.Burrow
	AvailableCount int
}

// calculateBurrowStats calculates statistics for all burrows
func (s *Scheduler) calculateBurrowStats(burrows []*ent.Burrow) BurrowStats {
	stats := BurrowStats{
		SmallestVolume: math.MaxFloat64,
	}

	for _, burrow := range burrows {
		stats.TotalDepth += burrow.Depth
		volume := utils.CalculateVolume(burrow)

		if volume >= stats.LargestVolume {
			stats.LargestVolume = volume
			stats.LargestBurrow = burrow
		}
		if volume < stats.SmallestVolume {
			stats.SmallestVolume = volume
			stats.SmallestBurrow = burrow
		}
		if !burrow.IsOccupied {
			stats.AvailableCount++
		}
	}
	return stats
}

// formatReport formats the report with the required information
func formatReport(stats BurrowStats, timestamp string) string {
	return fmt.Sprintf(`Burrow System Report
Generated at: %s

Total Depth: %.2f meters
Available Burrows: %d
Largest Burrow: %s (Volume: %.2f cubic meters)
Smallest Burrow: %s (Volume: %.2f cubic meters)
`, timestamp, stats.TotalDepth, stats.AvailableCount,
		stats.LargestBurrow.Name, stats.LargestVolume,
		stats.SmallestBurrow.Name, stats.SmallestVolume)
}

func (s *Scheduler) generateReport() error {
	log.Printf("Starting report generation")
	burrows, err := s.repo.GetAllBurrows(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get burrows: %w", err)
	}

	if len(burrows) == 0 {
		return fmt.Errorf("no burrows found")
	}

	// Calculate statistics
	stats := s.calculateBurrowStats(burrows)

	// Create reports directory
	if err := os.MkdirAll("reports", 0755); err != nil {
		return fmt.Errorf("failed to create reports directory: %v", err)
	}

	// Generate report content
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	report := formatReport(stats, timestamp)

	// Write report to file
	filename := filepath.Join("reports", fmt.Sprintf("burrow_report_%s.txt", time.Now().Format("2006-01-02_15-04-05")))
	if err := os.WriteFile(filename, []byte(report), 0644); err != nil {
		return fmt.Errorf("failed to write report: %v", err)
	}

	log.Printf("Report generated: %s", filename)
	return nil
}

// handleOldBurrow handles the deletion of old burrows
func (s *Scheduler) handleOldBurrow(ctx context.Context, b *ent.Burrow) error {
	if err := s.repo.DeleteBurrow(ctx, int64(b.ID)); err != nil {
		return fmt.Errorf("error deleting old burrow %d: %v", b.ID, err)
	}
	log.Printf("Deleted old burrow %d (age: %d)", b.ID, b.Age)
	return nil
}

// updateBurrowDepth updates a burrow's depth
func (s *Scheduler) updateBurrowDepth(ctx context.Context, b *ent.Burrow) error {
	newDepth := b.Depth + 0.009
	if err := s.repo.UpdateBurrowDepth(ctx, int64(b.ID), newDepth); err != nil {
		return fmt.Errorf("error updating burrow %d: %v", b.ID, err)
	}
	log.Printf("Updated burrow %d: depth %.2f -> %.2f", b.ID, b.Depth, newDepth)
	return nil
}

func (s *Scheduler) updateBurrows(ctx context.Context) error {
	burrows, err := s.repo.GetOccupiedBurrows(ctx)
	if err != nil {
		return err
	}

	log.Printf("Updating %d burrows...", len(burrows))

	for _, b := range burrows {
		if b.Age >= 1440 {
			if err := s.handleOldBurrow(ctx, b); err != nil {
				log.Printf("%v", err)
				continue
			}
			continue
		}
		if err := s.updateBurrowDepth(ctx, b); err != nil {
			log.Printf("%v", err)
			continue
		}
	}

	return nil
}

func (s *Scheduler) loadInitialBurrows(ctx context.Context) error {
	// Read the initial.json file
	data, err := os.ReadFile("data/initial.json")
	if err != nil {
		return fmt.Errorf("failed to read initial.json: %v", err)
	}

	// Unmarshal the JSON data into a slice of BurrowDto
	var burrows []dto.BurrowDto
	if err := json.Unmarshal(data, &burrows); err != nil {
		return fmt.Errorf("failed to parse initial.json: %v", err)
	}

	// Convert DTOs to ent.Burrow models
	var models []*ent.Burrow
	for _, burrow := range burrows {
		log.Printf("Loading burrow: %s", burrow.Name)
		models = append(models, burrow.ParseToModel())
	}

	// Bulk create burrows
	createdBurrows, err := s.repo.CreateBurrows(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create burrows: %v", err)
	}

	log.Printf("Loaded %d initial burrows", len(createdBurrows))
	return nil
}
