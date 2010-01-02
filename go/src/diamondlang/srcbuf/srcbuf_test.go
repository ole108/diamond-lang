package srcbuf

import (
  "testing";
  "diamondlang/common";
)


func TestAscii(t *testing.T) {
  var ascii [128]byte;
  for i := 0; i < len(ascii); i++ { ascii[i] = byte(i); }
  sb := NewSourceFromBuffer(&ascii);
  for i := 0; i < len(ascii); i++ {
    ch := sb.Getch();
    if ch != byte(i) {
      t.Fatalf("ASCII character not recognized (act %d != %d exp).", ch, i);
    }
  }
  if sb.Getch() != common.EOF {
    t.Error("EOF not recognized.");
  }
}

func TestAtStartOfLine(t *testing.T) {
  tstBuf := []byte{ 'a', '\n', 'b', '\n', 'c' };
  sb := NewSourceFromBuffer(tstBuf);
  if !sb.AtStartOfLine() {
    t.Error("Start of buffer isn't start of line.");
  }
  sb.Getch();  // got 'a'
  if !sb.AtStartOfLine() {
    t.Error("Start of first line isn't start of line.");
  }
  sb.Getch();  // got '\n'
  if sb.AtStartOfLine() {
    t.Error("Getch() should move away from start of line.");
  }
  sb.Getch();  // got 'b'
  if !sb.AtStartOfLine() {
    t.Error("Start of next line not recognized.");
  }
  sb.NotAtStartOfLine();
  if sb.AtStartOfLine() {
    t.Error("NotAtStartOfLine() didn't work.");
  }
  sb.Getch();  // got '\n'
  sb.Getch();  // got 'c'
  if !sb.AtStartOfLine() {
    t.Error("Start of last line not recognized.");
  }
  if sb.Getch() != common.EOF {
    t.Error("EOF not recognized.");
  }
}

func TestLineEndings(t *testing.T) {
  tstBuf := []byte{ 'a', '\n', 'b', '\r', '\n', '\r', '\n', 'c' };
  sb := NewSourceFromBuffer(tstBuf);
  sb.Getch();  // got 'a'
  if sb.curLine == nil {
    t.Fatal("1. line not read.");
  }
  if !sb.AtStartOfLine() {
    t.Error("Start of 1. line not recognized.");
  }
  if sb.curLine.num != 0 { t.Error("Line 0 not recognized."); }
  if sb.curLine.String() != "a" { t.Error("wholeLine != a."); }
  sb.Getch();  // got '\n'
  if sb.AtStartOfLine() {
    t.Error("Getch() should move away from start of line.");
  }
  if sb.curLine.num != 0 { t.Error("Line 0 not recognized (2)."); }
  if sb.curLine.String() != "a" { t.Error("wholeLine != a (2)."); }
  sb.Getch();  // got 'b'
  if !sb.AtStartOfLine() {
    t.Error("Start of 2. line not recognized.");
  }
  if sb.curLine.num != 1 { t.Error("Line 1 not recognized."); }
  if sb.curLine.String() != "b" { t.Error("wholeLine != b (1)."); }
  sb.Getch();  // got '\r'
  if sb.AtStartOfLine() {
    t.Error("'\\r' isn't start of line.");
  }
  if sb.curLine.num != 1 { t.Error("Line 1 not recognized. (2)"); }
  if sb.curLine.String() != "b" { t.Error("wholeLine != b (2)."); }
  sb.Getch();  // got '\n'
  if sb.AtStartOfLine() {
    t.Error("'\\n' isn't start of line.");
  }
  if sb.curLine.num != 1 { t.Error("Line 1 not recognized. (3)"); }
  if sb.curLine.String() != "b" { t.Error("wholeLine != b (3)."); }

  sb.Getch();  // got '\r'
  if !sb.AtStartOfLine() {
    t.Error("Start of 3. line not recognized.");
  }
  if sb.curLine.num != 2 { t.Error("Line 2 not recognized. (1)"); }
  if sb.curLine.String() != "" { t.Error("wholeLine !=  (1)."); }
  sb.Getch();  // got '\n'
  if sb.AtStartOfLine() {
    t.Error("'\\n' isn't start of line.");
  }
  if sb.curLine.num != 2 { t.Error("Line 2 not recognized. (2)"); }
  if sb.curLine.String() != "" { t.Error("wholeLine !=  (2)."); }
  sb.Getch();  // got 'c'
  if !sb.AtStartOfLine() {
    t.Error("Start of 4. line not recognized.");
  }
  if sb.curLine.num != 3 { t.Error("Line 3 not recognized."); }
  if sb.curLine.String() != "c" { t.Error("wholeLine != c (1)."); }
  if sb.Getch() != common.EOF {
    t.Error("EOF not recognized.");
  }
}

