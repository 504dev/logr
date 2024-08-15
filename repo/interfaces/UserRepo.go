package interfaces

import (
	"github.com/504dev/logr/types"
)

type UserRepo interface {
	GetAll() (types.Users, error)
	GetById(id int) (*types.User, error)
	GetByUsername(username string) (*types.User, error)
	Upsert(githubId int64, username string, role int) (*types.User, error)
	Create(githubId int64, username string, role int) (*types.User, error)
	LoginAt(id int) error
}
