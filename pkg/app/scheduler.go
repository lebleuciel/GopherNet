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

	"gophernet/pkg/config"
	"gophernet/pkg/db/ent"
	"gophernet/pkg/dto"
	"gophernet/pkg/repo"
	"gophernet/pkg/utils"
)

type IScheduler interface {
	Start(ctx context.Context)
	Stop()
}

// Scheduler manages periodic tasks for burrow maintenance and reporting
type Scheduler struct {
	repo         repo.IBurrowRepository
	updateTicker *time.Ticker
	reportTicker *time.Ticker
	config       *config.Scheduler
}

// BurrowStats holds the statistical information about the burrow system
type BurrowStats struct {
	TotalDepth     float64
	LargestVolume  float64
	SmallestVolume float64
	LargestBurrow  *ent.Burrow
	SmallestBurrow *ent.Burrow
	AvailableCount int
}

// NewScheduler creates a new scheduler instance
func NewScheduler(repo repo.IBurrowRepository, cfg *config.Scheduler) *Scheduler {
	scheduler := &Scheduler{
		repo:         repo,
		updateTicker: time.NewTicker(cfg.UpdateInterval),
		reportTicker: time.NewTicker(cfg.ReportInterval),
		config:       cfg,
	}
	return scheduler
}

// Start begins the scheduler's periodic tasks
func (s *Scheduler) Start(ctx context.Context) {
	if err := s.initializeSystem(ctx); err != nil {
		log.Printf("Error initializing scheduler system: %v", err)
	}

	go s.runPeriodicTasks(ctx)

	log.Println("Scheduler started")
}

// Stop gracefully shuts down the scheduler
func (s *Scheduler) Stop() {
	s.updateTicker.Stop()
	s.reportTicker.Stop()
}

// initializeSystem initializes the system with initial burrows if none exist
func (s *Scheduler) initializeSystem(ctx context.Context) error {
	// Check if we have any existing burrows
	existingBurrows, err := s.repo.GetAllBurrows(ctx)
	if err != nil {
		return fmt.Errorf("failed to check existing burrows: %w", err)
	}

	if len(existingBurrows) > 0 {
		s.handleExistingBurrowsOnStart(ctx, existingBurrows)
	}

	return nil
}

func (s *Scheduler) runPeriodicTasks(ctx context.Context) {
	for {
		select {
		case <-s.reportTicker.C:
			if err := s.generateReport(); err != nil {
				log.Printf("Error generating report: %v", err)
			}
		case <-s.updateTicker.C:
			if err := s.updateBurrows(ctx); err != nil {
				log.Printf("Error updating burrows: %v", err)
			}
		}
	}
}

// updateBurrows processes all occupied burrows
func (s *Scheduler) updateBurrows(ctx context.Context) error {
	burrows, err := s.repo.GetOccupiedBurrows(ctx)
	if err != nil {
		return fmt.Errorf("failed to get occupied burrows: %w", err)
	}

	for _, b := range burrows {
		if err := s.UpdateBurrow(ctx, b); err != nil {
			log.Printf("Failed to update burrow %d: %v", b.ID, err)
			continue
		}
	}

	return nil
}

// calculateBurrowStats computes statistical information about the burrow system
func (s *Scheduler) calculateBurrowStats(burrows []*ent.Burrow) BurrowStats {
	if len(burrows) == 0 {
		return BurrowStats{
			SmallestVolume: 0,
			LargestVolume:  0,
			AvailableCount: 0,
		}
	}

	stats := BurrowStats{
		SmallestVolume: math.MaxFloat64,
		LargestVolume:  0,
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

// formatReport creates a formatted report string
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

// generateReport creates and saves a new report
func (s *Scheduler) generateReport() error {
	burrows, err := s.repo.GetAllBurrows(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get burrows: %w", err)
	}

	if len(burrows) == 0 {
		return fmt.Errorf("no burrows found")
	}

	stats := s.calculateBurrowStats(burrows)
	if err := s.saveReport(stats); err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	return nil
}

// saveReport writes the report to a file
func (s *Scheduler) saveReport(stats BurrowStats) error {
	if err := os.MkdirAll("reports", 0755); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	report := formatReport(stats, timestamp)
	filename := filepath.Join("reports", fmt.Sprintf("burrow_report_%s.txt", time.Now().Format("2006-01-02_15-04-05")))

	if err := os.WriteFile(filename, []byte(report), 0644); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	return nil
}

// handleOldBurrow removes a burrow that has exceeded its maximum age
func (s *Scheduler) handleOldBurrow(ctx context.Context, b *ent.Burrow) error {
	if err := s.repo.DeleteBurrow(ctx, int64(b.ID)); err != nil {
		return fmt.Errorf("error deleting old burrow %d: %w", b.ID, err)
	}
	return nil
}

// loadInitialBurrows loads the initial burrow configuration
func (s *Scheduler) loadInitialBurrows(ctx context.Context) error {
	data, err := os.ReadFile("data/initial.json")
	if err != nil {
		return fmt.Errorf("failed to read initial.json: %w", err)
	}

	var burrows []dto.BurrowDto
	if err := json.Unmarshal(data, &burrows); err != nil {
		return fmt.Errorf("failed to parse initial.json: %w", err)
	}

	var models []*ent.Burrow
	for _, burrow := range burrows {
		models = append(models, burrow.ParseToModel())
	}

	createdBurrows, err := s.repo.CreateBurrows(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create burrows: %w", err)
	}

	log.Printf("Loaded %d initial burrows", len(createdBurrows))
	return nil
}

// handleExistingBurrowsOnStart processes existing burrows when the system starts
func (s *Scheduler) handleExistingBurrowsOnStart(ctx context.Context, burrows []*ent.Burrow) error {
	for _, b := range burrows {
		// Update the burrow with new age and depth
		if err := s.UpdateBurrow(ctx, b); err != nil {
			log.Printf("Failed to update burrow %d: %v", b.ID, err)
			continue
		}
	}
	return nil
}

func (s *Scheduler) UpdateBurrow(ctx context.Context, burrow *ent.Burrow) error {
	now := time.Now()
	timePassed := now.Sub(burrow.UpdatedAt)
	minutesPassed := int(timePassed.Minutes())
	newAge := burrow.Age + minutesPassed

	// If burrow is older than 25 days, delete it
	if newAge >= s.config.MaxBurrowAge {
		if err := s.handleOldBurrow(ctx, burrow); err != nil {
			log.Printf("Failed to delete old burrow %d: %v", burrow.ID, err)
		}
		return nil
	}

	// Calculate new depth based on time passed (0.009 meters per minute)
	depthIncrease := float64(minutesPassed) * s.config.DepthIncrementRate
	newDepth := burrow.Depth + depthIncrease

	if err := s.repo.UpdateBurrow(ctx, int64(burrow.ID), newDepth, newAge); err != nil {
		return fmt.Errorf("error updating burrow %d: %w", burrow.ID, err)
	}

	return nil
}
