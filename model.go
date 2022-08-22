// model.go

package main

import (
	"database/sql"
	"time"
	"sort"
)

type user struct {
	ID    int     `json:"id"`
	UserName  string  `json:"username"`
	FirstName  string  `json:"firstname"`
	LastName  string  `json:"lastname"`
	Token string `json:"token"`
	Image string `json:"image"`
	DateCreated string `json:"datecreated"`
	DateModified string `json:"datemodified"`
}

type screech struct {
	ID    int     `json:"id"`
	Content  string  `json:"content"`
	Creator  string  `json:"creator"`
	DateCreated string `json:"datecreated"`
	DateModified string `json:"datemodified"`
}

func (p *user) getUser(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM users WHERE id=$1",
		p.ID)
}

func (p *user) updateUser(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE users SET username=$2, firstname=$3, lastname=$4, image=$5, datemodified=$6 WHERE id=$1",
			p.ID, p.UserName, p.FirstName, p.LastName, p.Image, time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"))

	return err
}

func (p *user) updateProfilePicture(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE users SET image=$2, datemodified=$3 WHERE id=$1",
			p.ID, p.Image, time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"))

	return err
}

func (s *screech) createScreech(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO screechs(content, creator, datecreated, datemodified) VALUES($1, $2, $3, $4) RETURNING id",
		s.Content, s.Creator, time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"), time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")).Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *screech) updateScreech(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE screechs SET content=$2, datemodified=$3 WHERE id=$1",
			s.ID, s.Content, time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"))

	return err
}

func (s *screech) getScreech(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM screechs WHERE id=$1",
		s.ID)
}

func (s *screech) getScreechesByProfile(db *sql.DB, start, count int, ascending bool) ([]screech, error){
	rows, err := db.Query(
		"SELECT * FROM screechs WHERE creator=$1 ORDER BY datecreated ASC LIMIT $2 OFFSET $3",
		s.Creator, count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	screechs := []screech{}

	for rows.Next() {
		var s screech
		if err := rows.Scan(&s.ID, &s.Content, &s.Creator, &s.DateCreated, &s.DateModified); err != nil {
			return nil, err
		}
		screechs = append(screechs, s)
	}

	// if data is needed in the descending order, reverse the order and return
	if !ascending {
		return sort.Sort(sort.Reverse(screechs)), nil 
	}

	return screechs, nil
}
