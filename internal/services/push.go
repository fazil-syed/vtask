package services

import (
	"context"
	"encoding/json"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/syed.fazil/vtask/internal/config"
	"github.com/syed.fazil/vtask/internal/models"
)

type Payload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Data  struct {
		URL string `json:"url"`
	} `json:"data"`
	Vibrate bool `json:"vibrate"`
}
type PushService struct {
	ctx context.Context
	cfg config.Config
}

func NewPushService(ctx context.Context, cfg config.Config) *PushService {
	return &PushService{
		ctx: ctx,
		cfg: cfg,
	}
}
func (s *PushService) Send(sub models.PushSubscription, payload Payload) error {
	subscription := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.P256dh,
			Auth:   sub.Auth,
		},
	}
	body, _ := json.Marshal(payload)
	resp, err := webpush.SendNotificationWithContext(s.ctx, body, subscription, &webpush.Options{
		Subscriber:      "mailto:support@vtask.com",
		VAPIDPublicKey:  s.cfg.VAPID_PUBLIC_KEY,
		VAPIDPrivateKey: s.cfg.VAPID_PRIVATE_KEY,
		TTL:             60,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
