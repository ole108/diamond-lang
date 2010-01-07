#!/bin/bash

if ! fw +U +D "./diamondlang.fw" &> /dev/null ; then
  cat "./diamondlang.lis"
  exit 1
fi
rm -f *.lis

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

