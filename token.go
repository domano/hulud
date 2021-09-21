package main

import (
	"bufio"
	"unicode"
	"unicode/utf8"
)

const (
	Enum                         TokenType = "enum"
	Name                                   = ""
	CurlyStart                             = "{"
	CurlyEnd                               = "}"
	Comma                                  = ","
	Directive                              = "directive"
	DirectivePrefix                        = "@"
	BraceStart                             = "("
	BraceEnd                               = ")"
	On                                     = "on"
	UnionOr                                = "|"
	Union                                  = "union"
	LocationQuery                          = "QUERY"
	LocationMutation                       = "MUTATION"
	LocationSubscription                   = "SUBSCRIPTION"
	LocationField                          = "FIELD"
	LocationFragmentDefinition             = "FRAGMENT_DEFINITION"
	LocationFragmentSpread                 = "FRAGMENT_SPREAD"
	LocationInlineFragment                 = "INLINE_FRAGMENT"
	LocationSchema                         = "SCHEMA"
	LocationScalar                         = "SCALAR"
	LocationObject                         = "OBJECT"
	LocationFieldDefinition                = "FIELD_DEFINITION"
	LocationArgumentDefinition             = "ARGUMENT_DEFINITION"
	LocationInterface                      = "INTERFACE"
	LocationUnion                          = "UNION"
	LocationEnum                           = "ENUM"
	LocationEnumValue                      = "ENUM_VALUE"
	LocationInputObject                    = "INPUT_OBJECT"
	LocationInputFieldDefinition           = "INPUT_FIELD_DEFINITION"
	ParameterName                          = ""
	Colon                                  = ":"
	Type                                   = ""
	Required                               = "!"
	Scalar                                 = "scalar"
	Interface                              = "interface"
	FieldName                              = ""
	BracketStart                           = "["
	BracketEnd                             = "]"
	Query                                  = "query"
	Subscription                           = "subscription"
	Mutation                               = "Mutation"
	Equal                                  = "="
	True                                   = "true"
	False                                  = "false"
	String                                 = "String"
	Int                                    = "Int"
	Float                                  = "Float"
	Boolean                                = "Boolean"
	ID                                     = "ID"
	Input                                  = "input"
	Skip                                   = "skip"
	Include                                = "include"
	deprecated                             = "deprecated"
	DoubleUnderscore                       = "__"
	Comment                                = "#"
)

type TokenType string

// ScanGraphQLToken implements bufio.SplitFunc
func ScanGraphQLToken(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Scan first char and determine if it is a token
	advance, token, err, done := nextRune(data, atEOF)
	if done {
		return advance, token, err
	}
	if isToken(token) {
		return advance, token, err // We have a one-character token, we can return it and stop scanning for now
	}
	r, _ := utf8.DecodeRune(token)
	if unicode.IsSpace(r) {
		return advance, nil, err // We skip spaces
	}

	if len(data) <= advance { // If there is nothing else to scan
		return 0, nil, err // Continue scanning
	}

	// Loop over remaining chars until we hit another token
	var currentPos = advance
	for {
		a, t, err, done := nextRune(data[currentPos:], atEOF)
		if done {
			return a, t, err
		}

		if isToken(t) { // If this new rune is a token, then the bytes that came before must be a multi character token
			return currentPos, token, err
		}

		token = append(token, t...)

		if isToken(token) { // If the added bytes are a known token we can return it
			return currentPos + 1, token, err
		}

		currentPos += a // Lets keep looking in the slice for another token
	}

	return 0, nil, nil
}

func nextRune(data []byte, atEOF bool) (int, []byte, error, bool) {
	a, t, err := bufio.ScanRunes(data, atEOF)
	if err != nil || t == nil {
		return a, t, err, true
	}
	r, _ := utf8.DecodeRune(t)
	if r == utf8.RuneError {
		return a, t, err, true
	}
	return a, t, err, false
}

func isToken(token []byte) bool {
	tokenType := TokenType(token)
	switch tokenType {
	case Enum, CurlyStart, CurlyEnd, Comma, Directive, DirectivePrefix, BraceStart, BraceEnd, On, UnionOr, Union, LocationQuery, LocationMutation, LocationSubscription, LocationField, LocationFragmentDefinition, LocationFragmentSpread, LocationInlineFragment, LocationSchema, LocationScalar, LocationObject, LocationFieldDefinition, LocationArgumentDefinition, LocationInterface, LocationUnion, LocationEnum, LocationEnumValue, LocationInputObject, LocationInputFieldDefinition, Colon, Required, Scalar, Interface, BracketStart, BracketEnd, Query, Subscription, Mutation, Equal, True, False, String, Int, Float, Boolean, ID, Input, Skip, Include, deprecated, DoubleUnderscore, Comment:
		// For all types with predetermined names we now know it is a token
		return true
	default:
		return false
	}

}
