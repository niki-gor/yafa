package repository

import (
	"log"
	"strconv"
	"time"

	"yafa/internal/model"
	"yafa/internal/repository/queries"

	"github.com/jackc/pgx"
)

type ThreadRepo interface {
	GetBySlugOrId(slug string) (t model.Thread, err error)
	Create(new model.Thread) (t model.Thread, err error)
	Update(old model.Thread, new model.Thread) (t model.Thread, err error)
	CheckPost(parent int64, id int) (err error)
	CreatePosts(thread model.Thread, new model.Posts) (response *model.Posts, err error)
	CheckVotes(user int, thread int) (vote model.Vote, err error)
	CreateVote(userId int, vote model.VoteRequest, thread model.Thread) (err error)
	UpdateVote(vote model.VoteRequest, voteId int) (id int, err error)
	GetThreadPosts(thread model.Thread, since, sort string, limit int, desc bool) ([]model.Post, error)
}

type threadRepo struct {
	db *pgx.ConnPool
}

func NewThreadRepo(d *pgx.ConnPool) ThreadRepo {
	return &threadRepo{db: d}
}

func (r *threadRepo) Create(new model.Thread) (t model.Thread, err error) {
	if new.Created.String() == "" {
		new.Created = time.Now()
	}

	err = r.db.QueryRow(queries.CreateThread, new.Title, new.Author, new.Forum, new.Message, new.Slug, new.Created).Scan(
		&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
	return
}

func (r *threadRepo) Update(oldThread model.Thread, newThread model.Thread) (t model.Thread, err error) {
	err = r.db.QueryRow(queries.UpdateThread, newThread.Title, newThread.Message, oldThread.Id).Scan(&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
	return
}

func (r *threadRepo) GetBySlugOrId(slug string) (t model.Thread, err error) {
	id, _ := strconv.Atoi(slug)
	err = r.db.QueryRow(queries.GetThread, slug, id).Scan(&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
	return
}

func (r *threadRepo) CheckPost(parent int64, id int) (err error) {
	err = r.db.QueryRow(queries.CheckThreadPost, id, parent).Scan(&id)
	return
}

func (r *threadRepo) CreatePosts(thread model.Thread, new model.Posts) (response *model.Posts, err error) {
	values := make([]any, 0, 6*len(new.Posts))
	created := time.Now()
	for _, post := range new.Posts {
		values = append(values, post.Parent, post.Author, post.Message, thread.Forum, thread.Id, created)
	}

	rows, err := r.db.Query(queries.CreateThreadPosts(len(new.Posts)), values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response = &model.Posts{}

	for rows.Next() {
		var p model.Post
		err = rows.Scan(
			&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.Thread, &p.Created,
		)
		if err != nil {
			return nil, err
		}

		response.Posts = append(response.Posts, p)
	}

	return response, nil
}

func (r *threadRepo) CheckVotes(user, thread int) (vote model.Vote, err error) {
	err = r.db.QueryRow(queries.CheckVotes, user, thread).Scan(&vote.Id, &vote.User, &vote.Thread, &vote.Voice)
	return
}

func (r *threadRepo) CreateVote(userId int, vote model.VoteRequest, thread model.Thread) (err error) {
	err = r.db.QueryRow(queries.CreateVote, userId, thread.Id, vote.Voice).Scan(&userId)
	return
}

func (r *threadRepo) UpdateVote(vote model.VoteRequest, voteId int) (id int, err error) {
	err = r.db.QueryRow(queries.UpdateVote, vote.Voice, voteId).Scan(&id)
	return
}

func (r *threadRepo) GetThreadPosts(thread model.Thread, since string, sort string, limit int, desc bool) ([]model.Post, error) {
	posts := make([]model.Post, 0)

	query, err := queries.Render(queries.GetThreadPostsTemplate, queries.Params{
		Desc:  desc,
		Since: since,
		Limit: limit,
		Sort:  sort,
	})
	log.Print(query)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, thread.Id)
	if err != nil {
		return []model.Post{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Post
		err := rows.Scan(&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.Thread, &p.Created)
		if err != nil {
			return []model.Post{}, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}
