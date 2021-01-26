package dao

import (
	"Week02/src/model"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type UserDao struct {
	db *sql.DB
}

func (dao *UserDao) GetById(id string) (*model.User, error) {
	dao.db = connect()
	row := dao.db.QueryRow("SELECT * FROM t_user WHERE id = ?", id)
	var uid, username string
	err := row.Scan(&uid, &username)
	if err != nil {
		if err != sql.ErrNoRows {
			err = errors.Wrap(err, "query error")
		}
		return nil, err
	}
	disconnect(dao.db)
	return &model.User{
		Id:       uid,
		Username: username,
	}, nil
}

func (dao *UserDao) Exists(err error) bool {
	return err != sql.ErrNoRows
}

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:13643566666@/basic_platform")
	if err != nil {
		panic(errors.Wrap(err, "mysql connect failed"))
	}
	return db
}

func ping(db *sql.DB) error {
	err := db.Ping()
	return errors.Wrap(err, "mysql ping failed")
}

func disconnect(db *sql.DB) {
	err := db.Close()
	if err != nil {
		panic(errors.Wrap(err, "mysql disconnect failed"))
	}
}
