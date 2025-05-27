package app

import (
	"context"
	"errors"
	"testing"

	"gophernet/pkg/db/ent"
	"gophernet/pkg/mocks"

	"github.com/golang/mock/gomock"
)

func TestGetGopher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIBurrowRepository(ctrl)
	app := NewGopherApp(mockRepo)

	msg, err := app.GetGopher(context.Background())
	if err != nil {
		t.Errorf("GetGopher() error = %v", err)
	}
	if msg != "Gopher is ready to help!" {
		t.Errorf("GetGopher() = %v, want %v", msg, "Gopher is ready to help!")
	}
}

func TestRentBurrow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		burrowID      int
		initialBurrow *ent.Burrow
		expectedError error
		setupMock     func(*mocks.MockIBurrowRepository)
	}{
		{
			name:     "should rent unoccupied burrow",
			burrowID: 1,
			initialBurrow: &ent.Burrow{
				ID:         1,
				Name:       "Burrow 1",
				Depth:      5.0,
				Width:      2.0,
				IsOccupied: false,
				Age:        0,
			},
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetBurrowByID(gomock.Any(), 1).
					Return(&ent.Burrow{
						ID:         1,
						Name:       "Burrow 1",
						Depth:      5.0,
						Width:      2.0,
						IsOccupied: false,
						Age:        0,
					}, nil)
				mock.EXPECT().
					UpdateBurrowOccupancy(gomock.Any(), 1, true).
					Return(nil)
			},
		},
		{
			name:     "should fail when burrow is already occupied",
			burrowID: 2,
			initialBurrow: &ent.Burrow{
				ID:         2,
				Name:       "Burrow 2",
				Depth:      5.0,
				Width:      2.0,
				IsOccupied: true,
				Age:        0,
			},
			expectedError: errors.New("burrow is already occupied"),
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetBurrowByID(gomock.Any(), 2).
					Return(&ent.Burrow{
						ID:         2,
						Name:       "Burrow 2",
						Depth:      5.0,
						Width:      2.0,
						IsOccupied: true,
						Age:        0,
					}, nil)
			},
		},
		{
			name:          "should fail when burrow not found",
			burrowID:      3,
			expectedError: errors.New("failed to get burrow"),
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetBurrowByID(gomock.Any(), 3).
					Return(nil, errors.New("burrow not found"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIBurrowRepository(ctrl)
			tt.setupMock(mockRepo)
			app := NewGopherApp(mockRepo)

			result, err := app.RentBurrow(context.Background(), tt.burrowID)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("RentBurrow() expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("RentBurrow() error = %v, want %v", err, tt.expectedError)
				}
				return
			}

			if err != nil {
				t.Errorf("RentBurrow() unexpected error = %v", err)
				return
			}

			if !result.IsOccupied {
				t.Errorf("RentBurrow() burrow.IsOccupied = %v, want %v", result.IsOccupied, true)
			}

			// Verify all burrow fields are preserved
			if result.ID != tt.initialBurrow.ID ||
				result.Name != tt.initialBurrow.Name ||
				result.Depth != tt.initialBurrow.Depth ||
				result.Width != tt.initialBurrow.Width ||
				result.Age != tt.initialBurrow.Age {
				t.Errorf("RentBurrow() burrow = %+v, want %+v", result, tt.initialBurrow)
			}
		})
	}
}

