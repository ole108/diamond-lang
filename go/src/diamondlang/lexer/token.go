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

func Token2simple(tok common.Token) *SimpleToken {
  st, ok := tok.(*SimpleToken);
  if !ok { tok.Error("Not a simple token"); }
  return st;
}

func (lx *Lexer) newToken(typ common.TokEnum, mark common.SrcMark) *SimpleToken {
  return &SimpleToken{typ, lx.srcBuf.NewPiece(mark)};
}

func (lx *Lexer) NewCopyTok(typ common.TokEnum, tok common.Token) common.Token {
  return &SimpleToken{typ, tok.SourcePiece()};
}

func (lx *Lexer) NewAnyTok(typ common.TokEnum, start common.SrcMark, end common.SrcMark) common.Token {
  return &SimpleToken{typ, lx.srcBuf.NewAnyPiece(start, end)};
}

func (tok *SimpleToken) Type() common.TokEnum { return tok.typ; }
func (tok *SimpleToken) SourcePiece() common.SrcPiece { return tok.SrcPiece; }
func (tok *SimpleToken) String() string {
  return tok.typ.String() + ": `" + tok.SrcPiece.String() + "`";
}

/// HasSpaceAround - return:
///   -1 for front space only
///   +1 for back space only
///    0 for front and back being equal
func (tok *SimpleToken) HasSpaceAround() int {
  line  := tok.WholeLine();
  // look for space before token:
  front := 0;
  col   := tok.StartColumn() - 1;
  if col < 0 || common.IsSpace(line[col]) { front = 1 }

  // look for space after token:
  back := 0;
  col   = tok.StartColumn() + len(tok.Content());
  if col >= len(line) || common.IsSpace(line[col]) { back = 1 }

  return back - front;
}


// Special EOF token for better printing and easier creation
type EofTok struct {
  *SimpleToken;
}
func Token2eof(tok common.Token) *EofTok {
  et, ok := tok.(*EofTok);
  if !ok { tok.Error("Not an EOF token"); }
  return et;
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
func Token2int(tok common.Token) *IntTok {
  it, ok := tok.(*IntTok);
  if !ok { tok.Error("Not an integer token"); }
  return it;
}
func (lx *Lexer) newIntTok(mark common.SrcMark, val int64) *IntTok {
  tok := lx.newToken(common.TOK_INT, mark);
  return &IntTok{tok, val};
}
func (tok *IntTok) Value() int64 { return tok.value }

/// CharTok - Signal a character constant.
type CharTok struct {
  *SimpleToken;
  value byte;
}
func Token2char(tok common.Token) *CharTok {
  ct, ok := tok.(*CharTok);
  if !ok { tok.Error("Not a character token"); }
  return ct;
}
func (lx *Lexer) newCharTok(mark common.SrcMark, val byte) *CharTok {
  tok := lx.newToken(common.TOK_CHAR, mark);
  return &CharTok{tok, val};
}
func (tok *CharTok) Value() byte { return tok.value }

/// StringTok - Signal a string constant.
type StringTok struct {
  *SimpleToken;
  value string;
}
func Token2string(tok common.Token) *StringTok {
  st, ok := tok.(*StringTok);
  if !ok { tok.Error("Not a string token"); }
  return st;
}
func (lx *Lexer) newStringTok(mark common.SrcMark, val string) *StringTok {
  tok := &SimpleToken{common.TOK_STR, lx.srcBuf.NewPiece(mark)};
  return &StringTok{tok, val};
}
func (tok *StringTok) Value() string { return tok.value }

/// SpaceTok - Signal some space.
type SpaceTok struct {
  *SimpleToken;
  space int;
  atStartOfLine bool;
}
func Token2space(tok common.Token) *SpaceTok {
  st, ok := tok.(*SpaceTok);
  if !ok { tok.Error("Not a space token"); }
  return st;
}
func (lx *Lexer) newSpaceTok(mark common.SrcMark, space int, atStartOfLine bool) *SpaceTok {
  tok := lx.newToken(common.TOK_SPACE, mark);
  return &SpaceTok{tok, space, atStartOfLine};
}
func (lx *Lexer) NewSpaceTok(tok common.Token, space int, atStartOfLine bool) common.Token {
  newTok := Token2simple(lx.NewCopyTok(common.TOK_SPACE, tok));
  return &SpaceTok{newTok, space, atStartOfLine};
}
func (tok *SpaceTok) Space() int { return tok.space }
func (tok *SpaceTok) AtStartOfLine() bool { return tok.atStartOfLine }


type IdPart struct {
  typ         common.TokEnum;
  id          string;
  protected   bool;
}
func newIdPart(typ common.TokEnum, id string, protected bool) *IdPart {
  return &IdPart{typ, id, protected};
}
func (part *IdPart) Type() common.TokEnum { return part.typ }
func (part *IdPart) Id() string { return part.id }
func (part *IdPart) Protected() bool { return part.protected }

/// IdTok - Signal an ID.
type IdTok struct {
  *SimpleToken;
  parts       []*IdPart;
  halfApplied bool;
}
func Token2id(tok common.Token) *IdTok {
  it, ok := tok.(*IdTok);
  if !ok { tok.Error("Not an ID token"); }
  return it;
}
func (lx *Lexer) newIdTok(typ common.TokEnum, piece common.SrcPiece,
                          parts []*IdPart, halfApplied bool) *IdTok {
  if len(parts) <= 0 { piece.Error("ID has no parts"); }
  return &IdTok{&SimpleToken{typ, piece}, parts, halfApplied};
}
func (tok *IdTok) Parts() []*IdPart { return tok.parts }
func (tok *IdTok) HalfApplied() bool { return tok.halfApplied }

/// OperatorTok - Signal an operator.
type OperatorTok struct {
  *SimpleToken;
  halfApplied bool;
}
func Token2operator(tok common.Token) *OperatorTok {
  ot, ok := tok.(*OperatorTok);
  if !ok { tok.Error("Not an operator token"); }
  return ot;
}
func (lx *Lexer) newOperatorTok(mark common.SrcMark, halfApplied bool) *OperatorTok {
  return &OperatorTok{lx.newToken(common.TOK_OP_ID, mark), halfApplied};
}
func (tok *OperatorTok) HalfApplied() bool { return tok.halfApplied }

