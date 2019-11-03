package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hearkat/hearkat-go"
	"net/http"
	"regexp"
)

type Frame struct {
	Icon     string `json:"icon,omitempty"`
	Text     string `json:"text,omitempty"`
	GoalData struct {
		Start   int    `json:"start"`
		Current int    `json:"current"`
		End     int    `json:"end"`
		Unit    string `json:"unit"`
	} `json:"goalData,omitempty"`
	ChartData []int `json:"chartData,omitempty"`
}

type Sound struct {
	Category string `json:"category"`
	ID       string `json:"id"`
	Repeat   int    `json:"repeat"`
}

type Model struct {
	Frames []Frame `json:"frames"`
	Sound  Sound   `json:"sound"`
	Cycles int     `json:"cycles"`
}

type Notification struct {
	Priority string `json:"priority"`
	IconType string `json:"icon_type"`
	LifeTime int    `json:"lifeTime,omitempty"`
	Model    Model  `json:"model"`
}

func (p *Lametric) toLametric(channel string, n *hearkat.Notification) *Notification {
	icon := "a2899"
	sound := &Sound{
		Category: "notifications",
		ID:       "notification",
		Repeat:   1,
	}

	for _, v := range p.config.Icons {
		m, err := regexp.MatchString(v.Pattern, channel)
		if err != nil {
			panic(err)
		}
		if !m {
			continue
		}

		icon = v.Icon

		if v.Sound != nil {
			if *v.Sound == "none" {
				sound = nil
			} else {
				sound = &Sound{
					Category: "notifications",
					ID:       *v.Sound,
					Repeat:   1,
				}
			}
		}
	}

	return &Notification{
		Model: Model{
			Frames: []Frame{
				{
					Text: n.Title,
					Icon: icon,
				},
				{
					Text: n.Message,
				},
			},
			Cycles: 1,
			Sound:  *sound,
		},
		IconType: "info",
		Priority: "info",
	}
}

func (p *Lametric) sendNotification(notif *Notification) {
	d, err := json.Marshal(notif)
	if err != nil {
		panic(err)
	}

	c := http.Client{}
	addr := fmt.Sprintf("https://%s:4343/api/v2/device/notifications", p.config.IP)
	req, err := http.NewRequest("POST", addr, bytes.NewReader(d))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("dev", p.config.ApiKey)

	res, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
