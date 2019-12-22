// Provides utility functions.

package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	//"path/filepath"
)

// Generic error checking function.
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Alias to log.Fatalf function.
func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v)
}

// Alias to log.Fatal function.
func Fatal(v ...interface{}) {
	log.Fatal(v)
}

// Kills process and exits.
func Exit() {
	os.Exit(1)
}

// Chdir changes the current working directory to the named directory.
func Chdir(path string) {
	os.Chdir(path)
}

// Returns command line arguments.
func GetArg() int {
	inum := os.Args[1]
	return StrToInt(inum)
}

// Creates a new file.
func NewFile(filename string) *os.File {
	file, err := os.Create(filename)
	Check(err)
	return file
}

// Reads the file and returns byte data.
func ReadFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	Check(err)
	return data
}

// Exports given data to XML file.
func ExportXML(filename string) {

}

// Imports data from the given XML file.
func ImportXML(filename string) {

}

// Reports whether substr is within the string.
func StrContains(str, substr string) bool {
	return strings.Contains(str, substr)
}

// Converts integer to string.
func IntToStr(i int) string {
	str := strconv.Itoa(i)
	return str
}

// Converts string to float.
func StrToFloat(str string) float64 {
	val, _ := strconv.ParseFloat(str, 64)
	return val
}

// Converts string to integer.
func StrToInt(str string) int {
	val, _ := strconv.Atoi(str)
	return val
}

// Converts string to unsigned integer.
func StrToUint8(str string) uint8 {
	val, _ := strconv.ParseUint(str, 0, 8)
	return uint8(val)
}

// Converts string to unsigned integer.
func StrToUint16(str string) uint16 {
	val, _ := strconv.ParseUint(str, 0, 16)
	return uint16(val)
}

// Converts uint16 to string.
func UintToStr(vid, pid uint16) (vidstr, pidstr string) {
	vidstr = fmt.Sprintf("0x%04X", vid)
	pidstr = fmt.Sprintf("0x%04X", pid)
	return
}

// Converts uint8 to string.
func Uint8ToStr(x uint8) string {
	return fmt.Sprintf("0x%04X", x)
}

// Shows formatted representation of bits in a byte
func FmtBits(x uint8) string {
	bx := byte(x)
	return fmt.Sprintf("%08b", bx)
}

// Shows formatted representation of opts from bits.
func FmtOptStr(rx, tx, ledx, invert, hw_flow string) string {
	return fmt.Sprintf("%s%s%s000%s%s", rx, tx, ledx, invert, hw_flow)
}

// Shows formatted representation of pins from bits.
func FmtPinStr(suspend, usbcfg, rx, tx string) string {
	return fmt.Sprintf("%s%s00%s%s00", suspend, usbcfg, rx, tx)
}

// Finds the bit at given position in a byte.
func FindBit(x uint8, pos int) string {
	var bit string

	switch pos {
	case 0:
		bit = fmt.Sprintf("%01b", x&128)
	case 1:
		bit = fmt.Sprintf("%01b", x&64)
	case 2:
		bit = fmt.Sprintf("%01b", x&32)
	case 3:
		bit = fmt.Sprintf("%01b", x&16)
	case 4:
		bit = fmt.Sprintf("%01b", x&8)
	case 5:
		bit = fmt.Sprintf("%01b", x&4)
	case 6:
		bit = fmt.Sprintf("%01b", x&2)
	case 7:
		bit = fmt.Sprintf("%01b", x&1)
	}

	return bit
}

// Finds the bit at given position in a byte.
func GetBit(x uint8, pos int) string {
	return Uint8ToStr((x >> uint8(pos)) & 1)
}

// Concatenates given bytes.
func ConcatBytes(high, low uint8) uint16 {
	return uint16((high << 8) | low)
}

// Splits given uint16 value to uint8 values.
func SplitBytes(val uint16) (high, low uint8) {
	high = uint8(val & 0xff)
	low = uint8((val >> 8))
	return
}
