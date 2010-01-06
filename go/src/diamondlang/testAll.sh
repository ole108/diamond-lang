#!/bin/bash

BASEDIR=$(dirname "$0")
source "$BASEDIR/packages.sh"

for dir in $SUBDIRS ; do
  echo ""
  echo "$dir:"
  cd "$BASEDIR/$dir"
  make 'test' || exit 1
  make 'clean'
  cd -
done

