package main

import (
	"fmt"

	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
)

type hook struct {
	webhook api.WebhookClient
}

func (webh *hook) send_hook(content string) (bool, error) {
	_, err := webh.webhook.SendEmbeds(api.NewEmbedBuilder().
		SetFooter(fmt.Sprintf("ss-sniper ~ https://0xTokens.xyz ~ V%s", version), "").
		SetDescription(content).
		Build(),
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

func new_hook(token string, id string) (*hook, error) {
	webhook, err := disgohook.NewWebhookClientByToken(nil, nil, fmt.Sprintf("%s/%s", id, token))

	if err != nil {
		return nil, err
	}

	return &hook{
		webhook: webhook,
	}, nil
}
