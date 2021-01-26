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

const version string = "0.1.0"

func main() {
	var ipaddress_string string
	var port int
	var oneline bool
	flag.StringVar(&ipaddress_string, "ip", "127.0.0.1", "IP address of the device")
	flag.IntVar(&port, "port", 4001, "IP port of the device")
	flag.BoolVar(&oneline, "oneline", false, "Don't split the passed line")
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
	log.Printf("Connecting to %s:%d", ipaddress, port)

	// for _, arg := range flag.Args() {
	// 	linijka.LinijkaWriter(log.Writer(), arg)
	// }

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", ipaddress, port))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	var lines []string

	lines = append(lines, "<STATUS>")
	lines = append(lines, "<LEDS288>")
	lines = append(lines, "<CLOCK22:55:05>")

	if oneline {
		lines = append(lines, strings.Join(flag.Args(), " "))
	} else {
		lines = append(lines, flag.Args()...)
	}

	for _, arg := range lines {
		linijka.LinijkaWriter(conn, arg)
		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		log.Printf("Response: %s", status)
	}
}
