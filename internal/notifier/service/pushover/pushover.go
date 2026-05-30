package pushover

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"donetick.com/core/config"
	nModel "donetick.com/core/internal/notifier/model"
	"donetick.com/core/logging"
	"github.com/gregdel/pushover"
)

type Pushover struct {
	pushover *pushover.Pushover
	baseURL  string
}

func NewPushover(cfg *config.Config) *Pushover {

	pushoverApp := pushover.New(cfg.Pushover.Token)

	return &Pushover{
		pushover: pushoverApp,
		baseURL:  strings.TrimSuffix(cfg.Notifier.AppHost, "/"),
	}
}

func (p *Pushover) SendNotification(c context.Context, notification *nModel.NotificationDetails) error {
	if notification.TargetID == "" {
		return errors.New("unable to send notification, targetID is empty")
	}
	log := logging.FromContext(c)
	recipient := pushover.NewRecipient(notification.TargetID)
	message := pushover.NewMessageWithTitle(notification.Text, "Donetick")

	if p.baseURL != "" {
		if notification.ChoreID > 0 {
			message.URL = fmt.Sprintf("%s/chores/%d", p.baseURL, notification.ChoreID)
			message.URLTitle = "View Task"
		} else {
			message.URL = fmt.Sprintf("%s/chores", p.baseURL)
			message.URLTitle = "View Tasks"
		}
	}

	_, err := p.pushover.SendMessage(message, recipient)
	if err != nil {
		log.Debug("Error sending pushover notification", err)
		return err
	}

	return nil
}
