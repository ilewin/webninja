package main

import (
	"context"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"webp.ninja/slack/bot"
	"webp.ninja/utils"
)

func main() {

	log.Println("Slack Bot Started")
	config := utils.GetConfig()

	token := config.Slack_Auth_Token
	appToken := config.Slack_App_Token

	client := slack.New(token, slack.OptionDebug(false), slack.OptionAppLevelToken(appToken))

	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(false),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go bot.RunBot(ctx, client, socketClient)
	socketClient.Run()

}
