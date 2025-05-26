package app

import (
	"context"
	"testing"

	"gophernet/pkg/db/ent"
	"gophernet/pkg/mocks"

	"github.com/golang/mock/gomock"
)

func TestGetGopher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIBurrowRepository(ctrl)
	app := NewGopherApp(mockRepo, nil)

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

	mockRepo := mocks.NewMockIBurrowRepository(ctrl)
	app := NewGopherApp(mockRepo, nil)

	// Test renting an unoccupied burrow
	burrow := &ent.Burrow{ID: 1, Name: "Burrow 1", IsOccupied: false}
	mockRepo.EXPECT().GetBurrowByID(gomock.Any(), 1).Return(burrow, nil)
	mockRepo.EXPECT().UpdateBurrowOccupancy(gomock.Any(), 1, true).Return(nil)

	result, err := app.RentBurrow(context.Background(), 1)
	if err != nil {
		t.Errorf("RentBurrow() error = %v", err)
	}
	if !result.IsOccupied {
		t.Errorf("RentBurrow() burrow.IsOccupied = %v, want %v", result.IsOccupied, true)
	}

	// Test renting an already occupied burrow
	occupiedBurrow := &ent.Burrow{ID: 2, Name: "Burrow 2", IsOccupied: true}
	mockRepo.EXPECT().GetBurrowByID(gomock.Any(), 2).Return(occupiedBurrow, nil)

	_, err = app.RentBurrow(context.Background(), 2)
	if err == nil {
		t.Errorf("RentBurrow() expected error, got nil")
	}
}

func TestReleaseBurrow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIBurrowRepository(ctrl)
	app := NewGopherApp(mockRepo, nil)

	// Test releasing an occupied burrow
	burrow := &ent.Burrow{ID: 2, Name: "Burrow 2", IsOccupied: true}
	mockRepo.EXPECT().GetBurrowByID(gomock.Any(), 2).Return(burrow, nil)
	mockRepo.EXPECT().UpdateBurrowOccupancy(gomock.Any(), 2, false).Return(nil)

	result, err := app.ReleaseBurrow(context.Background(), 2)
	if err != nil {
		t.Errorf("ReleaseBurrow() error = %v", err)
	}
	if result.IsOccupied {
		t.Errorf("ReleaseBurrow() burrow.IsOccupied = %v, want %v", result.IsOccupied, false)
	}

	// Test releasing an unoccupied burrow
	unoccupiedBurrow := &ent.Burrow{ID: 1, Name: "Burrow 1", IsOccupied: false}
	mockRepo.EXPECT().GetBurrowByID(gomock.Any(), 1).Return(unoccupiedBurrow, nil)

	_, err = app.ReleaseBurrow(context.Background(), 1)
	if err == nil {
		t.Errorf("ReleaseBurrow() expected error, got nil")
	}
}

func TestGetBurrowStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIBurrowRepository(ctrl)
	app := NewGopherApp(mockRepo, nil)

	burrows := []*ent.Burrow{
		{ID: 1, Name: "Burrow 1", IsOccupied: false},
		{ID: 2, Name: "Burrow 2", IsOccupied: true},
	}
	mockRepo.EXPECT().GetAllBurrows(gomock.Any()).Return(burrows, nil)

	result, err := app.GetBurrowStatus(context.Background())
	if err != nil {
		t.Errorf("GetBurrowStatus() error = %v", err)
	}
	if len(result) != 2 {
		t.Errorf("GetBurrowStatus() len = %v, want %v", len(result), 2)
	}
}
