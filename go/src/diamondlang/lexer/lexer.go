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
  indentLevel int;     // current level of indentation
  curChar     byte;
}

func NewLexer(sb common.SrcBuffer) common.Lexer {
  lx := &Lexer{sb, new([MAX_PARENS]byte), 0, 0, 254};
  lx.nextChar();
  return lx;
}

func (lx *Lexer) nextChar() {
  lx.curChar = lx.srcBuf.Getch();
}

func openParen(ch byte, lx *Lexer) byte {
  var ret byte;
  switch ch {
    case ')': ret = '(';
    case ']': ret = '[';
    case '}': ret = '{';
    default: lx.srcBuf.Error("Unknown parenthesis type");
  }
  return ret;
}


// --------------------------------------------------------------------------
// Functions that emit tokens.
// --------------------------------------------------------------------------

func (lx *Lexer) GetToken() common.Token {
  tok := common.Token(nil);
  lxFuncs := []lexFunc{
      tryEof, tryIndent, skipComment, tryNewLine, trySemicolon, tryColon, tryParen,
      tryNumber, tryOperator, tryId, tryChar,
      skipSpaces, signalUndefined
  // NEXT: 
  };

  tok = lx.getFirstTok(lxFuncs);
//fmt.Println(">>> Got first token:", tok.Type(), tok);

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

