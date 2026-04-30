package model

import "time"

type guild struct {
	id           int       `db:"guild_id"`
	name         string    `db:"name"`
	wowaudit_key string    `db:"wowaudit_key"`
	updated      time.Time `db:"updated"`
}

type champion struct {
	id          int       `db:"champion_id"`
	name        string    `db:"name"`
	class       string    `db:"class"`
	realm       string    `db:"realm"`
	note        string    `db:"note"`
	guild_id    int       `db:"guild_id"`
	blizzard_id int       `db:"blizzard_id"`
	updated     time.Time `db:"updated"`
}
