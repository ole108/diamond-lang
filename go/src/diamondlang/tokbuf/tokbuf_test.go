package tokbuf

import (
  "testing";
  "diamondlang/common";
  "diamondlang/srcbuf";
  "diamondlang/lexer";
  "strings";
)

type tstTok struct {
  typ          common.TokEnum;
  content      string;
  checkContent bool;
  numVal       int64;
  strVal       string;
}


func TestIndent(t *testing.T) {
  testStr := `If
    mod
  Elif
    val`;

  testToks := []*tstTok{
    &tstTok{common.TOK_FUNC_ID, "If", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_INDENT, "    ", true, 0, ""},
    &tstTok{common.TOK_MODULE_ID, "mod", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_HALF_DEDENT, "  ", true, 0, ""},
    &tstTok{common.TOK_FUNC_ID, "Elif", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_HALF_INDENT, "    ", true, 0, ""},
    &tstTok{common.TOK_MODULE_ID, "val", true, 0, ""},
    &tstTok{common.TOK_DEDENT, "", true, 0, ""},
  };

  testStringVsTokens(t, testStr, testToks);
}

func TestIndentBlock(t *testing.T) {
  testStr := `If bla > 0:
    mod.Func mod.CONST mod.CONST.val Fn i
  Elif bla < 0:   # blue
    mod.FuncAli bla
  Else:   
    bla.val  # should work!
bla = 0
   # geschafft!`;

  testToks := []*tstTok{
    &tstTok{common.TOK_FUNC_ID, "If", true, 0, ""},             // 0
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_MODULE_ID, "bla", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_OP_ID, ">", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_INT, "0", false, 0, ""},
    &tstTok{common.TOK_BLOCK_START, ":\n", false, 0, ""},

    &tstTok{common.TOK_INDENT, "    ", true, 0, ""},            // 8
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

    &tstTok{common.TOK_HALF_DEDENT, "  ", true, 0, ""},         // 19
    &tstTok{common.TOK_FUNC_ID, "Elif", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_MODULE_ID, "bla", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_OP_ID, "<", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_INT, "0", false, 0, ""},
    &tstTok{common.TOK_BLOCK_START, ":\n", false, 0, ""},
    &tstTok{common.TOK_SPACE, "   ", true, 3, ""},
    &tstTok{common.TOK_COMMENT, "# blue", true, 0, ""},

    &tstTok{common.TOK_HALF_INDENT, "    ", true, 0, ""},       // 30
    &tstTok{common.TOK_FUNC_ID, "mod.FuncAli", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_MODULE_ID, "bla", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_HALF_DEDENT, "  ", true, 0, ""},         // 35
    &tstTok{common.TOK_FUNC_ID, "Else", true, 0, ""},
    &tstTok{common.TOK_BLOCK_START, ":\n", false, 0, ""},
    &tstTok{common.TOK_SPACE, "   ", true, 3, ""},

    &tstTok{common.TOK_HALF_INDENT, "    ", true, 0, ""},       // 39
    &tstTok{common.TOK_VAL_ID, "bla.val", true, 0, ""},
    &tstTok{common.TOK_SPACE, "  ", true, 2, ""},
    &tstTok{common.TOK_COMMENT, "# should work!", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_DEDENT, "", true, 0, ""},                // 44
    &tstTok{common.TOK_MODULE_ID, "bla", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_OP_ID, "=", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_INT, "0", false, 0, ""},
    &tstTok{common.TOK_NL, "\n", false, 0, ""},

    &tstTok{common.TOK_SPACE, "   ", true, 1003, ""},           // 51
    &tstTok{common.TOK_COMMENT, "# geschafft!", true, 0, ""},
  };

  testStringVsTokens(t, testStr, testToks);
}

func TestSemiNl(t *testing.T) {
  testStr := `If
    mod1;mod2
  Elif
    val1; val2`;

  testToks := []*tstTok{
    &tstTok{common.TOK_FUNC_ID, "If", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", true, 0, ""},

    &tstTok{common.TOK_INDENT, "    ", true, 0, ""},
    &tstTok{common.TOK_MODULE_ID, "mod1", true, 0, ""},
    &tstTok{common.TOK_NL, ";", true, 0, ""},
    &tstTok{common.TOK_MODULE_ID, "mod2", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", true, 0, ""},

    &tstTok{common.TOK_HALF_DEDENT, "  ", true, 0, ""},
    &tstTok{common.TOK_FUNC_ID, "Elif", true, 0, ""},
    &tstTok{common.TOK_NL, "\n", true, 0, ""},

    &tstTok{common.TOK_HALF_INDENT, "    ", true, 0, ""},
    &tstTok{common.TOK_MODULE_ID, "val1", true, 0, ""},
    &tstTok{common.TOK_NL, ";", true, 0, ""},
    &tstTok{common.TOK_SPACE, " ", true, 1, ""},
    &tstTok{common.TOK_MODULE_ID, "val2", true, 0, ""},
    &tstTok{common.TOK_DEDENT, "", true, 0, ""},
  };

  testStringVsTokens(t, testStr, testToks);
}


func testStringVsTokens(t *testing.T, str string, toks []*tstTok) {
  tb := NewTokenBuffer(lexer.NewLexer(srcbuf.NewSourceFromBuffer(strings.Bytes(str), "TestTokens")));
  var tok common.Token;
  var i   int;
  for tok, i = tb.GetToken(), 0; tok.Type() != common.TOK_EOF && i < len(toks);
      tok, i = tb.GetToken(), i+1 {
    switch typ := tok.(type) {
    case *lexer.IntTok:
      if toks[i].typ != typ.Type() {
        t.Errorf("%d: Expected token type %v, but got: %v.\n", i, toks[i].typ, typ.Type());
      }
      if toks[i].numVal != typ.Value() {
        t.Errorf("%d: Expected value %v, but got: %v.\n", i, toks[i].numVal, typ.Value());
      }
    case *lexer.CharTok:
      if toks[i].typ != typ.Type() {
        t.Errorf("%d: Expected token type %v, but got: %v.\n", i, toks[i].typ, typ.Type());
      }
      if toks[i].numVal != int64(typ.Value()) {
        t.Errorf("%d: Expected value %v, but got: %v.\n", i, toks[i].numVal, typ.Value());
      }
    case *lexer.StringTok:
      if toks[i].typ != typ.Type() {
        t.Errorf("%d: Expected token type %v, but got: %v.\n", i, toks[i].typ, typ.Type());
      }
      if toks[i].strVal != typ.Value() {
        t.Errorf("%d: Expected value %v, but got: %v.\n", i, toks[i].strVal, typ.Value());
      }
    case *lexer.SpaceTok:
      if toks[i].typ != typ.Type() {
        t.Errorf("%d: Expected token type %v, but got: %v.\n", i, toks[i].typ, typ.Type());
      }
      space := int64(typ.Space());
      if typ.AtStartOfLine() { space += 1000; }
      if toks[i].numVal != space {
        t.Errorf("%d: Expected space %v, but got: %v.\n", i, toks[i].numVal, space);
      }
    default:
      if toks[i].typ != typ.Type() {
        t.Errorf("%d: Expected token type %v, but got: %v.\n", i, toks[i].typ, typ.Type());
      }
      if toks[i].checkContent && toks[i].content != typ.Content() {
        t.Errorf("%d: Expected content %v, but got: %v.\n", i, toks[i].content, typ.Content());
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

