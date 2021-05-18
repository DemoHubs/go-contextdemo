// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gvalid implements powerful and useful data/form validation functionality.
package gvalid

import (
	"context"
	"regexp"
	"strings"

	"github.com/gogf/gf/text/gregex"
)

// Refer to Laravel validation: https://laravel.com/docs/5.5/validation#available-validation-rules
//
// All supported rules:
// required             format: required                              brief: Required.
// required-if          format: required-if:field,value,...           brief: Required unless all given field and its value are equal.
// required-unless      format: required-unless:field,value,...       brief: Required unless all given field and its value are not equal.
// required-with        format: required-with:field1,field2,...       brief: Required if any of given fields are not empty.
// required-with-all    format: required-with-all:field1,field2,...   brief: Required if all of given fields are not empty.
// required-without     format: required-without:field1,field2,...    brief: Required if any of given fields are empty.
// required-without-all format: required-without-all:field1,field2,...brief: Required if all of given fields are empty.
// date                 format: date                                  brief: Standard date, like: 2006-01-02, 20060102, 2006.01.02
// date-format          format: date-format:format                    brief: Custom date format.
// email                format: email                                 brief: Email address.
// phone                format: phone                                 brief: Phone number.
// telephone            format: telephone                             brief: Telephone number, like: "XXXX-XXXXXXX"、"XXXX-XXXXXXXX"、"XXX-XXXXXXX"、"XXX-XXXXXXXX"、"XXXXXXX"、"XXXXXXXX"
// passport             format: passport                              brief: Universal passport format rule: Starting with letter, containing only numbers or underscores, length between 6 and 18
// password             format: password                              brief: Universal password format rule1: Containing any visible chars, length between 6 and 18.
// password2            format: password2                             brief: Universal password format rule2: Must meet password rule1, must contain lower and upper letters and numbers.
// password3            format: password3                             brief: Universal password format rule3: Must meet password rule1, must contain lower and upper letters, numbers and special chars.
// postcode             format: postcode                              brief: Postcode number.
// resident-id          format: resident-id                           brief: Resident id number.
// bank-card            format: bank-card                             brief: Bank card nunber.
// qq                   format: qq                                    brief: Tencent QQ number.
// ip                   format: ip                                    brief: IPv4/IPv6.
// ipv4                 format: ipv4                                  brief: IPv4.
// ipv6                 format: ipv6                                  brief: IPv6.
// mac                  format: mac                                   brief: MAC.
// url                  format: url                                   brief: URL.
// domain               format: domain                                brief: Domain.
// length               format: length:min,max                        brief: Length between :min and :max. The length is calculated using unicode string, which means one chinese character or letter both has the length of 1.
// min-length           format: min-length:min                        brief: Length is equal or greater than :min. The length is calculated using unicode string, which means one chinese character or letter both has the length of 1.
// max-length           format: max-length:max                        brief: Length is equal or lesser than :max. The length is calculated using unicode string, which means one chinese character or letter both has the length of 1.
// between              format: between:min,max                       brief: Range between :min and :max. It supports both integer and float.
// min                  format: min:min                               brief: Equal or greater than :min. It supports both integer and float.
// max                  format: max:max                               brief: Equal or lesser than :max. It supports both integer and float.
// json                 format: json                                  brief: JSON.
// integer              format: integer                               brief: Integer.
// float                format: float                                 brief: Float. Note that an integer is actually a float number.
// boolean              format: boolean                               brief: Boolean(1,true,on,yes:true | 0,false,off,no,"":false)
// same                 format: same:field                            brief: Value should be the same as value of field.
// different            format: different:field                       brief: Value should be different from value of field.
// in                   format: in:value1,value2,...                  brief: Value should be in: value1,value2,...
// not-in               format: not-in:value1,value2,...              brief: Value should not be in: value1,value2,...
// regex                format: regex:pattern                         brief: Value should match custom regular expression pattern.

// CustomMsg is the custom error message type,
// like: map[field] => string|map[rule]string
type CustomMsg = map[string]interface{}

// doCheckStructWithParamMapInput is used for struct validation for internal function.
type doCheckStructWithParamMapInput struct {
	Object                          interface{} // Can be type of struct/*struct.
	ParamMap                        interface{} // Validation parameter map. Note that it acts different according attribute `UseParamMapInsteadOfObjectValue`.
	UseParamMapInsteadOfObjectValue bool        // Using `ParamMap` as its validation source instead of values from `Object`.
	CustomRules                     interface{} // Custom validation rules.
	CustomErrorMessageMap           CustomMsg   // Custom error message map for validation rules.
}

// apiNoValidation is an interface that marks current struct not validated by package `gvalid`.
type apiNoValidation interface {
	NoValidation()
}

const (
	// regular expression pattern for single validation rule.
	singleRulePattern   = `^([\w-]+):{0,1}(.*)`
	invalidRulesErrKey  = "invalid_rules"
	invalidParamsErrKey = "invalid_params"
	invalidObjectErrKey = "invalid_object"

	// no validation tag name for struct attribute.
	noValidationTagName = "nv"
)

