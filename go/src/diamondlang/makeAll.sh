#!/bin/bash

SUBDIRS="common srcbuf lexer tokbuf parser"

BASEDIR=$(dirname "$0")

for dir in $SUBDIRS ; do
  echo ""
  echo "$dir:"
  cd "$BASEDIR/$dir"
  make || exit 1
  cd -
done

