package core

type Migration struct {
	Version int64  `db:"version"`
	Up      string `db:"up"`
	Down    string `db:"down"`
}
