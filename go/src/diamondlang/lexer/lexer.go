package lexer

import (
  "diamondlang/common";
//"fmt";
)


// --------------------------------------------------------------------------
// The state of a lexer is held in a variable of this type:
// --------------------------------------------------------------------------
type Lexer struct {
  srcBuf      common.SrcBuffer; // our source for characters, ...
  parenStack  []byte;  // for handling nested parentheses
  inParens    int;     // (how deep) are we inside parentheses?
  curChar     byte;
}

func NewLexer(sb common.SrcBuffer) common.Lexer {
  lx := &Lexer{sb, new([MAX_PARENS]byte), 0, 254};
  lx.nextChar();
  return lx;
}

func (lx *Lexer) Error(msg string) {
  lx.srcBuf.Error(msg);
}

func (lx *Lexer) ClearUpTo(mark common.SrcMark) {
  lx.srcBuf.ClearUpTo(mark);
}

func (lx *Lexer) nextChar() {
  lx.curChar = lx.srcBuf.Getch();
}

func (lx *Lexer) prevChar() {
  lx.srcBuf.Ungetch();
  lx.srcBuf.Ungetch();
  lx.curChar = lx.srcBuf.Getch();
}

func openParen(ch byte, lx *Lexer) byte {
  var ret byte;
  switch ch {
    case ')': ret = '(';
    case ']': ret = '[';
    case '}': ret = '{';
    default: lx.Error("Unknown parenthesis type");
  }
  return ret;
}


// --------------------------------------------------------------------------
// Functions that emit tokens.
// --------------------------------------------------------------------------

func (lx *Lexer) GetToken() common.Token {
  tok := common.Token(nil);
  lxFuncs := []lexFunc{
      trySpace, tryEof, tryComment, tryNewLine, trySemicolon, tryColon, tryParen,
      tryNumber, tryOperator, tryId, tryChar, tryString, signalUndefined
  };

  tok = lx.getFirstTok(lxFuncs);
//fmt.Println(">>> Got token:", tok.Type(), tok);

  return tok;
}

func (lx *Lexer) getFirstTok(lxFuncs []lexFunc) common.Token {
  tok := common.Token(nil);
  moved := false;

  // move until a token is found:
  N := len(lxFuncs);
  for tok == nil {
    moved = false;
    for i := 0; !moved && i < N; i++ {
      tok, moved = lxFuncs[i](lx);
    }
  }
  return tok;
}

