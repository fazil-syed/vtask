package schemas

type SubscribeInput struct {
	Endpoint string `json:"endpoint" binding:"required"`

	Keys struct {
		P256dh string `json:"p256dh" binding:"required"`
		Auth   string `json:"auth" binding:"required"`
	} `json:"keys"`
}
