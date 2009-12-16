package lexer

import (
  "diamondlang/common";
)


// --------------------------------------------------------------------------
// Token types.
// --------------------------------------------------------------------------

/// SimpleToken - Base class for all tokens.
type SimpleToken struct {
  typ common.TokEnum;  // type of the token
  common.SrcPiece;
}

func (lx *Lexer) newToken(typ common.TokEnum, mark common.SrcMark) *SimpleToken {
  return &SimpleToken{typ, lx.srcBuf.NewPiece(mark)}
}

func (tok *SimpleToken) Type() common.TokEnum { return tok.typ; }

/// HasSpaceAround - return:
///   -1 for front space only
///   +1 for back space only
///    0 for front and back being equal
func (tok *SimpleToken) HasSpaceAround() int {
  line  := tok.WholeLine();
  // look for space before token:
  front := 0;
  col   := tok.Column() - 1;
  if col < 0 || common.IsSpace(line[col]) { front = 1 }

  // look for space after token:
  back := 0;
  col   = tok.Column() + len(tok.Content());
  if col >= len(line) || common.IsSpace(line[col]) { back = 1 }

  return back - front;
}

func (tok *SimpleToken) Error(msg string) {
  common.HandleFatal(common.MakeErrString(msg, tok.Line(), tok.WholeLine(),
      tok.Column(), len(tok.Content()) )
  );
}


// Special EOF token for better printing and easier creation
type EofTok struct {
  *SimpleToken;
}
func (lx *Lexer) newEofTok() *EofTok {
  return &EofTok{lx.newToken(common.TOK_EOF, lx.srcBuf.NewMark())}
}
func (tok *EofTok) String() string { return "<EOF>" }

/// IntTok - Signal an integer constant.
type IntTok struct {
  *SimpleToken;
  value int64;
}
func (lx *Lexer) newIntTok(mark common.SrcMark, val int64) *IntTok {
  tok := lx.newToken(common.TOK_INT, mark);
  return &IntTok{tok, val};
}
func (tok *IntTok) Value() int64 { return tok.value }

/// MultiDedentTok - Signal multiple dedentations.
type MultiDedentTok struct {
  *SimpleToken;
  dedent int;
}
func (lx *Lexer) newMultiDedentTok(dedent int) *MultiDedentTok {
  if (dedent & 1) > 0 || dedent <= 0 { lx.srcBuf.Error("Indentation error"); }

  // dedent is always even; so we don't loose information:
  return &MultiDedentTok{lx.newToken(common.TOK_MULTI_DEDENT, lx.srcBuf.NewMark()),
                         dedent/2};
}
func (tok *MultiDedentTok) Dedent() int { return tok.dedent }

type IdPart struct {
  typ common.TokEnum;
  id  string;
}
func newIdPart(typ common.TokEnum, id string) *IdPart {
  return &IdPart{typ, id};
}

/// IdTok - Signal an ID.
type IdTok struct {
  *SimpleToken;
  parts []*IdPart;
}
func (lx *Lexer) newIdTok(tok *SimpleToken, parts []*IdPart) *IdTok {
  if len(parts) <= 0 { tok.Error("ID has no parts"); }
  return &IdTok{tok, parts};
}
func (tok *IdTok) Parts() []*IdPart { return tok.parts }

