package main

import (
	"github.com/korayeyinc/microconfig/gui"
	"github.com/korayeyinc/microconfig/usb"
	"github.com/korayeyinc/microconfig/util"
)

var (
	win     gui.Window
	buffer  gui.TextBuffer
	micro   *usb.MCP
	gpio    *usb.AltPins
	opts    *usb.AltOpts
	conf    *Conf
	panel   *Panel
	console *Console
	button  *Button
	combo   *Combo
	icon    *Icon
	input   *Input
	radio   *Radio
	spin    *Spin
	toggle  *Toggle
)

// Represents device configuration for logging.
type Conf struct {
	VendID     string
	ProdID     string
	BaudRate   string
	IOConfig   string
	OutDefault string
	TxRxLeds   string
	CRTS       string
	USBCFG     string
	Suspend    string
	UARTPol    string
	LedFunc    string
	Blink      string
	Manufact   string
	Product    string
	Serial     string
}

type Button struct {
	Config  gui.Button
	Reset   gui.Button
	Export  gui.Button
	Import  gui.Button
	Reload  gui.Button
	Quit    gui.Button
	Console gui.Button
	Info    gui.TextBuffer
}

type Icon struct {
	Stat gui.Icon
}

type Combo struct {
	BaudRate gui.Combo
}

type Input struct {
	Manufacturer gui.Input
	Product      gui.Input
	Serial       gui.Input
	VendID       gui.Input
	ProdID       gui.Input
	IOConf       gui.Input
	OutDef       gui.Input
}

type Toggle struct {
	Leds    gui.Toggle
	Pins    gui.Toggle
	Usbcfg  gui.Toggle
	Suspend gui.Toggle
	UPol    gui.Toggle
}

type Radio struct {
	BlinkLeds  gui.RadioButton
	ToggleLeds gui.RadioButton
}

type Spin struct {
	Duration gui.Spin
}

type Console struct {
	View gui.TextView
}

type Panel struct {
	Header gui.Header
	Conf   gui.Grid
	Info   gui.Grid
}

// Stores device configuration to NVRAM.
func configDevice() {
	conf.BaudRate = combo.BaudRate.GetActiveText()
	micro.Data.Baud_Rate_H, micro.Data.Baud_Rate_L = micro.CalcHLBytes(conf.BaudRate)

	conf.IOConfig, _ = input.IOConf.GetText()
	micro.Data.IO_Bmap = util.StrToUint8(conf.IOConfig)

	conf.OutDefault, _ = input.OutDef.GetText()
	micro.Data.IO_Default = util.StrToUint8(conf.OutDefault)

	if toggle.Leds.GetActive() {
		conf.TxRxLeds = "1"
		gpio.TxLED = "1"
		gpio.RxLED = "1"
	} else {
		conf.TxRxLeds = "0"
		gpio.TxLED = "0"
		gpio.RxLED = "0"
	}

	if toggle.Pins.GetActive() {
		conf.CRTS = "1"
		opts.HW_Flow = "1"
	} else {
		conf.CRTS = "0"
		opts.HW_Flow = "0"
	}

	if toggle.Usbcfg.GetActive() {
		conf.USBCFG = "1"
		gpio.USBCFG = "1"
	} else {
		conf.USBCFG = "0"
		gpio.USBCFG = "0"
	}

	if toggle.Suspend.GetActive() {
		conf.Suspend = "1"
		gpio.SSPND = "1"
	} else {
		conf.Suspend = "0"
		gpio.SSPND = "0"
	}

	if toggle.UPol.GetActive() {
		conf.UARTPol = "1"
		opts.Invert = "1"
	} else {
		conf.UARTPol = "0"
		opts.Invert = "0"
	}

	if radio.BlinkLeds.GetActive() {
		opts.RxTGL = "0"
		opts.TxTGL = "0"
		if spin.Duration.GetValue() == 100.0 {
			opts.LEDX = "0"
		} else {
			opts.LEDX = "1"
		}
	} else if radio.ToggleLeds.GetActive() {
		opts.RxTGL = "1"
		opts.TxTGL = "1"
	}

	gpioStr := util.FmtPinStr(gpio.SSPND, gpio.USBCFG, gpio.RxLED, gpio.TxLED)
	micro.Data.Alt_Pins = util.StrToUint8(gpioStr)

	optsStr := util.FmtOptStr(opts.RxTGL, opts.TxTGL, opts.LEDX, opts.Invert, opts.HW_Flow)
	micro.Data.Alt_Opts = util.StrToUint8(optsStr)

	micro.ConfigCmd()
}

// Exports device configuration to XML.
func exportXML() {
	file := gui.ChooseXML(win)
	util.ExportXML(file, conf)
}

// Imports device configuration from XML.
func importXML() {
	gui.ChooseXML(win)
}

// Disconnects USB device and quits the application.
func quitApp() {
	gui.Quit()
}

// Logs info console output.
func logInfo() {
	buffer.SetText("")
}

// Sets LED configuration options.
func configLED() (ledfunc, duration string) {
	if opts.RxTGL == "0" || opts.TxTGL == "0" {
		ledfunc = "blink"
		if opts.LEDX == "0" {
			duration = "100"
		} else if opts.LEDX == "1" {
			duration = "200"
		}
	} else if opts.RxTGL == "1" || opts.TxTGL == "1" {
		ledfunc = "toggle"
	}
	return
}

// Sets info console.
func setConsole() {
	info := " [INFO]	USB Device (Microchip MCP2200) Connected!\n"
	info += " [Done]	Claiming HID Interface (2, 0)...\n"
	info += " [DONE]	Reading String Descriptors...\n"
	info += " [DONE]	Reading Device Serial Number...\n"
	info += " [DONE]	Reading Device Configurations...\n"
	buffer.SetText(info)
}

