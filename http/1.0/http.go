package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/tomocy/wataame/http"
	http0_9 "github.com/tomocy/wataame/http/0.9"
)

const (
	MethodGet  = http0_9.MethodGet
	MethodHead = "HEAD"
	MethodPost = "POST"
)

type Request interface {
	request()
}

type SimpleRequest struct {
	http0_9.Request
}

func (r *SimpleRequest) request() {}

type FullRequest struct {
	RequestLine *RequestLine
	Header      Header
	Body        io.ReadCloser
}

func (r *FullRequest) request() {}

func (r *FullRequest) WriteTo(dst io.Writer) (int64, error) {
	r.assureRequiredHeader()
	n, err := fmt.Fprint(dst, r)
	return int64(n), err
}

func (r *FullRequest) assureRequiredHeader() error {
	r.Header.assureRequired()
	if err := r.assureContentLength(); err != nil {
		return fmt.Errorf("failed to assure required header: %s", err)
	}

	return nil
}

func (r *FullRequest) assureContentLength() error {
	if r.Body == nil {
		return nil
	}

	l, err := r.measureBody()
	if err != nil {
		return fmt.Errorf("failed to assure content length: %s", err)
	}

	r.Header["Content-Length"] = []string{fmt.Sprint(l)}

	return nil
}

func (r *FullRequest) measureBody() (int64, error) {
	var w bytes.Buffer
	teed := r.teeBody()
	n, err := io.Copy(&w, teed)
	if err != nil {
		return 0, fmt.Errorf("failed to measure body: %s", err)
	}

	return n, nil
}

func (r *FullRequest) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, r.RequestLine)
	fmt.Fprintln(&b, r.Header)
	fmt.Fprintln(&b)
	if r.Body != nil {
		teed := r.teeBody()
		io.Copy(&b, teed)
	}

	return b.String()
}

func (r *FullRequest) teeBody() io.Reader {
	var b bytes.Buffer
	teed := io.TeeReader(r.Body, &b)
	r.Body = ioutil.NopCloser(&b)

	return teed
}

func (r *FullRequest) ReadFrom(src io.Reader) (int64, error) {
	n, err := fmt.Fscan(src, r)
	return int64(n), err
}

func (r *FullRequest) Scan(state fmt.ScanState, _ rune) error {
	r.RequestLine, r.Header = new(RequestLine), make(Header)
	if _, err := fmt.Fscanln(state, r.RequestLine); err != nil {
		return fmt.Errorf("failed to scan full request: %s", err)
	}
	if _, err := fmt.Fscan(state, r.Header); err != nil {
		return fmt.Errorf("failed to scan full request: %s", err)
	}
	if _, _, err := state.ReadRune(); err != nil {
		return fmt.Errorf("failed to scan full request: %s", err)
	}
	if !willBeEOF(state) {
		r.Body = &body{
			size: r.Header.contentLength(),
		}
		if _, err := fmt.Fscan(state, r.Body); err != nil {
			return fmt.Errorf("failed to scan full request: %s", err)
		}
	}

	return nil
}

type RequestLine struct {
	Method  string
	URI     *url.URL
	Version *Version
}

func (l RequestLine) String() string {
	return fmt.Sprintf("%s %s %s", l.Method, l.URI.Path, l.Version)
}

func (l *RequestLine) Scan(state fmt.ScanState, _ rune) error {
	var uri http.ScannableURL
	l.Version = new(Version)
	if _, err := fmt.Fscanf(state, "%s %v %v", &l.Method, &uri, l.Version); err != nil {
		return fmt.Errorf("failed to scan request line: %s", err)
	}
	l.URI = uri.URL()

	return nil
}

type Response interface {
	response()
}

type SimpleResponse struct {
	http0_9.Response
}

func (r *SimpleResponse) response() {}

type FullResponse struct {
	StatusLine *StatusLine
	Header     Header
	Body       io.ReadCloser
}

func (r *FullResponse) response() {}

func (r *FullResponse) WriteTo(dst io.Writer) (int64, error) {
	n, err := fmt.Fprint(dst, r)
	return int64(n), err
}

