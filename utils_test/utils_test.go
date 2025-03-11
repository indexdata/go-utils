package utils_test

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/indexdata/go-utils/utils"
	"github.com/stretchr/testify/assert"
)

type Root struct {
	XMLName        xml.Name           `xml:"http://some.com root" json:"-"`
	VersionNS      *utils.PrefixAttr  `xml:"http://some.com/z version,attr,omitempty" json:"@versionNS,omitempty"`
	Version        *utils.PrefixAttr  `xml:"version,attr,omitempty" json:"@version,omitempty"`
	Date           utils.XSDDateTime  `xml:"date,omitempty" json:"date,omitempty"`
	DatePtr        *utils.XSDDateTime `xml:"datePtr,omitempty" json:"datePtr,omitempty"`
	DateMissing    utils.XSDDateTime  `xml:"dateMissing,omitempty" json:"dateMissing,omitempty"`
	DateMissingPtr *utils.XSDDateTime `xml:"dateMissingPtr,omitempty" json:"dateMissingPtr,omitempty"`
	DateEmpty      utils.XSDDateTime  `xml:"dateEmpty,omitempty" json:"dateEmpty,omitempty"`
	DateEmptyPtr   *utils.XSDDateTime `xml:"dateEmptyPtr,omitempty" json:"dateEmptyPtr,omitempty"`
	Child          Child              `xml:"child" json:"child,omitempty"`
}

type Child struct {
	XMLName  xml.Name `xml:"http://other.com child" json:"-"`
	Text     string   `xml:",chardata" json:"#text,omitempty"`
	Revision string   `xml:"http://some.com/z revision,attr,omitempty" json:"@revision,omitempty"`
}

func TestGetEnvInt(t *testing.T) {
	os.Unsetenv("ENV_INT")
	//unset
	val, _ := utils.GetEnvInt("ENV_INT", 1)
	assert.Equal(t, 1, val)
	os.Setenv("ENV_INT", "2")
	//set
	val, _ = utils.GetEnvInt("ENV_INT", 1)
	assert.Equal(t, 2, val)
	//set to zero val
	os.Setenv("ENV_INT", "0")
	val, _ = utils.GetEnvInt("ENV_INT", 1)
	assert.Equal(t, 0, val)
	//empty
	os.Setenv("ENV_INT", "")
	val, _ = utils.GetEnvInt("ENV_INT", 1)
	assert.Equal(t, 1, val)
	//malformed
	os.Setenv("ENV_INT", "error")
	vale, err := utils.GetEnvInt("ENV_INT", 1)
	assert.Equal(t, 1, vale)
	assert.EqualError(t, err, "strconv.Atoi: parsing \"error\": invalid syntax")
	os.Unsetenv("ENV_INT")
}

func TestGetEnvBool(t *testing.T) {
	os.Unsetenv("ENV_BOOL")
	//unset
	val, _ := utils.GetEnvBool("ENV_BOOL", true)
	assert.Equal(t, true, val)
	//set
	os.Setenv("ENV_BOOL", "true")
	val, _ = utils.GetEnvBool("ENV_BOOL", true)
	assert.Equal(t, true, val)
	//set to zero value
	os.Setenv("ENV_BOOL", "false")
	val, _ = utils.GetEnvBool("ENV_BOOL", true)
	assert.Equal(t, false, val)
	//empty
	os.Setenv("ENV_BOOL", "")
	val, _ = utils.GetEnvBool("ENV_BOOL", true)
	assert.Equal(t, true, val)
	//malformed
	os.Setenv("ENV_BOOL", "error")
	vale, err := utils.GetEnvBool("ENV_BOOL", true)
	assert.Equal(t, true, vale)
	assert.EqualError(t, err, "strconv.ParseBool: parsing \"error\": invalid syntax")
	os.Unsetenv("ENV_BOOL")
}

