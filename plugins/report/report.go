package main

import (
	"fmt"
	"log"

	"strings"

	"github.com/mamur-rezeki/gowhatsplugins"
	"github.com/mamur-rezeki/gowhatsplugins/types"
	"github.com/mamur-rezeki/gowhatsplugins/whats"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

var Plugin = types.Plugin{
	Name:     "Reporter",
	Validate: gowhatsplugins.MasterValidator,
}

func init() {

	Plugin.CommandAddMany([]*types.Command{
		{
			Cmd:         []string{".report-app", ".ra"},
			Description: "Show App report",
			Usage:       "{cmd}",
			Execute:     reportApp,
		}, {
			Cmd:         []string{".report-auth", ".ru"},
			Description: "Show App report",
			Usage:       "{cmd}",
			Execute:     reportAuth,
		},
	})

}

func _(pattern string, args []string, cmd *types.Command, event *events.Message, ctx *waProto.ContextInfo, client *whatsmeow.Client) error {

	return nil
}

func reportApp(pattern string, args []string, cmd *types.Command, event *events.Message, ctx *waProto.ContextInfo, client *whatsmeow.Client) error {

	loadPlugins()
	var report_lines = []string{}

	for _, plg := range GoWhatsPlugins {
		var t = fmt.Sprintf("*[%d] %s*", len(plg.Commands), plg.Name)
		switch plg.Disabled {
		case true:
			t += " _*disabled_"
		}
		report_lines = append(report_lines, t)

		for _, cmd := range plg.Commands {
			var disabled = ""
			if cmd.Disabled {
				disabled = "_*disabled_"
			}
			report_lines = append(report_lines, fmt.Sprintf("  %s %s", strings.Join(cmd.Cmd, ", "), disabled))

		}
	}

	if len(report_lines) > 0 {
		report_lines = append([]string{"*📊 Report App*", strings.Repeat("-", 40)}, report_lines...)

		var sentctx, _ = whats.SanitizeContext(event, event.Info.IsGroup, client)
		if emsg, err := whats.NewExtendedMessage(strings.Join(report_lines, "\n"), sentctx); err != nil {
			return err
		} else {
			if resp, err := whats.SendMessage(event.Info.Chat, emsg, client); err != nil {
				return err
			} else {
				log.Println(pattern, resp.ID)
			}
		}
	}

	return nil
}

func reportAuth(pattern string, args []string, cmd *types.Command, event *events.Message, ctx *waProto.ContextInfo, client *whatsmeow.Client) error {
	var report_lines = []string{}

	for key, vals := range gowhatsplugins.Auths {
		report_lines = append(report_lines, fmt.Sprintf("*[%d]* for *%s*", len(vals), strings.ToTitle(key)))
		for _, u := range vals {
			report_lines = append(report_lines, "  - "+u)
		}
	}

	if len(report_lines) > 0 {
		report_lines = append([]string{"*📊 Report Auth*", strings.Repeat("-", 30)}, report_lines...)

		var sentctx, _ = whats.SanitizeContext(event, event.Info.IsGroup, client)
		if emsg, err := whats.NewExtendedMessage(strings.Join(report_lines, "\n"), sentctx); err != nil {
			return err
		} else {
			if resp, err := whats.SendMessage(event.Info.Chat, emsg, client); err != nil {
				return err
			} else {
				log.Println(pattern, resp.ID)
			}
		}
	}
	return nil
}
