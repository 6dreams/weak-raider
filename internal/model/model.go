package model

import "time"

type Guild struct {
	Id           int       `db:"guild_id"`
	Name         string    `db:"name"`
	Wowaudit_key string    `db:"wowaudit_key"`
	Updated      time.Time `db:"updated"`
}

type Champion struct {
	Id          int       `db:"champion_id"`
	Name        string    `db:"name"`
	Class       string    `db:"class"`
	Realm       string    `db:"realm"`
	Note        string    `db:"note"`
	Guild_id    int       `db:"guild_id"`
	Blizzard_id int       `db:"blizzard_id"`
	Updated     time.Time `db:"updated"`
}
