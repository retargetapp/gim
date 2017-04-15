package main

type Migration struct {
	Version uint64		`db:"version"`
	Up      string		`db:"up"`
	Down    string		`db:"down"`
	Done    bool		`db:"done"`
}
