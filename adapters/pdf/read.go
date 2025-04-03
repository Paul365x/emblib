package pdf

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func get_strings(path string) *bytes.Buffer {
	cmd := exec.Command("pdftotext", "-layout", path, "-")
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	return buf
}

type meta struct {
	title       string
	id          string
	size        string
	stitches    string
	colcnt      string
	description string
	chgs        string
}

type token struct {
	Tok     string
	Content string
}

func tokenise_EmbLib(buf *bytes.Buffer) []token {
	var tok []token
	words := bytes.Fields(buf.Bytes())

	state := "skip"

	for i := 0; i < len(words); {
		// this switch processes all key words except the CCs
		switch string(words[i]) {
		case "Name:", "Size:", "Stitches:":
			var t token
			t.Tok = string(words[i])
			tok = append(tok, t)
			state = "content"
			i++
			continue
		case "Product":
			switch string(words[i+1]) {
			case "ID:", "Comments:":
				var t token
				t.Tok = string(words[i]) + " " + string(words[i+1])
				tok = append(tok, t)
				state = "content"
				i += 2
				continue
			}
		case "Color":
			if string(words[i+1]) == "Changes:" {
				var t token
				t.Tok = string(words[i]) + " " + string(words[i+1])
				tok = append(tok, t)
				state = "content"
				i += 2
				continue
			}
		case "Colors":
			if string(words[i+1]) == "Used:" {
				var t token
				t.Tok = string(words[i]) + " " + string(words[i+1])
				tok = append(tok, t)
				state = "content"
				i += 2
				continue
			}
		case "COLOR":
			i += 4
			continue
		case "UNIQUE":
			i += 3
			var t token
			t.Tok = "UC"
			tok = append(tok, t)
			state = "unique"
			continue
		}

		// catch the unique colors
		if state == "unique" {
			str := string(words[i])
			index := len(tok) - 1
			if len(tok[index].Content) > 0 {
				tok[index].Content = tok[index].Content + " " + str
			} else {
				tok[index].Content = str
			}
			if str[0:1] == "#" {
				var t token
				t.Tok = "UC"
				tok = append(tok, t)
				state = "unique"
			}

		}

		// catch the CCs
		str := string(words[i])
		if len(str) > 3 {
			if str[:2] == "CC" {
				var t token
				t.Tok = str
				tok = append(tok, t)
				state = "content"
				i++
				continue
			}
		}

		// this processes all text strings
		if state == "content" {
			index := len(tok) - 1
			if len(tok[index].Content) > 0 {
				tok[index].Content = tok[index].Content + " " + string(words[i])
			} else {
				tok[index].Content = string(words[i])
			}
			i++
			continue
		}
		i++ // throw away any skips
	}
	index := len(tok) - 1
	if tok[index].Tok == "UC" {
		fmt.Println("|" + tok[index].Content[0:23] + "|")
		fmt.Println("All Embroidery Library")

		if strings.Compare(tok[index].Content[0:22], "All Embroidery Library") == 0 {
			tok = tok[0 : index-1]
		}
	}
	return tok
}

/*
	func tokeniser_EmbLib(buf *bytes.Buffer) func () ( string, bool ) {
		return func () string {
			var tok string
			state := "single"
			for word := range bytes.Fields(buf.Bytes()) {
				switch word {
				case "Product" :
					state = "multi"
					tok = word
				case "Color" :
					state = "multi"
					tok = word
				case "Colors":
					state = "multi"
					tok = word
				}


		}
	}

	func parse_EmbLib( buf *bytes.Buffer ) *meta {
		tokens := bytes.Fields(buf.Bytes())
		m := new(meta)
		for t := range tokens {
			tmp := ""
			switch t {
			case "Name:" :
			case "Product ID:" :
			case "Size:":
			case "Color Changes:":

			}
		}

}

func main() {
	str := get_strings("pdf/F3222.pdf")
	toks := tokenise_EmbLib(str)
	fmt.Println(toks)

}*/
