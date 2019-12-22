// Abstraction module for gousb library.
package usb

import (
	"github.com/google/gousb"
	"github.com/korayeyinc/microconfig/util"
)

// define alias for gousb.ID
type ID = gousb.ID

// Represents configuration data from READ_ALL and CONFIGURE opcode commands.
type Data struct {
	OpCmd       uint8
	EEP_Addr    uint8
	Reserved1   uint8
	EEP_Val     uint8
	IO_Bmap     uint8
	Alt_Pins    uint8
	IO_Default  uint8
	Alt_Opts    uint8
	Baud_Rate_H uint8
	Baud_Rate_L uint8
	IO_Port_Val uint8
}

// Represents Config_Alt_Options bitmap.
type AltOpts struct {
	HW_Flow string
	Invert  string
	LEDX    string
	TxTGL   string
	RxTGL   string
}

// Represents Config_Alt_Pins bitmap.
type AltPins struct {
	SSPND  string
	USBCFG string
	RxLED  string
	TxLED  string
}

// MCP struct represents all the data structures
// to interact with the USB device.
type MCP struct {
	Context   *gousb.Context
	Device    *gousb.Device
	Conf      *gousb.Config
	Interface *gousb.Interface
	InEP      *gousb.InEndpoint
	OutEP     *gousb.OutEndpoint
	VendID    ID
	ProdID    ID
	*Data
}

// define command opcodes for MCP2200
const (
	BASE_CONFIGURE = 0x01
	SET_CLEAR_OUT  = 0x08
	CONFIGURE      = 0x10
	READ_EEPROM    = 0x20
	WRITE_EEPROM   = 0x40
	READ_ALL       = 0x80
)

// Initializes a new USB context object.
func NewContext() *gousb.Context {
	return gousb.NewContext()
}

// Opens any device with a given VID/PID using a convenience function.
func (micro *MCP) OpenDevice(VendID, ProdID ID) *gousb.Device {
	device, err := micro.Context.OpenDeviceWithVIDPID(VendID, ProdID)
	if err != nil {
		util.Fatalf("Could not open device: %v", err)
	}

	return device
}

// Enables/disables automatic kernel driver detachment.
func (micro *MCP) AutoDetach() {
	micro.Device.SetAutoDetach(true)
}

// Select the device configuration number 0.
func (micro *MCP) SelectConfig() *gousb.Config {
	conf, err := micro.Device.Config(1)
	if err != nil {
		util.Fatalf("%s.Config(0): %v", micro.Device, err)
	}

	return conf
}

// Claims the specified HID interface using a convenience function.
// The default interface is always #0 alt #0 in the currently active config.
func (micro *MCP) ClaimHIDInterface() *gousb.Interface {
	intf, err := micro.Conf.Interface(2, 0)
	if err != nil {
		util.Fatalf("%s.Interface(2, 0): %v", micro.Conf, err)
	}

	return intf
}

//Prepares an IN endpoint for transfer.
func (micro *MCP) InEndpoint() *gousb.InEndpoint {
	input, err := micro.Interface.InEndpoint(1)
	if err != nil {
		util.Fatalf("%s.InEndpoint(1): %v", micro.Interface, err)
	}

	return input
}

//Prepares an OUT endpoint for transfer.
func (micro *MCP) OutEndpoint() *gousb.OutEndpoint {
	output, err := micro.Interface.OutEndpoint(1)
	if err != nil {
		util.Fatalf("%s.OutEndpoint(1): %v", micro.Interface, err)
	}

	return output
}

// Reload performs a USB port reset to reinitialize a device.
func (micro *MCP) Reload() bool {
	err := micro.Device.Reset()
	util.Check(err)
	return true
}

// Reads device manufacturer information.
func (micro *MCP) ReadManufacturer() string {
	manufacturer, err := micro.Device.Manufacturer()
	if err != nil {
		util.Fatalf("Could not read device manufacturer: %v", err)
	}
	return manufacturer
}

// Reads device's product name.
func (micro *MCP) ReadProduct() string {
	product, err := micro.Device.Product()
	if err != nil {
		util.Fatalf("Could not read device's product name: %v", err)
	}
	return product
}

