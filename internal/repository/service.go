package repository

import (
	"yafa/internal/model"
	"yafa/internal/repository/queries"

	"github.com/jackc/pgx"
)

type ServiceRepo interface {
	Status() (model.Status, error)
	Clear() error
}

type serviceRepo struct {
	db *pgx.ConnPool
}

func NewServiceRepo(d *pgx.ConnPool) ServiceRepo {
	return &serviceRepo{db: d}
}

func (r *serviceRepo) Status() (model.Status, error) {
	var status model.Status
	err := r.db.QueryRow(queries.GetDBInfo).Scan(
		&status.Forum, &status.Post, &status.Thread, &status.User,
	)
	return status, err
}

func (r *serviceRepo) Clear() error {
	_, err := r.db.Exec(queries.DeleteDB)
	return err
}
