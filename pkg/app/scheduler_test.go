package app

import (
	"context"
	"testing"
	"time"

	"gophernet/pkg/config"
	"gophernet/pkg/db/ent"
	"gophernet/pkg/logger"
	"gophernet/pkg/mocks"

	"github.com/golang/mock/gomock"
)

var testConfig = &config.Scheduler{
	UpdateInterval:     time.Minute,
	ReportInterval:     time.Hour,
	MaxBurrowAge:       25 * 24 * 60, // 25 days in minutes
	DepthIncrementRate: 0.009,        // meters per minute
}

func TestBulkBorrowUpdate(t *testing.T) {
	logger.InitTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		initialBurrows []*ent.Burrow
		expectedCount  int
		checkDepth     bool
		expectedDepth  float64
		setupMock      func(*mocks.MockIBurrowRepository)
	}{
		{
			name: "should delete old burrow",
			initialBurrows: []*ent.Burrow{
				{
					ID:        1,
					Name:      "Old Burrow",
					Depth:     10.0,
					Age:       25 * 24 * 60, // 25 days in minutes
					UpdatedAt: time.Now().Add(-24 * time.Hour),
				},
			},
			expectedCount: 0,
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					DeleteBurrow(gomock.Any(), int64(1)).
					Return(nil)
			},
		},
		{
			name: "should update occupied burrow depth and age",
			initialBurrows: []*ent.Burrow{
				{
					ID:         1,
					Name:       "Occupied Burrow",
					Depth:      5.0,
					IsOccupied: true,
					Age:        0,
					UpdatedAt:  time.Now().Add(-60 * time.Minute), // 1 hour ago
				},
			},
			expectedCount: 1,
			checkDepth:    true,
			expectedDepth: 5.0 + (60 * testConfig.DepthIncrementRate), // Initial depth + (minutes * rate)
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					UpdateBurrow(gomock.Any(), int64(1), 5.0+(60*testConfig.DepthIncrementRate), 60).
					Return(nil)
			},
		},
		{
			name: "should update unoccupied burrow age only",
			initialBurrows: []*ent.Burrow{
				{
					ID:         1,
					Name:       "Unoccupied Burrow",
					Depth:      5.0,
					IsOccupied: false,
					Age:        0,
					UpdatedAt:  time.Now().Add(-60 * time.Minute), // 1 hour ago
				},
			},
			expectedCount: 1,
			checkDepth:    true,
			expectedDepth: 5.0, // Depth should remain unchanged
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					UpdateBurrow(gomock.Any(), int64(1), 5.0, 1).
					Return(nil)
			},
		},
		{
			name: "should handle mixed burrows",
			initialBurrows: []*ent.Burrow{
				{
					ID:        1,
					Name:      "Old Burrow",
					Depth:     10.0,
					Age:       25 * 24 * 60,
					UpdatedAt: time.Now().Add(-24 * time.Hour),
				},
				{
					ID:         2,
					Name:       "Occupied Burrow",
					Depth:      5.0,
					IsOccupied: true,
					Age:        0,
					UpdatedAt:  time.Now().Add(-60 * time.Minute),
				},
				{
					ID:         3,
					Name:       "Unoccupied Burrow",
					Depth:      8.0,
					IsOccupied: false,
					Age:        30,
					UpdatedAt:  time.Now().Add(-60 * time.Minute),
				},
			},
			expectedCount: 2,
			checkDepth:    true,
			expectedDepth: 5.0 + (60 * testConfig.DepthIncrementRate),
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					DeleteBurrow(gomock.Any(), int64(1)).
					Return(nil)
				mock.EXPECT().
					UpdateBurrow(gomock.Any(), int64(2), 5.0+(60*testConfig.DepthIncrementRate), 60).
					Return(nil)
				mock.EXPECT().
					UpdateBurrow(gomock.Any(), int64(3), 8.0, 31).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockIBurrowRepository(ctrl)
			tt.setupMock(mockRepo)
			scheduler := NewScheduler(mockRepo, testConfig)

			// Execute
			err := scheduler.BulkBorrowUpdate(context.Background(), tt.initialBurrows)
			if err != nil {
				t.Errorf("BulkBorrowUpdate() error = %v", err)
				return
			}
		})
	}
}

