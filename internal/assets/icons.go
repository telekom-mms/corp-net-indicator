package assets

import (
	"context"
	"embed"
	"sync"

	gdkpixbuf "github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	gio "github.com/diamondburned/gotk4/pkg/gio/v2"
	glib "github.com/diamondburned/gotk4/pkg/glib/v2"
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
	SVGAlert Icon = "icons/close_circle.svg"
	SVGCheck Icon = "icons/check_circle.svg"
	SVGMinus Icon = "icons/minus_circle.svg"
)

var cache struct {
	mu sync.RWMutex
	m  map[Icon]*gdkpixbuf.Pixbuf
}

// returns known icons, otherwise panics
func GetIcon(file Icon) []byte {
	data, err := icons.ReadFile(string(file))
	if err != nil {
		panic(err)
	}
	return data
}

// returns icon as pixbuf, otherwise panics
func GetPixbuf(file Icon) *gdkpixbuf.Pixbuf {
	cache.mu.RLock()
	p, ok := cache.m[file]
	cache.mu.RUnlock()
	if ok {
		return p
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.m == nil {
		cache.m = make(map[Icon]*gdkpixbuf.Pixbuf)
	}

	p, err := gdkpixbuf.NewPixbufFromStream(
		context.Background(),
		gio.NewMemoryInputStreamFromBytes(glib.NewBytesWithGo(GetIcon(file))),
	)
	if err != nil {
		panic(err)
	}
	cache.m[file] = p

	return p
}
