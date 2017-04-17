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
	var p string

	p, err = filepath.Abs(path)
	if err != nil {
		return ml, errors.New(ERROR_INVALID_SRC_DIRECTORY)
	}

	vs, err := loadSrcVersions(p)
	if err != nil {
		return ml, err
	}

	for _, v := range vs {
		m, err := LoadSrcMigration(p, v)
		if err != nil {
			return ml, err
		}
		ml[m.Version] = m
	}

	return ml, nil
}

func LoadSrcMigration(path string, version string) (*Migration, error) {
	m := &Migration{}
	vInt, err := strconv.ParseInt(version, 10, 64)
	if err != nil {
		return m, NewSrcFileError(ERROR_INVALID_VERSION_FORMAT, version, "", err)
	}

	m.Version = vInt
	buf, err := ioutil.ReadFile(path + string(os.PathSeparator) + version + ".up.sql")
	if err != nil {
		return m, NewSrcFileError(ERROR_UNABLE_TO_OPEN_SRC_FILE, version, "up", err)
	}
	if (len(bytes.TrimSpace(buf))) <= 0 {
		return m, NewSrcFileError(ERROR_EMPTY_SRC_FILE, version, "up", err)
	}
	m.Up = string(buf)
	buf, err = ioutil.ReadFile(path + string(os.PathSeparator) + version + ".down.sql")
	if err != nil {
		return m, NewSrcFileError(ERROR_UNABLE_TO_OPEN_SRC_FILE, version, "down", err)
	}
	if (len(bytes.TrimSpace(buf))) <= 0 {
		return m, NewSrcFileError(ERROR_EMPTY_SRC_FILE, version, "down", err)
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

func loadSrcVersions(path string) ([]string, error) {
	var vs []string
	fl, err := ioutil.ReadDir(path)
	if err != nil {
		return vs, errors.New(ERROR_INVALID_SRC_DIRECTORY)
	}

	var exp, _ = regexp.Compile("^(\\d{10})\\.(up|down)\\.sql$")
	var vm = make(map[string]struct{})
	for _, f := range fl {
		if f.IsDir() {
			continue
		}
		if !exp.MatchString(f.Name()) {
			continue
		}
		v := exp.FindStringSubmatch(f.Name())[1]
		vm[v] = struct{}{}
	}

	for v, _ := range vm {
		vs = append(vs, v)
	}

	return vs, nil
}
