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
  return &SimpleToken{typ, lx.srcBuf.NewPiece(mark)};
}

func (lx *Lexer) NewCopyTok(typ common.TokEnum, tok common.Token) common.Token {
  return &SimpleToken{typ, lx.srcBuf.NewMultiPiece([]common.SrcPiece{tok})};
}

func (lx *Lexer) NewMultiTok(typ common.TokEnum, toks []common.Token) common.Token {
  pieces := make([]common.SrcPiece, len(toks));
  for i := 0; i < len(toks); i++ {
    pieces[i] = toks[i];
  }
  return &SimpleToken{typ, lx.srcBuf.NewMultiPiece(pieces)};
}

func (tok *SimpleToken) Type() common.TokEnum { return tok.typ; }
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

func (lx *Lexer) token2simple(tok common.Token) *SimpleToken {
  switch t := tok.(type) {
  case *SimpleToken:
    return t;
  default:
    lx.Error("Not a simple token");
  }

  return nil;
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

/// CharTok - Signal a character constant.
type CharTok struct {
  *SimpleToken;
  value byte;
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
func (lx *Lexer) newStringTok(mark common.MultiLineSrcMark, val string) *StringTok {
  tok := &SimpleToken{common.TOK_STR, lx.srcBuf.NewMultiLinePiece(mark)};
  return &StringTok{tok, val};
}
func (tok *StringTok) Value() string { return tok.value }

/// SpaceTok - Signal some space.
type SpaceTok struct {
  *SimpleToken;
  space int;
  atStartOfLine bool;
}
func (lx *Lexer) newSpaceTok(mark common.SrcMark, space int, atStartOfLine bool) *SpaceTok {
  tok := lx.newToken(common.TOK_SPACE, mark);
  return &SpaceTok{tok, space, atStartOfLine};
}
func (lx *Lexer) NewSpaceTok(tok common.Token, space int, atStartOfLine bool) common.Token {
  newTok := lx.token2simple(lx.NewCopyTok(common.TOK_SPACE, tok));
  return &SpaceTok{newTok, space, atStartOfLine};
}
func (tok *SpaceTok) Space() int { return tok.space }
func (tok *SpaceTok) AtStartOfLine() bool { return tok.atStartOfLine }


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

