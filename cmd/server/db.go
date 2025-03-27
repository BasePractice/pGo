package main

import (
	"crypto/sha3"
	"database/sql"
	"embed"
	"encoding/hex"
	"errors"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed sql/*.sql
var migrations embed.FS

type tokenManager interface {
	Token(token string) (*token, error)
	Access(username, password string) (*token, error)
}

type sqliteTokenManager struct {
	db *sql.DB
}

func (tokenManager *sqliteTokenManager) token(id int, username string) (*token, error) {
	rows, err := tokenManager.db.Query("SELECT t.token FROM tokens t WHERE t.user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tok string
	if rows.Next() {
		err = rows.Scan(&tok)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return newToken(tok, username), nil
}

func (tokenManager *sqliteTokenManager) Token(token string) (*token, error) {
	rows, err := tokenManager.db.Query("SELECT u.username FROM tokens t LEFT JOIN users u ON u.id = t.user_id WHERE token = ?", token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var username string
	if rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return newToken(token, username), nil
}

func (tokenManager *sqliteTokenManager) Access(username, password string) (*token, error) {
	rows, err := tokenManager.db.Query("SELECT u.password, u.id FROM users u WHERE u.username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bytes := sha3.Sum256([]byte(password))
	sha := hex.EncodeToString(bytes[:])
	if rows.Next() {
		var pass string
		var id int
		err = rows.Scan(&pass, &id)
		if err != nil {
			return nil, err
		}
		if sha != pass {
			return nil, errors.New("invalid password")
		}
		tok, err := tokenManager.token(id, username)
		if err != nil {
			return nil, err
		}
		return tok, nil
	} else {
		tx, err := tokenManager.db.Begin()
		if err != nil {
			return nil, err
		}
		result, err := tx.Exec("INSERT INTO users(username, password) VALUES(?, ?) RETURNING id", username, sha)
		if err != nil {
			return nil, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		tokenText := uuid.New().String()
		tokenText = strings.ReplaceAll(tokenText, "-", "")
		tokenText = strings.ToUpper(tokenText)
		result, err = tx.Exec("INSERT INTO tokens(user_id, token) VALUES(?, ?)", id, tokenText)
		if err != nil {
			return nil, err
		}
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
		return newToken(tokenText, username), nil
	}
}

func (tokenManager *sqliteTokenManager) Close() error {
	return tokenManager.db.Close()
}

func initScheme(db *sql.DB) {
	d, err := iofs.New(migrations, "sql")
	if err != nil {
		log.Fatal(err)
	}
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}
	instance, err := migrate.NewWithInstance("iofs", d, "tm", driver)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = instance.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
		return
	}
}

func newTokenManager() tokenManager {
	db, err := sql.Open("sqlite3", "tm.db3")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	initScheme(db)
	return &sqliteTokenManager{db}
}
