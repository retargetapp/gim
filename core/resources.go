package core

import (
	"log"
	"regexp"
	"path/filepath"
	"io/ioutil"
	"os"
	"bytes"
	"errors"
	"strconv"
)

func loadSources(path string) ([]*Migration, error) {
	var err error
	var migrations = []*Migration{}
	var dirPath string
	var versions = make(map[string]bool)

	dirPath, err = filepath.Abs(path)
	if err != nil {
		log.Panicln("No migraions resources directory")
		return migrations, err
	}

	filesInfo, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Panicln("Unable to read migartions files from resources directory")
		return migrations, err
	}

	var exp, _ = regexp.Compile("^(\\d{10})\\.(up|down)\\.sql$")

	for _, f := range filesInfo {
		if f.IsDir() {
			continue
		}
		if !exp.MatchString(f.Name()) {
			continue
		}
		v := exp.FindStringSubmatch(f.Name())[1]
		versions[v] = true
	}

	log.Println(versions)

	for v, _ := range versions {
		m, err := loadMigration(dirPath, v)
		if err != nil {
			log.Println("Unable to load migartion version " + v)
		}
		migrations = append(migrations, m)
	}

	return migrations, nil
}

func loadMigration(dirPath string, version string) (*Migration, error) {
	m := &Migration{}
	versionInt, err := strconv.ParseInt(version, 10, 64)
	if err != nil {
		log.Println("Unable to parse migation version as timestamp " + version)
		return m, err
	}
	m.Version = uint64(versionInt)
	buf, err := ioutil.ReadFile(dirPath + string(os.PathSeparator) + version + ".up.sql")
	if err != nil {
		log.Println("Unable to load up sql for migation version " + version)
		return m, err
	}
	if (len(bytes.TrimSpace(buf))) <= 0 {
		log.Println("Up migaration file for version " + version + " is empty")
		return m, errors.New("empty sql file")
	}
	m.Up = string(buf)
	buf, err = ioutil.ReadFile(dirPath + string(os.PathSeparator) + version + ".down.sql")
	if err != nil {
		log.Println("Unable to load down sql for migation version " + version)
		return m, err
	}
	if (len(bytes.TrimSpace(buf))) <= 0 {
		log.Println("Down migaration file for version " + version + " is empty")
		return m, errors.New("empty sql file")
	}
	m.Down = string(buf)
	return m, nil
}
