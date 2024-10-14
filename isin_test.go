// SPDX-FileCopyrightText: 2024 CmdScale GmbH
//
// SPDX-License-Identifier: BSD-3-Clause

package corfi

import (
	"errors"
	"testing"

	"code.pfad.fr/check"
)

func TestCleanISIN(t *testing.T) {
	f := func(s string) {
		t.Helper()

		isin, err := NewISIN(s)
		check.Equal(t, nil, err).Log(s)
		check.Equal(t, s, isin.String())
	}

	f("US0378331005")
	f("DE0006231004")
	f("DE000BAY0017")
	f("XF0000C14922")
	f("NL0000729408")
	f("CH0031240127")
	f("US5949181045")
	f("US38259P5089")
	f("JP3946600008")
	f("DE000DZ21632")
	f("DE000DB7HWY7")
	f("DE000CM7VX13")
	f("CH0031240127")
	f("CA9861913023")
}

func TestValidISIN(t *testing.T) {
	f := func(s string, country CountryCode, nsin string, checkDigit int) {
		t.Helper()

		isin, err := NewISIN(s)
		check.Equal(t, nil, err)
		check.Equal(t, country, isin.CountryCode).Log("failed CountryCode.")
		check.Equal(t, nsin, isin.NSIN).Log("failed NSIN.")
		check.Equal(t, checkDigit, isin.CheckDigit).Log("failed CheckDigit.")

		isin2, err := NewISIN(isin.String())
		check.Equal(t, nil, err)
		check.Equal(t, isin2, isin)
	}
	f("US0378331005", "US", "037833100", 5)
	f("us 0378331005", "US", "037833100", 5)

	f("AU0000XVGZA3", "AU", "0000XVGZA", 3)
	f("A u0000xVGZa 3", "AU", "0000XVGZA", 3)
	f("AU0000VXGZA3", "AU", "0000VXGZA", 3)

	f("GB0002634946", "GB", "000263494", 6)
}

func TestInvalidISIN(t *testing.T) {
	_, err := NewISIN("123")
	var lerr LenError
	check.Equal(t, true, errors.As(err, &lerr)).Logf("got error %T: %v", err, err)
	check.Equal(t, `expected 12 alphanumeric chars, got 3`, lerr.Error())

	_, err = NewISIN("12345678901A")
	var cerr CheckDigitError
	check.Equal(t, true, errors.As(err, &cerr)).Logf("got error %T: %v", err, err)
	check.Equal(t, `expected digit as last chart, got "A"`, cerr.Error())

	_, err = NewISIN("US0378331006")
	check.Equal(t, true, errors.As(err, &cerr)).Logf("got error %T: %v", err, err)
	check.Equal(t, `wrong check digit`, cerr.Error())
	t.Run("length", func(t *testing.T) {
		f := func(s string, length int) {
			t.Helper()

			_, err := NewISIN(s)
			var lerr LenError
			check.Equal(t, true, errors.As(err, &lerr)).Logf("got error %T: %v", err, err)
			check.Equal(t, length, int(lerr))
		}
		f("", 0)
		f("#", 0) // # gets removed
		f("1", 1)
		f("123", 3)
		f("1234567890", 10)
		f("12345678901", 11)
		f("12345678901#", 11) // # gets removed
		f("1234567890123", 13)
	})

	t.Run("check_digit", func(t *testing.T) {
		f := func(s string, given, computed int) {
			t.Helper()

			_, err := NewISIN(s)
			var cerr CheckDigitError
			check.Equal(t, true, errors.As(err, &cerr)).Logf("got error %T: %v", err, err)
			check.Equal(t, computed, cerr.Computed).Log("Computed")
			check.Equal(t, given, cerr.Given).Log("Given")
		}
		f("US0378331006", 6, 5)
		f("US0378331004", 4, 5)
		f("us 0378331000", 0, 5)

		f("AU0000XVGZA2", 2, 3)
		f("A u0000xVGZa 9", 9, 3)

		f("GB0002634947", 7, 6)
	})
}