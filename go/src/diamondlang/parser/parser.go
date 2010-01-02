package parser

import (
  "diamondlang/common";
  "diamondlang/lexer";
)


type parser struct {
  tb                 common.TokenBuffer; // our source for tokens
  curTok             common.Token;       // current token
  infixPrecedences   []int;
  halfIndentsAllowed bool;
}

func NewParser(tb common.TokenBuffer) common.Parser {
  p := &parser{tb, nil, infixPrecedences(), false};
  p.fetchNextToken();
  return p;
}

func (p *parser) Error(msg string) {
  p.tb.Error(msg);
}

/// fetchNextToken - Fetch the next meaningful token from the token buffer.
func (p *parser) fetchNextToken() {
  tok := p.tb.GetToken();
  for tok.Type() == common.TOK_SPACE || tok.Type() == common.TOK_COMMENT {
    tok = p.tb.GetToken();
  }
  p.curTok = tok;
}


func infixPrecedences() []int {
  var ip [128]int;
  for c := 'A'; c <= 'Z'; c++ { ip[c] = 10; }
  ip['_'] = 10;  // for protected functions
  ip['='] = 20;
  ip['?'] = 20;
  ip['!'] = 20;
  ip['<'] = 20;
  ip['>'] = 20;
  ip['|'] = 30;
  ip['&'] = 40;
  ip['+'] = 60;
  ip['-'] = 60;
  ip['*'] = 70;
  ip['/'] = 70;
  ip['%'] = 70;
  ip['^'] = 80;
  ip['~'] = 90;
  ip['$'] = 90;
  ip['@'] = 90;  // highest.
  return &ip;
}

/// calleePrecedence - Get the precedence of the pending binary operator token.
func (p *parser) calleePrecedence(callee string) int {
  tokPrec := p.infixPrecedences[callee[0]];
  // Make sure it's a declared binop.
  if tokPrec <= 0 {
    p.Error("Undefined start of a binary operator '" + callee + "'");
  }
  return tokPrec;
}

func (p *parser) ParseNumberExpr() common.ExprAst {
  it := lexer.Token2int(p.curTok);
  pi := new(int64);
  *pi = it.Value();
  p.fetchNextToken(); // consume the number
  return NewLiteralExprAst(it.SourcePiece(), common.TYPE_INT, pi);
}

func (p *parser) ParseCharExpr() common.ExprAst {
  ct := lexer.Token2char(p.curTok);
  pc := new(byte);
  *pc = ct.Value();
  p.fetchNextToken(); // consume the character
  return NewLiteralExprAst(ct.SourcePiece(), common.TYPE_CHAR, pc);
}

func (p *parser) ParseStringExpr() common.ExprAst {
  st := lexer.Token2string(p.curTok);
  ps := new(string);
  *ps = st.Value();
  p.fetchNextToken(); // consume the character
  return NewLiteralExprAst(st.SourcePiece(), common.TYPE_CHAR, ps);
}

func (p *parser) ParseValConstExpr() common.ExprAst {
  ret := common.ExprAst(nil);
  it := lexer.Token2id(p.curTok);
  parts := it.Parts();
  mainPart := parts[0];

  if it.Type() == common.TOK_CONST_ID {
    ret = parseConstExpr(it, parts, mainPart);
  } else {
    ret = parseValExpr(it, parts, mainPart);
  }

  p.fetchNextToken(); // consume the identifier
  return ret;
}

func parseConstExpr(it *lexer.IdTok, parts []*lexer.IdPart, mainPart *lexer.IdPart)
                        common.ExprAst {
  if protectedId(parts) {
    it.Error("Protected constants don't make sense");
  }
  module := "";
  subStart := 1;
  if mainPart.Type() == common.TOK_MODULE_ID {
    module = mainPart.Id();
    mainPart = parts[subStart];
    subStart++;
  }
  return NewConstantExprAst(it.SourcePiece(), module, mainPart.Id(),
                            parts2subs(parts[subStart:len(parts)]));
}

func parseValExpr(it *lexer.IdTok, parts []*lexer.IdPart, mainPart *lexer.IdPart)
                        common.ExprAst {
  if mainPart.Protected() {
    it.Error("Values protected at their first level don't make sense");
  }
  return NewValueExprAst(it.SourcePiece(), mainPart.Id(),
                         parts2subs(parts[1:len(parts)]));
}

func protectedId(parts []*lexer.IdPart) bool {
  for _, part := range parts {
    if part.Protected() { return true; }
  }
  return false;
}

func parts2subs(parts []*lexer.IdPart) []common.SubId {
  subs := make([]common.SubId, len(parts));
  for i, part := range parts {
    subs[i] = common.SubId{part.Id(), common.TYPE_UNKNOWN, part.Protected()};
  }
  return subs;
}