func (r FullResponse) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, r.StatusLine)
	if 0 < len(r.Header) {
		fmt.Fprintln(&b, r.Header)
	}
	fmt.Fprintln(&b)
	if r.Body != nil {
		teed := r.teeBody()
		io.Copy(&b, teed)
	}

	return b.String()
}

func (r *FullResponse) teeBody() io.Reader {
	var b bytes.Buffer
	teed := io.TeeReader(r.Body, &b)
	r.Body = ioutil.NopCloser(&b)

	return teed
}

func (r *FullResponse) ReadFrom(src io.Reader) (int64, error) {
	n, err := fmt.Fscan(src, r)
	return int64(n), err
}

func (r *FullResponse) Scan(state fmt.ScanState, _ rune) error {
	r.StatusLine, r.Header = new(StatusLine), make(Header)
	if _, err := fmt.Fscanln(state, r.StatusLine); err != nil {
		return fmt.Errorf("failed to scan full response: %s", err)
	}
	if _, err := fmt.Fscan(state, r.Header); err != nil {
		return fmt.Errorf("failed to scan full response: %s", err)
	}
	if _, _, err := state.ReadRune(); err != nil {
		return fmt.Errorf("failed to scan full response: %s", err)
	}
	if !willBeEOF(state) {
		r.Body = new(body)
		if _, err := fmt.Fscan(state, r.Body); err != nil {
			return fmt.Errorf("failed to scan full response: %s", err)
		}
	}

	return nil
}

type StatusLine struct {
	Version *Version
	Status  *Status
}

func (l StatusLine) String() string {
	return fmt.Sprintf("%s %s", l.Version, l.Status)
}

func (l *StatusLine) Scan(state fmt.ScanState, _ rune) error {
	l.Version, l.Status = new(Version), new(Status)
	if _, err := fmt.Fscan(state, l.Version, l.Status); err != nil {
		return fmt.Errorf("failed to scan status line: %s", err)
	}

	return nil
}

type Version struct {
	Major, Minor int
}

func (v Version) String() string {
	return fmt.Sprintf("HTTP/%d.%d", v.Major, v.Minor)
}

func (v *Version) Scan(state fmt.ScanState, _ rune) error {
	if _, err := fmt.Fscanf(state, "HTTP/%d.%d", &v.Major, &v.Minor); err != nil {
		return fmt.Errorf("failed to scan version: %s", err)
	}

	return nil
}

type Status struct {
	Code   int
	Phrase string
}

func (s Status) String() string {
	return fmt.Sprintf("%d %s", s.Code, s.Phrase)
}

func (s *Status) Scan(state fmt.ScanState, _ rune) error {
	if _, err := fmt.Fscanf(state, "%d %s", &s.Code, &s.Phrase); err != nil {
		return fmt.Errorf("failed to scan status: %s", err)
	}

	return nil
}

type Header map[string][]string

func (h Header) String() string {
	var b strings.Builder
	ns := h.names()
	sort.Sort(ns)
	for _, n := range ns {
		for _, v := range h[n.name] {
			fmt.Fprintf(&b, "%s: %s\n", n.name, v)
		}
	}

	return strings.TrimSuffix(b.String(), "\n")
}

func (h Header) names() headerFieldNames {
	ns := make(headerFieldNames, len(h))
	var i int
	for n := range h {
		var kind int
		switch {
		case search(generalHeaders, n):
			kind = headerGeneral
		case search(responseHeader, n):
			kind = headerResponse
		case search(requestHeader, n):
			kind = headerRequest
		default:
			kind = headerEntity
		}

		ns[i] = headerFieldName{
			name: n, kind: kind,
		}
		i++
	}

	return ns
}

func (h Header) assureRequired() {
	for k, vs := range requiredHeader {
		if _, ok := h[k]; ok {
			continue
		}

		h[k] = vs
	}
}

func (h Header) contentLength() int {
	var n int
	if ls, ok := h["Content-Length"]; ok {
		n, _ = strconv.Atoi(ls[0])
	}

	return n
}

