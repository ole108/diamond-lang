package common

import (
  "testing";
)

func TestIsSpace(t *testing.T) {
  if (!IsSpace(' ')) { t.Error("Blank not recognized as space."); }
  if (!IsSpace('\t')) { t.Error("Tab not recognized as space."); }
  if (IsSpace('\r')) { t.Error("<CR> recognized as space."); }
  if (IsSpace('\n')) { t.Error("<LF> recognized as space."); }
  if (IsSpace('\v')) { t.Error("<VTAB> recognized as space."); }
  if (IsSpace(0)) { t.Error("<NUL> recognized as space."); }
  if (IsSpace(12)) { t.Error("<^L> recognized as space."); }
}

func TestSpaceAmount(t *testing.T) {
  if (SpaceAmount(' ') != 1)        { t.Error("Blank not recognized as 1 space."); }
  if (SpaceAmount('\t') != TABSIZE) { t.Error("Tab not recognized as 4 spaces."); }
  if (SpaceAmount('\r') != 0)       { t.Error("<CR> recognized as space."); }
  if (SpaceAmount('\n') != 0)       { t.Error("<LF> recognized as space."); }
  if (SpaceAmount('\v') != 0)       { t.Error("<VTAB> recognized as space."); }
  if (SpaceAmount(0) != 0)          { t.Error("<NUL> recognized as space."); }
  if (SpaceAmount(12) != 0)         { t.Error("<^L> recognized as space."); }
}
