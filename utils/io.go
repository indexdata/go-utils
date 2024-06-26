package utils

import (
	"bytes"
	_ "embed"
	"io"
	"log"
	"os"
	"strconv"
)

type WriteFlusher interface {
	io.Writer
	Flush() error
}

func Warn[T any](res T, err error) T {
	if err != nil {
		log.Println(err)
	}
	return res
}

func Fail[T any](res T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return res
}

// Extract decimal number e.g `123.45` from the string `s`
// where param `places` specifies maximum number of decimal places to consider
// with -1 considering all numbers to the first right-hand delimiter as fractions
// returns:
// a string representing the decimal
// the integer base
// the exponent value (as in integer * 10^-exp)
func ExtractDecimal(s string, places int) (string, int, int) {
	buff := make([]byte, len(s))
	isFraction := false
	k := 0
	pow10 := 1
	integer := 0
	exp := 0
	//go from back
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c >= '0' && c <= '9' {
			if k == 0 { //fraction likely starts
				isFraction = true
			}
			if places > -1 && k == places { //fraction limit
				isFraction = false
			}
			buff[k] = c
			v := int(c - '0')
			integer += v * pow10
			pow10 = pow10 * 10
			k++
		}
		if isFraction && (c == '.' || c == ',') {
			buff[k] = '.' //float sep
			exp = k
			k++
			isFraction = false //fraction ends
		}
	}
	reverse(buff, k)
	return string(buff[:k]), integer, exp
}

func reverse(buff []byte, k int) {
	for i, j := 0, k-1; i < j; i, j = i+1, j-1 {
		a := buff[i]
		b := buff[j]
		buff[i] = b
		buff[j] = a
	}
}

func FormatDecimal(integer int, exp int) string {
	if integer == 0 {
		return "0"
	}
	num := integer
	var buff bytes.Buffer
	for num > 0 {
		digit := num % 10
		buff.WriteByte(byte(digit + '0'))
		if buff.Len() == exp {
			buff.WriteByte('.')
		}
		num = num / 10
	}
	bytes := buff.Bytes()
	reverse(bytes, buff.Len())
	return string(bytes)
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return fallback
}

func GetEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.ParseBool(value); err == nil {
			return v
		}
	}
	return fallback
}
