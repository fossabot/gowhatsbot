package main

import (
	"database/sql"
	"log"

	"github.com/mamur-rezeki/gowhatsplugins/errors"
	"github.com/mamur-rezeki/gowhatsplugins/helper"
	"github.com/mamur-rezeki/gowhatsplugins/logs"
	"github.com/mamur-rezeki/gowhatsplugins/types"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type logger struct {
	database *sql.DB
}

type iLogger interface {
	store(tags, content string) error
}

func newILogger(db *sql.DB) (iLogger, error) {
	var thelogger logger
	if db == nil {
		return nil, errors.NewPlugin("Logger", "Database nil")
	} else {
		thelogger.database = db
		thelogger.init()
		return &thelogger, nil
	}
}

func (i *logger) init() error {

	if ctx, err := i.database.Begin(); err != nil {
		return err
	} else {
		if _, err := ctx.Exec(`CREATE TABLE if not exists logs(
			"id"	INTEGER,
			"tags"	varchar(255) NOT NULL,
			"content"	text,
			PRIMARY KEY("id" AUTOINCREMENT)
	);
		CREATE TABLE if not exists sent(
			"id"	INTEGER,
			"jid"	varchar(255) NOT NULL,
			"message_id"	varchar(255) NOT NULL,
			"message"	text,
			PRIMARY KEY("id" AUTOINCREMENT)
	);`); err != nil {
			return err
		}
		ctx.Commit()

		return err
	}
}

func (i *logger) store(tags, content string) error {
	if ctx, err := i.database.Begin(); err != nil {
		return err
	} else {
		defer ctx.Commit()

		if _, err := ctx.Exec(`insert into logs(tags, content) values($1, $2)`, tags, content); err != nil {
			return err
		} else {
			return err
		}
	}
}

var mlogger iLogger

var Plugin = types.Plugin{
	Name:     "Logger",
	Validate: storeLog,
}

var Database *sql.DB

func init() {

	if db, err := sql.Open("sqlite3", "file:messages.db?_foreign_keys=on"); err == nil {
		Database = db
	} else {
		log.Println(logs.CodeFilename(1), err)
		panic(err)
	}

	if n, err := newILogger(Database); err != nil {
		log.Print(logs.CodeFilename(1), err)
	} else {
		mlogger = n
	}

}

func onlyLog(event interface{}, _ *whatsmeow.Client) (bool, error) {
	var the_event = helper.GetType(event)

	switch ee := event.(type) {
	case *events.Message:
		log.Println(the_event, ee.Info.Sender.ToNonAD(), "on", ee.Info.Chat.ToNonAD())
	case *events.Receipt:

	default:
		log.Println(the_event)
	}

	return false, nil
}

func storeLog(event interface{}, _ *whatsmeow.Client) (bool, error) {
	var the_event = helper.GetType(event)

	switch ee := event.(type) {
	case *events.Message:
		log.Println(the_event, ee.Info.Sender.ToNonAD(), "on", ee.Info.Chat.ToNonAD())
	case *events.Receipt:

	default:
		log.Println(the_event)
	}

	if j, err := helper.JsonMe(event); err == nil {
		if err := mlogger.store(the_event, j); err != nil {
			log.Println(logs.CodeFilename(1), err)
			return false, err

		}
	}

	return false, nil
}

func storeSent(jid waTypes.JID, id string, msg *waProto.Message) {
	if ctx, err := Database.Begin(); err != nil {
		log.Println("storeSent", err)
	} else {
		defer ctx.Commit()

		if jsn, err := helper.JsonMe(msg); err != nil {
			log.Println("storeSent", err)

		} else {

			if _, err := ctx.Exec(`insert into sent(jid, message_id, message) values($1, $2, $3)`, jid.String(), id, jsn); err != nil {
				log.Println("storeSent", err)
			}
		}
	}
}
