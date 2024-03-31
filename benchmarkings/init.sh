#!/usr/bin/env bash

hyperfine 'go run esolang.go ./benchmarkings/bench.eso' 'python3 ./benchmarkings/bench.py' -i --export-markdown ./benchmarkings/bench.md
