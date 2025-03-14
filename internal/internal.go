package internal

//go:generate sqlboiler --no-hooks psql

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var (
	PQ   *sql.DB
	SQLX *sqlx.DB

	Session *discordgo.Session
	Bot     *discordgo.User
)

const GuildConfigSchema = `
CREATE TABLE IF NOT EXISTS core_config (
	guild_id BIGINT PRIMARY KEY,
	guild_prefix TEXT
)
`

func Init() error {
	Session, err := discordgo.New(ConfigDgoBotToken())
	if err != nil {
		log.WithError(err).Fatal()
	}

	db := "asbwig"
	if ConfigPGDB != "" {
		db = ConfigPGDB
	}
	host := "localhost"
	if ConfigPGHost != "" {
		host = ConfigPGHost
	}

	err = postgresConnect(db, host, ConfigPGUsername, ConfigPGPassword)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}

	log.Infof("Initializing DB schema")
	initDB(GuildConfigSchema)

	run(Session)

	return err
}

func run(s *discordgo.Session) {
	s.Open()
	Session = s
	Bot = s.State.User
	log.Infoln("Bot is now running. Press CTRL-C to exit.")
}

func postgresConnect(database string, host string, username string, password string) error {
	if host == "" {
		host = "localhost"
	}

	if password != "" {
		password = " password='" + password + "'"
	}

	// Initialise database
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable%s", host, username, database, password))
	PQ = db
	SQLX = sqlx.NewDb(PQ, "postgres")
	boil.SetDB(PQ)
	return err
}

func initDB(schema string) {
	// Wait in case initial database creation. Bit starts too quickly
	_, err := PQ.Exec(schema)
	if err != nil {
		log.WithError(err).Fatal("Failed initializing postgres db schema")
	}
	log.Infoln("Database initialized")
}