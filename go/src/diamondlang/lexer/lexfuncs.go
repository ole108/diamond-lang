package lexer

import (
  "diamondlang/common";
  "strings";
  "strconv";
//"fmt";
)


type lexFunc func(*Lexer) (common.Token, bool)

func trySpace(lx *Lexer) (tok common.Token, moved bool) {
  atStart := lx.srcBuf.AtStartOfLine() && lx.inParens <= 0;
  if common.IsSpace(lx.curChar) || atStart {
    mark := lx.srcBuf.NewMark();
    tok, moved = lx.newSpaceTok(mark, countSpaces(lx), atStart), true;
    lx.srcBuf.NotAtStartOfLine();
  }
  return;
}

func countSpaces(lx *Lexer) int {
  spaces := 0;
  for ; common.IsSpace(lx.curChar); lx.nextChar() {
    spaces += common.SpaceAmount(lx.curChar);
  }
  return spaces;
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
    lx.Error("Too deeply nested parentheses");
  }
  mark := lx.srcBuf.NewMark();
  lx.parenStack[lx.inParens] = lx.curChar;
  lx.inParens++;
  lx.nextChar();
  return lx.newToken(common.TOK_PAREN_OPEN, mark);
}

func getParenClose(lx *Lexer) common.Token {
  if lx.inParens <= 0 {
    lx.Error("Too many closing parentheses");
  }
  mark := lx.srcBuf.NewMark();
  lx.inParens--;
  if lx.parenStack[lx.inParens] != openParen(lx.curChar, lx) {
    lx.Error("Parentheses don't fit together");
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
    if err != nil { lx.Error(err.String()) }
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
  if err != nil { lx.Error(err.String()) }
  if val > 36 || val < 2 { lx.Error("Invalid integer base") }
  return val;
}

func tryChar(lx *Lexer) (tok common.Token, moved bool) {
  if lx.curChar == '\'' {
    mark := lx.srcBuf.NewMark();
    lx.nextChar();
    char := readEscChar(lx);
    if lx.curChar != '\'' { lx.Error("Invalid character token"); }
    lx.nextChar();
    tok, moved = lx.newCharTok(mark, char), true;
  }
  return;
}

func readEscChar(lx *Lexer) byte {
  escaped := false;
  if lx.curChar == '\\' {
    escaped = true;
    lx.nextChar();
  }
  char := escaped2char(escaped, lx.curChar, lx);
  lx.nextChar();
  return char;
}

func escaped2char(escaped bool, char byte, lx *Lexer) byte {
  ret := char;
  if escaped && (isDigit(char) || isLower(char)) {
    switch char {
    case '0':  ret = 0;
    case 'a':  ret = '\a';
    case 'b':  ret = '\b';
    case 'd':  ret = 127;   // DEL
    case 'e':  ret = 27;    // ESC
    case 'f':  ret = '\f';
    case 'n':  ret = '\n';
    case 'r':  ret = '\r';
    case 't':  ret = '\t';
    case 'v':  ret = '\v';
    default:   lx.Error("Invalid escape character");
    }
  }
  return ret;
}

func tryString(lx *Lexer) (tok common.Token, moved bool) {
  if lx.curChar == '"' {
    mark := lx.srcBuf.NewMultiLineMark();
    str := readString(lx.curChar, lx, readEscString);
    tok, moved = lx.newStringTok(mark, str), true;
  } else if lx.curChar == '`' {
    mark := lx.srcBuf.NewMultiLineMark();
    str := readString(lx.curChar, lx, readRawString);
    tok, moved = lx.newStringTok(mark, str), true;
  }
  return;
}

func readString(delim byte, lx *Lexer, readTypString func(*Lexer,byte,int)string) string {
  cnt := readCharCount(lx.curChar, lx, 9);
  ret := "";
  if cnt == 1 {
    ret = readTypString(lx, delim, 1);
  } else if cnt == 2 {
    ret = "";
  } else if cnt == 3 {
    ret = readTypString(lx, delim, 3);
  } else if cnt < 6 {
    ret = strings.Repeat(string(delim), cnt - 3) + readTypString(lx, delim, 3);
  } else if cnt == 6 {
    ret = "";
  } else if cnt < 9 {
    ret = strings.Repeat(string(delim), cnt - 6);
  } else {
    lx.Error("Illegal number of consecutive string delimiters");
  }
  return ret;
}

func readRawString(lx *Lexer, delim byte, max int) string {
  ret := "";
  cnt := 0;
  for cnt < max && lx.curChar != common.EOF {
    for lx.curChar != delim && lx.curChar != common.EOF {
      ret += string(lx.curChar);
      lx.nextChar();
    }
    cnt = readCharCount(delim, lx, max);
    if cnt < max { ret += strings.Repeat(string(delim), cnt); }
  }
  return ret;
}

func readEscString(lx *Lexer, delim byte, max int) string {
  ret := "";
  cnt := 0;
  for cnt < max && lx.curChar != common.EOF {
    for lx.curChar != delim && lx.curChar != common.EOF {
      if max <= 1 && (lx.curChar == '\n' || lx.curChar == '\r') {
        lx.Error("Simple strings can't span multiple lines");
      }
      ret += string(readEscChar(lx));
    }
    cnt = readCharCount(delim, lx, max);
    if cnt < max { ret += strings.Repeat(string(delim), cnt); }
  }
  return ret;
}

func readCharCount(char byte, lx *Lexer, max int) int {
  ret := 0;
  for lx.curChar == char && ret < max {
    ret++;
    lx.nextChar();
  }
  return ret;
}

func tryNewLine(lx *Lexer) (tok common.Token, moved bool) {
  if lx.curChar == '\n' || lx.curChar == '\r' {
    mark := lx.srcBuf.NewMultiLineMark();
    readNewLine(lx);
    tok, moved = makeNewLineTok(lx, mark), true;
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
    mark := lx.srcBuf.NewMultiLineMark();
    for ; lx.curChar == ';'; lx.nextChar() { }
    moved = true;
    tok = makeNewLineTok(lx, mark);
  }
  return;
}

func makeNewLineTok(lx *Lexer, mark common.MultiLineSrcMark) common.Token {
  tok := common.Token(nil);
  if lx.inParens <= 0 {
    tok = &SimpleToken{common.TOK_NL, lx.srcBuf.NewMultiLinePiece(mark)};
  }
  return tok;
}

func tryComment(lx *Lexer) (tok common.Token, moved bool) {
  if lx.curChar == '#' {
    mark := lx.srcBuf.NewMark();
    for ; lx.curChar != '\n' && lx.curChar != '\r' && lx.curChar != common.EOF; lx.nextChar() { }
    tok, moved = lx.newToken(common.TOK_COMMENT, mark), true;
  }
  return;
}

func tryEof(lx *Lexer) (tok common.Token, moved bool) {
  if lx.curChar == common.EOF {
    tok, moved = lx.newEofTok(), true;
  }

  return;
}

func signalUndefined(lx *Lexer) (tok common.Token, moved bool) {
  lx.Error("Unknown token");
  return nil, false;
}

