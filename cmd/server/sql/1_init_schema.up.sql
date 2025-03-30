CREATE TABLE users
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(128) NOT NULL,
    password VARCHAR(128) NOT NULL,
    UNIQUE (username)
);

CREATE TABLE tokens
(
    token   VARCHAR(128) NOT NULL PRIMARY KEY,
    user_id INTEGER      NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE games
(
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    name   VARCHAR(28) NOT NULL,
    width  INTEGER     NOT NULL,
    height INTEGER     NOT NULL,
    data   BLOB        NOT NULL
);

CREATE TABLE games_token
(
    token   VARCHAR(128) NOT NULL PRIMARY KEY,
    game_id INTEGER      NOT NULL,
    FOREIGN KEY (game_id) REFERENCES games (id)
);

INSERT INTO games(name, width, height, data)
VALUES ('demo', '8', '9', '0,0,1,1,1,1,1,0,1,1,1,0,0,0,1,0,1,2,4,3,0,0,1,0,1,1,1,0,3,2,1,0,1,2,1,1,3,0,1,0,1,0,1,0,2,0,1,1,1,3,0,5,3,3,2,1,1,0,0,0,2,0,0,1,1,1,1,1,1,1,1,1');