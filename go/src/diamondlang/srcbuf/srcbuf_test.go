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
    t.Error("ASCII character not recognized.");
  }
}

func TestNonAscii(t *testing.T) {
}