var requiredHeader = Header{
	"Content-Length": []string{"0"},
}

func (h Header) Scan(state fmt.ScanState, _ rune) error {
	for {
		read, _, err := state.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to scan header: %s", err)
		}
		state.UnreadRune()
		if read == '\n' {
			break
		}

		var f headerField
		if _, err := fmt.Fscanln(state, &f); err != nil {
			return fmt.Errorf("failed to scan header: %s", err)
		}

		h[f.key] = append(h[f.key], f.vals...)
	}

	return nil
}

type headerField struct {
	key  string
	vals []string
}

func (f *headerField) Scan(state fmt.ScanState, _ rune) error {
	if err := f.scanKey(state); err != nil {
		return fmt.Errorf("failed to scan header field: %s", err)
	}
	state.ReadRune()
	state.ReadRune()
	if err := f.scanValues(state); err != nil {
		return fmt.Errorf("failed to scan header field: %s", err)
	}

	return nil
}

func (f *headerField) scanKey(r io.RuneScanner) error {
	var k []rune
	for {
		read, _, err := r.ReadRune()
		if err != nil {
			return fmt.Errorf("failed to scan key of header field: %s", err)
		}
		if read == ':' {
			r.UnreadRune()
			break
		}

		k = append(k, read)
	}

	f.key = string(k)

	return nil
}

func (f *headerField) scanValues(r io.RuneReader) error {
	var (
		vs [][]rune
		v  []rune
	)
	for {
		read, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to scan values of header field: %s", err)
		}
		if read == '\n' {
			break
		}
		if f.isListable() && read == ',' {
			vs, v = append(vs, v), nil
			continue
		}

		v = append(v, read)
	}
	vs = append(vs, v)

	f.vals = make([]string, len(vs))
	for i, v := range vs {
		f.vals[i] = strings.TrimLeft(string(v), " ")
	}

	return nil
}

func (f *headerField) isListable() bool {
	sort.Strings(listableHeaders)
	return search(listableHeaders, f.key)
}

var listableHeaders = []string{
	"Allow", "WWW-Authenticate",
}

func search(vs []string, x string) bool {
	begin, end := 0, len(vs)-1
	for begin <= end {
		mid := (begin + end) / 2
		if vs[mid] == x {
			return true
		}
		if vs[mid] < x {
			begin = mid + 1
		} else {
			end = mid - 1
		}
	}

	return false
}

type headerFieldNames []headerFieldName

func (ns headerFieldNames) Len() int {
	return len(ns)
}

func (ns headerFieldNames) Less(i, j int) bool {
	if ns[i].kind != ns[j].kind {
		return ns[j].kind < ns[i].kind
	}

	return ns[i].name < ns[j].name
}

func (ns headerFieldNames) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

func (ns headerFieldNames) names() []string {
	names := make([]string, len(ns))
	for i, n := range ns {
		names[i] = n.name
	}

	return names
}

type headerFieldName struct {
	name string
	kind int
}

const (
	headerEntity = iota
	headerRequest
	headerResponse
	headerGeneral
)

var entityHeader = []string{
	"Allow", "Content-Encoding", "Content-Length", "Content-Type", "Expires", "Last-Modified",
}

var requestHeader = []string{
	"Authorization", "From", "If-Modified-Since", "Referer", "User-Agent",
}

var responseHeader = []string{
	"Location", "Server", "WWW-Authenticate",
}

var generalHeaders = []string{
	"Date", "Pragma",
}

func willBeEOF(s io.RuneScanner) bool {
	_, _, err := s.ReadRune()
	s.UnreadRune()
	return err == io.EOF
}

type body struct {
	bytes.Buffer
	size int
}

func (b *body) Scan(state fmt.ScanState, _ rune) error {
	for i := 0; i < b.size; {
		read, n, err := state.ReadRune()
		i += n
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to scan body: %s", err)
		}

		if _, err := b.WriteRune(read); err != nil {
			return fmt.Errorf("failed to scan body: %s", err)
		}
	}

	return nil
}

func (b *body) Close() error {
	c := ioutil.NopCloser(b)
	return c.Close()
}
