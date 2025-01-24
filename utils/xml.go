package utils

import (
	"encoding/xml"
	"reflect"
	"strings"
	"sync"
	"time"
)

const XMLNS string = "xmlns"

var defaultNS string
var spaceToPrefix sync.Map
var prefixToSpace sync.Map
var attrDefaults sync.Map

func NSPrefix(prefix string, uri string) {
	prefixToSpace.Store(prefix, uri)
	spaceToPrefix.Store(uri, prefix)
}

func NSDefault(uri string) {
	defaultNS = uri
}

func AttrDefault(attr string, val string) {
	attrDefaults.Store(attr, val)
}

// Go's encoding/xml does not handle namespace prefixes and by extension cannot properly
// handle namespaces for qualified attributes. This implementation adds prefix support to attributes
// which allows to marshal valid XML documents.
// Note 1: Similar approach can be employed to deal with
// XML element prefixes but would require a struct for every simple string which is inefficient.
// Note 2: use as a pointer to properly handle "omitempty" tag
type PrefixAttr struct {
	xml.Attr
}

func NewPrefixAttr(name string, value string) *PrefixAttr {
	return &PrefixAttr{Attr: xml.Attr{Name: xml.Name{Local: name}, Value: value}}
}

func (pxAttr *PrefixAttr) UnmarshalXMLAttr(attr xml.Attr) error {
	ns := attr.Name.Space
	name := attr.Name.Local
	val := attr.Value
	//namespace registration
	if ns == XMLNS {
		spaceToPrefix.Store(val, name)
	}
	if ns == "" && name == XMLNS {
		defaultNS = val
	}
	pxAttr.Attr = attr
	return nil
}

func (pxAttr *PrefixAttr) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var ns, locName, value string
	if pxAttr != nil && !reflect.ValueOf(*pxAttr).IsZero() {
		ns = pxAttr.Name.Space
		locName = pxAttr.Name.Local
		value = pxAttr.Value
	} else { //value not set and omitempty=false
		ns = name.Space
		locName = name.Local
	}
	var qName string
	if ns == XMLNS {
		qName = "xmlns:" + locName
		if value == "" {
			v, ok := prefixToSpace.Load(locName)
			if ok {
				value = v.(string)
			}
		}
	} else {
		if ns == "" && locName != XMLNS {
			ns = defaultNS
		}
		prefix, ok := spaceToPrefix.Load(ns)
		if ok {
			qName = prefix.(string) + ":" + locName
		} else {
			qName = locName
		}
	}
	if value == "" {
		if locName == XMLNS {
			value = defaultNS
		} else {
			v, ok := attrDefaults.Load(locName)
			if ok {
				value = v.(string)
			}
		}
	}
	return xml.Attr{Name: xml.Name{Space: "", Local: qName}, Value: value}, nil
}

func (pxAttr *PrefixAttr) UnmarshalText(text []byte) error {
	pxAttr.Value = string(text)
	return nil
}

func (pxAttr *PrefixAttr) MarshalText() ([]byte, error) {
	return []byte(pxAttr.Value), nil
}

const ISO8601FormatInUTC = "2006-01-02T15:04:05.999Z"
const ISO8601FormatInLocal = "2006-01-02T15:04:05.999"
const ISO8601FormatInTZ = "2006-01-02T15:04:05.999-0700"
const ISO8601FormatOut = "2006-01-02T15:04:05.000Z"
const dateTimeLen = len("2006-01-02T15:04:05")

type XSDDateTime struct {
	time.Time
}

func hasTzOffset(timestamp string) bool {
	return strings.LastIndex(timestamp, "+") >= dateTimeLen ||
		strings.LastIndex(timestamp, "-") >= dateTimeLen
}

func tzOffsetPos(timestamp string) int {
	pos := strings.LastIndex(timestamp, "+")
	if pos >= dateTimeLen {
		return pos
	}
	pos = strings.LastIndex(timestamp, "-")
	if pos >= dateTimeLen {
		return pos
	}
	return -1
}

func normTzOffset(timestamp string) string {
	pos := tzOffsetPos(timestamp)
	if pos > -1 {
		return timestamp[:pos] + strings.ReplaceAll(timestamp[pos:], ":", "")
	}
	return timestamp
}

func ParseDateTime(timestamp string) (XSDDateTime, error) {
	var t time.Time
	var err error
	if strings.HasSuffix(timestamp, "Z") {
		t, err = time.Parse(ISO8601FormatInUTC, timestamp)
	} else if hasTzOffset(timestamp) {
		t, err = time.Parse(ISO8601FormatInTZ, normTzOffset(timestamp))
	} else {
		t, err = time.Parse(ISO8601FormatInLocal, timestamp)
	}
	if err != nil {
		return XSDDateTime{}, err
	}
	return XSDDateTime{Time: t}, nil
}

func (xdt *XSDDateTime) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*xdt = XSDDateTime{}
		return nil
	}
	s := string(text)
	//ignore parse errors and set zero time
	*xdt, _ = ParseDateTime(s)
	return nil
}

func (xdt *XSDDateTime) MarshalText() ([]byte, error) {
	if xdt == nil {
		xdt = &XSDDateTime{}
	}
	s := xdt.Format(ISO8601FormatOut)
	return []byte(s), nil
}

// /Decimal implementation for currencies, etc
type XSDDecimal struct {
	Base int
	Exp  int
}

func (xd *XSDDecimal) UnmarshalText(text []byte) error {
	*xd = XSDDecimal{}
	s := string(text)
	_, xd.Base, xd.Exp = ExtractDecimal(s, -1)
	return nil
}

func (xd *XSDDecimal) MarshalText() ([]byte, error) {
	return []byte(FormatDecimal(xd.Base, xd.Exp)), nil
}
