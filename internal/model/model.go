package model

import "time"

type Guild struct {
	Id          int       `db:"guild_id"`
	Name        string    `db:"name"`
	WowauditKey string    `db:"wowaudit_key"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Character struct {
	Id                int       `db:"character_id"`
	Name              string    `db:"name"`
	Class             string    `db:"class"`
	Rank              string    `db:"rank"`
	Role              string    `db:"role"`
	Realm             string    `db:"realm"`
	Note              string    `db:"note"`
	GuildId           int       `db:"guild_id"`
	UpdatedAt         time.Time `db:"updated_at"`
	WowauditUpdatedAt time.Time `db:"wowaudit_updated_at"`
}

type Season struct {
	Id         int    `db:"id"`
	BlizzardId *int   `db:"blizzard_id"`
	WowAuditId *int   `db:"wowaudit_id"`
	Name       string `db:"name"`
}

type Instance struct {
	Id        int       `db:"id"` // blizzard_id
	SeasonId  int       `db:"season_id"`
	Name      string    `db:"name"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Encounter struct {
	Id              int       `db:"id"`
	Name            string    `db:"name"`
	LogsEncounterId *int      `db:"logs_encounter_id"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type Slot struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

type Item struct {
	Id          int    `db:"id"` // ign id
	Name        string `db:"name"`
	InstanceId  *int   `db:"instance_id"`
	EncounterId *int   `db:"encounter_id"`
	IsUnique    bool   `db:"is_unique"`
	Slot        int    `db:"slot_id"`
	Sockets     int    `db:"sockets"`
}
