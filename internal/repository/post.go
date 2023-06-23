package repository

import (
	"yafa/internal/model"
	"yafa/internal/repository/queries"

	"github.com/jackc/pgx"
)

type PostRepo interface {
	Get(id int, related []string) (postInfo model.PostAll, err error)
	Update(id int, new model.Post) (p model.Post, err error)
}

type postRepo struct {
	db *pgx.ConnPool
}

func NewPostRepo(d *pgx.ConnPool) PostRepo {
	return &postRepo{db: d}
}

func (r *postRepo) Get(id int, related []string) (postInfo model.PostAll, err error) {
	var post model.Post
	err = r.db.QueryRow(queries.GetPost, id).Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)
	if err != nil {
		return
	}

	postInfo.Post = &post

	if len(related) != 0 {
		for _, q := range related {
			switch q {
			case "user":
				var u model.User
				err = r.db.QueryRow(queries.GetPostUser, post.Author).Scan(
					&u.Nickname, &u.Fullname, &u.About, &u.Email,
				)
				if err != nil {
					return
				}
				postInfo.Author = &u
			case "forum":
				var f model.Forum
				err = r.db.QueryRow(queries.GetPostForum, post.Forum).Scan(
					&f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads,
				)
				if err != nil {
					return
				}
				postInfo.Forum = &f
			case "thread":
				var t model.Thread
				err = r.db.QueryRow(queries.GetPostThread, post.Thread).Scan(
					&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created,
				)
				if err != nil {
					return
				}
				postInfo.Thread = &t
			}
		}
	}

	return
}

func (r *postRepo) Update(id int, new model.Post) (p model.Post, err error) {
	err = r.db.QueryRow(queries.UpdatePost, new.Message, id).Scan(
		&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.Thread, &p.Created,
	)
	return
}
