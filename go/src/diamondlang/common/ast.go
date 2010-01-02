package common

import (
  "fmt";
)


// Define call types as 'enumeration':
type CallTypeEnum int
const (
  UNKNOWN_CALL = iota;
  FREE_CALL;
  BIND_CALL;
)

// Define data types as 'enumeration':
type DataTypeEnum int
const (
  TYPE_UNKNOWN = iota;
  TYPE_BOOL;
  TYPE_INT;
  TYPE_CHAR;
  TYPE_STRING;
)


func Any2bool(val interface{}) bool {
  pb, ok := val.(*bool);
  if !ok { HandleFatal(fmt.Sprint("Unable to convert to boolean:", val)); }
  return *pb;
}
func Any2int(val interface{}) int64 {
  pi, ok := val.(*int64);
  if !ok { HandleFatal(fmt.Sprint("Unable to convert to integer:", val)); }
  return *pi;
}
func Any2char(val interface{}) byte {
  pc, ok := val.(*byte);
  if !ok { HandleFatal(fmt.Sprint("Unable to convert to character:", val)); }
  return *pc;
}
func Any2string(val interface{}) string {
  ps, ok := val.(*string);
  if !ok { HandleFatal(fmt.Sprint("Unable to convert to string:", val)); }
  return *ps;
}


/// AstNode - Basic interface of all AST nodes.
type AstNode interface {
  SourcePiece() SrcPiece;
}

/// ExprAST - Basic interface of all expression AST nodes.
type ExprAst interface {
  AstNode;
  DataType() DataTypeEnum;
}

/// LiteralExprAst - Interface of expressions for literals like 123, 'c' or "abc".
type LiteralExprAst interface {
  ExprAst;
  Value() interface{};
}

type SubId struct {
  Name      string;
  DataType  DataTypeEnum;
  Protected bool;
}
/// ValueExprAst - Interface of expressions for referencing a value, like 'a'.
type ValueExprAst interface {
  ExprAst;
  ValueName() string;
  SubIds()    []SubId;
}

/// ConstantExprAst - Interface of expressions for referencing a constant, like 'PI'.
type ConstantExprAst interface {
  ExprAst;
  Module()       string;
  ConstantName() string;
  SubIds()       []SubId;
}

/// CallExprAst - Interface of expressions for function calls.
type CallExprAst interface {
  ExprAst;
  Module()       string;
  FuncName()     string;
  CallType()     CallTypeEnum;
  HalfApplied()  bool;
  Args()         []ExprAst;
}

/// AssignmentAst - Interface of assignment statements line: value = expr
type AssignmentAst interface {
  AstNode;
  Value() ValueExprAst;
  Expr() ExprAst;
}

/// BlockExprAst - Interface of expressions for function calls.
type BlockExprAst interface {
  ExprAst;
  Assignments() []AssignmentAst;
  Expr() ExprAst;
}

type Arg struct {
  Name     string;
  DataType DataTypeEnum;
}
// PrototypeAst - Interface of a 'prototype' for a function.
type PrototypeAst interface {
  AstNode;
  FuncName() string;
  FuncDataType() DataTypeEnum;
  Args() []Arg;
}

// FunctionAst - Interface of a function definition itself.
type FunctionAst interface {
  PrototypeAst;
  Body() ExprAst;
}

