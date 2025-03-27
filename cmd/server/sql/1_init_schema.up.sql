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