// Reads device's serial number.
func (micro *MCP) ReadSerial() string {
	serial, err := micro.Device.SerialNumber()
	if err != nil {
		util.Fatalf("Could not read device's serial number: %v", err)
	}
	return serial
}

// Sends READ_ALL command to MCP2200.
func (micro *MCP) ReadAllCmd() int {
	buf := make([]byte, 16)
	buf[0] = READ_ALL

	// write READ_ALL command opcode via OutEndpoint.
	val, err := micro.OutEP.Write(buf)

	if err != nil {
		util.Fatalf("%s.Write: got error %v:", micro.OutEP, err)
	}

	return val
}

// Sends the CONFIGURE command to MCP2200.
func (micro *MCP) ConfigCmd() int {
	data := micro.NewReqData()

	// write CONFIGURE command opcode via OutEndpoint.
	val, err := micro.OutEP.Write(data)

	if err != nil {
		util.Fatalf("%s.Write: got error %v:", micro.OutEP, err)
	}

	return val
}

// Parses READ_ALL command response.
func (micro *MCP) ParseResponse() *Data {
	buf := make([]byte, 16)

	// read READ_ALL command response via InEndpoint.
	_, err := micro.InEP.Read(buf)

	if err != nil {
		util.Fatalf("%s.Read: got error %v:", micro.InEP, err)
	}

	data := new(Data)
	data.OpCmd = buf[0]
	data.EEP_Addr = buf[1]
	data.Reserved1 = buf[2]
	data.EEP_Val = buf[3]
	data.IO_Bmap = buf[4]
	data.Alt_Pins = buf[5]
	data.IO_Default = buf[6]
	data.Alt_Opts = buf[7]
	data.Baud_Rate_H = buf[8]
	data.Baud_Rate_L = buf[9]
	data.IO_Port_Val = buf[10]

	return data
}

// Parses Alt_Opts bitmap.
func (micro *MCP) ParseAltOpts(bitmap uint8) *AltOpts {
	opts := new(AltOpts)
	opts.HW_Flow = util.GetBit(bitmap, 0)
	opts.Invert = util.GetBit(bitmap, 1)
	opts.LEDX = util.GetBit(bitmap, 5)
	opts.TxTGL = util.GetBit(bitmap, 6)
	opts.RxTGL = util.GetBit(bitmap, 7)

	return opts
}

// Parses Alt_Pins bitmap.
func (micro *MCP) ParseAltPins(bitmap uint8) *AltPins {
	gpio := new(AltPins)
	gpio.SSPND = util.GetBit(bitmap, 7)
	gpio.USBCFG = util.GetBit(bitmap, 6)
	gpio.RxLED = util.GetBit(bitmap, 3)
	gpio.TxLED = util.GetBit(bitmap, 2)

	return gpio
}

// Creates new request data for CONFIGURE command.
func (micro *MCP) NewReqData() []byte {
	buf := make([]byte, 16)
	buf[0] = CONFIGURE
	buf[4] = micro.Data.IO_Bmap
	buf[5] = micro.Data.Alt_Pins
	buf[6] = micro.Data.IO_Default
	buf[7] = micro.Data.Alt_Opts
	buf[8] = micro.Data.Baud_Rate_H
	buf[9] = micro.Data.Baud_Rate_L

	return buf
}

// Gets baud rate value.
func (micro *MCP) GetBaudRate(high, low uint8) int {
	baud_rate := util.ConcatBytes(high, low)
	return 12000000 / (int(baud_rate) + 1)
}

// Calculates high/low byte values.
func (micro *MCP) CalcHLBytes(baud string) (high, low uint8) {
	baud_rate := util.StrToInt(baud)
	divisor := 12000000/baud_rate - 1
	high, low = util.SplitBytes(uint16(divisor))
	return
}

// Returns active baud rate index
func (micro *MCP) GetBaudRateIndex(input string) int {
	baudrates := [12]string{
		"300", "1200", "2400", "4800", "9600", "19200",
		"38400", "57600", "115200", "230400", "460800", "921600",
	}

	var index int
	var baud string

	for index, baud = range baudrates {
		if baud == input {
			break
		}
	}

	return index
}
