# syntax=docker/dockerfile:1

# Install ucihttp.
FROM golang:1.23.5 AS ucihttp

WORKDIR /tmp
COPY . ./
RUN go build ./cmd/ucihttp

# Install Stockfish.
FROM ubuntu:24.04 AS stockfish

WORKDIR /tmp
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    wget=1.21.4-1ubuntu4.1 \
    build-essential=12.10ubuntu1
RUN wget "https://github.com/official-stockfish/Stockfish/archive/refs/tags/sf_17.tar.gz" \
  && tar xf sf_17.tar.gz \
  && rm sf_17.tar.gz
WORKDIR /tmp/Stockfish-sf_17/src
RUN make net \
  && make -j profile-build

# Run ucihttp.
FROM scratch

COPY --from=perftfuzz /tmp/cmd/ucihttp /ucihttp
COPY --from=stockfish /tmp/Stockfish-sf_17/src/stockfish /stockfish
CMD ["/ucihttp", "--path", "/stockfish"]
