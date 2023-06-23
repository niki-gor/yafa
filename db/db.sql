CREATE EXTENSION citext;


CREATE UNLOGGED TABLE users (
    id       BIGSERIAL PRIMARY KEY,
    nickname CITEXT UNIQUE,
    fullname CITEXT,
    about    TEXT,
    email    CITEXT UNIQUE
);

CREATE UNLOGGED TABLE forum (
    id      BIGSERIAL PRIMARY KEY,
    title   TEXT,
    "user"  CITEXT,
    slug    CITEXT UNIQUE,
    posts   BIGINT DEFAULT 0,
    threads INTEGER DEFAULT 0
);

CREATE UNLOGGED TABLE forum_user (
    id     BIGSERIAL PRIMARY KEY,
    "user" BIGINT REFERENCES users (id),
    forum  BIGINT REFERENCES forum (id)  
);

CREATE UNLOGGED TABLE thread (
    id      BIGSERIAL PRIMARY KEY,
    title   TEXT,
    author  CITEXT,
    forum   CITEXT,
    message TEXT,
    votes   INTEGER DEFAULT 0,
    slug    CITEXT,
    created TIMESTAMPTZ DEFAULT now()
);

CREATE UNLOGGED TABLE post (
    id        BIGSERIAL PRIMARY KEY,
    parent    BIGINT DEFAULT 0,
    author    CITEXT,
    message   TEXT,
    is_edited BOOLEAN DEFAULT false,
    forum     CITEXT,
    thread    INT,
    created   TIMESTAMPTZ DEFAULT now(),
    path      BIGINT[] DEFAULT '{0}'
);

CREATE UNLOGGED TABLE vote (
    id     BIGSERIAL PRIMARY KEY,
    "user" BIGINT REFERENCES users (id),
    thread BIGINT REFERENCES thread (id),
    voice  INT,
    UNIQUE ("user", thread)
);


CREATE INDEX users__info_idx 
on users (nickname, fullname, about, email);

CREATE INDEX forum__slug_idx 
ON forum (slug);

CREATE INDEX forum__user_idx 
ON forum_user (forum, "user");

CREATE INDEX post__thread_idx 
ON post (thread);

CREATE INDEX post__thread_path_idx 
ON post (thread, path);

CREATE INDEX post__path_parent_idx 
ON post (thread, id, (path[1]), parent);

CREATE INDEX thread__slug_idx 
ON thread (slug);

CREATE INDEX thread__author_idx 
ON thread (author);

CREATE INDEX thread__forum_idx 
ON thread (forum);

CREATE INDEX thread__created_idx 
ON thread (created);

CREATE INDEX vote__user_thread_idx 
ON vote ("user", thread);


CREATE OR REPLACE FUNCTION thread_vote() 
RETURNS TRIGGER AS $$
    BEGIN
        UPDATE "thread"
        SET votes=(votes + new.voice)
        WHERE id = new.thread;
        RETURN new;
    end;
$$ language plpgsql;

CREATE TRIGGER vote_insert
AFTER INSERT ON vote
FOR EACH ROW EXECUTE PROCEDURE thread_vote();


CREATE OR REPLACE FUNCTION thread_vote_update() 
RETURNS TRIGGER AS $$
    BEGIN
        UPDATE "thread"
        SET "votes" = (votes + 2 * new.voice)
        WHERE "id" = new.thread;
        RETURN new;
    END;
$$ language plpgsql;

CREATE TRIGGER vote_update
AFTER UPDATE ON vote
FOR EACH ROW EXECUTE PROCEDURE thread_vote_update();


CREATE OR REPLACE FUNCTION create_post() 
RETURNS TRIGGER AS $$
    DECLARE
        _id bigint;
    BEGIN
        SELECT u.id, u.nickname, u.fullname, u.about, u.email
        FROM users u
        WHERE u.nickname = new.author
        INTO _id;

        UPDATE forum
        SET posts = posts + 1
        WHERE slug = new.forum;

        new.path = (
            SELECT path 
            FROM post 
            WHERE id = new.parent LIMIT 1
        ) || new.id;

        INSERT INTO forum_user ("user", forum)
        VALUES (_id, ( 
                SELECT id 
                FROM forum
                WHERE new.forum = slug
            )
        );

        RETURN new;
    END
$$ language plpgsql;

CREATE TRIGGER create_post
BEFORE INSERT ON post
FOR EACH ROW EXECUTE PROCEDURE create_post();


CREATE OR REPLACE FUNCTION create_thread() 
RETURNS TRIGGER AS $$
    DECLARE
        _id bigint;
    BEGIN
        SELECT u.id
        FROM users u
        WHERE u.nickname = new.author
        INTO _id;

        UPDATE forum
        SET threads = threads + 1
        WHERE slug = new.forum;
        INSERT INTO forum_user ("user", forum)
        VALUES (_id, (
                SELECT id
                FROM forum
                WHERE new.forum = slug
            )
        );
        RETURN new;
    END
$$ language plpgsql;

CREATE TRIGGER create_thread
BEFORE INSERT ON thread
FOR EACH ROW EXECUTE PROCEDURE create_thread();
