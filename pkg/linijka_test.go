package linijka

import (
	"net"
	"testing"
)

func TestWrap(t *testing.T) {
	var s string

	s = Wrapincrc("<START2>Strona 2 ")
	if s != "<START2>Strona 2 <STOP2377>\r\n" {
		t.Errorf("Got: %s", s)
	}

	s = Wrapincrc("<START9>Przykład strona 9")
	if s != "<START9>Przykład strona 9<STOPE9AB>\r\n" {
		t.Errorf("Got: %s", s)
	}

	s = Wrapincrc("<START1>Przykładowa <PAUSE1>strona linijki dynamicznej")
	if s != "<START1>Przykładowa <PAUSE1>strona linijki dynamicznej<STOPBEA6>\r\n" {
		t.Errorf("Got: %s", s)
	}

	s = Wrapincrc("<START1>Przykładowa <PAUSE1>strona linijki dynamicznej")
	if s != "<START1>Przykładowa <PAUSE1>strona linijki dynamicznej<STOPBEA6>\r\n" {
		t.Errorf("Got: %s", s)
	}

	s = Wrapincrc("Przykładowa <PAUSE1>strona linijki dynamicznej")
	if s != "<START1>Przykładowa <PAUSE1>strona linijki dynamicznej<STOPBEA6>\r\n" {
		t.Errorf("Got: %s", s)
	}
}

func TestWriter(t *testing.T) {
	server, client := net.Pipe()
	text := "<START1>Przykładowa <PAUSE1>strona linijki dynamicznej<STOPBEA6>\r\n"
	go func() {
		// b, err := ioutil.ReadAll(client)
		n, err := client.Read()
		t.Log(string(b))
		if err != nil {
			t.Error(err)
		}
		s := string(b)
		if s != text {
			t.Errorf("Got: %s", s)
		}
		client.Close()
	}()
	LinijkaWriter(server, "<START9>Przykład strona 9")
	server.Close()
}
