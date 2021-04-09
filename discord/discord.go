package discord

import (
	"context"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type Discord struct {
	channelID        string
	client           *discordgo.Session
	isDevelopment    bool
	messagesIdCache  map[string]string
	ctx              context.Context
	sendMessageMutex sync.Mutex
	timeLocation     *time.Location
}

func InitDiscord(token string, channelID string, isDevelopment bool, ctx context.Context, timeLocation *time.Location) (*Discord, error) {
	d := &Discord{channelID, nil, isDevelopment, nil, ctx, sync.Mutex{}, timeLocation}
	d.messagesIdCache = make(map[string]string, 200)

	var err error
	d.client, err = discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	err = d.client.Open()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Discord) cleanChannel() error {
	messages, err := d.client.ChannelMessages(d.channelID, 100, "", "", "")
	if err != nil {
		return err
	}
	ids := make([]string, 0, 100)
	for _, m := range messages {
		if m.Author.Bot {
			ids = append(ids, m.ID)
		}
	}
	err = d.client.ChannelMessagesBulkDelete(d.channelID, ids)
	if err != nil {
		return err
	}
	return nil
}

func (d *Discord) SendMessageToUser(message string, userId string) error {
	d.sendMessageMutex.Lock()
	defer d.sendMessageMutex.Unlock()
	channel, err := d.client.UserChannelCreate(userId)
	if err != nil {
		logrus.WithError(err).Error("Error sending message")
	}
	_, err = d.client.ChannelMessageSend(channel.ID, message)
	return err
}

func (d *Discord) SendMessage(message string) error {
	d.sendMessageMutex.Lock()
	defer d.sendMessageMutex.Unlock()
	_, err := d.client.ChannelMessageSend(d.channelID, message)
	return err
}

func (d *Discord) SendDebugMessage(message string) error {
	if !d.isDevelopment {
		return nil
	}
	return d.SendMessage(message)
}

func (d *Discord) Close() {
	d.client.Close()
}
