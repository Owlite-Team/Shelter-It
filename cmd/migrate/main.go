package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func runMigration(m *migrate.Migrate, cmd string) error {
	var err error
	switch cmd {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "force":
		// Check if a version is provided for force
		if len(os.Args) < 4 {
			return errors.New("force command requires a version number")
		}
		version, err := parseVersion(os.Args[len(os.Args)-2])
		//version, _, err := m.Version()
		if err != nil {
			return err
		}
		err = m.Force(int(version))
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				log.Println("No migration has been applied yet.")
				return nil
			}
			return err
		}
		log.Printf("Current migration version: %d, Dirty: %t", version, dirty)
		return nil
	default:
		return errors.New("invalid command")
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func parseVersion(versionStr string) (uint, error) {
	var version uint
	_, err := fmt.Sscan(versionStr, &version)
	if err != nil {
		return 0, fmt.Errorf("invalid version format: %w", err)
	}
	return version, nil
}

func main() {
	m, err := migrate.New(
		"file://cmd/migrate/migrations",
		"postgres://novita:root@localhost:5432/shelter_it?sslmode=disable")
	if err != nil {
		log.Fatal("Error initialize migration: ", err.Error())
	}
	cmd := os.Args[len(os.Args)-1]
	if err := runMigration(m, cmd); err != nil {
		log.Fatalf("Error migrate %s: %s", cmd, err.Error())
	} else {
		log.Printf("Success migrate %s", cmd)
	}
}
