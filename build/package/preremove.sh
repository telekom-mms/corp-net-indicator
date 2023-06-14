#!/bin/sh -e
# taken from https://git.launchpad.net/ubuntu/+source/debhelper/tree/autoscripts/prerm-systemd?h=applied/13.6ubuntu1
# adapted to user units

UNIT='corp-net-indicator.service'

case "$1" in
  'remove' | 'purge')
    if [ -z "${DPKG_ROOT:-}" ] && [ -d /run/systemd/system ] ; then
      deb-systemd-invoke --global stop $UNIT >/dev/null || true
    fi
    ;;
esac