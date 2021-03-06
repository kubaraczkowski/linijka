package linijka

import (
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func encodeWindows1250(inp string) string {
	enc := charmap.Windows1250.NewEncoder()
	out, _ := enc.String(inp)
	return out
}

func tobytes(inp string) []byte {
	return []byte(encodeWindows1250(inp))
}

func sumxor(text []byte) (int, int) {
	var s byte
	var x byte
	for _, ch := range text {
		s += ch
		//s = s
		x = x ^ ch
	}
	return int(s), int(x)
}

func addstart(text string) string {
	if !strings.HasPrefix(text, "<START") {
		return fmt.Sprintf("<START1>%s", text)
	} else {
		return text
	}
}

func Wrapincrc(text string) string {
	text = addstart(text)
	s, x := sumxor(tobytes(text))
	return fmt.Sprintf("%s<STOP%X%X>\r\n", text, s, x)

}

func checkspecial(s []string, e string) bool {
	for _, a := range s {
		if strings.HasPrefix(e, a) {
			return true
		}
	}
	return false
}

func InjectFlag(text string, flag string) string {
	var header string
	if strings.HasPrefix(text, "<START") {
		split_text := strings.SplitN(text, ">", 2)
		if len(split_text) == 2 {
			header = fmt.Sprintf("%s>", split_text[0])
			text = split_text[1]
		} else {
			log.Fatalf("Can't split, got :%s, len(split): %v", text, split_text)
		}
	}
	output := fmt.Sprintf("%s%s%s", header, flag, text)
	return output
}

func LinijkaWriter(w io.Writer, text string) (n int, err error) {
	var specials = []string{"<STATUS>", "<LEDS", "<CLOCK", "<TIME", "<SETP", "<RESETP"}
	if checkspecial(specials, text) {
	} else {
		text = Wrapincrc(text)
	}
	text = fmt.Sprintf("%s\r\n", strings.TrimSpace(text))
	n, err = w.Write(tobytes(text))
	return n, err
}
