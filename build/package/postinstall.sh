#!/bin/sh -e
# taken from https://git.launchpad.net/ubuntu/+source/debhelper/tree/autoscripts/postinst-systemd-user-enable?h=applied/13.6ubuntu1

UNIT='corp-net-indicator.service'

case "$1" in
  'configure' | 'abort-upgrade' | 'abort-deconfigure' | 'abort-remove')
    # systemctl daemon-reload
    # systemctl --global enable $UNIT
    # This will only remove masks created by d-s-h on package removal.
    deb-systemd-helper --user unmask $UNIT >/dev/null || true

    # was-enabled defaults to true, so new installations run enable.
    if deb-systemd-helper --quiet --user was-enabled $UNIT ; then
      # Enables the unit on first installation, creates new
      # symlinks on upgrades if the unit file has changed.
      deb-systemd-helper --user enable $UNIT >/dev/null || true
    else
      # Update the statefile to add new symlinks (if any), which need to be
      # cleaned up on purge. Also remove old symlinks.
      deb-systemd-helper --user update-state $UNIT >/dev/null || true
    fi
    ;;
esac
