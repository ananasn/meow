package webterm

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"io/ioutil"

	"github.com/tarm/serial"
	"golang.org/x/net/html"
)

const (
	inBufferSize  = 1
	outBufferSize = 1024
	colMax        = 114
	rowsDefault   = 29
)

var Doc *html.Node = nil              // contains tree of html pattern after HTMLToNodesTree was called
var NewLineChan = make(chan struct{}) // channel, throught that command "move to new line" can be sent
var rowCurrent int = rowsDefault      // number of the latest line; increments if current position is close to rowsDefault
var row, col = 0, 0                   // current cursor position
var re_html = regexp.MustCompile(`<tbody[\s\S]*</tbody>`)
var re = regexp.MustCompile(`(\x1b\[[\d{0,3};{0,1}]{0,3}[a-zA-Z]|\x0d\x0a)`)
var regMap = map[string]func(string, *html.Node){
	`\x0d\x0a`:                  SetCursToNewLine,
	`\x1b\[3\dm`:                SetForColor,
	`\x1b\[4\dm`:                SetBgColor,
	`\x1b\[0m`:                  ResetFormat,
	`\x1b\[\d{0,1}K`:            Erase,
	`\x1b\[2J`:                  EraseDisp,
	`\x1b\[\d{0,3};\d{0,3}[Hf]`: SetCursPos,
	`\x1b\[\d{0,3}A`:            SetCursUp,
	`\x1b\[\d{0,3}B`:            SetCursDown,
	`\x1b\[\d{0,3}D`:            SetCursBack,
	`\x1b\[\d{0,3}C`:            SetCursFow,
}

type TextArr []string

// Resets current cursor position.
func resetValues() {
	col, row = 0, 0
}

// Opens serial port, where console interface is set.
func StartSerial(dev string, baud int) (*serial.Port, error) {
	config := &serial.Config{Name: dev, Baud: baud}
	s, err := serial.OpenPort(config)
	if err != nil {
		fmt.Println("Error, cant open serial port")
		return nil, err
	}
	return s, nil
}

//Closes serial port, where console interface is set.
func CloseSerial(s *serial.Port) error {
	return s.Close()
}

// Reads message from serial port in buffer with size outBufferSize
// and returns a slice with size of recieved message.
func SerialRead(s *serial.Port) (string, error) {
	buf := make([]byte, outBufferSize)
	n, err := s.Read(buf)
	if err != nil {
		fmt.Println("Error, cant read from serial port")
		return "", err
	}
	return string(buf[:n]), nil
}

// Writes one byte in serial port.
func SerialWrite(msg byte, s *serial.Port) error {
	byte_arr := make([]byte, inBufferSize)
	byte_arr[0] = msg
	_, err := s.Write(byte_arr)
	if err != nil {
		return err
	}
	return nil
}

// Reads file, that contains html console pattern, with path,
// specified in name param and convert it to *html.Node tree.
func HTMLToNodesTree(name string) (*html.Node, error) {
	resetValues()
	dat, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(dat)
	Doc, e := html.Parse(r)
	if e != nil {
		return nil, e
	}
	return Doc, nil
}

// Cuts output html console pattern,
// and returns content of table tag.
func truncate(s string) string {
	return re_html.FindString(s)
}

