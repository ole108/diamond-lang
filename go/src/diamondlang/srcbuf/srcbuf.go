package srcbuf

import (
  "diamondlang/common";
  "container/list";
  "bufio";
  "bytes";
  "os";
  "io";
  "fmt";
)

const (
  MAX_BUFFER_LINES = 1024;
)

type readByter interface {
  ReadByte() (c byte, err os.Error);
}


// --------------------------------------------------------------------------
// The state of a source reading buffer is held in a variable of this type:
// --------------------------------------------------------------------------
type SrcBuffer struct {
  source       readByter;    // our source for bytes
  buf         *list.List;    // the real buffer of lines
  curElem     *list.Element; // the current element in the buffer
  curLine     *line;         // the current line in the buffer
  curCol       int;          // current column in the current line
  atLineStart  bool;         // are we *really* at the start of the line?
  eof          bool;
}


// This function doesn't handle any errors but returns it instead.
func NewSourceFromFile(filename string) (srcBuf *SrcBuffer, err os.Error) {
  file, e := os.Open(filename, os.O_RDONLY, 0444);
  if e != nil { return nil, e }
  return NewSourceFromReader(file), nil;
}

func NewSourceFromBuffer(buf []byte) *SrcBuffer {
  return NewSourceFromReader(bytes.NewBuffer(buf));
}

func NewSourceFromReader(rd io.Reader) *SrcBuffer {
  src, ok := rd.(readByter);
  if !ok {
    src = bufio.NewReader(rd);
  }
  ret := &SrcBuffer{src, list.New(), nil, nil, -1, true, false};
  return ret;
}


// Error - Handle errors by writing a description to STDERR and exiting.
func (sb *SrcBuffer) Error(msg string) {
  common.HandleFatal(
      common.MakeErrString(msg, sb.curLine.num, sb.curLine.String(), sb.curCol, 1)
  );
}

func (sb *SrcBuffer) AtStartOfLine() bool {
  return sb.atLineStart && sb.curCol <= 0;
}

func (sb *SrcBuffer) NotAtStartOfLine() {
  sb.atLineStart = false;
}

// remove old lines from the buffer
func (sb *SrcBuffer) ClearUpTo(mark common.SrcMark) {
  for mark.Elem != sb.buf.Front() {
    sb.buf.Remove(sb.buf.Front());
  }
}

// unget a character from the source code
func (sb *SrcBuffer) Ungetch() {
    if sb.curCol < 0 {
      sb.Error("Unable to unget characters beyond the beginning of the current line");
    }
    sb.curCol--;
}

// get the next character from the source code
func (sb *SrcBuffer) Getch() byte {
  // handle EOF
  if sb.eof { return common.EOF; }

  // read a new line if necessary
  if sb.mustGotoNextLine() {
    sb.gotoNextLine();
  }

  // read new character
  sb.curCol++;
  ch := sb.curLine.buf[sb.curCol];
  if ch == common.EOF { sb.eof = true; }  // handle EOF

  return ch;
}
func (sb *SrcBuffer) mustGotoNextLine() bool {
  return sb.curLine == nil || sb.curCol+1 >= len(sb.curLine.buf);
}
func (sb *SrcBuffer) gotoNextLine() {
  if sb.curElem == nil || sb.curElem.Next() == nil {
    sb.readNewLine();
  } else {
    sb.curElem = sb.curElem.Next();
    sb.curLine = any2line(sb.curElem.Value);
  }
  sb.curCol = -1;
  sb.atLineStart = true;
}

// read a new line and append it to the source buffer
func (sb *SrcBuffer) readNewLine() {
  if sb.buf.Len() > MAX_BUFFER_LINES {
    sb.Error("Buffer overflow; please preak up your code into smaller pieces");
  }

  num := 0;
  if sb.curLine != nil {
    num = sb.curLine.num + 1;
  }
  line, err := newLine(sb.source, num);
  if err != nil {
    common.HandleFatal(fmt.Sprintf("While reading line %d: %v", num+1, err.String()));
  }
  sb.curLine = line;
  sb.curElem = sb.buf.PushBack(line);
}


// --------------------------------------------------------------------------
// Special types that help the lexer.
// --------------------------------------------------------------------------

func (sb *SrcBuffer) NewMark() common.SrcMark {
  return common.SrcMark{sb.curElem, sb.curCol};
}

type SrcPiece struct {
  start common.SrcMark;
  end   common.SrcMark;
}

// Create a new SrcPiece.
// The given mark is the start of the piece.
// The piece ends one character before the current reading position.
func (sb *SrcBuffer) NewPiece(start common.SrcMark) common.SrcPiece {
  return &SrcPiece{start, sb.NewMark()};
}

func (sb *SrcBuffer) NewAnyPiece(start common.SrcMark, end common.SrcMark) common.SrcPiece {
  return &SrcPiece{start, end};
}

func (piece *SrcPiece) Start() common.SrcMark { return piece.start; }
func (piece *SrcPiece) End() common.SrcMark { return piece.end; }
func (piece *SrcPiece) StartLine() int { return any2line(piece.start.Elem.Value).num; }
func (piece *SrcPiece) StartColumn() int { return piece.start.Col; }
func (piece *SrcPiece) String() string { return piece.Content() }

func (piece *SrcPiece) Error(msg string) {
  common.HandleFatal(common.MakeErrString(msg, piece.StartLine(), piece.WholeLine(),
      piece.StartColumn(), len(piece.Content()) )
  );
}

func (piece *SrcPiece) WholeLine() string {
  ret := "";
  first := true;
  for elem := piece.start.Elem; elem != piece.end.Elem; elem = elem.Next() {
    if !first { ret += "\n"; }
    first = false;
    ret += any2line(elem.Value).String();
  }
  if !first { ret += "\n"; }
  return ret + any2line(piece.end.Elem.Value).String();
}

func (piece *SrcPiece) Content() string {
  if piece.start.Elem == piece.end.Elem {
    return string(any2line(piece.start.Elem.Value).buf[piece.start.Col : piece.end.Col]);
  }

  line := any2line(piece.start.Elem.Value);
  start := string(line.buf[piece.start.Col : len(line.buf)]);

  mid := "";
  for elem := piece.start.Elem.Next(); elem != piece.end.Elem; elem = elem.Next() {
    mid += string(any2line(elem.Value).buf);
  }

  line = any2line(piece.end.Elem.Value);
  end := string(line.buf[0 : piece.end.Col]);

  return start + mid + end;
}

