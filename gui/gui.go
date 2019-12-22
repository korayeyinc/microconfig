// GTK3 GUI toolbox module.

package gui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/korayeyinc/microconfig/util"
)

type Window = *gtk.Window
type Header = *gtk.HeaderBar
type FlowBox = *gtk.FlowBox
type Scroll = *gtk.ScrolledWindow
type Grid = *gtk.Grid
type Icon = *gtk.Image
type Input = *gtk.Entry
type Button = *gtk.Button
type RadioButton = *gtk.RadioButton
type Combo = *gtk.ComboBoxText
type Spin = *gtk.SpinButton
type Toggle = *gtk.Switch
type TextView = *gtk.TextView
type TextBuffer = *gtk.TextBuffer

// Quit event handler.
func Quit() {
	gtk.MainQuit()
}

// Creates a new window.
func NewWin() *gtk.Window {
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	util.Check(err)
	win.SetTitle("Microconfig")
	win.Connect("destroy", Quit)
	return win
}

// Creates a new image from icon.
func NewIcon(icon string) *gtk.Image {
	img, err := gtk.ImageNewFromIconName(icon, gtk.ICON_SIZE_LARGE_TOOLBAR)
	util.Check(err)
	return img
}

// Adds a new button widget.
func NewButton(icon string) *gtk.Button {
	button, err := gtk.ButtonNewFromIconName(icon, gtk.ICON_SIZE_LARGE_TOOLBAR)
	util.Check(err)
	return button
}

// Adds a new ComboBox widget.
func ComboBox() *gtk.ComboBoxText {
	combo, err := gtk.ComboBoxTextNew()
	util.Check(err)
	return combo
}

// Adds a new switch button widget.
func NewToggle() *gtk.Switch {
	toggle, err := gtk.SwitchNew()
	util.Check(err)
	return toggle
}

// Adds a new inputbox widget.
func InputBox() *gtk.Entry {
	inbox, err := gtk.EntryNew()
	util.Check(err)
	return inbox
}

// Adds a new label widget.
func Label(text string) *gtk.Label {
	label, err := gtk.LabelNew(text)
	util.Check(err)
	return label
}

// Adds a new radio button.
func RadioButtons(label1, label2 string) (opt1, opt2 *gtk.RadioButton) {
	var err error
	opt1, err = gtk.RadioButtonNewWithLabel(nil, label1)
	util.Check(err)
	group, err := opt1.GetGroup()
	util.Check(err)
	opt2, err = gtk.RadioButtonNewWithLabel(group, label2)
	util.Check(err)
	//radiobutton.SetMode(false);
	opt1.SetActive(true)
	return
}

// Adds a new spin button widget.
func SpinButton(min, max, step float64) *gtk.SpinButton {
	spinButton, err := gtk.SpinButtonNewWithRange(min, max, step)
	util.Check(err)
	return spinButton
}

// Adds a new textview widget.
func TxtView() *gtk.TextView {
	tview, err := gtk.TextViewNew()
	util.Check(err)
	//tp := gtk.TextWindowType(gtk.TEXT_WINDOW_WIDGET)
	//txtbox.SetBorderWindowSize(tp, 2)
	tview.SetEditable(false)
	tview.SetCursorVisible(false)
	return tview
}

func GetBuffer(tview *gtk.TextView) *gtk.TextBuffer {
	buffer, err := tview.GetBuffer()
	if err != nil {
		util.Fatal("Unable to get buffer:", err)
	}

	return buffer
}

// Adds a new flowbox widget.
func NewFlowBox(width, height int) FlowBox {
	fbox, err := gtk.FlowBoxNew()
	util.Check(err)
	fbox.SetVAlign(gtk.ALIGN_START)
	fbox.SetMaxChildrenPerLine(2)
	fbox.SetColumnSpacing(20)
	fbox.SetRowSpacing(20)
	fbox.SetSelectionMode(gtk.SELECTION_SINGLE)
	fbox.SetSizeRequest(width, height)
	return fbox
}

