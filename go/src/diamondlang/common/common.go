package common

import (
  "strings";
  "fmt";
  "os";
  "container/list";
)

var TABSIZE = 4;  // 'almost constant' should be set only by main!
const EOF = 255;  // used between SrcBuffer and Lexer

// --------------------------------------------------------------------------
// Constants for token types:
// --------------------------------------------------------------------------

// Define token types as 'enumeration':
type TokEnum int
const (
  // basic source file structure:
  TOK_EOF = iota;
  TOK_NL;
  TOK_COLON;

  // indentation:
  TOK_INDENT;
  TOK_HALF_INDENT;
  TOK_DEDENT;
  TOK_HALF_DEDENT;

  // parentheses:
  TOK_PAREN_OPEN;
  TOK_PAREN_CLOSE;

  // beginning of block ::= ':' '\n'
  TOK_BLOCK_START;

  // comment ::= '#' ...
  TOK_COMMENT;
  TOK_SPACE;

  // identifiers:
  TOK_VAL_ID;
  TOK_FUNC_ID;
  TOK_OP_ID;
  TOK_CONST_ID;
  TOK_MODULE_ID;

  // constant values:
  TOK_INT;
  TOK_STR;
  TOK_CHAR;

  // keywords:
  TOK_DEF;
  TOK_EXTERN;
  TOK_IMPORT;
  TOK_SHELF;
  TOK_BIND;
  TOK_SHADOWED;
)
func (te TokEnum) String() string {
  ret := "";
  switch te {
  case TOK_EOF:          ret = "<TOK EOF>";
  case TOK_NL:           ret = "<TOK NL>";
  case TOK_COLON:        ret = "<TOK COLON>";
  case TOK_INDENT:       ret = "<TOK INDENT>";
  case TOK_HALF_INDENT:  ret = "<TOK HALF INDENT>";
  case TOK_DEDENT:       ret = "<TOK DEDENT>";
  case TOK_HALF_DEDENT:  ret = "<TOK HALF DEDENT>";
  case TOK_PAREN_OPEN:   ret = "<TOK PAREN OPEN>";
  case TOK_PAREN_CLOSE:  ret = "<TOK PAREN CLOSE>";
  case TOK_BLOCK_START:  ret = "<TOK BLOCK START>";
  case TOK_COMMENT:      ret = "<TOK COMMENT>";
  case TOK_SPACE:        ret = "<TOK SPACE>";
  case TOK_CONST_ID:     ret = "<TOK CONST ID>";
  case TOK_MODULE_ID:    ret = "<TOK MODULE ID>";
  case TOK_VAL_ID:       ret = "<TOK VAL ID>";
  case TOK_FUNC_ID:      ret = "<TOK FUNC ID>";
  case TOK_OP_ID:        ret = "<TOK OP ID>";
  case TOK_INT:          ret = "<TOK INT>";
  case TOK_STR:          ret = "<TOK STR>";
  case TOK_CHAR:         ret = "<TOK CHAR>";
  case TOK_DEF:          ret = "<TOK DEF>";
  case TOK_EXTERN:       ret = "<TOK EXTERN>";
  case TOK_IMPORT:       ret = "<TOK IMPORT>";
  case TOK_SHELF:        ret = "<TOK SHELF>";
  case TOK_BIND:         ret = "<TOK BIND>";
  case TOK_SHADOWED:     ret = "<TOK SHADOWED>";
  default:               ret = fmt.Sprintf("<TOK %d>", te);
  }
  return ret;
}


// --------------------------------------------------------------------------
// Free functions:
// --------------------------------------------------------------------------

func IsSpace(ch byte) bool {
  return ch == ' ' || ch == '\t';
}
func SpaceAmount(ch byte) int {
  ret := 0;
  if ch == ' '       { ret = 1; }
  else if ch == '\t' { ret = TABSIZE; }
  return ret;
}

func MakeErrString(msg string, lineNum int, lineStr string, markStart int, markLen int) string {
  return fmt.Sprintf("%s at line %d near:\n", msg, lineNum+1) +
         lineStr + "\n" +
         strings.Repeat(" ", markStart) + strings.Repeat("^", markLen) + "\n";
}

func HandleFatal(msg string) {
  fmt.Fprint(os.Stderr, "FATAL ERROR: " + msg);
  os.Exit(1);
}


// --------------------------------------------------------------------------
// Interfaces:
// --------------------------------------------------------------------------

// a marker inside a source buffer
type SrcMark struct {
  Elem  *list.Element;
  Col    int;
}

// a piece of a source buffer
type SrcPiece interface {
  Start() SrcMark;
  End() SrcMark;
  StartLine() int;
  StartColumn() int;
  Content() string;
  WholeLine() string;
  String() string;
  Error(msg string);
}

// the interface the Lexer needs as a source
type SrcBuffer interface {
  Error(msg string);
  AtStartOfLine() bool;
  NotAtStartOfLine();
  Ungetch();
  Getch() byte;
  NewMark() SrcMark;
  NewPiece(start SrcMark) SrcPiece;
  NewAnyPiece(start SrcMark, end SrcMark) SrcPiece;
  ClearUpTo(mark SrcMark);
}

// The interface all tokens returned from the Lexer implement
type Token interface {
  SourcePiece() SrcPiece;
  HasSpaceAround() int;
  Type() TokEnum;
  Error(msg string);
  StartLine() int;
  StartColumn() int;
  WholeLine() string;
  Content() string;
  String() string;
}

type Lexer interface {
  GetToken() Token;
  NewCopyTok(typ TokEnum, tok Token) Token;
  NewAnyTok(typ TokEnum, start SrcMark, end SrcMark) Token;
  NewSpaceTok(tok Token, space int, atStartOfLine bool) Token;
  Error(msg string);
  ClearUpTo(mark SrcMark);
}

type TokenBuffer interface {
  GetToken() Token;
  ClearUpTo(tok Token);
  Error(msg string);
}

type Parser interface {
}
