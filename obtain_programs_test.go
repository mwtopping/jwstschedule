package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPIparsing(t *testing.T) {
	PIstring := "PI: Andrew Swan"
	PI := get_PI(PIstring)
	assert.Equal(t, "Andrew Swan", PI)

	PIstring = "PI: Charles Beichman<br />                       Co-PI: Dimitri Mawet"
	PI = get_PI(PIstring)
	assert.Equal(t, "Charles Beichman", PI)

	PIstring = "Bob George"
	PI = get_PI(PIstring)
	assert.Equal(t, "", PI)

}

func TestExpTimeParse(t *testing.T) {
	testString := "0/591"
	priT, parT := parse_exptime(testString)
	assert.Equal(t, float32(0.0), priT)
	assert.Equal(t, float32(591.0), parT)

	testString = "25"
	priT, parT = parse_exptime(testString)
	assert.Equal(t, float32(25.0), priT)
	assert.Equal(t, float32(0.0), parT)

	testString = "10/25"
	priT, parT = parse_exptime(testString)
	assert.Equal(t, float32(10.0), priT)
	assert.Equal(t, float32(25.0), parT)

}

func TestInstrumentParse(t *testing.T) {
	testString := "NIRSpec/IFU<br />                        MIRI/MRS"
	result := parse_instrument_mode(testString)
	expected := []string{"NIRSpec/IFU", "MIRI/MRS"}
	assert.Equal(t, expected, result)

	testString = "NIRSpec/MOS"
	result2 := parse_instrument_mode(testString)
	expected2 := []string{"NIRSpec/MOS"}
	assert.Equal(t, expected2, result2)

}

func TestModeParse(t *testing.T) {
	testString := "GO, Calibration"
	result := parse_program_type(testString)
	assert.Equal(t, []string{"GO", "Calibration"}, result)

	testString = "GTO"
	result2 := parse_program_type(testString)
	assert.Equal(t, []string{"GTO"}, result2)
}
