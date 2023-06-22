# Corporate Network Indicator
Provides a tray icon and a status window to display corporate network status data, such as identity agent, vpn details.

## Installation

For installation you can chose between 2 options:

### Using Debian/Ubuntu package

Download the package from releases page and use the following instructions to install and activate the agent:

```console
$ sudo apt install ./corp-net-indicator.deb
$ sudo systemctl --user start corp-net-indicator.service
```

### Using tar.gz archive

Download the archive from releases page and use the following instructions to install and activate the agent:

```console
$ tar -xf corp-net-indicator.tar.gz && cd <extracted directory>
$ sudo cp corp-net-indicator /usr/bin/
$ sudo cp corp-net-indicator-win /usr/bin/
$ sudo cp corp-net-indicator.service /usr/lib/systemd/user/
$ sudo cp corp-net-indicator-win.desktop /usr/share/applications/
$ sudo systemctl --user enable corp-net-indicator.service
$ sudo systemctl --user start corp-net-indicator.service
```

## Interface
![Screenshot](/screenshot.png)
# Credits
[Eva Icons](https://github.com/akveo/eva-icons)