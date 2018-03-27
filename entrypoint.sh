#!/bin/sh

if [ "$#" -ne 0 ]; then
  exec $@
fi

exec ipxe-builder
