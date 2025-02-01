package assets

import (
	"embed"
)

//go:embed icons/*
var icons embed.FS

type Icon string

const (
	ShieldOff  Icon = "icons/shield_off.png"
	ShieldOn   Icon = "icons/shield.png"
	Umbrella   Icon = "icons/umbrella.png"
	Connect    Icon = "icons/connect.png"
	Disconnect Icon = "icons/disconnect.png"
	Status     Icon = "icons/activity.png"
)

const (
	SVGAlert   Icon = "icons/close_circle.svg"
	SVGCheck   Icon = "icons/check_circle.svg"
	SVGMinus   Icon = "icons/minus_circle.svg"
	SVGWarning Icon = "icons/alert_triangle.svg"
)

// returns known icons, otherwise panics
func GetIcon(file Icon) []byte {
	data, err := icons.ReadFile(string(file))
	if err != nil {
		panic(err)
	}
	return data
}
