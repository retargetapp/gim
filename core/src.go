package core

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	ERROR_INVALID_SRC_DIRECTORY = "invalid_sources_directory"
)

func LoadSrcMigrations(path string) (map[int64]*Migration, error) {
	var err error
	var ml = make(map[int64]*Migration)

	vs, err := LoadSrcVersions(path)
	if err != nil {
		return ml, err
	}

	for fn, v := range vs {
		m, err := LoadSrcMigration(path, v, fn)
		if err != nil {
			return ml, err
		}
		ml[m.Version] = m
	}

	return ml, nil
}

func LoadSrcMigration(path string, fName string, version int64) (*Migration, error) {
	m := &Migration{}
	m.Version = version

	p, err := filepath.Abs(path)
	if err != nil {
		return m, errors.New(ERROR_INVALID_SRC_DIRECTORY)
	}

	buf, err := ioutil.ReadFile(p + string(os.PathSeparator) + fName + ".up.sql")
	if err != nil {
		return m, NewSrcFileError(ERROR_UNABLE_TO_OPEN_SRC_FILE, fName, "up", err)
	}
	if (len(bytes.TrimSpace(buf))) <= 0 {
		return m, NewSrcFileError(ERROR_EMPTY_SRC_FILE, fName, "up", err)
	}
	m.Up = string(buf)
	buf, err = ioutil.ReadFile(p + string(os.PathSeparator) + fName + ".down.sql")
	if err != nil {
		return m, NewSrcFileError(ERROR_UNABLE_TO_OPEN_SRC_FILE, fName, "down", err)
	}
	if (len(bytes.TrimSpace(buf))) <= 0 {
		return m, NewSrcFileError(ERROR_EMPTY_SRC_FILE, fName, "down", err)
	}
	m.Down = string(buf)
	return m, nil
}

func CreateSrcVersionTpl(path string, v string) error {
	f, err := os.OpenFile(path+string(os.PathSeparator)+v+".up.sql", os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	f, err = os.OpenFile(path+string(os.PathSeparator)+v+".down.sql", os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func LoadSrcVersions(path string) (map[int64]string, error) {
	var vs map[int64]string

	p, err := filepath.Abs(path)
	if err != nil {
		return vs, errors.New(ERROR_INVALID_SRC_DIRECTORY)
	}

	fl, err := ioutil.ReadDir(p)
	if err != nil {
		return vs, errors.New(ERROR_INVALID_SRC_DIRECTORY)
	}

	vs = make(map[int64]string)

	var exp, _ = regexp.Compile("^((\\d{1,10})(_.*)?)\\.(up|down)\\.sql$")
	for _, f := range fl {
		if f.IsDir() {
			continue
		}
		if !exp.MatchString(f.Name()) {
			continue
		}
		p := exp.FindStringSubmatch(f.Name())
		vInt, err := strconv.ParseInt(p[2], 10, 64)
		if err != nil {
			return vs, NewSrcFileError(ERROR_INVALID_VERSION_FORMAT, p[2], p[4], err)
		}
		vs[vInt] = p[1]
	}

	return vs, nil
}