func TestGetEnvAny(t *testing.T) {
	os.Unsetenv("ENV_BOOL")
	//unset
	val, _ := utils.GetEnvAny("ENV_BOOL", true, func(env string) (bool, error) {
		return strconv.ParseBool(env)
	})
	assert.Equal(t, true, val)
	//set
	os.Setenv("ENV_BOOL", "true")
	val, _ = utils.GetEnvAny("ENV_BOOL", true, func(env string) (bool, error) {
		return strconv.ParseBool(env)
	})
	assert.Equal(t, true, val)
	//set to zero value
	os.Setenv("ENV_BOOL", "false")
	val, _ = utils.GetEnvAny("ENV_BOOL", true, func(env string) (bool, error) {
		return strconv.ParseBool(env)
	})
	assert.Equal(t, false, val)
	//empty
	os.Setenv("ENV_BOOL", "")
	val, _ = utils.GetEnvAny("ENV_BOOL", true, func(env string) (bool, error) {
		return strconv.ParseBool(env)
	})
	assert.Equal(t, true, val)
	//malformed
	os.Setenv("ENV_BOOL", "error")
	vale, err := utils.GetEnvAny("ENV_BOOL", true, func(env string) (bool, error) {
		return strconv.ParseBool(env)
	})
	assert.Equal(t, true, vale)
	assert.EqualError(t, err, "strconv.ParseBool: parsing \"error\": invalid syntax")
	os.Unsetenv("ENV_BOOL")
}

func TestXMLUtils(t *testing.T) {
	xf, err := os.Open("test.xml")
	if err != nil {
		t.Fatalf("Failure parsing test file: %v", err)
	}
	xd, _ := io.ReadAll(xf)
	var doc Root
	err = xml.Unmarshal(xd, &doc)
	if err != nil {
		t.Fatalf("Failure unmarshaling test file %v", err)
	}
	actual, _ := xml.MarshalIndent(&doc, "", "  ")
	actual = append(actual, '\n')
	expectedFile, err := os.Open("test.xml.out")
	if err != nil {
		t.Fatalf("Failure parsing test file: %v", err)
	}
	expected, _ := io.ReadAll(expectedFile)
	assert.Equal(t, string(expected), string(actual))
}

func TestXMLToJSON(t *testing.T) {
	xf, err := os.Open("test.xml")
	if err != nil {
		t.Fatalf("Failure parsing test file: %v", err)
	}
	xd, _ := io.ReadAll(xf)
	var doc Root
	err = xml.Unmarshal(xd, &doc)
	if err != nil {
		t.Fatalf("Failure marshaling test file %v", err)
	}
	actual, _ := json.MarshalIndent(&doc, "", "  ")
	actual = append(actual, '\n')
	expectedFile, err := os.Open("test.json")
	if err != nil {
		t.Fatalf("Failure parsing test file: %v", err)
	}
	expected, _ := io.ReadAll(expectedFile)
	assert.Equal(t, string(expected), string(actual))
}

func TestJSONUnThenMarshal(t *testing.T) {
	xf, err := os.Open("test.json")
	if err != nil {
		t.Fatalf("Failure parsing test file: %v", err)
	}
	xd, _ := io.ReadAll(xf)
	var doc Root
	err = json.Unmarshal(xd, &doc)
	if err != nil {
		t.Fatalf("Failure marshaling test file %v", err)
	}
	actual, _ := json.MarshalIndent(&doc, "", "  ")
	actual = append(actual, '\n')
	expectedFile, err := os.Open("test.json")
	if err != nil {
		t.Fatalf("Failure parsing test file: %v", err)
	}
	expected, _ := io.ReadAll(expectedFile)
	assert.Equal(t, string(expected), string(actual))
}