func TestSrcMark(t *testing.T) {
  tstBuf := []byte{ 'a', '\n', 'b', '\r', '\n', '\r', '\n', 'e' };
  sb := NewSourceFromBuffer(tstBuf);
  sb.Getch();  // got 'a'
  mark1 := sb.NewMark();
  testSrcPiece(sb.NewPiece(mark1), 0, 0, "", "a", t, 11);
  sb.Getch();  // got '\n'
  mark2 := sb.NewMark();
  testSrcPiece(sb.NewPiece(mark1), 0, 0, "a", "a", t, 21);
  testSrcPiece(sb.NewPiece(mark2), 0, 1, "", "a", t, 22);
  sb.Getch();  // got 'b'
  mark3 := sb.NewMark();
  testSrcPiece(sb.NewPiece(mark1), 0, 0, "a\n", "a\nb", t, 31);
  testSrcPiece(sb.NewPiece(mark2), 0, 1, "\n", "a\nb", t, 32);
  testSrcPiece(sb.NewPiece(mark3), 1, 0, "", "b", t, 33);
  sb.Getch();  // got '\r'
  testSrcPiece(sb.NewPiece(mark1), 0, 0, "a\nb", "a\nb", t, 41);
  testSrcPiece(sb.NewPiece(mark2), 0, 1, "\nb", "a\nb", t, 42);
  testSrcPiece(sb.NewPiece(mark3), 1, 0, "b", "b", t, 43);
  sb.Getch();  // got '\n'
  testSrcPiece(sb.NewPiece(mark1), 0, 0, "a\nb\r", "a\nb", t, 51);
  testSrcPiece(sb.NewPiece(mark2), 0, 1, "\nb\r", "a\nb", t, 52);
  testSrcPiece(sb.NewPiece(mark3), 1, 0, "b\r", "b", t, 53);
  sb.Getch();  // got '\r'
  testSrcPiece(sb.NewPiece(mark1), 0, 0, "a\nb\r\n", "a\nb\n", t, 61);
  testSrcPiece(sb.NewPiece(mark2), 0, 1, "\nb\r\n", "a\nb\n", t, 62);
  testSrcPiece(sb.NewPiece(mark3), 1, 0, "b\r\n", "b\n", t, 63);
  sb.Getch();  // got '\n'
  testSrcPiece(sb.NewPiece(mark1), 0, 0, "a\nb\r\n\r", "a\nb\n", t, 71);
  testSrcPiece(sb.NewPiece(mark2), 0, 1, "\nb\r\n\r", "a\nb\n", t, 72);
  testSrcPiece(sb.NewPiece(mark3), 1, 0, "b\r\n\r", "b\n", t, 73);
  sb.Getch();  // got 'e'
  testSrcPiece(sb.NewPiece(mark1), 0, 0, "a\nb\r\n\r\n", "a\nb\n\ne", t, 81);
  testSrcPiece(sb.NewPiece(mark2), 0, 1, "\nb\r\n\r\n", "a\nb\n\ne", t, 82);
  testSrcPiece(sb.NewPiece(mark3), 1, 0, "b\r\n\r\n", "b\n\ne", t, 83);
  sb.Getch();  // got EOF
  testSrcPiece(sb.NewPiece(mark1), 0, 0, "a\nb\r\n\r\ne", "a\nb\n\ne", t, 91);
  testSrcPiece(sb.NewPiece(mark2), 0, 1, "\nb\r\n\r\ne", "a\nb\n\ne", t, 92);
  testSrcPiece(sb.NewPiece(mark3), 1, 0, "b\r\n\r\ne", "b\n\ne", t, 93);
  sb.Getch();  // got EOF
  testSrcPiece(sb.NewPiece(mark1), 0, 0, "a\nb\r\n\r\ne", "a\nb\n\ne", t, 101);
  testSrcPiece(sb.NewPiece(mark2), 0, 1, "\nb\r\n\r\ne", "a\nb\n\ne", t, 102);
  testSrcPiece(sb.NewPiece(mark3), 1, 0, "b\r\n\r\ne", "b\n\ne", t, 103);
}

func testSrcPiece(piece common.SrcPiece, line int, column int,
                  content string, wholeLine string, t *testing.T, num int) {
  failed := false;
  if piece.StartLine() != line {
    failed = true;
    t.Errorf("(%d) Expected line %d but got %d.", num, line, piece.StartLine());
  }
  if piece.StartColumn() != column {
    failed = true;
    t.Errorf("(%d) Expected column %d but got %d.", num, column, piece.StartColumn());
  }
  if piece.Content() != content {
    failed = true;
    t.Errorf("(%d) Expected content `%v` but got `%v`.", num, content, piece.Content());
  }
  if piece.WholeLine() != wholeLine {
    failed = true;
    t.Errorf("(%d) Expected wholeLine `%v` but got `%v`.", num, wholeLine, piece.WholeLine());
  }
  if failed { t.FailNow(); }
}
