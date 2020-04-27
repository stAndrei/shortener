package main

import (
	"time"
)

// URL ...
type URL struct {
	ID        string    `storm:"unique"`
	URL       string    `storm:"unique"`
	Name      string    `storm:"index"`
	CreatedAt time.Time `storm:"index"`
	UpdatedAt time.Time `storm:"index"`
}

func GenerateID() string {
	for {
		// TODO: Make length (5) configurable
		id := RandomString(10)
		err := db.One("ID", id, nil)
		if err != nil {
			return id
		}
	}
}

func NewURL(target string) (url *URL, err error) {
	var u URL

	err = db.One("URL", target, &u)

	if err != nil {

		url = &URL{ID: GenerateID(), URL: target, CreatedAt: time.Now()}
		err = db.Save(url)
		return url, err
	} else {
		return &u, err
	}
}

// SetName ...
func (u *URL) SetName(name string) error {
	u.Name = name
	u.UpdatedAt = time.Now()
	return db.Save(&u)
}