// Reconnects USB device and reloads application
func reloadApp() {
	micro.Reload()
	logInfo()
	setConsole()
}

func main() {
	// create new usb.MCP object
	micro = new(usb.MCP)

	// initialize a new usb context
	micro.Context = usb.NewContext()
	defer micro.Context.Close()

	// set Vendor/Product IDs for MCP2200 device
	micro.VendID, micro.ProdID = 0x04D8, 0x00DF

	// open USB device with Vendor/Product ID
	micro.Device = micro.OpenDevice(micro.VendID, micro.ProdID)

	// enable Linux kernel driver auto detachment
	micro.AutoDetach()

	// ***TO BE REMOVED***
	//var vendID, prodID usb.ID  = 0x12D1, 0x1039

	conf = new(Conf)
	vid, pid := uint16(micro.VendID), uint16(micro.ProdID)
	conf.VendID, conf.ProdID = util.UintToStr(vid, pid)

	// check if the USB device is connected
	if micro.Device == nil {
		util.Fatalf("Matching USB device not found!\n[VendorID ProductID] %s %s", conf.VendID, conf.ProdID)
	}
	defer micro.Device.Close()

	// initialize device configuration
	micro.Conf = micro.SelectConfig()
	defer micro.Conf.Close()

	// claim HID interface
	micro.Interface = micro.ClaimHIDInterface()
	defer micro.Interface.Close()

	// set In/Out Endpoints
	micro.InEP = micro.InEndpoint()
	micro.OutEP = micro.OutEndpoint()

	// send READ_ALL command request to MCP2200
	micro.ReadAllCmd()
	// parse READ_ALL command response from MCP2200
	micro.Data = micro.ParseResponse()

	// read string descriptors
	conf.Manufact = micro.ReadManufacturer()
	conf.Product = micro.ReadProduct()
	conf.Serial = micro.ReadSerial()

	// init widget objects
	win = gui.NewWin()
	panel = new(Panel)
	button = new(Button)
	icon = new(Icon)
	input = new(Input)
	combo = new(Combo)
	radio = new(Radio)
	spin = new(Spin)
	toggle = new(Toggle)
	console = new(Console)

	// set headerbar widgets
	panel.Header, button.Import, button.Export, button.Reload, button.Quit = gui.HeaderBar()

	// set config panel widgets
	panel.Conf, input.VendID, input.ProdID, combo.BaudRate, input.IOConf, input.OutDef,
		toggle.Leds, toggle.Pins, toggle.Usbcfg, toggle.Suspend, toggle.UPol, radio.BlinkLeds,
		radio.ToggleLeds, spin.Duration, button.Config, button.Reset = gui.ConfigPanel()

	// set IDs
	input.VendID.SetText(conf.VendID)
	input.ProdID.SetText(conf.ProdID)

	// parse Alt_Opts and Alt_Pins data
	opts = micro.ParseAltOpts(micro.Data.Alt_Opts)
	gpio = micro.ParseAltPins(micro.Data.Alt_Pins)

	// set Conf
	baud := micro.GetBaudRate(micro.Data.Baud_Rate_H, micro.Data.Baud_Rate_L)
	conf.BaudRate = util.IntToStr(baud)
	conf.IOConfig = util.FmtBits(micro.Data.IO_Bmap)
	conf.OutDefault = util.FmtBits(micro.Data.IO_Default)
	conf.TxRxLeds = gpio.TxLED
	conf.CRTS = opts.HW_Flow
	conf.USBCFG = gpio.USBCFG
	conf.Suspend = gpio.SSPND
	conf.UARTPol = opts.Invert

	// set LED configuration options
	conf.LedFunc, conf.Blink = configLED()

	// set active vals
	index := micro.GetBaudRateIndex(conf.BaudRate)
	combo.BaudRate.SetActive(index)
	input.IOConf.SetText(conf.IOConfig)
	input.OutDef.SetText(conf.OutDefault)

	if conf.TxRxLeds == "1" {
		toggle.Leds.SetActive(true)
	} else {
		toggle.Leds.SetActive(false)
	}

	if conf.USBCFG == "1" {
		toggle.Usbcfg.SetActive(true)
	} else {
		toggle.Usbcfg.SetActive(false)
	}

	if conf.Suspend == "1" {
		toggle.Suspend.SetActive(true)
	} else {
		toggle.Suspend.SetActive(false)
	}

	// set active radio buttons
	if conf.LedFunc == "blink" {
		radio.BlinkLeds.SetActive(true)
	} else if conf.LedFunc == "toggle" {
		radio.ToggleLeds.SetActive(true)
	}

	// set info panel widgets
	panel.Info, icon.Stat, input.Manufacturer, input.Product, input.Serial, button.Console, console.View = gui.InfoPanel(conf.Manufact, conf.Product, conf.Serial)
	buffer = gui.GetBuffer(console.View)
	setConsole()

	// wrap panels inside a root box
	rootBox := gui.RootBox(panel.Conf, panel.Info)

	// handle button click events
	button.Config.Connect("clicked", configDevice)
	button.Reset.Connect("clicked", micro.Reload)
	button.Import.Connect("clicked", importXML)
	button.Export.Connect("clicked", exportXML)
	button.Reload.Connect("clicked", reloadApp)
	button.Quit.Connect("clicked", quitApp)
	button.Console.Connect("clicked", logInfo)

	button.Reset.SetSensitive(false)
	button.Import.SetSensitive(false)
	button.Export.SetSensitive(false)

	// render window with the widgets
	gui.Render(win, panel.Header, rootBox)
}
