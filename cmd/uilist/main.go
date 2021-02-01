package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	linijka "github.com/kubaraczkowski/linijka/pkg"
	"github.com/lxn/walk"

	. "github.com/lxn/walk/declarative"
)

const version = "0.1.0"

var ipaddress = "192.168.0.50"

var ipport = 4001

var conn net.Conn
var connected bool
var printonly bool

func main() {
	var openAction, showAboutBoxAction *walk.Action
	mw := &MyMainWindow{model: NewLiniaModel()}
	toolbar := &walk.ToolBar{}

	var ipaddress_string string

	flag.StringVar(&ipaddress_string, "ip", "192.168.0.50", "IP address of the device")
	flag.IntVar(&ipport, "port", 4001, "IP port of the device")
	flag.BoolVar(&printonly, "printonly", false, "Only print the messages to be sent, don't connect")
	version_flag := flag.Bool("v", false, "Display program version and exit")

	flag.Parse()

	if *version_flag {
		fmt.Println(version)
		os.Exit(0)
	}

	ipaddress := net.ParseIP(ipaddress_string)
	if ipaddress == nil {
		log.Fatalf("Could not parse IP address: %s", ipaddress_string)
	}

	var filename string
	if len(flag.Args()) >= 1 {
		filename = flag.Args()[0]
	}

	if filename != "" {
		mw.loadFile(filename)
	}

	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Linijka",
		MinSize:  Size{240, 320},
		Size:     Size{300, 400},
		Layout:   VBox{MarginsZero: false},
		OnDropFiles: func(files []string) {
			if len(files) >= 1 {
				mw.loadFile(files[0])
			}
		},
		// Children: []Widget{
		// 	VSplitter{
		Children: []Widget{
			ListBox{
				AssignTo:              &mw.lb,
				Model:                 mw.model,
				OnCurrentIndexChanged: mw.lb_CurrentIndexChanged,
				OnItemActivated:       mw.lb_ItemActivated,
				AlwaysConsumeSpace:    true,
			},
			// 	},
			// },
		},
		MenuItems: []MenuItem{
			Menu{
				Text: "&File",
				Items: []MenuItem{
					Action{
						AssignTo:    &openAction,
						Text:        "&Open",
						Shortcut:    Shortcut{walk.ModControl, walk.KeyO},
						OnTriggered: mw.openAction_Triggered,
					},
					Separator{},
					Action{
						Text:        "E&xit",
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						AssignTo:    &showAboutBoxAction,
						Text:        "About",
						OnTriggered: mw.showAboutBoxAction_Triggered,
					},
				},
			},
		},
		ToolBar: ToolBar{AssignTo: &toolbar,
			ButtonStyle: ToolBarButtonImageBeforeText,
			Items: []MenuItem{
				Action{Text: "Open", Visible: true, Enabled: true, OnTriggered: mw.openAction_Triggered},
				Action{AssignTo: &mw.connectAction, Text: "Connect", Visible: true, Enabled: true, OnTriggered: mw.connectAction_Triggered},
				Action{AssignTo: &mw.disconnectAction, Text: "Disconnect", Enabled: false, OnTriggered: mw.disconnectAction_Triggered},
				Separator{},
				Action{AssignTo: &mw.permanentAction, Text: "Permanent", Enabled: true, Checkable: true, Checked: false},
			},
		},
		StatusBarItems: []StatusBarItem{
			StatusBarItem{
				AssignTo: &mw.sbi,
				Text:     "Disconnected",
				Width:    80,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

type MyMainWindow struct {
	*walk.MainWindow
	model            *LiniaModel
	prevFilePath     string
	lb               *walk.ListBox
	connectAction    *walk.Action
	disconnectAction *walk.Action
	sbi              *walk.StatusBarItem
	permanentAction  *walk.Action
}

func (mw *MyMainWindow) lb_CurrentIndexChanged() {
	i := mw.lb.CurrentIndex()
	if i >= 0 && i <= len(mw.model.items) {
		item := mw.model.items[i]

		// fmt.Println("CurrentIndex: ", i)
		// fmt.Println("CurrentEnvVarName: ", item.Bar)
		mw.sendLine(item.Bar)
	}
}

func (mw *MyMainWindow) lb_ItemActivated() {
	value := mw.model.items[mw.lb.CurrentIndex()].Bar

	walk.MsgBox(mw, "Value", value, walk.MsgBoxIconInformation)
}

type Linia struct {
	Bar string
}

type LiniaModel struct {
	walk.ListModelBase
	items []*Linia
}

func NewLiniaModel() *LiniaModel {
	m := new(LiniaModel)
	m.ExampleRows()
	return m
}

func (m *LiniaModel) ItemCount() int {
	return len(m.items)
}

func (m *LiniaModel) Value(index int) interface{} {
	return m.items[index].Bar
}

func (m *LiniaModel) ExampleRows() {
	m.items = make([]*Linia, 10)

	for i := range m.items {
		m.items[i] = &Linia{
			Bar: fmt.Sprintf("Linia %d", i),
		}
	}

}

func (mw *MyMainWindow) openAction_Triggered() {
	// walk.MsgBox(mw, "Open", "Pretend to open a file...", walk.MsgBoxIconInformation)
	mw.openFile()
}

func (mw *MyMainWindow) showAboutBoxAction_Triggered() {
	walk.MsgBox(mw, "About", fmt.Sprintf("Linijka, version %v", version), walk.MsgBoxIconInformation)
}

func (mw *MyMainWindow) connectAction_Triggered() {
	mw.connect()
	mw.updateConnectedStatus()
}

func (mw *MyMainWindow) disconnectAction_Triggered() {
	mw.disconnect()
	mw.updateConnectedStatus()
}

func (mw *MyMainWindow) updateConnectedStatus() {
	if connected {
		mw.sbi.SetText(fmt.Sprintf("Connected to: %s:%d", ipaddress, ipport))
		mw.connectAction.SetEnabled(false)
		mw.disconnectAction.SetEnabled(true)
	} else {
		mw.sbi.SetText("Disconnected")
		mw.connectAction.SetEnabled(true)
		mw.disconnectAction.SetEnabled(false)
	}
}

func (mw *MyMainWindow) openFile() error {
	var err error
	dlg := new(walk.FileDialog)

	dlg.FilePath = mw.prevFilePath
	dlg.Filter = "Text files (*.txt)"
	dlg.Title = "Select a Linijka file"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		return err
	} else if !ok {
		return nil
	}
	mw.prevFilePath = dlg.FilePath
	log.Println(dlg.FilePath)

	mw.loadFile(dlg.FilePath)
	return err
}

func (mw *MyMainWindow) loadFile(path string) error {
	var err error

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var line string

	var new []*Linia
	mw.model.items = new

	for {
		line, err = reader.ReadString('\n')
		log.Println(line)
		if err == io.EOF {
			break
		}
		lin := Linia{Bar: strings.TrimSpace(line)}
		if len(line) == 0 {
			break
		}
		mw.model.items = append(mw.model.items, &lin)
		if err != nil {
			break
		}
	}
	mw.model.PublishItemsReset()

	return nil
}

func (mw *MyMainWindow) connect() error {
	var err error
	if printonly {
		log.Println("Printonly enabled, not connecting")
		return nil
	}
	if connected == false {

		log.Printf("Connecting to %s:%d", ipaddress, ipport)
		var d net.Dialer
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		conn, err = d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", ipaddress, ipport))
		if err != nil {
			log.Printf("Failed to connect: %v", err)
			walk.MsgBox(mw, "Error connecting", fmt.Sprintf("Failed to connect to: %s:%d\n %v", ipaddress, ipport, err), walk.MsgBoxIconWarning)
		}

		if err == nil {
			connected = true
			log.Printf("Connected to %s:%d", ipaddress, ipport)
		}
		send("<STATUS>")
		send("<LEDS288>")
		send("<CLOCK22:55:05>")

	} else {
		log.Printf("Already connected - skipping")

	}
	return err
	// defer conn.Close()
}

func (mw *MyMainWindow) disconnect() error {
	var err error
	if printonly {
		log.Println("Printonly enabled, not disconnecting")
		return nil
	}
	err = conn.Close()
	if err != nil {
		log.Printf("Error disconnecting: %v", err)
	} else {
		connected = false
		log.Println("Disconnected")
	}
	return err
}

func (mw *MyMainWindow) sendLine(line string) error {
	var err error
	if !printonly {
		if !connected {
			err = mw.connect()
			if err != nil {
				return err
			}
			defer mw.disconnect()
		}
	}

	if mw.permanentAction.Checked() {
		line = linijka.InjectFlag(line, "<PERMANENT>")
	}

	send(line)

	return err
}

func send(line string) {
	linijka.LinijkaWriter(log.Writer(), line)
	if !printonly {
		linijka.LinijkaWriter(conn, line)
		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to send: %v", err)
		}
		log.Printf("Response: %s", status)
	}
}
