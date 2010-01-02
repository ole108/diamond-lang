package parser

import (
  "diamondlang/common";
)


type AstNode struct {
  srcPiece  common.SrcPiece;
}
func (an *AstNode) SourcePiece() common.SrcPiece { return an.srcPiece; }


type ExprAst struct {
  *AstNode;
  dataType  common.DataTypeEnum;
}
func (an *ExprAst) DataType() common.DataTypeEnum { return an.dataType; }


type LiteralExprAst struct {
  *ExprAst;
  value  interface{};
}
func (an *LiteralExprAst) Value() interface{} { return an.value; }
func NewLiteralExprAst(piece common.SrcPiece, dataType common.DataTypeEnum,
                       value interface{}) common.LiteralExprAst {
  return &LiteralExprAst{&ExprAst{&AstNode{piece}, dataType}, value};
}


type ValueExprAst struct {
  *ExprAst;
  value  string;
  subs   []common.SubId;
}
func (an *ValueExprAst) ValueName() string { return an.value; }
func (an *ValueExprAst) SubIds() []common.SubId { return an.subs; }
func NewValueExprAst(piece common.SrcPiece, value string,
                     subs []common.SubId) common.ValueExprAst {
  return &ValueExprAst{&ExprAst{&AstNode{piece}, common.TYPE_UNKNOWN}, value, subs};
}


type ConstantExprAst struct {
  *ExprAst;
  module   string;
  constant string;
  subs     []common.SubId;
}
func (an *ConstantExprAst) Module() string { return an.module; }
func (an *ConstantExprAst) ConstantName() string { return an.constant; }
func (an *ConstantExprAst) SubIds() []common.SubId { return an.subs; }
func NewConstantExprAst(piece common.SrcPiece, module string, constant string,
                        subs []common.SubId) common.ConstantExprAst {
  return &ConstantExprAst{&ExprAst{&AstNode{piece}, common.TYPE_UNKNOWN},
                          module, constant, subs};
}


type CallExprAst struct {
  *ExprAst;
  module      string;
  function    string;
  typ         common.CallTypeEnum;
  halfApplied bool;
  args        []common.ExprAst;
}
func (an *CallExprAst) Module() string { return an.module; }
func (an *CallExprAst) FuncName() string { return an.function; }
func (an *CallExprAst) CallType() common.CallTypeEnum { return an.typ; }
func (an *CallExprAst) HalfApplied() bool { return an.halfApplied; }
func (an *CallExprAst) Args() []common.ExprAst { return an.args; }
func NewCallExprAst(piece common.SrcPiece, module string, function string,
          typ common.CallTypeEnum, halfApplied bool, args []common.ExprAst)
                 common.CallExprAst {
  return &CallExprAst{&ExprAst{&AstNode{piece}, common.TYPE_UNKNOWN},
                      module, function, typ, halfApplied, args};
}


type AssignmentAst struct {
  *AstNode;
  value common.ValueExprAst;
  expr  common.ExprAst;
}
func (an *AssignmentAst) Value() common.ValueExprAst { return an.value; }
func (an *AssignmentAst) Expr() common.ExprAst { return an.expr; }
func NewAssignmentAst(piece common.SrcPiece, value common.ValueExprAst,
                      expr common.ExprAst) common.AssignmentAst {
  return &AssignmentAst{&AstNode{piece}, value, expr};
}


type BlockExprAst struct {
  *ExprAst;
  assignments []common.AssignmentAst;
  expr  common.ExprAst;
}
func (an *BlockExprAst) Assignments() []common.AssignmentAst { return an.assignments; }
func (an *BlockExprAst) Expr() common.ExprAst { return an.expr; }
func NewBlockExprAst(piece common.SrcPiece, assignments []common.AssignmentAst,
                     expr common.ExprAst) common.BlockExprAst {
  return &BlockExprAst{&ExprAst{&AstNode{piece}, common.TYPE_UNKNOWN}, assignments, expr};
}


type PrototypeAst struct {
  *AstNode;
  function string;
  dataType  common.DataTypeEnum;
  args     []common.Arg;
}
func (an *PrototypeAst) FuncName() string { return an.function; }
func (an *PrototypeAst) FuncDataType() common.DataTypeEnum { return an.dataType; }
func (an *PrototypeAst) Args() []common.Arg { return an.args; }
func NewPrototypeAst(piece common.SrcPiece, function string,
        dataType common.DataTypeEnum, args []common.Arg) common.PrototypeAst {
  return &PrototypeAst{&AstNode{piece}, function, dataType, args};
}


type FunctionAst struct {
  common.PrototypeAst;
  body  common.ExprAst;
}
func (an *FunctionAst) Body() common.ExprAst { return an.body; }
func NewFunctionAst(piece common.SrcPiece, prototype common.PrototypeAst,
                    body common.ExprAst) common.FunctionAst {
  return &FunctionAst{prototype, body};
}

