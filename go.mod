module github.com/telekom-mms/corp-net-indicator

go 1.23.0

toolchain go1.23.5

require (
	github.com/diamondburned/gotk4/pkg v0.2.2
	github.com/slytomcat/systray v1.10.2
	golang.org/x/text v0.21.0
)

require (
	github.com/KarpelesLab/weak v0.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/telekom-mms/tnd v0.5.1 // indirect
	github.com/vishvananda/netlink v1.3.0 // indirect
	github.com/vishvananda/netns v0.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/godbus/dbus/v5 v5.1.0
	github.com/stretchr/testify v1.10.0
	github.com/telekom-mms/fw-id-agent v1.1.2
	github.com/telekom-mms/oc-daemon v1.2.0
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20231121144256-b99613f794b6 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

replace github.com/godbus/dbus/v5 v5.1.0 => github.com/malaupa/dbus/v5 v5.2.0-next
