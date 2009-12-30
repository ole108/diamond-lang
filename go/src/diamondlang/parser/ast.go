package parser

import (
  "diamondlang/common";
)


type AstNode struct {
  line      int;
  column    int;
  wholeLine string;
  content   string;
}
// Error - Handle errors by writing a description to STDERR and exiting.
func (an *AstNode) Error(msg string) {
  common.HandleFatal(
      common.MakeErrString(msg, an.line, an.wholeLine, an.column, len(an.content))
  );
}
func (an *AstNode) String() string { return an.content; }


type ExprAst struct {
  *AstNode
}
func (an *ExprAst) IsExpr() bool { return true; }


type NumberExprAst struct {
  *ExprAst;
  value int64;
}
func (an *NumberExprAst) Value() int64 { return an.value; }
func NewNumberExprAst(line int, column int, wholeLine string, content string,
                      value int64) common.NumberExprAst                       {
  return &NumberExprAst{&ExprAst{&AstNode{line, column, wholeLine, content}}, value};
}


type BoolExprAst struct {
  *ExprAst;
  value bool;
}
func (an *BoolExprAst) Value() bool { return an.value; }
func NewBoolExprAst(line int, column int, wholeLine string, content string,
                      value bool) common.BoolExprAst                       {
  return &BoolExprAst{&ExprAst{&AstNode{line, column, wholeLine, content}}, value};
}


type CharExprAst struct {
  *ExprAst;
  value byte;
}
func (an *CharExprAst) Value() byte { return an.value; }
func NewCharExprAst(line int, column int, wholeLine string, content string,
                      value byte) common.CharExprAst                       {
  return &CharExprAst{&ExprAst{&AstNode{line, column, wholeLine, content}}, value};
}


type StringExprAst struct {
  *ExprAst;
  value string;
}
func (an *StringExprAst) Value() string { return an.value; }
func NewStringExprAst(line int, column int, wholeLine string, content string,
                      value string) common.StringExprAst                       {
  return &StringExprAst{&ExprAst{&AstNode{line, column, wholeLine, content}}, value};
}


type ValueExprAst struct {
  *ExprAst;
  module string;
  value  string;
  subs   []string;
}
func (an *ValueExprAst) Module() string { return an.module; }
func (an *ValueExprAst) Value() string { return an.value; }
func (an *ValueExprAst) Subs() []string { return an.subs; }
func NewValueExprAst(line int, column int, wholeLine string, content string,
             module string, value string, subs []string) common.ValueExprAst {
  return &ValueExprAst{&ExprAst{&AstNode{line, column, wholeLine, content}},
                         module, value, subs};
}


type ConstantExprAst struct {
  *ExprAst;
  module   string;
  constant string;
  subs     []string;
}
func (an *ConstantExprAst) Module() string { return an.module; }
func (an *ConstantExprAst) Constant() string { return an.constant; }
func (an *ConstantExprAst) Subs() []string { return an.subs; }
func NewConstantExprAst(line int, column int, wholeLine string, content string,
          module string, constant string, subs []string) common.ConstantExprAst {
  return &ConstantExprAst{&ExprAst{&AstNode{line, column, wholeLine, content}},
                         module, constant, subs};
}


type CallExprAst struct {
  *ExprAst;
  module   string;
  function string;
  typ      common.CallTypeEnum;
  args     []common.ExprAst;
}
func (an *CallExprAst) Module() string { return an.module; }
func (an *CallExprAst) Func() string { return an.function; }
func (an *CallExprAst) Type() common.CallTypeEnum { return an.typ; }
func (an *CallExprAst) Args() []common.ExprAst { return an.args; }
func NewCallExprAst(line int, column int, wholeLine string, content string,
     module string, function string, typ common.CallTypeEnum, args []common.ExprAst)
         common.CallExprAst {
  return &CallExprAst{&ExprAst{&AstNode{line, column, wholeLine, content}},
                         module, function, typ, args};
}


type AssignmentAst struct {
  *AstNode;
  value common.ValueExprAst;
  expr  common.ExprAst;
}
func (an *AssignmentAst) Value() common.ValueExprAst { return an.value; }
func (an *AssignmentAst) Expr() common.ExprAst { return an.expr; }
func NewAssignmentAst(line int, column int, wholeLine string, content string,
         value common.ValueExprAst, expr common.ExprAst) common.AssignmentAst {
  return &AssignmentAst{&AstNode{line, column, wholeLine, content}, value, expr};
}


type BlockExprAst struct {
  *ExprAst;
  assignments []common.AssignmentAst;
  expr  common.ExprAst;
}
func (an *BlockExprAst) Assignments() []common.AssignmentAst { return an.assignments; }
func (an *BlockExprAst) Expr() common.ExprAst { return an.expr; }
func NewBlockExprAst(line int, column int, wholeLine string, content string,
     assignments []common.AssignmentAst, expr common.ExprAst) common.BlockExprAst {
  return &BlockExprAst{&ExprAst{&AstNode{line, column, wholeLine, content}},
                         assignments, expr};
}


type PrototypeAst struct {
  *AstNode;
  function string;
  args     []string;
}
func (an *PrototypeAst) Func() string { return an.function; }
func (an *PrototypeAst) Args() []string { return an.args; }
func NewPrototypeAst(line int, column int, wholeLine string, content string,
                     function string, args []string) common.PrototypeAst {
  return &PrototypeAst{&AstNode{line, column, wholeLine, content}, function, args};
}


type FunctionAst struct {
  *AstNode;
  prototype common.PrototypeAst;
  body      common.ExprAst;
}
func (an *FunctionAst) Prototype() common.PrototypeAst { return an.prototype; }
func (an *FunctionAst) Body() common.ExprAst { return an.body; }
func NewFunctionAst(line int, column int, wholeLine string, content string,
     prototype common.PrototypeAst, body common.ExprAst) common.FunctionAst {
  return &FunctionAst{&AstNode{line, column, wholeLine, content}, prototype, body};
}

