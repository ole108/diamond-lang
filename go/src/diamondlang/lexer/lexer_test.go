package lexer

import (
  "testing";
  "diamondlang/common";
  "diamondlang/srcbuf";
  "strings";
)

type tstTok struct {
  typ          common.TokEnum;
  content      string;
  checkContent bool;
  numVal       int64;
  strVal       string;
}

func TestComments(t *testing.T) {
  testStr := `#bla 0b0110
BLA:

 # nope
TRUE`;

  testToks := []*tstTok{
    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_COMMENT, "#bla 0b0110", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_CONST_ID, "BLA", true, 0, ""},
    &tstTok{common.TOK_COLON, ":", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, " ", true, 1001, ""},
    &tstTok{common.TOK_COMMENT, "# nope", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_CONST_ID, "TRUE", true, 0, ""},
  };

  testStringVsTokens(t, testStr, testToks);
}

func TestNumsParensOps(t *testing.T) {
  testStr := `(E = 2*3)&[{TRUE?(0x01_F)+PI_HOCH_2 ^  001230};0c17 * 0b0110] -/- 0r036_10`;

  testToks := []*tstTok{
    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_PAREN_OPEN, "(", true, 0, ""},
    &tstTok{common.TOK_CONST_ID, "E", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_OP_ID, "=", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_INT, "2", false, 2, ""},
    &tstTok{common.TOK_OP_ID, "*", true, 0, ""},
    &tstTok{common.TOK_INT, "3", false, 3, ""},
    &tstTok{common.TOK_PAREN_CLOSE, ")", true, 0, ""},
    &tstTok{common.TOK_OP_ID, "&", true, 0, ""},
    &tstTok{common.TOK_PAREN_OPEN, "[", true, 0, ""},
    &tstTok{common.TOK_PAREN_OPEN, "{", true, 0, ""},
    &tstTok{common.TOK_CONST_ID, "TRUE", true, 0, ""},
    &tstTok{common.TOK_OP_ID, "?", true, 0, ""},
    &tstTok{common.TOK_PAREN_OPEN, "(", true, 0, ""},
    &tstTok{common.TOK_INT, "0x01_F", false, 31, ""},
    &tstTok{common.TOK_PAREN_CLOSE, ")", true, 0, ""},
    &tstTok{common.TOK_OP_ID, "+", true, 0, ""},
    &tstTok{common.TOK_CONST_ID, "PI_HOCH_2", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_OP_ID, "^", true, 0, ""},
    &tstTok{common.TOK_SPACE, "  ", true, 2, ""},
    &tstTok{common.TOK_INT, "001230", false, 1230, ""},
    &tstTok{common.TOK_PAREN_CLOSE, "}", true, 0, ""},
    &tstTok{common.TOK_INT, "0c17", false, 15, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_OP_ID, "*", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_INT, "0b0110", false, 6, ""},
    &tstTok{common.TOK_PAREN_CLOSE, "]", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_OP_ID, "-/-", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_INT, "0r036_10", false, 36, ""},
  };

  testStringVsTokens(t, testStr, testToks);
}

func TestIdsIndent(t *testing.T) {
  testStr := `If bla > 0:
    mod.Func mod.CONST mod.CONST.val Fn i
  Elif bla < 0:
    mod.FuncAli bla
  Else:
    bla.val  # should work!
bla = 0
   # geschafft!`;

  testToks := []*tstTok{
    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_FUNC_ID, "If", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_MODULE_ID, "bla", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_OP_ID, ">", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_INT, "0", false, 0, ""},
    &tstTok{common.TOK_COLON, ":", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "    ", true, 1004, ""},
    &tstTok{common.TOK_FUNC_ID, "mod.Func", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_CONST_ID, "mod.CONST", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_CONST_ID, "mod.CONST.val", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_FUNC_ID, "Fn", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_MODULE_ID, "i", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "  ", true, 1002, ""},
    &tstTok{common.TOK_FUNC_ID, "Elif", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_MODULE_ID, "bla", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_OP_ID, "<", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_INT, "0", false, 0, ""},
    &tstTok{common.TOK_COLON, ":", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "    ", true, 1004, ""},
    &tstTok{common.TOK_FUNC_ID, "mod.FuncAli", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_MODULE_ID, "bla", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "  ", true, 1002, ""},
    &tstTok{common.TOK_FUNC_ID, "Else", true, 0, ""},
    &tstTok{common.TOK_COLON, ":", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "    ", true, 1004, ""},
    &tstTok{common.TOK_VAL_ID, "bla.val", true, 0, ""},
    &tstTok{common.TOK_SPACE, "  ", true, 2, ""},
    &tstTok{common.TOK_COMMENT, "# should work!", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_MODULE_ID, "bla", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_OP_ID, "=", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_INT, "0", false, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "   ", true, 1003, ""},
    &tstTok{common.TOK_COMMENT, "# geschafft!", true, 0, ""},
  };

  testStringVsTokens(t, testStr, testToks);
}

