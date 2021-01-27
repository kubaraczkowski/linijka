package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	linijka "github.com/kubaraczkowski/linijka/pkg"
)

const version string = "0.2.0"

type param struct {
	permanent  *bool
	font       *int
	speed      *int
	pojawienie *bool
	gora       *bool
	dol        *bool
	flash      *int
}

func (p param) injectparam(text string) string {
	if *p.permanent {
		text = linijka.InjectFlag(text, "<PERMANENT>")
	}
	if *p.pojawienie {
		text = linijka.InjectFlag(text, "<POJAWIENIE>")
	}
	if *p.gora {
		text = linijka.InjectFlag(text, "<GORA>")
	}
	if *p.dol {
		text = linijka.InjectFlag(text, "<DOL>")
	}
	if *p.font != 0 {
		text = linijka.InjectFlag(text, fmt.Sprintf("<FONT%d>", *p.font))
	}
	if *p.speed != 0 {
		text = linijka.InjectFlag(text, fmt.Sprintf("<SPEED%d>", *p.speed))
	}
	if *p.flash != 0 {
		text = linijka.InjectFlag(text, fmt.Sprintf("<FLASH%d>", *p.speed))
	}
	return text
}

func main() {
	var ipaddress_string string
	var port int
	var oneline bool
	var printonly bool
	p := &param{}
	flag.StringVar(&ipaddress_string, "ip", "127.0.0.1", "IP address of the device")
	flag.IntVar(&port, "port", 4001, "IP port of the device")
	flag.BoolVar(&oneline, "oneline", false, "Don't split the passed line")
	flag.BoolVar(&printonly, "printonly", false, "Only print the messages to be sent, don't connect")
	version_flag := flag.Bool("v", false, "Display program version and exit")
	p.permanent = flag.Bool("permanent", false, "set <PERMANTENT> flag")
	p.font = flag.Int("font", 0, "set <FONT#> flag")
	p.speed = flag.Int("speed", 0, "set <SPEED#> flag")
	p.pojawienie = flag.Bool("pojawienie", false, "set <POJAWIENIE> flag")
	p.gora = flag.Bool("gora", false, "set <GORA> flag")
	p.dol = flag.Bool("dol", false, "set <DOL> flag")
	p.flash = flag.Int("flash", 0, "set <FLASH#> flag")

	flag.Parse()

	if *version_flag {
		fmt.Println(version)
		os.Exit(0)
	}

	ipaddress := net.ParseIP(ipaddress_string)
	if ipaddress == nil {
		log.Fatalf("Could not parse IP address: %s", ipaddress_string)
	}

	// for _, arg := range flag.Args() {
	// 	linijka.LinijkaWriter(log.Writer(), arg)
	// }

	var lines []string

	lines = append(lines, "<STATUS>")
	lines = append(lines, "<LEDS288>")
	lines = append(lines, "<CLOCK22:55:05>")

	var text string
	if oneline {
		text = strings.Join(flag.Args(), " ")
		text = p.injectparam(text)
		lines = append(lines, text)
	} else {
		for _, line := range flag.Args() {
			line = p.injectparam(line)
			lines = append(lines, line)
		}
	}

	for _, arg := range lines {
		linijka.LinijkaWriter(log.Writer(), arg)
		if !printonly {
			log.Printf("Connecting to %s:%d", ipaddress, port)
			var d net.Dialer
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", ipaddress, port))
			if err != nil {
				log.Fatalf("Failed to connect: %v", err)
			}
			defer conn.Close()
			linijka.LinijkaWriter(conn, arg)
			status, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Fatalf("Failed to connect: %v", err)
			}
			log.Printf("Response: %s", status)
		}
	}
}
