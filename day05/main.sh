#!/bin/sh

nums=$(cat input.txt \
  | tr 'FLBR' '0011' \
  | while read -r line; do printf '%d\n' "$((2#$line))"; done \
  | sort -n)

printf 'Part 1: %s\n' $(printf '%s' "$nums" | tail -n 1)

first=$(printf '%s' "$nums" | head -n 1)
printf 'Part 2: '
printf '%s' "$nums" | awk -v i="$first" '$1 != i {while (i < $1) {print i; i++}}; {i++}'
