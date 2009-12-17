package lexer

import (
  "diamondlang/common";
  "strings";
  "strconv";
//"fmt";
)


type lexFunc func(*Lexer) (common.Token, bool)

func tryIndent(lx *Lexer) (tok common.Token, moved bool) {
  if lx.srcBuf.AtStartOfLine() && lx.inParens <= 0 {
    mark := lx.srcBuf.NewMark();
    // count spaces:
    spaces := 0;
    for ; common.IsSpace(lx.curChar); lx.nextChar() {
      spaces += common.SpaceAmount(lx.curChar);
    }

    // ignore comment and empty lines
    // (EOF will be handled later):
    _, mov1 := skipComment(lx);
    _, mov2 := tryNewLine(lx);

    // did we move in any way?
    if spaces > 0 || mov1 || mov2 { moved = true; }

    if !mov1 && !mov2 {
      // return indent token
      tok = spaces2tok(lx, mark, spaces);
      if tok != nil {
        moved = true;
        lx.srcBuf.NotAtStartOfLine();
      }
    }
  }

  return;
}

func spaces2tok(lx *Lexer, mark common.SrcMark, spaces int) common.Token {
  tok := common.Token(nil);

  indent := spaces - lx.indentLevel*2;
  if indent < 0 {
    tok = lx.newMultiDedentTok(-indent);
    lx.indentLevel += indent/2;
  } else {
    tok = indent2tok(lx, mark, indent);
  }

  return tok;
}

func indent2tok(lx *Lexer, mark common.SrcMark, indent int) common.Token {
  tok := common.Token(nil);

  // we can have half indentations (2 spaces) and
  //             full indentations (4 spaces)
  switch indent {
  case  0:
    tok = nil;  // no change in indentation
  case  2:
    lx.indentLevel++;
    tok = lx.newToken(common.TOK_HALF_INDENT, mark);
  case  4:
    lx.indentLevel+=2;
    tok = lx.newToken(common.TOK_INDENT, mark);
  default:
    lx.srcBuf.Error("Indentation error");
  }

  return tok;
}

func tryId(lx *Lexer) (tok common.Token, moved bool) {
  if isAlpha(lx.curChar) {
    fullId := readFullId(lx);
    parts := fullId2parts(fullId);
    fullId.typ = setIdTypes(parts, fullId);
    tok, moved = lx.newIdTok(fullId, parts), true;
  }
  return;
}

func readFullId(lx *Lexer) *SimpleToken {
  mark := lx.srcBuf.NewMark();
  for ; isIdChar(lx.curChar); lx.nextChar() { }
  return lx.newToken(common.TOK_VAL_ID, mark);
}

func fullId2parts(fullId common.Token) []*IdPart {
  strParts := strings.Split(fullId.Content(), ".", 0);
  idParts := make([]*IdPart, len(strParts));
  for i, s := range strParts {
    if len(s) <= 0 {
      fullId.Error("Illegal identifier");
    }
    idParts[i] = newIdPart(common.TOK_MODULE_ID, s);
  }
  return idParts;
}

func setIdTypes(parts []*IdPart, tok common.Token) common.TokEnum {
  for _, part := range parts {
    part.typ = getIdType(part.id, tok);
  }

  var tokTyp common.TokEnum = common.TOK_MODULE_ID;
  for i, part := range parts {
    switch {
    case part.typ == common.TOK_CONST_ID:
      if i <= 1 && tokTyp == common.TOK_MODULE_ID {
        tokTyp = common.TOK_CONST_ID;
      } else {
        tok.Error("Illegal constant identifier part");
      }
    case tokTyp == common.TOK_CONST_ID:
      if part.typ != common.TOK_MODULE_ID && part.typ != common.TOK_VAL_ID {
        tok.Error("Illegal constant identifier part");
      }
      part.typ = common.TOK_VAL_ID;
    case part.typ == common.TOK_MODULE_ID:
      if tokTyp == common.TOK_FUNC_ID {
        tok.Error("Illegal value after function identifier part");
      }
      if i >= 1 {
        part.typ = common.TOK_VAL_ID;
        tokTyp = common.TOK_VAL_ID;
      }
    case part.typ == common.TOK_FUNC_ID:
      if i > 1 || tokTyp != common.TOK_MODULE_ID {
        tok.Error("Illegal function identifier");
      }
      tokTyp = common.TOK_FUNC_ID;
    case part.typ == common.TOK_VAL_ID:
      if tokTyp == common.TOK_FUNC_ID {
        tok.Error("Illegal value after function identifier part");
      }
      if tokTyp == common.TOK_MODULE_ID {
        tokTyp = common.TOK_VAL_ID;
      }
    default:
      tok.Error("Illegal identifier part");
    }
  }

  return tokTyp;
}

