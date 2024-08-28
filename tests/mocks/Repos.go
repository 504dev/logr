package mocks

import (
	"github.com/504dev/logr/repo"
	"sync"
)

var once sync.Once
var repos *repo.Repos

func GetRepos() *repo.Repos {
	once.Do(func() {
		repos = &repo.Repos{
			User:            nil,
			Dashboard:       nil,
			DashboardMember: nil,
			DashboardKey:    nil,
			Count:           nil,
			Log:             nil,
		}
	})
	return repos
}
