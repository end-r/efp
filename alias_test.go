package efp

import (
	"fmt"
	"testing"

	"github.com/end-r/goutil"
)

func TestAliasingStandardAliases(t *testing.T) {
	p := createPrototypeParserString("")
	goutil.Assert(t, len(p.prototype.valueAliases) == len(standards), "wrong number of standard aliases")
}

func TestAliasingStandardInt(t *testing.T) {
	p := createPrototypeParserString("")
	goutil.AssertNow(t, p.prototype.valueAliases != nil, "text aliases is nil")
	regex := p.prototype.valueAliases["int"].values[0].value
	goutil.Assert(t, regex.MatchString("99"), "int regex didn't match")
	goutil.Assert(t, regex.MatchString("0"), "int regex didn't match")
	goutil.Assert(t, regex.MatchString("-100"), "int regex didn't match")

	goutil.Assert(t, !regex.MatchString("aaa"), "int regex did match")
	goutil.Assert(t, !regex.MatchString("-0"), "int regex did match")
}

func TestAliasingStandardUInt(t *testing.T) {
	p := createPrototypeParserString("")
	regex := p.prototype.valueAliases["uint"].values[0].value
	goutil.Assert(t, regex.MatchString("99"), "uint regex didn't match")
	goutil.Assert(t, regex.MatchString("0"), "uint regex didn't match")

	goutil.Assert(t, !regex.MatchString("-100"), "uint regex did match")
	goutil.Assert(t, !regex.MatchString("-0"), "uint regex did match")
}

func TestAliasingStandardFloat(t *testing.T) {
	p := createPrototypeParserString("")
	regex := p.prototype.valueAliases["float"].values[0].value
	goutil.Assert(t, regex.MatchString("99"), "float regex didn't match")
	goutil.Assert(t, regex.MatchString("0"), "float regex didn't match")
	goutil.Assert(t, regex.MatchString("-100"), "float regex didn't match")
	goutil.Assert(t, regex.MatchString("0.5"), "float regex didn't match")

	goutil.Assert(t, !regex.MatchString("aaa"), "float regex did match")
	goutil.Assert(t, !regex.MatchString("-0"), "float regex did match")
}

func TestAliasingStandardBool(t *testing.T) {
	p := createPrototypeParserString("")
	regex := p.prototype.valueAliases["bool"].values[0].value
	goutil.Assert(t, regex.MatchString("true"), "bool regex didn't match")
	goutil.Assert(t, regex.MatchString("false"), "bool regex didn't match")

	goutil.Assert(t, !regex.MatchString("tru"), "bool regex did match")
	goutil.Assert(t, !regex.MatchString("flase"), "bool regex did match")
}

func TestAliasingFieldAlias(t *testing.T) {
	p, errs := PrototypeString(`
        alias x = name : string
        x`)
	goutil.AssertNow(t, errs == nil, "errs are not nil")
	goutil.AssertNow(t, p.fields["name"] != nil, "name is nil")
}

func TestAliasingElementAlias(t *testing.T) {
	p, errs := PrototypeString(`
        alias x = name {

		}
        x`)
	goutil.Assert(t, errs == nil, "errs are not nil")
	goutil.Assert(t, p.elementAliases["x"] != nil, "alias not found")
	goutil.AssertNow(t, p.elements != nil, "elements is nil")
	goutil.AssertNow(t, p.elements["name"] != nil, "name is nil")
}

/*func TestAliasingTextAliasMax(t *testing.T) {
	const limit = 2
	const regex = "[a-z]{3}"
	p, errs := PrototypeString(`
        alias LIMIT = 2
			<LIMIT:"[a-z]{3}":LIMIT> : [LIMIT:string:LIMIT]
		`)
	goutil.Assert(t, len(p.valueAliases) == 1+len(standards), "wrong text alias length")
	goutil.Assert(t, errs == nil, "errs should be nil")
	goutil.Assert(t, p.valueAliases["LIMIT"].value == "2", "wrong limit value")
	goutil.AssertNow(t, p.fields[regex] != nil, "field is nil")
	goutil.AssertNow(t, len(p.fields[regex].types) == 1, "types must not be nil")
	goutil.AssertNow(t, p.fields[regex].types[0].isArray, "not array")
	goutil.AssertNow(t, p.fields[regex].types[0].types[0].value.String() == standards["string"].value, "incorrect regex")
	goutil.AssertNow(t, p.fields[regex].types[0].max == limit, "incorrect value max")
	goutil.AssertNow(t, p.fields[regex].types[0].min == limit, "incorrect value min")
	goutil.AssertNow(t, p.fields[regex].key.min == limit, "incorrect key min")
	goutil.AssertNow(t, p.fields[regex].key.max == limit, "incorrect key max")
}*/

func TestAliasingTextAliasValue(t *testing.T) {
	p, errs := PrototypeString(`
        alias x = string
        name : x`)
	_, ok := p.valueAliases["x"]
	goutil.AssertNow(t, ok, "alias is nil")
	goutil.AssertNow(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.fields["name"] != nil, "name is nil")
	goutil.AssertNow(t, p.fields["name"].types[0].value.String() == standards["string"], "wrong value")
}

func TestAliasingDoubleIndirection(t *testing.T) {
	p, errs := PrototypeString(`
        alias y = string
        alias x = name : y
        x`)
	goutil.Assert(t, len(p.valueAliases) == 1+len(standards), fmt.Sprintf(`wrong text alias length %d (expected %d)`,
		len(p.valueAliases), 1+len(standards)))
	goutil.Assert(t, errs == nil, "errs must be nil")
	goutil.AssertNow(t, p.fields["name"] != nil, "name is nil")
	goutil.AssertNow(t, len(p.fields["name"].types) == 1, "wrong type length")
	goutil.AssertNow(t, p.fields["name"].types[0].value != nil, "regex is nil")
	goutil.AssertNow(t, p.fields["name"].types[0].value.String() == standards["string"], "wrong value")
}

func TestAliasingDuplicateTextAlias(t *testing.T) {
	_, errs := PrototypeString(`alias y = 25 alias y = 30`)
	goutil.Assert(t, errs != nil, "there should be an error")
}

func TestAliasingDuplicateFieldAlias(t *testing.T) {
	_, errs := PrototypeString(`alias y = name : string alias y = name : string`)
	goutil.Assert(t, errs != nil, "there should be an error")
}

func TestAliasingDuplicateElementAlias(t *testing.T) {
	_, errs := PrototypeString(`alias y = name{} alias y = name{}`)
	goutil.Assert(t, errs != nil, "there should be an error")
}

func TestAliasingDuplicateMixedAlias(t *testing.T) {
	_, errs := PrototypeString(`alias y = name{} alias y = 10`)
	goutil.Assert(t, errs != nil, "there should be an error")
}

// test that element recursion is allowed
func TestAliasingRecursionValid(t *testing.T) {
	p, errs := PrototypeString(`
        alias p = x {
            p
        }
        p`)
	goutil.AssertNow(t, errs == nil, "errs should be nil")
	goutil.AssertNow(t, p.elements["x"] != nil, "x is nil")
}

// test that field recursion is disallowed
func TestAliasingRecursionInvalid(t *testing.T) {
	_, errs := PrototypeString(`alias x = x
		x`)
	goutil.AssertNow(t, errs != nil, "errs should not be nil")
}
