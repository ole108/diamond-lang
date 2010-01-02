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
  indentLevel int;       // current level of indentation
}

func NewTokenBuffer(lx common.Lexer) common.TokenBuffer {
  return &tokBuf{lx, list.New(), nil, 0};
}

func (tb *tokBuf) Error(msg string) {
  tb.lx.Error(msg);
}

func (tb *tokBuf) GetToken() common.Token {
  if tb.atEnd() {
    tb.readToken();
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
  tok := tb.lx.GetToken();
  caseHandlers := []caseHandler{
    handleSpace, handleEof, handleColon, handleDefault,
  };
  handled := false;
  for i := 0; i < len(caseHandlers) && !handled; i++ {
    handled = caseHandlers[i](tok, tb);
  }
  tb.ensureSize();
}

func (tb *tokBuf) getFilteredToken() common.Token {
  tok := tb.lx.GetToken();
  for ; tok.Type() == common.TOK_COMMENT || tok.Type() == common.TOK_SPACE
      ; tok = tb.lx.GetToken()                                             {

      tb.tokBuf.PushBack(tok);
  }
  return tok;
}

func (tb *tokBuf) ensureSize() {
  for tb.tokBuf.Len() > MAX_TOK {
    if tb.curTok == nil || tb.curTok.Prev() == nil {
      tb.Error("Unable to remove current token from buffer");
    }
    tb.tokBuf.Remove(tb.tokBuf.Front());
  }
}

// special handling of EOF (so we have valid code if possible)
func handleEof(tok common.Token, tb *tokBuf) bool {
  if tok.Type() == common.TOK_EOF {
    if tb.indentLevel > 0 {
      handleIndent(tb.lx.NewSpaceTok(tok, 0, true), tok, tb);
    } else {
      tb.curTok = tb.tokBuf.PushBack(tok);
    }
    return true;
  }
  return false;
}

func handleSpace(tok common.Token, tb *tokBuf) bool {
  if tok.Type() == common.TOK_SPACE {
    spaceTok := tb.token2space(tok);
    if spaceTok.AtStartOfLine() {
      handlePossibleIndent(tok, tb);
    } else {
      tb.curTok = tb.tokBuf.PushBack(tok);
    }
    return true;
  }
  return false;
}

func handlePossibleIndent(tok common.Token, tb *tokBuf) {
  tok2 := tb.lx.GetToken();
  if tok2.Type() == common.TOK_COMMENT ||
     tok2.Type() == common.TOK_NL        {
    tb.curTok = tb.tokBuf.PushBack(tok);
    tb.tokBuf.PushBack(tok2);
  } else if tok2.Type() == common.TOK_EOF       {
    curTok := tb.tokBuf.PushBack(tok);
    handleEof(tok2, tb);
    tb.curTok = curTok;
  } else {
    handleIndent(tok, tok2, tb);
  }
}

func handleIndent(tok common.Token, tok2 common.Token, tb *tokBuf) {
  handled := recordAnyIndent(tb.token2space(tok).Space(), tok, tb);
  if handled {
    tb.tokBuf.PushBack(tok2);
  } else {
    tb.curTok = tb.tokBuf.PushBack(tok2);
  }
}

func recordAnyIndent(spaces int, tok common.Token, tb *tokBuf) bool {
  indent := spaces - tb.indentLevel*2;
  if indent < 0 {
    recordDedent(-indent, tok, tb);
  } else if indent > 0 {
    recordIndent(indent, tok, tb);
  } else {
    return false;
  }
  return true;
}

func recordIndent(indent int, tok common.Token, tb *tokBuf) {
  // we can have half indentations (2 spaces) and
  //             full indentations (4 spaces)
  switch indent {
  case  2:
    tb.indentLevel++;
    tb.curTok = tb.tokBuf.PushBack(tb.lx.NewCopyTok(common.TOK_HALF_INDENT, tok));
  case  4:
    tb.indentLevel += 2;
    tb.curTok = tb.tokBuf.PushBack(tb.lx.NewCopyTok(common.TOK_INDENT, tok));
  default:
    tok.Error("Indentation error");
  }
}

func recordDedent(dedent int, tok common.Token, tb *tokBuf) {
  if dedent & 1 != 0 { tok.Error("Uneven dedentation error"); }
  recordOtherDedents(recordFirstDedent(dedent/2, tok, tb), tok, tb);
}

func recordFirstDedent(dedent int, tok common.Token, tb *tokBuf) int {
  ret := 0;
  if dedent & 1 != 0 {
    tb.indentLevel--;
    tb.curTok = tb.tokBuf.PushBack(tb.lx.NewCopyTok(common.TOK_HALF_DEDENT, tok));
    ret = dedent - 1;
  } else {
    tb.indentLevel -= 2;
    tb.curTok = tb.tokBuf.PushBack(tb.lx.NewCopyTok(common.TOK_DEDENT, tok));
    ret = dedent - 2;
  }
  return ret;
}

func recordOtherDedents(dedent int, tok common.Token, tb *tokBuf) {
  for d := 0; dedent > d; d += 2 {
    tb.indentLevel -= 2;
    tb.tokBuf.PushBack(tb.lx.NewCopyTok(common.TOK_DEDENT, tok));
  }
}

func handleColon(tok common.Token, tb *tokBuf) bool {
  if tok.Type() == common.TOK_COLON {
    curTok := tb.tokBuf.PushBack(tok);
    tok2 := tb.getFilteredToken();  // whats behind the colon?
    if tok2.Type() != common.TOK_NL {
      // No TOK_NL: just 2 normal tokens then
      tb.tokBuf.PushBack(tok2);
    } else {
      curTok.Value = tb.lx.NewAnyTok(common.TOK_BLOCK_START,
          tok.SourcePiece().Start(), tok2.SourcePiece().End());
    }
    tb.curTok = curTok;
    return true;
  }
  return false;
}

func handleDefault(tok common.Token, tb *tokBuf) bool {
  tb.curTok = tb.tokBuf.PushBack(tok);
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

func (tb *tokBuf) token2space(tok common.Token) *lexer.SpaceTok {
  switch t := tok.(type) {
  case *lexer.SpaceTok:
    return t;
  default:
    tb.Error("Not a space token");
  }

  return nil;
}

