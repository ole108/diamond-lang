#!/bin/bash

BASEDIR=$(dirname "$0")
source "$BASEDIR/packages.sh"

for dir in $SUBDIRS ; do
  echo ""
  echo "$dir:"
  cd "$BASEDIR/$dir"
  make || exit 1
  cd -
done
