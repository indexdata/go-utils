package utils_test

import (
	"encoding/xml"
	"io"
	"os"
	"testing"

	"github.com/indexdata/go-utils/utils"
	"github.com/stretchr/testify/assert"
)

type Root struct {
	XMLName        xml.Name           `xml:"root"`
	Date           utils.XSDDateTime  `xml:"date,omitempty"`
	DatePtr        *utils.XSDDateTime `xml:"datePtr,omitempty"`
	DateMissing    utils.XSDDateTime  `xml:"dateMissing,omitempty"`
	DateMissingPtr *utils.XSDDateTime `xml:"dateMissingPtr,omitempty"`
	DateEmpty      utils.XSDDateTime  `xml:"dateEmpty,omitempty"`
	DateEmptyPtr   *utils.XSDDateTime `xml:"dateEmptyPtr,omitempty"`
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
		t.Fatalf("Failure marshaling test file %v", err)
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
