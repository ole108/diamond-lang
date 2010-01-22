include ../../../Make.$(GOARCH)

TARG=diamondlang/llvm
CGOFILES=\
  llvm.go\

CGO_CFLAGS=`llvm-config --cflags`
CGO_LDFLAGS=`llvm-config --ldflags core analysis executionengine jit interpreter native` `llvm-config --libs`

include ../../../Make.pkg

