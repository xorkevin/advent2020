#!/bin/sh

set -e

awkprg1=$(cat <<'EOF'
/^bin/ {printf "%-8s %-8s %32s %32s %32s %32s %32s\n", "go", $2, $3, $4, $5, $6, $7; fflush()}
/^target/ {printf "%-8s %-8s %32s %32s %32s %32s %32s\n", "rs", $2, $3, $4, $5, $6, $7; fflush()}
EOF
)

awkprg2=$(cat <<'EOF'
{mean[$1] += $3; variance[$1] += $4^2; median[$1] += $5; min[$1] += $6; max[$1] += $7}
END {
  for (i in mean) {
    printf "%-8s %-8s %32s %32s %32s %32s %32s\n", i, "total", mean[i], sqrt(variance[i]), median[i], min[i], max[i]
  }
}
EOF
)

days=$(find . -maxdepth 1 -type d -name 'day*' -printf '%P\n' | sort)

statsdir=$(mktemp --tmpdir -d adventbench.XXXXXXXX)

printf '%-8s %-8s %32s %32s %32s %32s %32s\n' 'lang' 'day' 'mean' 'stddev' 'median' 'min' 'max' 1>&2

for i in $days; do
  statsfile="$statsdir/$i.json"
  cd "$i" && hyperfine -s basic --export-json "$statsfile" -w 8 -L bin "./bin/$i,./target/release/$i" '{bin}' 2>/dev/null && cd ".."
  jq -r '.results[] | (.command | capture("(?<kind>(bin|target/release))/day(?<day>[0-9]+)$")) as $cmd | "\($cmd.kind) \($cmd.day) \(.mean) \(.stddev) \(.median) \(.min) \(.max)"' "$statsfile"
done | awk "$awkprg1" | tee /dev/stderr | awk "$awkprg2"

rm -r "$statsdir"
