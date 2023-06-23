package repository

import (
	"yafa/internal/model"
	"yafa/internal/repository/queries"

	"github.com/jackc/pgx"
)

type UserRepo interface {
	Create(newUser model.User) (model.User, error)
	GetByNickname(nickname string) (user model.User, err error)
	GetByEmail(email string) (user model.User, err error)
	GetByEmailOrNick(email, nickname string) (users []*model.User, err error)
	Update(user model.User) (NewUser model.User, err error)
}

type userRepo struct {
	db *pgx.ConnPool
}

func NewUserRepo(d *pgx.ConnPool) UserRepo {
	return &userRepo{db: d}
}

func (r *userRepo) Create(newUser model.User) (user model.User, err error) {
	err = r.db.QueryRow(queries.CreateUser, newUser.Nickname, newUser.Fullname, newUser.About, newUser.Email).Scan(&user.Id, &user.Nickname, &user.Fullname, &user.About, &user.Email)
	return
}

func (r *userRepo) GetByNickname(nickname string) (user model.User, err error) {
	err = r.db.QueryRow(queries.GetUserByNickname, nickname).Scan(&user.Id, &user.Nickname, &user.Fullname, &user.About, &user.Email)
	return
}

func (r *userRepo) GetByEmail(email string) (user model.User, err error) {
	err = r.db.QueryRow(queries.GetUserByEmail, email).Scan(&user.Id, &user.Nickname, &user.Fullname, &user.About, &user.Email)
	return
}

func (r *userRepo) GetByEmailOrNick(email, nickname string) (users []*model.User, err error) {
	rows, err := r.db.Query(queries.GetUserByEmailOrNick, nickname, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := &model.User{}
		err = rows.Scan(&u.Id, &u.Nickname, &u.Fullname, &u.About, &u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepo) Update(user model.User) (NewUser model.User, err error) {
	err = r.db.QueryRow(queries.UpdateUser, user.Nickname, user.Fullname, user.About, user.Email).Scan(&NewUser.Nickname, &NewUser.Fullname, &NewUser.About, &NewUser.Email)
	return
}
