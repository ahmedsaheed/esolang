#!/usr/bin/env bash

hyperfine './dist/esolang_darwin_arm64/esolang ./benchmarkings/bench.eso' 'python3 ./benchmarkings/bench.py' -i --export-markdown ./benchmarkings/bench.md 
