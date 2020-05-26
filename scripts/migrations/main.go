package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/PedPet/user/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]

	settings, err := config.LoadSettings()
	if err != nil {
		log.Fatalf("goose: failed to load settings: %v\n", err)
	}

	dbSource := settings.DB.User + ":" + settings.DB.Password +
		"@tcp(" + settings.DB.Host + ")/" + settings.DB.Database + "?parseTime=true"
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("goose: failed to ping DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 2 {
		arguments = append(arguments, args[3:]...)
	}
	err = goose.Run(command, db, *dir, arguments...)
	if err != nil {
		log.Fatalf("goose: %v: %v", command, err)
	}
}