func TestExtractDecimal(t *testing.T) {
	in := []string{"1800.25", "1800,25", "$1,800.25", "1,800.25 USD", "€ 1.800,25", "1.800,25 EUR"}
	expected := "1800.25"
	exp1 := 180025
	exp2 := 2
	//autodetect
	for _, s := range in {
		actual, act1, act2 := utils.ExtractDecimal(s, -1)
		assert.Equal(t, expected, actual)
		assert.Equal(t, exp1, act1)
		assert.Equal(t, exp2, act2)
		formatted := utils.FormatDecimal(act1, act2)
		assert.Equal(t, expected, formatted)
	}
	//up to two decimal places
	for _, s := range in {
		actual, act1, act2 := utils.ExtractDecimal(s, 2)
		assert.Equal(t, expected, actual)
		assert.Equal(t, exp1, act1)
		assert.Equal(t, exp2, act2)
		formatted := utils.FormatDecimal(act1, act2)
		assert.Equal(t, expected, formatted)
	}
	expected = "180025"
	exp1 = 180025
	exp2 = 0
	//up to one decimal place
	for _, s := range in {
		actual, act1, act2 := utils.ExtractDecimal(s, 1)
		assert.Equal(t, expected, actual)
		assert.Equal(t, exp1, act1)
		assert.Equal(t, exp2, act2)
		formatted := utils.FormatDecimal(act1, act2)
		assert.Equal(t, expected, formatted)
	}
	//autodetect
	in = []string{"180.025", "180,025.", ".180,025"}
	expected = "180.025"
	exp1 = 180025
	exp2 = 3
	for _, s := range in {
		actual, act1, act2 := utils.ExtractDecimal(s, -1)
		assert.Equal(t, expected, actual)
		assert.Equal(t, exp1, act1)
		assert.Equal(t, exp2, act2)
		formatted := utils.FormatDecimal(act1, act2)
		assert.Equal(t, expected, formatted)
	}
	//up to to two decimal places
	expected = "180025"
	exp1 = 180025
	exp2 = 0
	for _, s := range in {
		actual, act1, act2 := utils.ExtractDecimal(s, 2)
		assert.Equal(t, expected, actual)
		assert.Equal(t, exp1, act1)
		assert.Equal(t, exp2, act2)
		formatted := utils.FormatDecimal(act1, act2)
		assert.Equal(t, expected, formatted)
	}
	//up to three decimal place
	expected = "180.025"
	exp1 = 180025
	exp2 = 3
	for _, s := range in {
		actual, act1, act2 := utils.ExtractDecimal(s, 3)
		assert.Equal(t, expected, actual)
		assert.Equal(t, exp1, act1)
		assert.Equal(t, exp2, act2)
		formatted := utils.FormatDecimal(act1, act2)
		assert.Equal(t, expected, formatted)
	}
	//mixed
	in = []string{"180.025", "180,.025.", ".180 025"}
	expected = "180025"
	exp1 = 180025
	exp2 = 0
	for _, s := range in {
		actual, act1, act2 := utils.ExtractDecimal(s, 2)
		assert.Equal(t, expected, actual)
		assert.Equal(t, exp1, act1)
		assert.Equal(t, exp2, act2)
		formatted := utils.FormatDecimal(act1, act2)
		assert.Equal(t, expected, formatted)
	}
	//simple tests
	out := utils.FormatDecimal(0, 0)
	assert.Equal(t, out, "0")
}

func TestNewPrefixAttrNS(t *testing.T) {
	prefixAttr := utils.NewPrefixAttrNS("http://www.mygoofy.org/", "x", "v")
	assert.Equal(t, "http://www.mygoofy.org/", prefixAttr.Name.Space)
	assert.Equal(t, "x", prefixAttr.Name.Local)
	assert.Equal(t, "v", prefixAttr.Value)
}

func TestNewPrefixAttr(t *testing.T) {
	prefixAttr := utils.NewPrefixAttr("x", "v")
	assert.Empty(t, prefixAttr.Name.Space)
	assert.Equal(t, "x", prefixAttr.Name.Local)
	assert.Equal(t, "v", prefixAttr.Value)
}

func TestUrlWithQuery(t *testing.T) {
	in, _ := url.Parse("http://example.com")
	exp, _ := url.Parse("http://example.com?a=x&b=y")
	actual := utils.UrlWithQuery(*in, "a", "x", "b", "y")
	assert.Equal(t, exp.String(), actual.String())
	actual = utils.UrlWithQuery(*in, "a", "x", "b", "y", "c")
	assert.Equal(t, exp.String(), actual.String())
	actual = utils.UrlWithQuery(*in, "a", "x", "b", "y", "", "")
	assert.Equal(t, exp.String(), actual.String())
}