func TestCharsStrings(t *testing.T) {
  testStr := "'c'    'h' '\\t' '\\a' '\\r' '\\n' \n"
             "\"bla\\r\\n\" \"blue\\t\\a\\0\\n\" \n"
             "```\r\n dia\\nbla``` \n"
             "\"\"\"\r\n dia\\nbla\"\"\" \n"
             "```````` \n"
             "\"\"\"\"\"\"\"\"" ;

  testToks := []*tstTok{
    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_CHAR, "'c'", true, int64('c'), ""},
    &tstTok{common.TOK_SPACE, "    ", true, 4, ""},
    &tstTok{common.TOK_CHAR, "'h'", true, int64('h'), ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_CHAR, "'\\t'", true, int64('\t'), ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_CHAR, "'\\a'", true, int64('\a'), ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_CHAR, "'\\r'", true, int64('\r'), ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_CHAR, "'\\n'", true, int64('\n'), ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_STR, "\"bla\\r\\n\"", true, 0, "bla\r\n"},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_STR, "\"blue\\t\\a\\0\\n\"", true, 0, "blue\t\a\x00\n"},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_STR, "```\r\n dia\\nbla```", true, 0, "\r\n dia\\nbla"},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_STR, "\"\"\"\r\n dia\\nbla\"\"\"", true, 0, "\r\n dia\nbla"},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_STR, "````````", true, 0, "``"},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "", true, 1000, ""},
    &tstTok{common.TOK_STR, "\"\"\"\"\"\"\"\"", true, 0, "\"\""},
  };

  testStringVsTokens(t, testStr, testToks);
}

func testStringVsTokens(t *testing.T, str string, toks []*tstTok) {
  lx := NewLexer(srcbuf.NewSourceFromBuffer(strings.Bytes(str), "TestTokens"));
  var tok common.Token;
  var i   int;
  for tok, i = lx.GetToken(), 0; tok.Type() != common.TOK_EOF && i < len(toks);
      tok, i = lx.GetToken(), i+1 {
    switch typ := tok.(type) {
    case *IntTok:
      if toks[i].typ != typ.Type() {
        t.Errorf("Expected token type %v, but got: %v.\n", toks[i].typ, typ.Type());
      }
      if toks[i].numVal != typ.Value() {
        t.Errorf("Expected value %v, but got: %v.\n", toks[i].numVal, typ.Value());
      }
    case *CharTok:
      if toks[i].typ != typ.Type() {
        t.Errorf("Expected token type %v, but got: %v.\n", toks[i].typ, typ.Type());
      }
      if toks[i].numVal != int64(typ.Value()) {
        t.Errorf("Expected value %v, but got: %v.\n", toks[i].numVal, typ.Value());
      }
    case *StringTok:
      if toks[i].typ != typ.Type() {
        t.Errorf("Expected token type %v, but got: %v.\n", toks[i].typ, typ.Type());
      }
      if toks[i].strVal != typ.Value() {
        t.Errorf("Expected value %v, but got: %v.\n", toks[i].strVal, typ.Value());
      }
    case *SpaceTok:
      if toks[i].typ != typ.Type() {
        t.Errorf("Expected token type %v, but got: %v.\n", toks[i].typ, typ.Type());
      }
      space := int64(typ.Space());
      if typ.AtStartOfLine() { space += 1000; }
      if toks[i].numVal != space {
        t.Errorf("Expected space %v, but got: %v.\n", toks[i].numVal, space);
      }
    default:
      if toks[i].typ != typ.Type() {
        t.Errorf("Expected token type %v, but got: %v.\n", toks[i].typ, typ.Type());
      }
      if toks[i].checkContent && toks[i].content != typ.Content() {
        t.Errorf("Expected content %v, but got: %v.\n", toks[i].content, typ.Content());
      }
    }
  }
  if i != len(toks) {
    t.Error("Got too few tokens!");
  }
  if tok.Type() != common.TOK_EOF {
    t.Error("Got too many tokens!");
  }
}