// Returns node with specified id from nodes tree, i.e. the
// function allows to get a pointer to any available position in
// console at the moment. Id has a special format "%row_%colmn".
func getElementById(id string, n *html.Node) (element *html.Node, ok bool) {
	if n == nil {
		return nil, false
	}
	for _, a := range n.Attr {
		if a.Key == "id" && a.Val == id {
			return n, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if element, ok = getElementById(id, c); ok {
			return
		}
	}
	return
}

//Deletes text in node n and in it's children.
func DeleteAllTextNodes(n *html.Node) {
	if n == nil {
		fmt.Println("elem is nil\n")
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		DeleteAllTextNodes(c)
		if c.Type == html.TextNode {
			c.Data = ""
		}
	}
}

// Rerurns pointers to a current cursor node and to a table cell,
// where cursor is placed at the moment; param is a piece of tree,
// where cursor is must be found.
func getCursor(n *html.Node) (c *html.Node, p *html.Node, ok bool) {
	c, ok = getElementById("cursor", n)
	if ok == false {
		fmt.Println("No cursor")
		return nil, nil, false
	}
	p = c.Parent
	return c, p, true
}

// Deletes cursor node. Params are pointers to cursor and
// its's parent.
func deleteCursor(c *html.Node, p *html.Node) bool {
	if c == nil || p == nil || c.Parent != p {
		return false
	}
	p.RemoveChild(c)
	return true
}

// Secondary function. Sets cursor to specified position in the tree.
//Id has a special format "%row_%colmn".
func setCursor(position_id string, cursor *html.Node, Doc *html.Node) bool {
	if cursor.Parent != nil {
		return false
	}
	p, ok := getElementById(position_id, Doc)
	if ok == false {
		fmt.Printf("No such id %s!", position_id)
		return false
	}
	p.AppendChild(cursor)
	return true
}

// Creates id string.
func createId(row int, col int) string {
	return fmt.Sprintf("%v_%v", row, col)
}

// Deletes cursor from current position and sets it to new specified position.
func setCursorToPos(row int, col int, Doc *html.Node) (ok bool) {
	c, p, _ := getCursor(Doc)
	ok = deleteCursor(c, p)
	if ok == false {
		return false
	}
	ok = setCursor(createId(row, col), c, Doc)
	if ok == false {
		return false
	}
	return true
}

// Sets text to current position.
func setTextToCurrentPos(text string, Doc *html.Node) (ok bool) {
	for _, letter := range text { // Add text node to cursor position if it doesn't have any
		_, p, ok := getCursor(Doc)
		if ok == false {
			fmt.Printf("error in setTextToCurrentPos() %v\n", p)
			return false
		}
		found := false
		for c := p.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				c.Data = string(letter)
				found = true
			}
		}
		if !found {
			var textNode html.Node
			textNode.Type = html.TextNode
			textNode.Data = string(letter)
			p.AppendChild(&textNode)
		}
		col += 1
		setCursorToPos(row, col, Doc)
	}
	fmt.Printf("setTextToCurrentPos() %v%v %v\n", text, row, col)
	return true
}

// Creates the iterator, that iterates over string array,
// that was obtained by spliting input string at esc characters.
func (s TextArr) Next() func() string {
	i := 0
	return func() string {
		res := s[i]
		i += 1
		return res
	}
}

// Converts string with esc characters to html page.
func EscStringToHTML(s string) string {
	escArr := re.FindAllString(s, -1)   // Finds all esc characters
	textArr := TextArr(re.Split(s, -1)) // Finds text between esc characters
	next := textArr.Next()
	if textArr[0] == "" { //if esc code was at the begining
		fmt.Printf("First pos empty\n")
		next()

	} else {
		setTextToCurrentPos(next(), Doc)
	}
	for _, escSeq := range escArr {
		for key, value := range regMap {
			re := regexp.MustCompile(key)
			if re.MatchString(escSeq) { // if esc character is recognized
				value(escSeq, Doc) // call esc character handler
				if str := next(); str != "" {
					setTextToCurrentPos(str, Doc)
				}
				break
			}
		}
	}
	buf := new(bytes.Buffer)
	html.Render(buf, Doc)
	return truncate(buf.String())
}

// Set of functions for each escape sequense type

func SetForColor(s string, Doc *html.Node) {
	fmt.Println("SetFontColor()")
}

func SetBgColor(s string, Doc *html.Node) {
	fmt.Println("SetBgColor()")
}

func ResetFormat(s string, Doc *html.Node) {
	fmt.Println("ResetFormat()")
}

