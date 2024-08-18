package user

import (
	"fmt"
	"github.com/504dev/logr/dbs/mysql"
	"github.com/504dev/logr/types"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	conn *sqlx.DB
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		conn: mysql.Conn(),
	}
}

func (repo *UserRepo) findAllByField(fieldname string, val interface{}, limit int) (types.Users, error) {
	users := types.Users{}
	sql := fmt.Sprintf("SELECT id, github_id, username, role, login_at, created_at FROM users WHERE %v = ?", fieldname)

	if limit > 0 {
		sql = fmt.Sprintf("%v LIMIT %v", sql, limit)
	}

	err := repo.conn.Select(&users, sql, val)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *UserRepo) findOneByField(fieldname string, val interface{}) (*types.User, error) {
	users, err := repo.findAllByField(fieldname, val, 1)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

func (repo *UserRepo) GetAll() (types.Users, error) {
	users := types.Users{}
	err := repo.conn.Select(&users, "SELECT id, github_id, username, role, login_at, created_at FROM users")
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepo) GetById(id int) (*types.User, error) {
	return repo.findOneByField("id", id)
}

func (repo *UserRepo) GetByUsername(username string) (*types.User, error) {
	return repo.findOneByField("username", username)
}

func (repo *UserRepo) GetByGithubId(id int64) (*types.User, error) {
	return repo.findOneByField("github_id", id)
}

func (repo *UserRepo) Upsert(githubId int64, username string, role int) (*types.User, error) {
	user, err := repo.GetByGithubId(githubId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user, err = repo.Create(githubId, username, role)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (repo *UserRepo) Create(githubId int64, username string, role int) (*types.User, error) {
	sql := `
		INSERT INTO users (github_id, username, role)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE role=VALUES(role)
	`
	values := []interface{}{githubId, username, role}

	_, err := repo.conn.Exec(sql, values...)
	if err != nil {
		return nil, err
	}

	user, err := repo.GetByGithubId(githubId)

	return user, err
}

func (repo *UserRepo) LoginAt(id int) error {
	_, err := repo.conn.Exec("UPDATE users SET login_at = NOW() WHERE id = ?", id)

	return err
}
