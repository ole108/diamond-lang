package srcbuf

import (
  "diamondlang/common";
  "fmt";
  "os";
)

type line struct {
  num   int;
  buf   []byte;
}

func any2line(any interface{}) *line {
  line, ok := any.(*line);
  if !ok {
    common.HandleFatal("Internal error (not a line)!");
  }
  return line;
}

// String -- returns the content of the line as a string without line ending.
func (l *line) String() string {
  n := len(l.buf) - 1;
  if n >= 0 && l.buf[n] == common.EOF {
    n--;
  } else {
    for n >= 0 && (l.buf[n] == '\n' || l.buf[n] == '\r') {
      n--;
    }
  }
  if n >= 0 {
    return string(l.buf[0 : n+1]);
  }
  return "";
}

// read a line from the readByter (only "\n" and "\r\n" are recognized as line endings)
func newLine(rb readByter, num int) (l *line, err os.Error) {
  buf := new([256]byte);
  pos := 0;

  // read almost a full line into the buffer
  c,e := rb.ReadByte();
  for ; e == nil && c != '\n';
      c,e = rb.ReadByte() {
    if !legalChar(c) {
      common.HandleFatal(fmt.Sprintf("Unable to handle non-ASCII character %d in line %d!", c, num));
    }
    putChar(c, buf, pos, num);
    pos++;
  }

  // handle last character
  if e == os.EOF {
    putChar(common.EOF, buf, pos, num);
  } else if e != nil {
    return nil, e;
  } else {
    putChar(c, buf, pos, num);
  }

  return &line{num, buf[0:pos+1]}, nil;
}

func putChar(c byte, buf []byte, pos int, num int) {
  if pos >= len(buf) {
    common.HandleFatal(fmt.Sprintf("Line %d is longer than %d bytes!", num, len(buf)));
  }
  buf[pos] = c;
}

func legalChar(c byte) bool {
  return c >= 0 && c < 128;
}

