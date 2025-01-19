package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/clfs/perftfuzz/engine"
	_ "github.com/clfs/perftfuzz/engine/stockfish"
)

func main() {
	var pathFlag = flag.String("path", "", "absolute path to the chess engine")
	flag.Parse()

	if *pathFlag == "" {
		flag.Usage()
		return
	}

	ctx := context.Background()

	eng, err := engine.New(ctx, *pathFlag)
	if err != nil {
		log.Fatalf("new engine: %v", err)
	}
	defer eng.Close()

	fn := func(w http.ResponseWriter, httpReq *http.Request) {
		var req engine.Request

		if err := json.NewDecoder(httpReq.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			return
		}

		resp, err := eng.Do(&req)
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

	http.Handle("/", http.HandlerFunc(fn))
	http.ListenAndServe(":8080", nil)
}
