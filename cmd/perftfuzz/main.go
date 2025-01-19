package main

import (
	"context"
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

	http.ListenAndServe(":8080", eng)
}
