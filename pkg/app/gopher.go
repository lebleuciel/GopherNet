package app

import "gophernet/pkg/repo"

type IGopherApp interface {
	GetGopher() string
}

type GopherApp struct {
	repo repo.IBurrowRepository
}

func NewGopherApp(repo repo.IBurrowRepository) *GopherApp {
	return &GopherApp{
		repo: repo,
	}
}

func (g *GopherApp) GetGopher() string {
	return "gopher"
}
