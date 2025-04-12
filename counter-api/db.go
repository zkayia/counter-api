package main

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ostafen/clover/v2"
)

func dbCreateCollection(db *clover.DB, name string) error {
	if res, err := db.HasCollection(name); err != nil {
		return err
	} else if !res {
		if err := db.CreateCollection(name); err != nil {
			return err
		}
	}
	return nil
}

func dbExecuteOperation(counter string, operation func(db *clover.DB) (any, error)) (any, error) {

	db, err := clover.Open(*dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err := dbCreateCollection(db, counter); err != nil {
		return nil, err
	}

	return operation(db)
}

func dbCreateBackup() error {

	if _, err := os.Stat(*dbPath); err != nil || errors.Is(err, os.ErrNotExist) {
		log.Println("[WARN] Database doesn't exists/no access; Skipping backup")
		return nil
	}

	archive, err := os.Create(filepath.Join(
		*backupsPath,
		time.Now().Format("2006-01-02_15-04-05")+".zip",
	))
	if err != nil {
		return err
	}
	defer archive.Close()

	writer := zip.NewWriter(archive)
	defer writer.Close()

	return filepath.Walk(*dbPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		var zippablePath string
		if filepath.IsAbs(path) {
			relPath, _ := strings.CutPrefix(path, filepath.Dir(*dbPath))
			zippablePath = filepath.Join(".", relPath)
		} else {
			zippablePath = path
		}

		f, err := writer.Create(zippablePath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	})
}
