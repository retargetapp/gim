package core

// TODO: Refactor version type to int32
type Migration struct {
	Version uint64 `db:"version"`
	Up      string `db:"up"`
	Down    string `db:"down"`
}
