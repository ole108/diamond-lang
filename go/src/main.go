package main

import (
  "diamondlang/common";
  "diamondlang/srcbuf";
  "diamondlang/lexer";
  "os";
  "flag";
  "fmt";
  "strings";
)

var useCommandLine = flag.Bool("c", false, "use command line as source")


func main() {
  flag.Parse(); // Scans the arg list and sets up flags
  var sb common.SrcBuffer;
  // fill the source buffer either from the command line or from a file:
  if *useCommandLine {
    if flag.NArg() <= 0 {
      fmt.Fprintln(os.Stderr, "FATAL ERROR: Need source line(s) as argument(s)!");
      os.Exit(1);
    }

    srcStr := "";
    for i := 0; i < flag.NArg(); i++ {
      srcStr += flag.Arg(i) + "\n";
    }
    sb = srcbuf.NewSourceFromBuffer(strings.Bytes(srcStr), "command line");
  } else {
    if flag.NArg() <= 0 {
      fmt.Fprintln(os.Stderr, "FATAL ERROR: Need name of source file as argument!");
      os.Exit(1);
    }

    // Initialize the source buffer:
    var err os.Error;
    sb, err = srcbuf.NewSourceFromFile(flag.Arg(0));
    if err != nil {
      fmt.Fprintf(os.Stderr, "FATAL ERROR: Unable to read source file '%s': %s\n",
          flag.Arg(0), err);
      os.Exit(1);
    }
  }

//for ch := sb.Getch(); ch != common.EOF; ch = sb.Getch() {
//  fmt.Println("Found char:", ch, string(ch));
//}

  // Initialize the lexer:
  lx := lexer.NewLexer(sb);

  // Test output:
  for tok := lx.GetToken(); tok.Type() != common.TOK_EOF; tok = lx.GetToken() {
    switch t := tok.(type) {
    case *lexer.IntTok:
      fmt.Println("Got int:", t.Value(), t);
    case *lexer.MultiDedentTok:
      fmt.Println("Got dedentation:", t.Dedent());
    default:
      fmt.Println("Got token:", t.Type(), t);
    }
  }

}
