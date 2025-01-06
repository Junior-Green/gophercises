package devserver

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

type recoverPanic struct {
	handler    http.Handler
	production bool
}

func (rp *recoverPanic) catchPanic(w http.ResponseWriter, _ *http.Request) {
	e := recover()
	if e == nil {
		return
	}
	if bw, ok := w.(*bufferedResponseWriter); ok {
		bw.Clear()
	}
	fmt.Fprintln(os.Stderr, e)

	if rp.production {
		http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("<h1>Panic: %v</h1>\n", e))

	trace := debug.Stack()
	fmt.Fprintln(os.Stderr, string(trace))
	buffer.Write(trace)

	if paths, err := parseStackTrace(trace); err == nil {
		buffer.WriteString("<ul>")
		for _, path := range paths {
			path, lineNum := path[0], path[1]
			params := url.Values{}
			params.Add("line", lineNum)
			href := fmt.Sprintf("/debug%s?%s", path, params.Encode())
			html := fmt.Sprintf("<li><a href=\"%s\">%s:%s</a></li>", href, path, lineNum)
			buffer.WriteString(html)
		}
		buffer.WriteString("</ul>")
	}

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, buffer.String())
}

func (rp *recoverPanic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bw := &bufferedResponseWriter{writer: w}
	defer bw.Flush()
	defer rp.catchPanic(bw, r)

	rp.handler.ServeHTTP(bw, r)
}

type bufferedResponseWriter struct {
	writer http.ResponseWriter
	buffer bytes.Buffer
	status int
}

func (brw *bufferedResponseWriter) Header() http.Header {
	return brw.writer.Header()
}

func (brw *bufferedResponseWriter) Write(data []byte) (int, error) {
	return brw.buffer.Write(data)
}

func (brw *bufferedResponseWriter) Clear() {
	brw.buffer.Reset()
}

func (brw *bufferedResponseWriter) WriteHeader(statusCode int) {
	brw.status = statusCode
}

func (brw *bufferedResponseWriter) Flush() {
	if brw.status != 0 {
		brw.writer.WriteHeader(brw.status)
	}
	brw.writer.Write(brw.buffer.Bytes())
}
func Start() {
	var prod bool

	flag.BoolVar(&prod, "p", false, "run in production mode")
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/debug/", srcCodeHandler)
	mux.HandleFunc("/", hello)

	wrapped := newRecoverPanic(mux, prod)

	log.Fatal(http.ListenAndServe(":3000", wrapped))
}
func newRecoverPanic(toWrap http.Handler, prod bool) *recoverPanic {
	return &recoverPanic{toWrap, prod}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func srcCodeHandler(w http.ResponseWriter, r *http.Request) {
	filePath, ok := strings.CutPrefix(r.URL.Path, "/debug/")
	if !ok {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	file, err := openSourceFile(filePath)
	if err != nil {
		errMsg := fmt.Sprintf("error opening: %q", filePath)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	src, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "error occured retrieving source code", http.StatusInternalServerError)
		return
	}

	var ranges [][2]int

	line := r.URL.Query().Get("line")
	if lineNum, err := strconv.Atoi(line); err == nil {
		ranges = append(ranges, [2]int{lineNum, lineNum})
	}

	newSrc, err := beautifySource(src, path.Ext(filePath), ranges...)
	if err != nil {
		http.Error(w, "error occured beautifying source code", http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(newSrc); err != nil {
		http.Error(w, "error writing to body", http.StatusInternalServerError)
		return
	}
}

func openSourceFile(path string) (*os.File, error) {
	path, _ = strings.CutPrefix(path, "/")
	var (
		file *os.File
		err  error
	)
	//try relative
	file, err = os.Open(path)
	if err == nil {
		return file, nil
	}
	//try absolute
	file, err = os.Open("/" + path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func parseStackTrace(trace []byte) ([][2]string, error) {
	regex, err := regexp.Compile("(/.+.[a-zA-Z]+):([0-9]+)")
	if err != nil {
		return nil, err
	}

	matches := regex.FindAllSubmatch(trace, -1)
	paths := make([][2]string, 0, len(matches))

	for _, match := range matches {
		if len(match) != 3 {
			return nil, fmt.Errorf("something unexpected occured while parsing stack trace")
		}
		path, line := string(match[1]), string(match[2])
		paths = append(paths, [2]string{path, line})
	}

	return paths, nil
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func beautifySource(src []byte, extentsion string, highlightRanges ...[2]int) ([]byte, error) {
	lexer := lexers.Get(extentsion)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.Standalone(true), html.WithLineNumbers(true), html.HighlightLines(highlightRanges))

	it, err := lexer.Tokenise(nil, string(src))
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	if err = formatter.Format(&buffer, style, it); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