func TestUpdateBurrows(t *testing.T) {
	logger.InitTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		initialBurrows []*ent.Burrow
		expectedDepth  float64
		expectedAge    int
		setupMock      func(*mocks.MockIBurrowRepository)
	}{
		{
			name: "should update occupied burrow",
			initialBurrows: []*ent.Burrow{
				{
					ID:         1,
					Name:       "Occupied Burrow",
					Depth:      5.0,
					IsOccupied: true,
					Age:        0,
					UpdatedAt:  time.Now(),
				},
			},
			expectedDepth: 5.0,
			expectedAge:   0,
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetAllBurrows(gomock.Any()).
					Return([]*ent.Burrow{
						{
							ID:         1,
							Name:       "Occupied Burrow",
							Depth:      5.0,
							IsOccupied: true,
							Age:        0,
							UpdatedAt:  time.Now(),
						},
					}, nil)
				mock.EXPECT().
					UpdateBurrow(gomock.Any(), int64(1), 5.0, 0).
					Return(nil)
			},
		},
		{
			name: "should update unoccupied burrow age",
			initialBurrows: []*ent.Burrow{
				{
					ID:         1,
					Name:       "Unoccupied Burrow",
					Depth:      5.0,
					IsOccupied: false,
					Age:        0,
					UpdatedAt:  time.Now(),
				},
			},
			expectedDepth: 5.0,
			expectedAge:   1,
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetAllBurrows(gomock.Any()).
					Return([]*ent.Burrow{
						{
							ID:         1,
							Name:       "Unoccupied Burrow",
							Depth:      5.0,
							IsOccupied: false,
							Age:        0,
							UpdatedAt:  time.Now(),
						},
					}, nil)
				mock.EXPECT().
					UpdateBurrow(gomock.Any(), int64(1), 5.0, 1).
					Return(nil)
			},
		},
		{
			name: "should handle mixed burrows",
			initialBurrows: []*ent.Burrow{
				{
					ID:         1,
					Name:       "Occupied Burrow",
					Depth:      5.0,
					IsOccupied: true,
					Age:        0,
					UpdatedAt:  time.Now(),
				},
				{
					ID:         2,
					Name:       "Unoccupied Burrow",
					Depth:      10.0,
					IsOccupied: false,
					Age:        5,
					UpdatedAt:  time.Now(),
				},
			},
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetAllBurrows(gomock.Any()).
					Return([]*ent.Burrow{
						{
							ID:         1,
							Name:       "Occupied Burrow",
							Depth:      5.0,
							IsOccupied: true,
							Age:        0,
							UpdatedAt:  time.Now(),
						},
						{
							ID:         2,
							Name:       "Unoccupied Burrow",
							Depth:      10.0,
							IsOccupied: false,
							Age:        5,
							UpdatedAt:  time.Now(),
						},
					}, nil)
				mock.EXPECT().
					UpdateBurrow(gomock.Any(), int64(1), 5.0, 0).
					Return(nil)
				mock.EXPECT().
					UpdateBurrow(gomock.Any(), int64(2), 10.0, 6).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockIBurrowRepository(ctrl)
			tt.setupMock(mockRepo)
			scheduler := NewScheduler(mockRepo, testConfig)

			// Execute
			err := scheduler.updateBurrows(context.Background())
			if err != nil {
				t.Errorf("updateBurrows() error = %v", err)
				return
			}
		})
	}
}

func TestCalculateBurrowStats(t *testing.T) {
	logger.InitTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		burrows       []*ent.Burrow
		expectedStats BurrowStats
	}{
		{
			name:    "should handle empty burrows",
			burrows: []*ent.Burrow{},
			expectedStats: BurrowStats{
				SmallestVolume: 0,
				LargestVolume:  0,
				AvailableCount: 0,
			},
		},
		{
			name: "should calculate stats for valid burrows",
			burrows: []*ent.Burrow{
				{
					ID:         1,
					Name:       "Burrow 1",
					Depth:      5.0,
					Width:      2.0,
					IsOccupied: true,
				},
				{
					ID:         2,
					Name:       "Burrow 2",
					Depth:      10.0,
					Width:      3.0,
					IsOccupied: false,
				},
			},
			expectedStats: BurrowStats{
				TotalDepth:     15.0,
				AvailableCount: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIBurrowRepository(ctrl)
			scheduler := NewScheduler(mockRepo, testConfig)

			stats := scheduler.calculateBurrowStats(tt.burrows)

			if tt.name == "should handle empty burrows" {
				if stats.SmallestVolume != tt.expectedStats.SmallestVolume ||
					stats.LargestVolume != tt.expectedStats.LargestVolume ||
					stats.AvailableCount != tt.expectedStats.AvailableCount {
					t.Errorf("calculateBurrowStats() = %+v, want %+v", stats, tt.expectedStats)
				}
			} else {
				if stats.TotalDepth != tt.expectedStats.TotalDepth ||
					stats.AvailableCount != tt.expectedStats.AvailableCount {
					t.Errorf("calculateBurrowStats() = %+v, want %+v", stats, tt.expectedStats)
				}
				if stats.LargestBurrow == nil || stats.SmallestBurrow == nil {
					t.Error("calculateBurrowStats() should set largest and smallest burrows")
				}
			}
		})
	}
}