// Adds a new toolbar widget.
func HeaderBar() (header *gtk.HeaderBar, importBtn, exportBtn, reloadBtn, quitBtn *gtk.Button) {
	header, err := gtk.HeaderBarNew()
	util.Check(err)
	header.SetShowCloseButton(false)
	header.SetTitle("Microconfig v1.0")

	quitBtn = NewButton("system-shutdown-symbolic")
	quitBtn.SetLabel("Quit")
	header.PackEnd(quitBtn)
	reloadBtn = NewButton("view-refresh-symbolic")
	reloadBtn.SetLabel("Reload")
	header.PackEnd(reloadBtn)

	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	util.Check(err)
	hbox.SetSpacing(20)

	importBtn = NewButton("view-sort-ascending-symbolic")
	importBtn.SetLabel("Import")
	exportBtn = NewButton("view-sort-descending-symbolic")
	exportBtn.SetLabel("Export")

	hbox.Add(importBtn)
	hbox.Add(exportBtn)
	header.PackStart(hbox)
	return
}

// Adds a configuration panel widget for device settings.
func ConfigPanel() (grid Grid, vendID, prodID Input, baudrate Combo, ioconf, outdef Input, leds, pins, usbcfg, suspend, upol Toggle, opt1, opt2 RadioButton, duration Spin, config, reset Button) {
	var err error
	grid, err = gtk.GridNew()
	util.Check(err)
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	grid.SetMarginStart(20)
	grid.SetMarginTop(20)
	grid.SetMarginBottom(20)
	grid.SetColumnSpacing(20)
	grid.SetRowSpacing(20)

	vendlab := Label("Vendor ID:")
	vendID = InputBox()
	vendID.SetMaxLength(8)
	vendID.SetWidthChars(8)

	prodlab := Label("Product ID:")
	prodID = InputBox()
	prodID.SetMaxLength(8)
	prodID.SetWidthChars(8)

	baudrates := [12]string{
		"300", "1200", "2400", "4800", "9600", "19200",
		"38400", "57600", "115200", "230400", "460800", "921600",
	}

	baudlab := Label("Baud Rate:")
	baudrate = ComboBox()

	for _, baud := range baudrates {
		baudrate.AppendText(baud)
	}

	//baudrate.SetActive(1)

	iolab := Label("IO Config:")
	ioconf = InputBox()
	ioconf.SetMaxLength(8)
	ioconf.SetWidthChars(8)
	ioconf.SetText("00000000")

	outlab := Label("Output Default:")
	outdef = InputBox()
	outdef.SetMaxLength(8)
	outdef.SetWidthChars(8)
	outdef.SetText("00000000")

	ledslab := Label("Enable Tx/Rx LEDs:")
	leds = NewToggle()
	leds.SetActive(false)

	pinslab := Label("Enable CTS/RTS Pins:")
	pins = NewToggle()
	pins.SetActive(false)

	usbcfglab := Label("Enable USBCFG Pin:")
	usbcfg = NewToggle()
	usbcfg.SetActive(false)

	susplab := Label("Enable Suspend Pin:")
	suspend = NewToggle()
	suspend.SetActive(false)

	upollab := Label("Enable UART Polarity:")
	upol = NewToggle()
	upol.SetActive(false)

	config = NewButton("preferences-system-symbolic")
	config.SetLabel("Configure")

	reset = NewButton("document-revert-symbolic")
	reset.SetLabel("Reset")

	funclab := Label("LED Function:")
	opt1, opt2 = RadioButtons("Blink LEDs", "Toggle LEDs")

	blinklab := Label("Blink Duration:")
	duration = SpinButton(100, 200, 100)
	duration.SetValue(100.0)

	//save := NewButton("gtk-apply")
	//save.SetLabel("Kaydet")

	grid.Attach(vendlab, 0, 0, 1, 1)
	grid.Attach(vendID, 1, 0, 1, 1)
	grid.Attach(prodlab, 0, 1, 1, 1)
	grid.Attach(prodID, 1, 1, 1, 1)
	grid.Attach(baudlab, 0, 2, 1, 1)
	grid.Attach(baudrate, 1, 2, 1, 1)
	grid.Attach(iolab, 0, 3, 1, 1)
	grid.Attach(ioconf, 1, 3, 1, 1)
	grid.Attach(outlab, 0, 4, 1, 1)
	grid.Attach(outdef, 1, 4, 1, 1)
	grid.Attach(ledslab, 0, 5, 1, 1)
	grid.Attach(leds, 1, 5, 1, 1)
	grid.Attach(pinslab, 0, 6, 1, 1)
	grid.Attach(pins, 1, 6, 1, 1)
	grid.Attach(usbcfglab, 0, 7, 1, 1)
	grid.Attach(usbcfg, 1, 7, 1, 1)
	grid.Attach(susplab, 0, 8, 1, 1)
	grid.Attach(suspend, 1, 8, 1, 1)
	grid.Attach(upollab, 0, 9, 1, 1)
	grid.Attach(upol, 1, 9, 1, 1)
	grid.Attach(funclab, 0, 10, 1, 1)
	grid.Attach(opt1, 1, 10, 1, 1)
	grid.Attach(opt2, 1, 11, 1, 1)
	grid.Attach(blinklab, 0, 12, 1, 1)
	grid.Attach(duration, 1, 12, 1, 1)
	grid.Attach(config, 0, 14, 1, 1)
	grid.Attach(reset, 1, 14, 1, 1)

	return
}

