// SPDX-FileCopyrightText: 2024 CmdScale GmbH
//
// SPDX-License-Identifier: BSD-3-Clause

// package corfin provides methods to check the validity of an International Securities Identification Number (ISIN)
package corfin

import (
	"regexp"
	"strconv"
	"strings"
)

// ISINValidator makes additional checks on a given ISIN.
type ISINValidator func(ISIN) error

// NewISIN sanitizes, parses and checks an ISIN (length, check digit).
func NewISIN(s string, additionalRules ...ISINValidator) (ISIN, error) {
	s = strings.ToUpper(nonAlphanumeric.ReplaceAllString(s, ""))

	var isin ISIN
	if len(s) != isin_len {
		return isin, LenError(len(s))
	}
	isin.CountryCode = CountryCode(s[0:2])
	isin.NSIN = s[2:11]
	isin.CheckDigit = int(s[11] - '0')
	if isin.CheckDigit < 0 || isin.CheckDigit > 9 {
		return isin, CheckDigitError{Given: isin.CheckDigit, Computed: -1}
	}
	c := isin.checkDigit()
	if c != isin.CheckDigit {
		return isin, CheckDigitError{Given: isin.CheckDigit, Computed: c}
	}

	for _, v := range additionalRules {
		if err := v(isin); err != nil {
			return isin, err
		}
	}

	return isin, nil
}

const isin_len = 12

var nonAlphanumeric = regexp.MustCompile(`[^0-9A-Za-z]`)

// ISIN represents an International Securities Identification Number
type ISIN struct {
	CountryCode CountryCode
	// National Securities Identifying Number
	NSIN       string
	CheckDigit int
}

// CountryCode is the ISO 3166-1 alpha-2 code of the country (uppercase)
type CountryCode string

// String returns the ISIN
func (isin ISIN) String() string {
	return string(isin.CountryCode) + isin.NSIN + strconv.Itoa(isin.CheckDigit)
}

func (isin ISIN) checkDigit() int {
	nsinDigit, multiply := luhnDigit(isin.NSIN, true)
	countryDigit, _ := luhnDigit(string(isin.CountryCode), multiply)
	sum := nsinDigit + countryDigit
	return (10 - (sum % 10)) % 10
}

func luhnDigit(s string, m bool) (sum int, multiply bool) {
	multiply = m
	ingest := func(i int) {
		if multiply {
			if i > 4 {
				i = 2*i - 9
			} else {
				i *= 2
			}
		}
		sum += i
		multiply = !multiply
	}
	for i := len(s) - 1; i >= 0; i-- {
		c := int(s[i] - '0')

		if c > 9 {
			// alpha: adjust the number and ingest both digits
			c = c + '0' - 55
			ingest(c % 10)
			c = c / 10
		}
		ingest(c)
	}
	return sum, multiply
}

// LenError indicates that the provided ISIN is of the wrong length
type LenError int

func (e LenError) Error() string {
	return "expected " + strconv.Itoa(isin_len) + " alphanumeric chars, got " + strconv.Itoa(int(e))
}

// CheckDigitError indicates that the checksum did not match the expected check digit.
type CheckDigitError struct {
	Given    int
	Computed int
}

func (e CheckDigitError) Error() string {
	if e.Given < 0 || e.Given > 9 {
		return "expected digit as last chart, got " + strconv.Quote(string(byte(e.Given+'0')))
	}
	// don't display the expected/actual digits, because the error is most likely somewhere else in the NSIN
	return "wrong check digit"
}
