package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

type GoWhatsBot struct {
	Config       Config
	Client       *whatsmeow.Client
	Container    *sqlstore.Container
	DeviceStore  *store.Device
	DeviceProps  *waProto.DeviceProps
	EventHandler whatsmeow.EventHandler
	DLog         waLog.Logger
	CLog         waLog.Logger
	QRCode       string
	Connected    bool
	Channel      chan os.Signal
}

type IGoWhatsBot interface {
	GetConfig(key string) string
	GetDeviceProps() *waProto.DeviceProps
	Run()
	Stop()
	GetClient() *whatsmeow.Client
	AddEventHandler(e whatsmeow.EventHandler) uint32
}

func NewGoWhatsBot(cfg Config) IGoWhatsBot {
	var newGWB = GoWhatsBot{}

	newGWB.DLog = waLog.Stdout("Database", "ERROR", true)
	newGWB.CLog = waLog.Stdout("Client", "ERROR", true)
	newGWB.Config = cfg

	dbdriver := cfg.GetByKey("driver", "sqlite3")
	dbaddress := cfg.GetByKey(dbdriver, "file:gowhatsbot.db?_foreign_keys=on")
	if ctr, err := sqlstore.New(dbdriver, dbaddress, newGWB.DLog); err != nil {
		panic(err)
	} else {
		newGWB.Container = ctr
	}

	if dvc, err := newGWB.Container.GetFirstDevice(); err != nil {
		panic(err)
	} else {
		newGWB.DeviceStore = dvc
	}

	store.DeviceProps.Os = proto.String("GoWhatsBot")
	store.DeviceProps.PlatformType = waProto.DeviceProps_DESKTOP.Enum()

	newGWB.DeviceProps = store.DeviceProps
	newGWB.Client = whatsmeow.NewClient(newGWB.DeviceStore, newGWB.CLog)

	return &newGWB
}
func (g *GoWhatsBot) AddEventHandler(e whatsmeow.EventHandler) uint32 {
	return g.Client.AddEventHandler(e)
}
func (g *GoWhatsBot) GetDeviceProps() *waProto.DeviceProps {
	return g.DeviceProps
}

func (g *GoWhatsBot) GetConfig(key string) string {

	return g.Config.GetByKey(key, "")

}

func (g *GoWhatsBot) Run() {
	if g.Client.Store.ID == nil {
		if qrchan, err := g.Client.GetQRChannel(context.Background()); err == nil {
			if err := g.Client.Connect(); err != nil {
				panic(err)
			}

			for qritem := range qrchan {
				if qritem.Event == "code" {
					g.QRCode = qritem.Code
					qrterminal.GenerateHalfBlock(qritem.Code, qrterminal.M, os.Stdout)
				} else {
					g.Connected = qritem.Event == "success"
				}
			}
		}
	} else {
		if err := g.Client.Connect(); err != nil {
			panic(err)
		}
	}

	g.Channel = make(chan os.Signal, 1)
	signal.Notify(g.Channel, os.Interrupt, syscall.SIGTERM)

	<-g.Channel

	g.Client.Disconnect()
}

func (g *GoWhatsBot) Stop() {
	close(g.Channel)
}

func (g *GoWhatsBot) GetClient() *whatsmeow.Client {
	return g.Client
}
