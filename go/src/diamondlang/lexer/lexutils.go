package lexer

import (
  "strings";
)


// --------------------------------------------------------------------------
// Constants for the lexer
// --------------------------------------------------------------------------

const MAX_PARENS = 8
const OPERATOR_CHARS = "+-*/%^<>!=&|,?@$~"
const NUM_CHARS = "_0123456789abcdefghijklmnopqrstuvwxyz"


// --------------------------------------------------------------------------
// Free functions:
// --------------------------------------------------------------------------

// look for char in string
func IsIn(ch byte, str string) bool {
  return strings.Index(str, string(ch)) >= 0;
}

func isOpChar(ch byte) bool { return IsIn(ch, OPERATOR_CHARS) }
func isDigit(ch byte) bool { return (ch >= '0' && ch <= '9') }
func isUpper(ch byte) bool { return (ch >= 'A' && ch <= 'Z') }
func isLower(ch byte) bool { return (ch >= 'a' && ch <= 'z') }
func isAlpha(ch byte) bool { return isLower(ch) || isUpper(ch) }
func isConstChar(ch byte) bool { return (isUpper(ch) || ch == '_' || isDigit(ch)) }

func isNumChar(ch byte, base int) bool {
  idx := strings.Index(NUM_CHARS, string(lower(ch)));
  return idx >= 0 && idx <= base;
}


func lower(ch byte) byte {
  if ch >= 'A' && ch <= 'Z' {
    return 'a' + (ch - 'A');
  }
  return ch;
}

func isIdChar(ch byte) bool {
  return (isAlpha(ch) || isDigit(ch) || ch == '_' || ch == '.');
}

func strRemove(s string, ch byte) string {
  ret := "";
  delRune := int(ch);
  for _, rune := range s {
    if rune != delRune { ret += string(rune); }
  }
  return ret;
}

