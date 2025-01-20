# syntax=docker/dockerfile:1

# Install perftfuzz.
FROM golang:1.23.5 AS perftfuzz

WORKDIR /tmp
COPY . ./
RUN go build ./cmd/perftfuzz

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

# Run perftfuzz.
FROM scratch

COPY --from=perftfuzz /tmp/cmd/perftfuzz /perftfuzz
COPY --from=stockfish /tmp/Stockfish-sf_17/src/stockfish /stockfish
CMD ["/perftfuzz", "--path", "/stockfish"]