func getIdType(id string, tok common.Token) common.TokEnum {
  // ordinary flags:
  gotUpper := false;
  gotLower := false;
  gotUnder := false;
  got2Uppr := false;

  // first (and sometimes second) character are special:
  i := 0;
  firstUnder := (id[i] == '_');
  if firstUnder { i++; }
  firstUpper := isUpper(id[i]);
  firstLower := isLower(id[i]);
  i++;

  if !firstUpper && !firstLower {
    tok.Error("Illegal start of identifier part");
  }

  // set flags:
  lastUpper := firstUpper;
  for ; i < len(id); i++ {
    b := id[i];    // SrcBuf guaranties 7 bit clean!
    switch {
    case isUpper(b):
      gotUpper = true;
      if lastUpper { got2Uppr = true; }
      lastUpper = true;
    case isLower(b):
      gotLower = true;
      lastUpper = false;
    case isDigit(b):
      // digits are allowed in any ID
      lastUpper = false;
    case b == '_':
      gotUnder = true;
      lastUpper = false;
    default:
      tok.Error("Illegal character in identifier part");
    }
  }

  // evaluate flags:
  var typ common.TokEnum;
  switch {
  case !gotUpper && firstLower && !gotUnder && !firstUnder:
    typ = common.TOK_MODULE_ID;
  case firstUpper && !gotLower && !firstUnder:
    typ = common.TOK_CONST_ID;
  case firstLower && !gotUnder && !got2Uppr:
    typ = common.TOK_VAL_ID;
  case firstUpper && gotLower && !gotUnder && !got2Uppr:
    typ = common.TOK_FUNC_ID;
  default:
//  fmt.Println("firstUpper:", firstUpper, ", gotLower:", gotLower, ", gotUnder:", gotUnder, ", got2Uppr:", got2Uppr);
    tok.Error("Illegal identifier part");
  }

  return typ;
}

func tryOperator(lx *Lexer) (tok common.Token, moved bool) {
  if isOpChar(lx.curChar) {
    mark := lx.srcBuf.NewMark();
    for ; isOpChar(lx.curChar); lx.nextChar() { }
    tok, moved = lx.newToken(common.TOK_OP_ID, mark), true;
  }
  return;
}

func tryParen(lx *Lexer) (tok common.Token, moved bool) {
  if lx.curChar == '(' || lx.curChar == '[' || lx.curChar == '{' {
    tok, moved = getParenOpen(lx), true;
  } else if lx.curChar == ')' || lx.curChar == ']' || lx.curChar == '}' {
    tok, moved = getParenClose(lx), true;
  }
  return;
}

func getParenOpen(lx *Lexer) common.Token {
  if lx.inParens >= len(lx.parenStack) {
    lx.srcBuf.Error("Too deeply nested parentheses");
  }
  mark := lx.srcBuf.NewMark();
  lx.parenStack[lx.inParens] = lx.curChar;
  lx.inParens++;
  lx.nextChar();
  return lx.newToken(common.TOK_PAREN_OPEN, mark);
}

func getParenClose(lx *Lexer) common.Token {
  if lx.inParens <= 0 {
    lx.srcBuf.Error("Too many closing parentheses");
  }
  mark := lx.srcBuf.NewMark();
  lx.inParens--;
  if lx.parenStack[lx.inParens] != openParen(lx.curChar, lx) {
    lx.srcBuf.Error("Parentheses don't fit together");
  }
  lx.nextChar();
  return lx.newToken(common.TOK_PAREN_CLOSE, mark);
}

func tryNumber(lx *Lexer) (tok common.Token, moved bool) {
  if isDigit(lx.curChar) {
    mark := lx.srcBuf.NewMark();
    tok, moved = lx.newIntTok(mark, readInt(lx)), true;
  }
  return;
}

