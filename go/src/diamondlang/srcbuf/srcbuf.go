package srcbuf

import (
  "diamondlang/common";
  "io/ioutil";
  "os";
//"fmt";
)

// --------------------------------------------------------------------------
// The state of a source reading buffer is held in a variable of this type:
// --------------------------------------------------------------------------
type SrcBuffer struct {
  buf          []byte;  // the whole input source file is stored here
  size         int;     // len(buf)
  pos          int;     // current 'reading' position in the buffer
  line         int;     // current line number
  lineStartPos int;     // starting position of current line in buf
  atLineStart  bool;    // are we *really* at the start of the line?
  wholeLine    string;  // the whole current line
  bufName      string;  // for messages to the user only
}


// This function doesn't handle any errors but returns it instead.
func NewSourceFromFile(filename string) (srcBuf *SrcBuffer, err os.Error) {
  var buf []byte;
  buf, err = ioutil.ReadFile(filename);
  if err != nil { return nil, err }
  return NewSourceFromBuffer(buf, filename), err;
}

func NewSourceFromBuffer(buf []byte, name string) *SrcBuffer {
  ret := &SrcBuffer{buf, len(buf), -1, 0, 0, true, "", name};
  ret.updateSrcWholeLine();
  return ret;
}


// returns current column in current line.
func (sb *SrcBuffer) curCol() int { return sb.pos - sb.lineStartPos }

// update the global variable 'sb.wholeLine' to the current line of source code
func (sb *SrcBuffer) updateSrcWholeLine() {
  n := sb.lineStartPos;
  for ; n<sb.size && sb.buf[n]!='\n' && sb.buf[n]!='\r'; n++ { }
  sb.wholeLine = string(sb.buf[sb.lineStartPos : n]);
}

// Error - Handle errors by writing a description to STDERR and exiting.
func (sb *SrcBuffer) Error(msg string) {
  common.HandleFatal(
      common.MakeErrString(msg, sb.line, sb.wholeLine, sb.curCol(), 1)
  );
}

func (sb *SrcBuffer) AtStartOfLine() bool {
  return sb.atLineStart && sb.pos <= sb.lineStartPos;
}

func (sb *SrcBuffer) NotAtStartOfLine() {
  sb.atLineStart = false;
}

// unget a character from the source code
func (sb *SrcBuffer) Ungetch() {
    sb.pos--;
    if sb.pos < sb.lineStartPos {
      sb.Error("Unable to unget characters beyond the beginning of the current line");
    }
}

// get the next character from the source code
func (sb *SrcBuffer) Getch() byte {
  // handle EOF
  if sb.pos+1 >= sb.size  {
    sb.pos = sb.size;
    return common.EOF
  }

  // read and check new character
  sb.pos++;
  ch := sb.buf[sb.pos];
  sb.handleNewLine(ch);
  sb.ensureAscii(ch);

  return ch;
}

// choke on multibyte UTF8 characters
func (sb *SrcBuffer) ensureAscii(ch byte) {
  if ch > 127  { sb.Error("Unable to handle non-ASCII character") }
}

// handle new lines (can be done here since they can't be escaped :-)
func (sb *SrcBuffer) handleNewLine(ch byte) {
  oldChar := byte(0);
  if sb.pos > 0  { oldChar = sb.buf[sb.pos-1] }
  if (oldChar == '\n' || (oldChar == '\r' && ch != '\n')) {
    sb.updateCurNewLine();
  }
}

// update global variables because of a new line
func (sb *SrcBuffer) updateCurNewLine() {
  sb.line++;
  sb.lineStartPos = sb.pos;
  sb.atLineStart = true;
  sb.updateSrcWholeLine();
}


// --------------------------------------------------------------------------
// Special types that help the lexer.
// --------------------------------------------------------------------------

type SrcMark  int
func (sb *SrcBuffer) NewMark() common.SrcMark {
  ret := new(SrcMark);
  *ret = SrcMark(sb.pos);
  return ret;
}
func (mark *SrcMark) Pos() int { return int(*mark); }

type MultiLineSrcMark struct {
  pos       int;
  line      int;
  column    int;
  wholeLine string;
}
func (sb *SrcBuffer) NewMultiLineMark() common.MultiLineSrcMark {
  return &MultiLineSrcMark{sb.pos, sb.line, sb.curCol(), sb.wholeLine};
}
func (mark *MultiLineSrcMark) Pos() int { return mark.pos; }
func (mark *MultiLineSrcMark) Line() int { return mark.line; }
func (mark *MultiLineSrcMark) Column() int { return mark.column; }
func (mark *MultiLineSrcMark) WholeLine() string { return mark.wholeLine; }

type SrcPiece struct {
  startLine    int;
  startColumn  int;
  wholeLine    string;
  content      string;
}

// Create a new SrcPiece.
// The given mark is the start of the piece.
// The piece ends one character before the current reading position.
func (sb *SrcBuffer) NewPiece(start common.SrcMark) common.SrcPiece {
  col := start.Pos() - sb.lineStartPos;
  if col < 0 { sb.Error("Unable to start before start of line"); }
  return &SrcPiece{sb.line, col, sb.wholeLine, string(sb.buf[start.Pos():sb.pos])};
}

func (sb *SrcBuffer) NewMultiPiece(pieces []common.SrcPiece) common.SrcPiece {
  if len(pieces) <= 0 {
    sb.Error("Too few source pieces");
  }
  cnt  := "";
  wl   := pieces[0].WholeLine();
  line := pieces[0].Line();
  for _, piece := range pieces {
    if line < piece.Line() {
      wl += "\n" + piece.WholeLine();
      line = piece.Line();
    }
    cnt += piece.Content();
  }
  return &SrcPiece{pieces[0].Line(), pieces[0].Column(), wl, cnt};
}

func (sb *SrcBuffer) NewMultiLinePiece(start common.MultiLineSrcMark) common.SrcPiece {
  content := string(sb.buf[start.Pos():sb.pos]);
  return &SrcPiece{start.Line(), start.Column(),
      makeWholeLine(start, content, sb), content};
}

func (piece *SrcPiece) Line() int { return piece.startLine }
func (piece *SrcPiece) Column() int { return piece.startColumn }
func (piece *SrcPiece) Content() string { return piece.content }
func (piece *SrcPiece) WholeLine() string { return piece.wholeLine }
func (piece *SrcPiece) String() string { return piece.content }

func makeWholeLine(start common.MultiLineSrcMark, content string, sb *SrcBuffer) string {
  beg := start.WholeLine()[0 : start.Column()];
  mid := content;
  end := "";
  if sb.curCol() < len(sb.wholeLine) {
    end = sb.wholeLine[sb.curCol() : len(sb.wholeLine)];
  }
  return beg + mid + end;
}

