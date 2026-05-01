package model

import "time"

type Guild struct {
	Id          int       `db:"guild_id"`
	Name        string    `db:"name"`
	WowauditKey string    `db:"wowaudit_key"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Character struct {
	Id         int       `db:"character_id"`
	Name       string    `db:"name"`
	Class      string    `db:"class"`
	Rank       string    `db:"rank"`
	Role       string    `db:"role"`
	Realm      string    `db:"realm"`
	Note       string    `db:"note"`
	GuildId    int       `db:"guild_id"`
	BlizzardId int       `db:"blizzard_id"`
	UpdatedAt  time.Time `db:"updated_at"`
}