func readInt(lx *Lexer) int64 {
  base := 10;
  strVal := "";

  if lx.curChar == '0' {
    lx.nextChar();
    base = readBase(lx);
  }
  for ; isNumChar(lx.curChar, base); lx.nextChar() {
    if (lx.curChar != '_') { strVal += string(lx.curChar); }
  }
  var ret int64 = 0;
  if len(strVal) > 0 {
    val, err := strconv.Btoi64(strVal, base);
    if err != nil { lx.srcBuf.Error(err.String()) }
    ret = val;
  }
  return ret;
}

func readBase(lx *Lexer) int {
  base := 10;
  switch lx.curChar {
  case 'b':
    lx.nextChar();
    base = 2;
  case 'c':
    lx.nextChar();
    base = 8;
  case 'x':
    lx.nextChar();
    base = 16;
  case 'r':
    lx.nextChar();
    base = readExplicitBase(lx);
  }
  return base;
}

func readExplicitBase(lx *Lexer) int {
  strVal := "";
  for ; isDigit(lx.curChar); lx.nextChar() {
    strVal += string(lx.curChar);
  }
  val, err := strconv.Atoi(strVal);
  if err != nil { lx.srcBuf.Error(err.String()) }
  if val > 36 || val < 2 { lx.srcBuf.Error("Invalid integer base") }
  return val;
}

func tryChar(lx *Lexer) (tok common.Token, moved bool) {
  if lx.curChar == '\'' {
    mark := lx.srcBuf.NewMark();
    lx.nextChar();
    if lx.curChar == '\\' { lx.nextChar(); }
    lx.nextChar();
    if lx.curChar != '\'' { lx.srcBuf.Error("Invalid character token"); }
    tok, moved = lx.newToken(common.TOK_CHAR, mark), true;
  }
  return;
}

func tryNewLine(lx *Lexer) (tok common.Token, moved bool) {
  if lx.curChar == '\n' || lx.curChar == '\r' {
    mark := lx.srcBuf.NewMark();
    readNewLine(lx);
    moved = true;
    tok = makeNewLineTok(lx, mark);
  }
  return;
}

func readNewLine(lx *Lexer) {
  oldChar := lx.curChar;
  lx.nextChar();
  if oldChar == '\r' && lx.curChar == '\n' { lx.nextChar(); }
}

func tryColon(lx *Lexer) (tok common.Token, moved bool) {
  if lx.curChar == ':' {
    mark := lx.srcBuf.NewMark();
    lx.nextChar();
    tok, moved = lx.newToken(common.TOK_COLON, mark), true;
  }
  return;
}

func trySemicolon(lx *Lexer) (tok common.Token, moved bool) {
  if lx.curChar == ';' {
    mark := lx.srcBuf.NewMark();
    for ; lx.curChar == ';'; lx.nextChar() { }
    moved = true;
    tok = makeNewLineTok(lx, mark);
  }
  return;
}

func makeNewLineTok(lx *Lexer, mark common.SrcMark) common.Token {
  tok := common.Token(nil);
  if lx.inParens <= 0 {
    tok = lx.newToken(common.TOK_NL, mark);
  }
  return tok;
}

func skipComment(lx *Lexer) (tok common.Token, moved bool) {
  if lx.curChar == '#' {
    for ; lx.curChar != '\n' && lx.curChar != '\r' && lx.curChar != common.EOF; lx.nextChar() { }
    moved = true;
  }
  return;
}

// special handling of EOF (so we have valid code if possible)
func tryEof(lx *Lexer) (tok common.Token, moved bool) {
  if lx.curChar == common.EOF {
    moved = true;
    if lx.indentLevel <= 0 {
      tok = lx.newEofTok();
    } else {
      tok = lx.newMultiDedentTok(lx.indentLevel*2);
      lx.indentLevel = 0;
    }
  }

  return;
}

func skipSpaces(lx *Lexer) (tok common.Token, moved bool) {
  for ; common.IsSpace(lx.curChar); lx.nextChar() {
    moved = true;
  }
  return;
}

func signalUndefined(lx *Lexer) (tok common.Token, moved bool) {
  lx.srcBuf.Error("Unknown token");
  return nil, false;
}