// Adds a new panel widget.
func InfoPanel(manufacturer, product, serial string) (grid Grid, statico *gtk.Image, manufact, prod, serinum Input, infoico Button, console TextView) {
	var err error
	grid, err = gtk.GridNew()
	util.Check(err)
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	grid.SetMarginStart(20)
	grid.SetMarginTop(20)
	grid.SetMarginBottom(20)
	grid.SetColumnSpacing(20)
	grid.SetRowSpacing(20)

	statlab := Label("Connection Status:")
	statico = NewIcon("object-select-symbolic")

	manulab := Label("Manufacturer:")
	manufact = InputBox()
	manufact.SetMaxLength(50)
	manufact.SetWidthChars(50)
	manufact.SetText(manufacturer)

	prodlab := Label("Product:")
	prod = InputBox()
	prod.SetMaxLength(50)
	prod.SetWidthChars(50)
	prod.SetText(product)

	serilab := Label("Serial Number:")
	serinum = InputBox()
	serinum.SetMaxLength(50)
	serinum.SetWidthChars(50)
	serinum.SetText(serial)

	infoico = NewButton("dialog-information-symbolic")
	infoico.SetLabel("Info Console")

	scroll, err := gtk.ScrolledWindowNew(nil, nil)
	util.Check(err)
	scroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	console = TxtView()
	scroll.Add(console)

	//connlab = Label("")
	//connlab.SetUseMarkup(true)
	//connlab.SetMarkup("<b>Connected</b>")

	grid.Attach(statlab, 0, 0, 1, 1)
	grid.Attach(statico, 1, 0, 1, 1)
	grid.Attach(manulab, 0, 1, 1, 1)
	grid.Attach(manufact, 1, 1, 1, 1)
	grid.Attach(prodlab, 0, 2, 1, 1)
	grid.Attach(prod, 1, 2, 1, 1)
	grid.Attach(serilab, 0, 3, 1, 1)
	grid.Attach(serinum, 1, 3, 1, 1)
	grid.Attach(infoico, 0, 4, 1, 1)
	grid.Attach(scroll, 0, 5, 2, 15)

	return
}

// Adds root box containing other GTK widgets.
func RootBox(confPanel, infoPanel Grid) *gtk.Box {
	rootBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	util.Check(err)
	rootBox.SetSpacing(0)
	vsep, err := gtk.SeparatorNew(gtk.ORIENTATION_VERTICAL)
	util.Check(err)
	rootBox.PackStart(confPanel, true, true, 0)
	rootBox.PackStart(vsep, false, false, 0)
	rootBox.PackEnd(infoPanel, true, true, 0)
	return rootBox
}

// Renders the window with the container widgets.
func Render(win *gtk.Window, header Header, rootBox *gtk.Box) {
	win.Add(rootBox)
	win.SetTitlebar(header)
	win.SetDefaultSize(800, 600)
	//win.Maximize()
	win.ShowAll()
	gtk.Main()
}

func ChooseXML(win *gtk.Window) string {
	title := "Import From XML"
	xml := gtk.OpenFileChooserNative(title, win)
	return *xml
}
