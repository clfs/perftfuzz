package engine

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
)

type translator struct {
	name   string
	encode func(*Request) (string, error)
	decode func(io.Reader) (*Response, error)
}

var (
	translatorsMutex  sync.Mutex
	atomicTranslators atomic.Value
)

type EncodeFunc func(*Request) (string, error)

type DecodeFunc func(io.Reader) (*Response, error)

func Register(name string, e EncodeFunc, d DecodeFunc) {
	translatorsMutex.Lock()
	defer translatorsMutex.Unlock()
	translators, _ := atomicTranslators.Load().([]translator)
	atomicTranslators.Store(append(translators, translator{name, e, d}))
}

type Request struct {
	Position string   `json:"position"`
	Moves    []string `json:"moves"`
	Depth    int      `json:"depth"`
}

type Response struct {
	Perft map[string]int `json:"perft"`
}

type Engine struct {
	cmd *exec.Cmd
	t   translator
	w   io.Writer // write to the underlying engine
	r   io.Reader // read from the underlying engine
}

func New(ctx context.Context, path string) (*Engine, error) {
	cmd := exec.CommandContext(ctx, path)

	w, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %v", err)
	}
	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %v", err)
	}

	eng := &Engine{cmd: cmd, w: w, r: r}

	name, err := eng.name()
	if err != nil {
		return nil, fmt.Errorf("name: %v", err)
	}

	translators, _ := atomicTranslators.Load().([]translator)
	i := slices.IndexFunc(translators, func(t translator) bool { return t.name == name })
	if i == -1 {
		return nil, fmt.Errorf("unknown name: %q", name)
	}

	eng.t = translators[i]

	return eng, nil
}

func (e *Engine) do(req *Request) (*Response, error) {
	commands, err := e.t.encode(req)
	if err != nil {
		return nil, fmt.Errorf("encode: %v", err)
	}

	fmt.Fprintf(e.w, "%s", commands)

	resp, err := e.t.decode(e.r)
	if err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	return resp, nil
}

func (e *Engine) name() (string, error) {
	var res string

	fmt.Fprintln(e.w, "uci")

	// TODO(clfs): Replace with bytes.Lines in Go 1.24.

	s := bufio.NewScanner(e.r)
	for s.Scan() {
		line := s.Text()
		if name, ok := strings.CutPrefix(line, "id name "); ok {
			res = name
			break
		}
	}

	if err := s.Err(); err != nil {
		return "", fmt.Errorf("scanner: %v", err)
	}

	if res == "" {
		return "", errors.New("name not found")
	}

	return res, nil
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, httpReq *http.Request) {
	var req Request

	if err := json.NewDecoder(httpReq.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	resp, err := e.do(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway) // 502
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusBadGateway) // 502
		return
	}

	w.WriteHeader(http.StatusOK) // 200
}

func (e *Engine) Close() error {
	fmt.Fprintln(e.w, "quit")
	return e.cmd.Wait()
}