var (
	defaultValidator     = New()                            // defaultValidator is the default validator for package functions.
	structTagPriority    = []string{"gvalid", "valid", "v"} // structTagPriority specifies the validation tag priority array.
	aliasNameTagPriority = []string{"param", "params", "p"} // aliasNameTagPriority specifies the alias tag priority array.

	// all internal error keys.
	internalErrKeyMap = map[string]string{
		invalidRulesErrKey:  invalidRulesErrKey,
		invalidParamsErrKey: invalidParamsErrKey,
		invalidObjectErrKey: invalidObjectErrKey,
	}
	// regular expression object for single rule
	// which is compiled just once and of repeatable usage.
	ruleRegex, _ = regexp.Compile(singleRulePattern)

	// mustCheckRulesEvenValueEmpty specifies some rules that must be validated
	// even the value is empty (nil or empty).
	mustCheckRulesEvenValueEmpty = map[string]struct{}{
		"required":             {},
		"required-if":          {},
		"required-unless":      {},
		"required-with":        {},
		"required-with-all":    {},
		"required-without":     {},
		"required-without-all": {},
		//"same":                 {},
		//"different":            {},
		//"in":                   {},
		//"not-in":               {},
		//"regex":                {},
	}
	// allSupportedRules defines all supported rules that is used for quick checks.
	allSupportedRules = map[string]struct{}{
		"required":             {},
		"required-if":          {},
		"required-unless":      {},
		"required-with":        {},
		"required-with-all":    {},
		"required-without":     {},
		"required-without-all": {},
		"date":                 {},
		"date-format":          {},
		"email":                {},
		"phone":                {},
		"phone-loose":          {},
		"telephone":            {},
		"passport":             {},
		"password":             {},
		"password2":            {},
		"password3":            {},
		"postcode":             {},
		"resident-id":          {},
		"bank-card":            {},
		"qq":                   {},
		"ip":                   {},
		"ipv4":                 {},
		"ipv6":                 {},
		"mac":                  {},
		"url":                  {},
		"domain":               {},
		"length":               {},
		"min-length":           {},
		"max-length":           {},
		"between":              {},
		"min":                  {},
		"max":                  {},
		"json":                 {},
		"integer":              {},
		"float":                {},
		"boolean":              {},
		"same":                 {},
		"different":            {},
		"in":                   {},
		"not-in":               {},
		"regex":                {},
	}
	// boolMap defines the boolean values.
	boolMap = map[string]struct{}{
		"1":     {},
		"true":  {},
		"on":    {},
		"yes":   {},
		"":      {},
		"0":     {},
		"false": {},
		"off":   {},
		"no":    {},
	}
)

// Check checks single value with specified rules.
// It returns nil if successful validation.
//
// The parameter `value` can be any type of variable, which will be converted to string
// for validation.
// The parameter `rules` can be one or more rules, multiple rules joined using char '|'.
// The parameter `messages` specifies the custom error messages, which can be type of:
// string/map/struct/*struct.
// The optional parameter `params` specifies the extra validation parameters for some rules
// like: required-*、same、different, etc.
func Check(ctx context.Context, value interface{}, rules string, messages interface{}, params ...interface{}) Error {
	return defaultValidator.Ctx(ctx).Check(value, rules, messages, params...)
}

// CheckMap validates map and returns the error result. It returns nil if with successful validation.
//
// The parameter `rules` can be type of []string/map[string]string. It supports sequence in error result
// if `rules` is type of []string.
// The optional parameter `messages` specifies the custom error messages for specified keys and rules.
func CheckMap(ctx context.Context, params interface{}, rules interface{}, messages ...CustomMsg) Error {
	return defaultValidator.Ctx(ctx).CheckMap(params, rules, messages...)
}

// CheckStruct validates strcut and returns the error result.
//
// The parameter `object` should be type of struct/*struct.
// The parameter `rules` can be type of []string/map[string]string. It supports sequence in error result
// if `rules` is type of []string.
// The optional parameter `messages` specifies the custom error messages for specified keys and rules.
func CheckStruct(ctx context.Context, object interface{}, rules interface{}, messages ...CustomMsg) Error {
	return defaultValidator.Ctx(ctx).CheckStruct(object, rules, messages...)
}

// CheckStructWithParamMap validates struct with given parameter map and returns the error result.
//
// The parameter `object` should be type of struct/*struct.
// The parameter `rules` can be type of []string/map[string]string. It supports sequence in error result
// if `rules` is type of []string.
// The optional parameter `messages` specifies the custom error messages for specified keys and rules.
func CheckStructWithParamMap(ctx context.Context, object interface{}, paramMap interface{}, rules interface{}, messages ...CustomMsg) Error {
	return defaultValidator.Ctx(ctx).CheckStructWithParamMap(object, paramMap, rules, messages...)
}

// parseSequenceTag parses one sequence tag to field, rule and error message.
// The sequence tag is like: [alias@]rule[...#msg...]
func parseSequenceTag(tag string) (field, rule, msg string) {
	// Complete sequence tag.
	// Example: name@required|length:2,20|password3|same:password1#||密码强度不足|两次密码不一致
	match, _ := gregex.MatchString(`\s*((\w+)\s*@){0,1}\s*([^#]+)\s*(#\s*(.*)){0,1}\s*`, tag)
	return strings.TrimSpace(match[2]), strings.TrimSpace(match[3]), strings.TrimSpace(match[5])
}
