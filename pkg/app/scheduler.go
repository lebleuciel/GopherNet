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

// Scheduler manages periodic tasks for burrow maintenance and reporting
type Scheduler struct {
	repo   repo.IBurrowRepository
	ticker *time.Ticker
	done   chan bool
	config *config.Scheduler
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
	if cfg == nil {
		cfg = &config.DefaultScheduler
	}
	return &Scheduler{
		repo:   repo,
		ticker: time.NewTicker(cfg.UpdateInterval),
		done:   make(chan bool),
		config: cfg,
	}
}

// Start begins the scheduler's periodic tasks
func (s *Scheduler) Start(ctx context.Context) {
	if err := s.initializeSystem(ctx); err != nil {
		log.Printf("Error initializing system: %v", err)
	}

	reportTicker := time.NewTicker(s.config.ReportInterval)
	go s.runPeriodicTasks(ctx, reportTicker)
	log.Println("Scheduler started")
}

// Stop gracefully shuts down the scheduler
func (s *Scheduler) Stop() {
	s.done <- true
	log.Println("Scheduler stopped")
}

// initializeSystem sets up the initial state of the system
func (s *Scheduler) initializeSystem(ctx context.Context) error {
	if err := s.repo.DeleteAllBurrows(ctx); err != nil {
		return fmt.Errorf("failed to clear database: %w", err)
	}
	return s.loadInitialBurrows(ctx)
}

// runPeriodicTasks executes the scheduled tasks
func (s *Scheduler) runPeriodicTasks(ctx context.Context, reportTicker *time.Ticker) {
	defer reportTicker.Stop()
	for {
		select {
		case <-reportTicker.C:
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
}

// calculateBurrowStats computes statistical information about the burrow system
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
	log.Printf("Starting report generation")
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

	log.Printf("Report generated: %s", filename)
	return nil
}

// handleOldBurrow removes a burrow that has exceeded its maximum age
func (s *Scheduler) handleOldBurrow(ctx context.Context, b *ent.Burrow) error {
	if err := s.repo.DeleteBurrow(ctx, int64(b.ID)); err != nil {
		return fmt.Errorf("error deleting old burrow %d: %w", b.ID, err)
	}
	log.Printf("Deleted old burrow %d (age: %d)", b.ID, b.Age)
	return nil
}

// updateBurrowDepth increases a burrow's depth
func (s *Scheduler) updateBurrowDepth(ctx context.Context, b *ent.Burrow) error {
	newDepth := b.Depth + s.config.DepthIncrement
	if err := s.repo.UpdateBurrowDepth(ctx, int64(b.ID), newDepth); err != nil {
		return fmt.Errorf("error updating burrow %d: %w", b.ID, err)
	}
	log.Printf("Updated burrow %d: depth %.2f -> %.2f", b.ID, b.Depth, newDepth)
	return nil
}

// updateBurrows processes all occupied burrows
func (s *Scheduler) updateBurrows(ctx context.Context) error {
	burrows, err := s.repo.GetOccupiedBurrows(ctx)
	if err != nil {
		return fmt.Errorf("failed to get occupied burrows: %w", err)
	}

	log.Printf("Updating %d burrows...", len(burrows))

	for _, b := range burrows {
		if b.Age >= s.config.MaxBurrowAge {
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
		log.Printf("Loading burrow: %s", burrow.Name)
		models = append(models, burrow.ParseToModel())
	}

	createdBurrows, err := s.repo.CreateBurrows(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create burrows: %w", err)
	}

	log.Printf("Loaded %d initial burrows", len(createdBurrows))
	return nil
}
