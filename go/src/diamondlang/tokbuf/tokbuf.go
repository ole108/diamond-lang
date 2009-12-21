package tokbuf

import (
  "diamondlang/common";
  "diamondlang/lexer";
  "container/list";
//"fmt";
)

const MAX_TOK = 128;

type tokBuf struct {
  lx      common.Lexer;  // our source for tokens
  tokBuf  *list.List;    // the real token buffer
  curTok  *list.Element; // position of the current token in the buffer
  emitNl  bool;          // should a new line token be returned?
}

func NewTokenBuffer(lx common.Lexer) common.TokenBuffer {
  return &tokBuf{lx, list.New(), nil, false};
}

func (tb *tokBuf) Error(msg string) {
  tb.lx.Error(msg);
}

// unget a token (simply go back in buffer if still possible).
func (tb *tokBuf) ungetToken() {
  if tb.curTok == nil || tb.curTok.Prev() == nil {
    tb.Error("Unable to go beyond the start of the token buffer");
  }
  tb.curTok = tb.curTok.Prev();
}

func (tb *tokBuf) GetToken() common.Token {
  if tb.atEnd() {
    tb.readToken();
    tb.curTok = tb.tokBuf.Back();
  } else {
    tb.curTok = tb.curTok.Next();
  }
  return tb.any2token(tb.curTok.Value);
}

func (tb *tokBuf) atEnd() bool {
  return tb.curTok == nil || tb.curTok.Next() == nil;
}

type caseHandler func(common.Token, *tokBuf) bool

func (tb *tokBuf) readToken() {
  tok := tb.getFilteredToken();
  caseHandlers := []caseHandler{ handleColon, handleMultiDedent, handleDefault };
  handled := false;
  for i := 0; i < len(caseHandlers) && !handled; i++ {
    handled = caseHandlers[i](tok, tb);
  }
  tb.ensureSize();
}

func (tb *tokBuf) getFilteredToken() common.Token {
  for tok := tb.lx.GetToken(); tok == nil; tok = tb.lx.GetToken() {
    if tok.Type() == common.TOK_COMMENT || tok.Type() == common.TOK_SPACE {
      tok = nil;
    } else {
      return tok;
    }
  }
  return nil;
}

func (tb *tokBuf) ensureSize() {
  for tb.tokBuf.Len() > MAX_TOK {
    if tb.curTok == nil || tb.curTok.Prev() == nil {
      tb.Error("Unable to remove current token from buffer");
    }
    tb.tokBuf.Remove(tb.tokBuf.Front());
  }
}

func handleColon(tok common.Token, tb *tokBuf) bool {
  if tok.Type() == common.TOK_COLON {
    tok2 := tb.getFilteredToken();  // whats behind the colon?
    if tok2.Type() != common.TOK_NL {
      // No TOK_NL: just 2 normal tokens then
      tb.curTok = tb.tokBuf.PushBack(tok);
      tb.tokBuf.PushBack(tok2);
    } else {
      newTok := tb.lx.NewMultiTok(common.TOK_BLOCK_START, []common.Token{tok, tok2});
      tb.curTok = tb.tokBuf.PushBack(newTok);
    }
    return true;
  }
  return false;
}

func handleMultiDedent(tok common.Token, tb *tokBuf) bool {
  if tok.Type() == common.TOK_MULTI_DEDENT {
    dedent := lexer.Tok2multiDedent(tok).Dedent();
    dedent = handleFirstDedent(dedent, tok, tb);

    // handle all other dedentations:
    return true;
  }
  return false;
}

func handleFirstDedent(dedent int, tok common.Token, tb *tokBuf) int {
  ret := 0;
  if dedent & 1 != 0 {
    tb.curTok = tb.tokBuf.PushBack(tb.lx.NewCopyTok(common.TOK_HALF_DEDENT, tok));
    ret = dedent - 1;
  } else {
    tb.curTok = tb.tokBuf.PushBack(tb.lx.NewCopyTok(common.TOK_DEDENT, tok));
    ret = dedent - 2;
  }
  return ret;
}

func handleOtherDedents(dedent int, tok common.Token, tb *tokBuf) {
  for d := 0; dedent > d; d += 2 {
    tb.tokBuf.PushBack(tb.lx.NewCopyTok(common.TOK_DEDENT, tok));
  }
}

func handleDefault(tok common.Token, tb *tokBuf) bool {
  tb.tokBuf.PushBack(tok);
  return true;
}

func (tb *tokBuf) any2token(val interface{}) common.Token {
  switch t := val.(type) {
  case common.Token:
    return t;
  default:
    tb.Error("Not a token type");
  }

  return nil;
}

