#!/bin/bash

BASEDIR=$(dirname "$0")
source "$BASEDIR/packages.sh"

for dir in $SUBDIRS ; do
  echo ""
  echo "$dir:"
  cd "$BASEDIR/$dir"
  if [[ "${dir}.go.fw" -nt "${dir}_test.go" ]] ; then
    fw "${dir}.go.fw" || exit 1
  fi
  rm -f *.lis
  make || exit 1
  cd -
done

