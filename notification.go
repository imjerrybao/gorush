package main

import (
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	_ "github.com/google/go-gcm"
	"log"
)

type ExtendJSON struct {
	Key   string `json:"key"`
	Value string `json:"val"`
}

type alert struct {
	Action       string   `json:"action,omitempty"`
	ActionLocKey string   `json:"action-loc-key,omitempty"`
	Body         string   `json:"body,omitempty"`
	LaunchImage  string   `json:"launch-image,omitempty"`
	LocArgs      []string `json:"loc-args,omitempty"`
	LocKey       string   `json:"loc-key,omitempty"`
	Title        string   `json:"title,omitempty"`
	TitleLocArgs []string `json:"title-loc-args,omitempty"`
	TitleLocKey  string   `json:"title-loc-key,omitempty"`
}

type RequestPushNotification struct {
	// Common
	Tokens   []string `json:"tokens" binding:"required"`
	Platform int      `json:"platform" binding:"required"`
	Message  string   `json:"message" binding:"required"`
	Priority string   `json:"priority,omitempty"`
	ContentAvailable bool `json:"content_available,omitempty"`

	// Android
	CollapseKey    string `json:"collapse_key,omitempty"`
	DelayWhileIdle bool   `json:"delay_while_idle,omitempty"`
	TimeToLive     int    `json:"time_to_live,omitempty"`

	// iOS
	ApnsID string       `json:"apns_id,omitempty"`
	Topic  string       `json:"topic,omitempty"`
	Badge  int          `json:"badge,omitempty"`
	Sound  string       `json:"sound,omitempty"`
	Expiry int          `json:"expiry,omitempty"`
	Retry  int          `json:"retry,omitempty"`
	Category string     `json:"category,omitempty"`
	URLArgs          []string    `json:"url-args,omitempty"`
	Extend []ExtendJSON `json:"extend,omitempty"`
	Alert alert `json:"alert,omitempty"`
	// meta
	IDs []uint64 `json:"seq_id,omitempty"`
}

func pushNotification(notification RequestPushNotification) bool {
	var (
		success bool
	)

	cert, err := certificate.FromPemFile("./key.pem", "")
	if err != nil {
		log.Println("Cert Error:", err)
	}

	apnsClient := apns.NewClient(cert).Development()

	switch notification.Platform {
	case PlatFormIos:
		success = pushNotificationIos(notification, apnsClient)
		if !success {
			apnsClient = nil
		}
	case PlatFormAndroid:
		success = pushNotificationAndroid(notification)
	}

	return success
}

func pushNotificationIos(req RequestPushNotification, client *apns.Client) bool {

	for _, token := range req.Tokens {
		notification := &apns.Notification{}
		notification.DeviceToken = token

		if len(req.ApnsID) > 0 {
			notification.ApnsID = req.ApnsID
		}

		if len(req.Topic) > 0 {
			notification.Topic = req.Topic
		}

		if len(req.Priority) > 0 && req.Priority == "normal" {
			notification.Priority = apns.PriorityLow
		}

		payload := payload.NewPayload().Alert(req.Message)

		if req.Badge > 0 {
			payload.Badge(req.Badge)
		}

		if len(req.Sound) > 0 {
			payload.Sound(req.Sound)
		}

		if req.ContentAvailable {
			payload.ContentAvailable()
		}

		if len(req.Extend) > 0 {
			for _, extend := range req.Extend {
				payload.Custom(extend.Key, extend.Value)
			}
		}

		// Alert dictionary

		if len(req.Alert.Title) > 0 {
			payload.AlertTitle(req.Alert.Title)
		}

		if len(req.Alert.TitleLocKey) > 0 {
			payload.AlertTitleLocKey(req.Alert.TitleLocKey)
		}

		if len(req.Alert.LocArgs) > 0 {
			payload.AlertTitleLocArgs(req.Alert.LocArgs)
		}

		if len(req.Alert.Body) > 0 {
			payload.AlertBody(req.Alert.Body)
		}

		if len(req.Alert.LaunchImage) > 0 {
			payload.AlertLaunchImage(req.Alert.LaunchImage)
		}

		if len(req.Alert.LocKey) > 0 {
			payload.AlertLocKey(req.Alert.LocKey)
		}

		if len(req.Alert.Action) > 0 {
			payload.AlertAction(req.Alert.Action)
		}

		if len(req.Alert.ActionLocKey) > 0 {
			payload.AlertActionLocKey(req.Alert.ActionLocKey)
		}

		// General

		if len(req.Category) > 0 {
			payload.Category(req.Category)
		}

		if len(req.URLArgs) > 0 {
			payload.URLArgs(req.URLArgs)
		}

		notification.Payload = payload

		// send ios notification
		res, err := client.Push(notification)

		if err != nil {
			log.Println("There was an error", err)
			return false
		}

		if res.Sent() {
			log.Println("APNs ID:", res.ApnsID)
		}
	}

	client = nil

	return true
}

func pushNotificationAndroid(req RequestPushNotification) bool {

	return true
}
