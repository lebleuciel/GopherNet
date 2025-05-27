package app

import (
	"context"
	"testing"
	"time"

	"gophernet/pkg/config"
	"gophernet/pkg/db/ent"
	"gophernet/pkg/mocks"

	"github.com/golang/mock/gomock"
)

func TestHandleExistingBurrowsOnStart(t *testing.T) {
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
					GetAllBurrows(gomock.Any()).
					Return([]*ent.Burrow{}, nil)
				mock.EXPECT().
					DeleteBurrow(gomock.Any(), int64(1)).
					Return(nil)
			},
		},
		{
			name: "should update burrow depth and age",
			initialBurrows: []*ent.Burrow{
				{
					ID:        1,
					Name:      "Young Burrow",
					Depth:     5.0,
					Age:       0,
					UpdatedAt: time.Now().Add(-60 * time.Minute), // 1 hour ago
				},
			},
			expectedCount: 1,
			checkDepth:    true,
			expectedDepth: 5.0 + (60 * 0.09), // Initial depth + (minutes * rate)
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetAllBurrows(gomock.Any()).
					Return([]*ent.Burrow{
						{
							ID:        1,
							Name:      "Young Burrow",
							Depth:     5.0 + (60 * 0.09),
							Age:       60,
							UpdatedAt: time.Now(),
						},
					}, nil)
				mock.EXPECT().
					UpdateBurrowDepth(gomock.Any(), int64(1), 5.0+(60*0.09)).
					Return(nil)
			},
		},
		{
			name: "should handle multiple burrows",
			initialBurrows: []*ent.Burrow{
				{
					ID:        1,
					Name:      "Old Burrow",
					Depth:     10.0,
					Age:       25 * 24 * 60,
					UpdatedAt: time.Now().Add(-24 * time.Hour),
				},
				{
					ID:        2,
					Name:      "Young Burrow",
					Depth:     5.0,
					Age:       0,
					UpdatedAt: time.Now().Add(-60 * time.Minute),
				},
			},
			expectedCount: 1,
			checkDepth:    true,
			expectedDepth: 5.0 + (60 * 0.09),
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetAllBurrows(gomock.Any()).
					Return([]*ent.Burrow{
						{
							ID:        2,
							Name:      "Young Burrow",
							Depth:     5.0 + (60 * 0.09),
							Age:       60,
							UpdatedAt: time.Now(),
						},
					}, nil)
				mock.EXPECT().
					DeleteBurrow(gomock.Any(), int64(1)).
					Return(nil)
				mock.EXPECT().
					UpdateBurrowDepth(gomock.Any(), int64(2), 5.0+(60*0.09)).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockIBurrowRepository(ctrl)
			tt.setupMock(mockRepo)
			scheduler := NewScheduler(mockRepo, &config.DefaultScheduler)

			// Execute
			err := scheduler.handleExistingBurrowsOnStart(context.Background(), tt.initialBurrows)
			if err != nil {
				t.Errorf("handleExistingBurrowsOnStart() error = %v", err)
				return
			}
		})
	}
}

func TestUpdateBurrows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		initialBurrows []*ent.Burrow
		expectedDepth  float64
		setupMock      func(*mocks.MockIBurrowRepository)
	}{
		{
			name: "should update occupied burrow depth",
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
			expectedDepth: 5.0 + config.DefaultScheduler.DepthIncrement,
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetOccupiedBurrows(gomock.Any()).
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
					UpdateBurrowDepth(gomock.Any(), int64(1), 5.0+config.DefaultScheduler.DepthIncrement).
					Return(nil)
			},
		},
		{
			name: "should not update unoccupied burrow",
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
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetOccupiedBurrows(gomock.Any()).
					Return([]*ent.Burrow{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockIBurrowRepository(ctrl)
			tt.setupMock(mockRepo)
			scheduler := NewScheduler(mockRepo, &config.DefaultScheduler)

			// Execute
			err := scheduler.updateBurrows(context.Background())
			if err != nil {
				t.Errorf("updateBurrows() error = %v", err)
				return
			}
		})
	}
}
