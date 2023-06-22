package cmp

import (
	"context"
	"sync"

	gdkpixbuf "github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	gio "github.com/diamondburned/gotk4/pkg/gio/v2"
	glib "github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/telekom-mms/corp-net-indicator/internal/assets"
)

var cache struct {
	mu sync.RWMutex
	m  map[assets.Icon]*gdkpixbuf.Pixbuf
}

// returns icon as pixbuf, otherwise panics
func getPixbuf(file assets.Icon) *gdkpixbuf.Pixbuf {
	cache.mu.RLock()
	p, ok := cache.m[file]
	cache.mu.RUnlock()
	if ok {
		return p
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.m == nil {
		cache.m = make(map[assets.Icon]*gdkpixbuf.Pixbuf)
	}

	p, err := gdkpixbuf.NewPixbufFromStream(
		context.Background(),
		gio.NewMemoryInputStreamFromBytes(glib.NewBytesWithGo(assets.GetIcon(file))),
	)
	if err != nil {
		panic(err)
	}
	cache.m[file] = p

	return p
}
