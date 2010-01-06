#!/bin/bash

fw "./diamondlang.fw" || exit 1
rm -f *.lis

BASEDIR=$(dirname "$0")
source "$BASEDIR/packages.sh"

for dir in $SUBDIRS ; do
  echo ""
  echo "$dir:"
  cd "$BASEDIR/$dir"
  make || exit 2
  cd -
done

