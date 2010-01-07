package parser

import (
  "diamondlang/common";
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
  ip['$'] = 90;  // highest.
  return &ip;
}



/// infixPrecedence - Get the precedence of the pending binary operator token.
func (p *parser) infixPrecedence(operator string) int {
  tokPrec := p.infixPrecedences[operator[0]];
  // Make sure it's a declared binop.
  if tokPrec <= 0 {
    p.Error("Undefined start of a binary operator '" + operator + "'");
  }
  return tokPrec;
}


