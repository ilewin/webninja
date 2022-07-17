package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"webp.ninja/services"
	"webp.ninja/utils"
)

func handleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client, socketClient *socketmode.Client) error {
	switch event.Type {
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			handleAppMentionEvent(ev, client)
		case *slackevents.MessageEvent:
			handleDirectMessageEvent(ev, client, socketClient)
		}
	case slackevents.Message:

	default:
		return errors.New("Unsupported event type")
	}

	return nil
}

func handleDirectMessageEvent(event *slackevents.MessageEvent, client *slack.Client, socketClient *socketmode.Client) error {
	
	files := event.Files
	
	if event.User == "U0327KPCSDP" || len(files) < 1 {
		return nil
	}

	config := utils.GetConfig()

	p := services.NewProcessor(len(files))

	sid := uuid.New().String()

	os.Mkdir(fmt.Sprintf(config.App_Storage+"%s", sid), 0755)
	defer services.Clean(fmt.Sprintf(config.App_Storage+"%s", sid))

	for _, f := range files {

		lp := fmt.Sprintf(config.App_Storage+"%s/%s", sid, f.Name)
		fname, oserr := os.Create(lp)
		if oserr != nil {
			log.Fatalf("\nCouldnt create file: %s\n", lp)
		}
		defer fname.Close()

		er := client.GetFile(f.URLPrivateDownload, fname)
		if er != nil {
			fmt.Printf("\nCould not save file: %s to %s", f.Permalink, fname.Name())
		}

		nfrmt, err := utils.GetFormat("WEBP")

		if err != nil {
			log.Fatalf("Could not get format: %v", err)
		}

		p.Files = append(p.Files, services.File{
			Sid:      sid,
			Name:     f.Name,
			Size:     int64(f.Size),
			NSize:    0,
			Format:   f.Mimetype,
			EncodeTo: nfrmt,
			Path:     lp,
		})

	}

	p.Convert()
	go services.UpdateMeta(&p.Files)

	for _, cf := range p.Files {

		abdfpath, _ := filepath.Abs(cf.Path)

		reader, err := os.Open(abdfpath)
		if err != nil {
			log.Printf("Can read file: %s", err)
		}
		defer reader.Close()
		params := slack.FileUploadParameters{
			Title:          cf.Name,
			Filename:       cf.Name,
			Filetype:       cf.Format,
			File:           abdfpath,
			InitialComment: cf.Name,
			Channels:       []string{event.Channel},
		}
		_, uperr := socketClient.UploadFile(params)
		if uperr != nil {
			fmt.Printf("%s\n", err)
			continue
		}

	}

	return nil
}

func handleAppMentionEvent(event *slackevents.AppMentionEvent, client *slack.Client) error {
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}

	text := strings.ToLower(event.Text)

	log.Printf("User %s mentioned WebP Ninja with: \"%s\"", user.Name, text)

	attachment := slack.Attachment{}

	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date",
			Value: time.Now().String(),
		},
		{
			Title: "Initializer",
			Value: user.Name,
		},
	}

	if strings.Contains(text, "hello") {
		attachment.Text = fmt.Sprintf("Hello %s", user.Name)
		attachment.Pretext = "Greetings"
		attachment.Color = "#4af030"
	} else {
		attachment.Text = fmt.Sprintf("How can I help you %s?", user.Name)
		attachment.Pretext = "How can I be of service"
		attachment.Color = "#3d3d3d"
	}

	_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("Failed to send Message: %v", err)
	}
	return nil
}

func handleSlashCommand(command slack.SlashCommand, client *slack.Client) (interface{}, error) {
	switch command.Command {
	case "/convert":
		return nil, handleConvertCommand(command, client)
	case "/ninja-help":
		return handleNinjaHelp(command, client)
	}

	return nil, nil
}

func handleConvertCommand(command slack.SlashCommand, client *slack.Client) error {

	attachment := slack.Attachment{}

	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date",
			Value: time.Now().String(),
		},
		{
			Title: "Initializer",
			Value: command.UserName,
		},
	}

	attachment.Text = fmt.Sprintf("Hello %s", command.Text)
	attachment.Color = "#4af030"
	_, _, err := client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	return nil
}

func handleNinjaHelp(command slack.SlashCommand, client *slack.Client) (interface{}, error) {
	attachment := slack.Attachment{}

	checkbox := slack.NewCheckboxGroupsBlockElement(
		"answer",
		slack.NewOptionBlockObject("yes",
			&slack.TextBlockObject{
				Text: "Yes", Type: slack.MarkdownType,
			},
			&slack.TextBlockObject{
				Text: "Did you enjoy it?", Type: slack.MarkdownType,
			},
		),
		slack.NewOptionBlockObject("no",
			&slack.TextBlockObject{
				Text: "No", Type: slack.MarkdownType,
			},
			&slack.TextBlockObject{
				Text: "Did you dislike it?", Type: slack.MarkdownType,
			},
		),
	)

	// Create the Accessory that will be included in the Block and add the checkbox to it
	accessory := slack.NewAccessory(checkbox)
	// Add Blocks to the attachment
	attachment.Blocks = slack.Blocks{
		BlockSet: []slack.Block{
			// Create a new section block element and add some text and the accessory to it
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: "Did you think this article was helpful?",
				},
				nil,
				accessory,
			),
		},
	}

	attachment.Text = "Rate the tutorial"
	attachment.Color = "#4af030"
	return attachment, nil

}

func handleInteractionEvent(interaction slack.InteractionCallback, client *slack.Client) error {

	log.Printf("The action called is: %s\n", interaction.ActionID)
	log.Printf("The response was of type: %s\n", interaction.Type)
	switch interaction.Type {
	case slack.InteractionTypeBlockActions:
		// This is a block action, so we need to handle it

		for _, action := range interaction.ActionCallback.BlockActions {
			log.Printf("%+v", action)
			log.Println("Selected option: ", action.SelectedOptions)

		}

	default:

	}

	return nil
}

func RunBot(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {

	for {

		select {
		case <-ctx.Done():
			log.Println("Shutting down ninja bot")
			return
		case event := <-socketClient.Events:
			switch event.Type {
			case socketmode.EventTypeEventsAPI:
				eventsApiEvent, ok := event.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Printf("Couldnt typecast event")
					continue
				}
				socketClient.Ack(*event.Request)
				err := handleEventMessage(eventsApiEvent, client, socketClient)
				if err != nil {
					log.Printf("%v", err)
				}

			case socketmode.EventTypeSlashCommand:
				command, ok := event.Data.(slack.SlashCommand)
				if !ok {
					log.Printf("Could not type cast Slash Command")
					continue
				}

				payload, err := handleSlashCommand(command, client)
				if err != nil {
					log.Fatal(err)
				}
				socketClient.Ack(*event.Request, payload)

			case socketmode.EventTypeInteractive:
				interaction, ok := event.Data.(slack.InteractionCallback)
				if !ok {
					log.Printf("Could not type cast the message to interaction callback")
					continue
				}

				err := handleInteractionEvent(interaction, client)
				if err != nil {
					log.Fatal(err)
				}
				socketClient.Ack(*event.Request)
			default:

			}

		}
	}
}
