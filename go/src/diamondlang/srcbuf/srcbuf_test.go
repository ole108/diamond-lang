package srcbuf

import (
  "testing";
  "diamondlang/common";
)


func TestAscii(t *testing.T) {
  var ascii [128]byte;
  for i := 0; i < len(ascii); i++ { ascii[i] = byte(i); }
  sb := NewSourceFromBuffer(&ascii, "TestAscii");
  for i := 0; i < len(ascii); i++ {
    if sb.Getch() != byte(i) {
      t.Error("ASCII character not recognized.");
    }
  }
  if sb.Getch() != common.EOF {
    t.Error("EOF not recognized.");
  }
}

func TestAtStartOfLine(t *testing.T) {
  tstBuf := []byte{ 'a', '\n', 'b', '\n', 'c' };
  sb := NewSourceFromBuffer(tstBuf, "TestAtStartOfLine");
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
  tstBuf := []byte{ 'a', '\n', 'b', '\r', '\n', 'c', '\r', 'd', '\n', '\r', 'e' };
  sb := NewSourceFromBuffer(tstBuf, "TestLineEndings");
  sb.Getch();  // got 'a'
  if !sb.AtStartOfLine() {
    t.Error("Start of 1. line not recognized.");
  }
  if sb.line != 0 { t.Error("Line 0 not recognized."); }
  if sb.wholeLine != "a" { t.Error("wholeLine != a."); }
  sb.Getch();  // got '\n'
  if sb.AtStartOfLine() {
    t.Error("Getch() should move away from start of line.");
  }
  if sb.line != 0 { t.Error("Line 0 not recognized (2)."); }
  if sb.wholeLine != "a" { t.Error("wholeLine != a (2)."); }
  sb.Getch();  // got 'b'
  if !sb.AtStartOfLine() {
    t.Error("Start of 2. line not recognized.");
  }
  if sb.line != 1 { t.Error("Line 1 not recognized."); }
  if sb.wholeLine != "b" { t.Error("wholeLine != b (1)."); }
  sb.Getch();  // got '\r'
  if sb.AtStartOfLine() {
    t.Error("'\\r' isn't start of line.");
  }
  if sb.line != 1 { t.Error("Line 1 not recognized. (2)"); }
  if sb.wholeLine != "b" { t.Error("wholeLine != b (2)."); }
  sb.Getch();  // got '\n'
  if sb.AtStartOfLine() {
    t.Error("'\\n' isn't start of line.");
  }
  if sb.line != 1 { t.Error("Line 1 not recognized. (3)"); }
  if sb.wholeLine != "b" { t.Error("wholeLine != b (3)."); }
  sb.Getch();  // got 'c'
  if !sb.AtStartOfLine() {
    t.Error("Start of 3. line not recognized.");
  }
  if sb.line != 2 { t.Error("Line 2 not recognized."); }
  if sb.wholeLine != "c" { t.Error("wholeLine != c (1)."); }
  sb.Getch();  // got '\r'
  if sb.AtStartOfLine() {
    t.Error("'\\r' isn't start of line (2).");
  }
  if sb.line != 2 { t.Error("Line 2 not recognized (2)."); }
  if sb.wholeLine != "c" { t.Error("wholeLine != c (2)."); }
  sb.Getch();  // got 'd'
  if !sb.AtStartOfLine() {
    t.Error("Start of 4. line not recognized.");
  }
  if sb.line != 3 { t.Error("Line 3 not recognized."); }
  if sb.wholeLine != "d" { t.Error("wholeLine != d (1)."); }
  sb.Getch();  // got '\n'
  if sb.AtStartOfLine() {
    t.Error("'\\n' isn't start of line.");
  }
  if sb.line != 3 { t.Error("Line 3 not recognized (2)."); }
  if sb.wholeLine != "d" { t.Error("wholeLine != d (2)."); }
  sb.Getch();  // got '\r'
  if !sb.AtStartOfLine() {
    t.Error("Start of 5. line not recognized.");
  }
  if sb.line != 4 { t.Error("Line 4 not recognized."); }
  if sb.wholeLine != "" { t.Error("wholeLine !=  (1)."); }
  sb.Getch();  // got 'e'
  if !sb.AtStartOfLine() {
    t.Error("Start of 6. line not recognized.");
  }
  if sb.line != 5 { t.Error("Line 5 not recognized."); }
  if sb.wholeLine != "e" { t.Error("wholeLine != e (1)."); }
  if sb.Getch() != common.EOF {
    t.Error("EOF not recognized.");
  }
}

