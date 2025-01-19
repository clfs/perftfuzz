package stockfish

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/clfs/perftfuzz/engine"
)

func init() {
	engine.Register(Name, Encode, Decode)
}

const Name = "Stockfish 17"

func Encode(req *engine.Request) (string, error) {
	var sb strings.Builder

	fmt.Fprintf(&sb, "position %s", req.Position)
	if len(req.Moves) > 0 {
		fmt.Fprintf(&sb, " moves %s", strings.Join(req.Moves, " "))
	}
	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "go perft %d\n", req.Depth)

	return sb.String(), nil
}

var perftRegexp = regexp.MustCompile(`^(\w+): (\d+)$`)

func Decode(r io.Reader) (*engine.Response, error) {
	perft := make(map[string]int)

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()

		m := perftRegexp.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		k := m[0]
		v, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, fmt.Errorf("bad line: %q", line)
		}

		perft[k] = v
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("scanner: %v", err)
	}

	return &engine.Response{Perft: perft}, nil
}
