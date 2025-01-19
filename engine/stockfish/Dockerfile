FROM ubuntu:24.04

WORKDIR /engine

RUN apt update \
  && apt install -y wget build-essential \
  && rm -rf /var/lib/apt/lists/*

RUN wget "https://github.com/official-stockfish/Stockfish/releases/download/sf_17/stockfish-ubuntu-x86-64.tar" \
  && tar xf stockfish-ubuntu-x86-64.tar stockfish/stockfish-ubuntu-x86-64 \
  && rm stockfish-ubuntu-x86-64.tar

WORKDIR /app

