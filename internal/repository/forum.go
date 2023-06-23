package repository

import (
	"yafa/internal/model"
	"yafa/internal/repository/queries"

	"github.com/jackc/pgx"
)

type ForumRepo interface {
	Create(new model.Forum) (model.Forum, error)
	GetBySlug(slug string) (forum model.Forum, err error)
	GetThreads(slug, since string, limit int, desc bool) ([]model.Thread, error)
	GetUsers(forum model.Forum, since string, limit int, desc bool) ([]model.User, error)
}

type forumRepo struct {
	db *pgx.ConnPool
}

func NewForumRepo(d *pgx.ConnPool) ForumRepo {
	return &forumRepo{db: d}
}

func (r *forumRepo) Create(new model.Forum) (forum model.Forum, err error) {
	err = r.db.QueryRow(queries.CreateForum, new.Title, new.User, new.Slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	return
}

func (r *forumRepo) GetBySlug(slug string) (forum model.Forum, err error) {
	err = r.db.QueryRow(queries.GetForumBySlug, slug).Scan(&forum.Id, &forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	return
}

func (r *forumRepo) GetThreads(slug string, since string, limit int, desc bool) ([]model.Thread, error) {
	query, err := queries.Render(queries.GetForumThreadsTemplate, queries.Params{
		Desc:  desc,
		Since: since,
		Limit: limit,
	})

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, slug)
	if err != nil {
		return []model.Thread{}, err
	}
	defer rows.Close()

	threads := make([]model.Thread, 0)
	for rows.Next() {
		var t model.Thread
		err = rows.Scan(
			&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created,
		)
		if err != nil {
			return []model.Thread{}, err
		}

		threads = append(threads, t)
	}
	return threads, nil
}

func (r *forumRepo) GetUsers(forum model.Forum, since string, limit int, desc bool) ([]model.User, error) {
	query, err := queries.Render(queries.GetForumUsersTemplate, queries.Params{
		Desc:  desc,
		Since: since,
		Limit: limit,
	})

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, forum.Id)
	if err != nil {
		return []model.User{}, err
	}
	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {
		var u model.User
		err := rows.Scan(
			&u.Nickname,
			&u.About,
			&u.Email,
			&u.Fullname,
		)
		if err != nil {
			return []model.User{}, err
		}

		users = append(users, u)
	}
	return users, nil
}
