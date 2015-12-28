package main

import (
        "fmt"
        "github.com/nlopes/slack"
        "os"
)

func main() {
        apiKey := os.Getenv("SLACK_API_KEY")
        slackChannel := os.Getenv("SLACK_CHANNEL")

        api := slack.New(apiKey)
        api.SetDebug(false)

        rtm := api.NewRTM()
        go rtm.ManageConnection()

mainEventLoop:
        for {
                select {
                case msg := <-rtm.IncomingEvents:
                        switch slackEvent := msg.Data.(type) {
                        case *slack.MessageEvent:
                                if slackEvent.SubType == "channel_leave" && slackEvent.Channel == slackChannel {
                                        api.InviteUserToChannel(slackEvent.Channel, slackEvent.User)
                                        userData, _ := api.GetUserInfo(slackEvent.User)
                                        fmt.Printf("User %s attempted to leave!\n", userData.Name)
                                }
                        case *slack.InvalidAuthEvent:
                                fmt.Printf("Invalid Authentication: %+v\n", slackEvent)
                                break mainEventLoop // Exit main event loop, The API key is invalid
                        }
                }
        }
}