func Erase(s string, Doc *html.Node) {
	_, p, _ := getCursor(Doc)
	for n := p; n != nil; n = n.NextSibling {
		if n.FirstChild == nil {
			continue
		}
		if n.FirstChild.Type == html.TextNode {
			n.FirstChild.Data = ""
		}
	}
	//setCursorToPos(row, col, Doc)
	fmt.Println("Erase()")
}

func EraseDisp(s string, Doc *html.Node) {
	DeleteAllTextNodes(Doc)
	row = 0
	col = 0
	setCursorToPos(row, col, Doc)
	fmt.Println("EraseDisp()")
}

func SetCursPos(s string, Doc *html.Node) {
	res := strings.Split(string(s[2:len(s)-1]), ";")
	if len(res) != 2 {
		row = 0
		col = 0
		setCursorToPos(row, col, Doc)
		fmt.Println("SetCursPos()")
		return
	}
	if row, _ := strconv.Atoi(res[0]); row < 0 || row > rowCurrent {
		fmt.Printf("SetCursPos() row is %v", row)
		row = 0
	}
	if col, _ := strconv.Atoi(res[1]); col < 0 || col > colMax {
		fmt.Printf("SetCursPos() col is %v", col)
		col = 0
	}
	setCursorToPos(row, col, Doc)
	fmt.Println("SetCursPos()")
}

func SetCursUp(s string, Doc *html.Node) {
	offset, _ := strconv.Atoi(string(s[2 : len(s)-1]))
	if row = row - offset; row < 0 {
		row = 0
	}
	setCursorToPos(row, col, Doc)
	fmt.Println("SetCursUp()")
}

func SetCursDown(s string, Doc *html.Node) {
	offset, _ := strconv.Atoi(string(s[2 : len(s)-1])) //TODO: Change to something more appropriate
	if row = row + offset; row > rowCurrent {
		fmt.Println("SetCursDown) Overflow col")
	}
	setCursorToPos(row, col, Doc)
	fmt.Println("SetCursDown()")
}

func SetCursBack(s string, Doc *html.Node) {
	offset, _ := strconv.Atoi(string(s[2 : len(s)-1]))
	if col = col - offset; col < 0 {
		col = 0
	}
	setCursorToPos(row, col, Doc)
	fmt.Println("SetCursBack()")
}

func SetCursFow(s string, Doc *html.Node) {
	offset, _ := strconv.Atoi(string(s[2 : len(s)-1]))
	if col = col + offset; col > colMax { //TODO: Change to something more appropriate
		fmt.Println("SetCursFow() Overflow row")
		NewLineChan <- struct{}{}
		return
	}
	setCursorToPos(row, col, Doc)
	fmt.Println("SetCursFow()")
}

func SetCursToNewLine(s string, Doc *html.Node) {
	fmt.Println("SetCursToNewLine()")
	if row == rowCurrent-2 {

		//Copy last <tr> in newNode
		var newNode html.Node
		newNode.Type = html.ElementNode
		newNode.Data = "tr"
		//Add childs to newNode
		fmt.Printf("Add childs to newNode")
		for c := 0; c <= colMax; c++ {
			fmt.Println("Set Child!")
			var newChild html.Node
			newChild.Type = html.ElementNode
			newChild.Data = "td"
			var attr html.Attribute
			newChild.Attr = append(newChild.Attr, attr)
			newChild.Attr[0].Key = "id"
			//fmt.Println(string(rowCurrent+1) + "_" + strings.Split(a.Val, "_")[1])
			newChild.Attr[0].Val = fmt.Sprintf("%v_%v", rowCurrent+1, c)
			fmt.Println(fmt.Sprintf("%v_%v", rowCurrent+1, c))
			newNode.AppendChild(&newChild)

		}
		_, p, _ := getCursor(Doc)
		fmt.Println(p.Parent.Parent.Data)
		p.Parent.Parent.AppendChild(&newNode)
		rowCurrent += 1
		fmt.Println("SetCursToNewLine() Overflow col")
	}
	row += 1
	col = 0
	setCursorToPos(row, col, Doc)
}
