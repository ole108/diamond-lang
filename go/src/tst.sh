#!/bin/bash

../main -c \
  '#bla 0b0110' \
  'BLA:' \
  '(E = 2*3)&[{TRUE?(0x01_F)+PI_HOCH_2 ^  001230};0c17 * 0b0110] -/- 0r036_10' \
  'If bla > 0:' \
  '    mod.Func mod.CONST mod.CONST.val Fn i' \
  '  Elif bla < 0:' \
  '    mod.FuncAli bla' \
  '  Else:' \
  '    bla.val  # should work!' \
  'bla = 0' \
  '   # geschafft!'  &> tst.out

if ( cmp -s tst.out tst.right ) ; then
  echo "OK."
  rm -f tst.out
else
  diff tst.out tst.right
fi