func TestReleaseBurrow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		burrowID      int
		initialBurrow *ent.Burrow
		expectedError error
		setupMock     func(*mocks.MockIBurrowRepository)
	}{
		{
			name:     "should release occupied burrow",
			burrowID: 2,
			initialBurrow: &ent.Burrow{
				ID:         2,
				Name:       "Burrow 2",
				Depth:      5.0,
				Width:      2.0,
				IsOccupied: true,
				Age:        0,
			},
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetBurrowByID(gomock.Any(), 2).
					Return(&ent.Burrow{
						ID:         2,
						Name:       "Burrow 2",
						Depth:      5.0,
						Width:      2.0,
						IsOccupied: true,
						Age:        0,
					}, nil)
				mock.EXPECT().
					UpdateBurrowOccupancy(gomock.Any(), 2, false).
					Return(nil)
			},
		},
		{
			name:     "should fail when burrow is not occupied",
			burrowID: 1,
			initialBurrow: &ent.Burrow{
				ID:         1,
				Name:       "Burrow 1",
				Depth:      5.0,
				Width:      2.0,
				IsOccupied: false,
				Age:        0,
			},
			expectedError: errors.New("burrow is not occupied"),
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetBurrowByID(gomock.Any(), 1).
					Return(&ent.Burrow{
						ID:         1,
						Name:       "Burrow 1",
						Depth:      5.0,
						Width:      2.0,
						IsOccupied: false,
						Age:        0,
					}, nil)
			},
		},
		{
			name:          "should fail when burrow not found",
			burrowID:      3,
			expectedError: errors.New("failed to get burrow"),
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetBurrowByID(gomock.Any(), 3).
					Return(nil, errors.New("burrow not found"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIBurrowRepository(ctrl)
			tt.setupMock(mockRepo)
			app := NewGopherApp(mockRepo)

			result, err := app.ReleaseBurrow(context.Background(), tt.burrowID)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("ReleaseBurrow() expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("ReleaseBurrow() error = %v, want %v", err, tt.expectedError)
				}
				return
			}

			if err != nil {
				t.Errorf("ReleaseBurrow() unexpected error = %v", err)
				return
			}

			if result.IsOccupied {
				t.Errorf("ReleaseBurrow() burrow.IsOccupied = %v, want %v", result.IsOccupied, false)
			}

			// Verify all burrow fields are preserved
			if result.ID != tt.initialBurrow.ID ||
				result.Name != tt.initialBurrow.Name ||
				result.Depth != tt.initialBurrow.Depth ||
				result.Width != tt.initialBurrow.Width ||
				result.Age != tt.initialBurrow.Age {
				t.Errorf("ReleaseBurrow() burrow = %+v, want %+v", result, tt.initialBurrow)
			}
		})
	}
}

func TestGetBurrowStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		burrows       []*ent.Burrow
		expectedError error
		setupMock     func(*mocks.MockIBurrowRepository)
	}{
		{
			name: "should return all burrows",
			burrows: []*ent.Burrow{
				{
					ID:         1,
					Name:       "Burrow 1",
					Depth:      5.0,
					Width:      2.0,
					IsOccupied: false,
					Age:        0,
				},
				{
					ID:         2,
					Name:       "Burrow 2",
					Depth:      10.0,
					Width:      3.0,
					IsOccupied: true,
					Age:        0,
				},
			},
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetAllBurrows(gomock.Any()).
					Return([]*ent.Burrow{
						{
							ID:         1,
							Name:       "Burrow 1",
							Depth:      5.0,
							Width:      2.0,
							IsOccupied: false,
							Age:        0,
						},
						{
							ID:         2,
							Name:       "Burrow 2",
							Depth:      10.0,
							Width:      3.0,
							IsOccupied: true,
							Age:        0,
						},
					}, nil)
			},
		},
		{
			name:    "should handle empty burrow list",
			burrows: []*ent.Burrow{},
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetAllBurrows(gomock.Any()).
					Return([]*ent.Burrow{}, nil)
			},
		},
		{
			name:          "should handle repository error",
			expectedError: errors.New("failed to get burrows"),
			setupMock: func(mock *mocks.MockIBurrowRepository) {
				mock.EXPECT().
					GetAllBurrows(gomock.Any()).
					Return(nil, errors.New("database error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIBurrowRepository(ctrl)
			tt.setupMock(mockRepo)
			app := NewGopherApp(mockRepo)

			result, err := app.GetBurrowStatus(context.Background())

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("GetBurrowStatus() expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("GetBurrowStatus() error = %v, want %v", err, tt.expectedError)
				}
				return
			}

			if err != nil {
				t.Errorf("GetBurrowStatus() unexpected error = %v", err)
				return
			}

			if len(result) != len(tt.burrows) {
				t.Errorf("GetBurrowStatus() len = %v, want %v", len(result), len(tt.burrows))
				return
			}

			// Verify all burrow fields are preserved
			for i, burrow := range result {
				expected := tt.burrows[i]
				if burrow.ID != expected.ID ||
					burrow.Name != expected.Name ||
					burrow.Depth != expected.Depth ||
					burrow.Width != expected.Width ||
					burrow.IsOccupied != expected.IsOccupied ||
					burrow.Age != expected.Age {
					t.Errorf("GetBurrowStatus() burrow[%d] = %+v, want %+v", i, burrow, expected)
				}
			}
		})
	}
}
