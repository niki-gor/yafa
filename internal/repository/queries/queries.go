package queries

import (
	"bytes"
	"text/template"
)

const (
	CreateForum    = `INSERT INTO forum (title, "user", slug) values ($1, $2, $3) RETURNING title, "user", slug, posts, threads;`
	GetForumBySlug = `SELECT id, title, "user", slug, posts, threads FROM forum WHERE slug = $1;`
)

const (
	GetPost       = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM post WHERE id = $1;`
	GetPostUser   = `SELECT nickname, fullname, about, email FROM users WHERE nickname = $1;`
	GetPostForum  = `SELECT title, "user", slug, posts, threads FROM forum WHERE slug = $1;`
	GetPostThread = `SELECT id, title, author, forum, message, votes, slug, created FROM thread WHERE id = $1;`
	UpdatePost    = `UPDATE post SET message = $1, is_edited = TRUE WHERE id = $2 RETURNING id, parent, author, message, is_edited, forum, thread, created;`
)

const (
	GetDBInfo = `SELECT (SELECT count(*) from forum), (SELECT count(*) from post), (SELECT count(*) from thread), (SELECT count(*) from users);`
	DeleteDB  = `TRUNCATE users, forum, thread, post, vote, forum_user CASCADE;`
)

const (
	CreateThread      = `INSERT INTO thread (title, author, forum, message, slug, created) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, title, author, forum, message, votes, slug, created;`
	UpdateThread      = `UPDATE thread SET title = $1, message = $2 WHERE id = $3 RETURNING id, title, author, forum, message, votes, slug, created;`
	GetThread         = `SELECT id, title, author, forum, message, votes, slug, created FROM thread WHERE slug = $1 OR id = $2;`
	CheckThreadPost   = `SELECT id FROM post WHERE thread = $1 AND id = $2;`
	CreateThreadPosts = `INSERT INTO post (parent, author, message, forum, thread, created) values `
	CheckVotes        = `SELECT id, "user", thread, voice from vote where "user" = $1 and thread = $2;`
	CreateVote        = `INSERT INTO vote ("user", thread, voice)  VALUES ($1, $2, $3)  RETURNING "user";`
	UpdateVote        = `UPDATE vote SET voice = $1 WHERE id = $2 RETURNING id;`
)

const (
	CreateUser           = `INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4) RETURNING id, nickname, fullname, about, email;`
	GetUserByNickname    = `SELECT id, nickname, fullname, about, email FROM users WHERE nickname = $1;`
	GetUserByEmail       = `SELECT id, nickname, fullname, about, email  FROM users WHERE email = $1;`
	GetUserByEmailOrNick = `SELECT id, nickname, fullname, about, email FROM users WHERE nickname = $1 OR email = $2;`
	UpdateUser           = `UPDATE users SET fullname = $2, about = $3, email = $4 WHERE nickname = $1 RETURNING nickname, fullname, about, email;`
)

type Params struct {
	Desc  bool
	Since string
	Limit int
	Sort  string
}

var out = bytes.NewBuffer(nil)

func Render(t *template.Template, v any) (string, error) {
	out.Reset()
	err := t.Execute(out, v)
	return out.String(), err
}

var GetForumUsersTemplate = template.Must(template.New("").Parse(`
{{$less := ""}}{{$desc := ""}}
{{if .Desc}}{{$less = "<"}}{{$desc = "DESC"}}{{else}}{{$less = ">"}}{{$desc = "ASC"}}{{end}}

SELECT nickname, about, email, fullname
FROM users
WHERE id IN (SELECT "user" FROM forum_user WHERE forum = $1)
{{if .Since}} AND nickname {{$less}} '{{.Since}}' {{end}}
ORDER BY "nickname" {{$desc}}
{{if .Limit}} LIMIT {{.Limit}} {{end}};
`))

var GetForumThreadsTemplate = template.Must(template.New("").Parse(`
{{$less := ""}}{{$desc := ""}}
{{if .Desc}}{{$less = "<"}}{{$desc = "DESC"}}{{else}}{{$less = ">"}}{{$desc = "ASC"}}{{end}}

SELECT id, title, author, forum, message, votes, slug, created 
FROM thread WHERE forum = $1
{{if .Since}} AND created {{$less}}= '{{.Since}}' {{end}}
ORDER BY created {{$desc}}
{{if .Limit}} LIMIT {{.Limit}} {{end}};
`))

var GetThreadPostsTemplate = template.Must(template.New("").Parse(`
{{$less := ""}}{{$desc := ""}}
{{if .Desc}}{{$less = "<"}}{{$desc = "DESC"}}{{else}}{{$less = ">"}}{{$desc = "ASC"}}{{end}}

SELECT 
	id, parent, author, message, is_edited, forum, thread, created 
FROM 
	post 
WHERE 
	thread = $1

{{if eq .Sort "flat"}} 
{{if .Since}} 
	AND id {{$less}} '{{.Since}}' 
{{end}}
ORDER BY created {{$desc}}, id {{$desc}}
LIMIT {{.Limit}}

{{else if eq .Sort "tree"}}
{{if .Since}} 
	AND path {{$less}} (SELECT path FROM post WHERE id = {{.Since}}) 
{{end}}
ORDER BY path[1] {{$desc}}, path {{$desc}}
LIMIT {{.Limit}}

{{else if eq .Sort "parent_tree"}}
AND path && (
	SELECT ARRAY (
		SELECT id FROM post WHERE thread = $1 AND parent = 0
		{{if .Since}} 
			AND path {{$less}} (SELECT path[1:1] FROM post WHERE id = {{.Since}}) 
		{{end}}
		ORDER BY path[1] {{$desc}}, path
		LIMIT {{.Limit}}
	)
)
ORDER BY path[1] {{$desc}}, path
{{end}};
`))

var CreateThreadPostsTemplate = template.Must(template.New("").Parse(`
INSERT INTO post (parent, author, message, forum, thread, created) VALUES 
{{.}}
RETURNING id, parent, author, message, is_edited, forum, thread, created;
`))
