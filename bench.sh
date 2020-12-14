#!/bin/sh

export benchargs="-w 256 -s basic -u millisecond"

days=$(find . -maxdepth 1 -type d -name 'day*' -printf '%P,%f\n' | sort)

awkprg1=$(cat <<'EOF'
/^Group/ {printf "%s", $2}
/^Benchmark #1/ {printf "##go##"}
/^Benchmark #2/ {printf "##rs##"}
/^Time/ {printf "%s %s ", $2, $3}
/^Range/ {printf "%s %s", $2, $3}
/^Summary/ {printf "\n"; fflush()}
EOF
)

awkprg2=$(cat <<'EOF'
{printf "%s %s %s\n%s %s %s\n", $1, $2, $3, $1, $4, $5; fflush()}
EOF
)

awkprg3=$(cat <<'EOF'
{printf "%-8s %-8s %16s %16s %16s\n", $1, $2, $3, $5, $6; fflush()}
EOF
)

awkprg4=$(cat <<'EOF'
{mean[$2] += $3; min[$2] += $4; max[$2] += $5}
END {
  for (i in mean) {
    printf "%-8s %-8s %16s %16s %16s\n", "total", i, mean[i], min[i], max[i]
  }
}
EOF
)

printf '%-8s %-8s %16s %16s %16s\n' "day" "lang" "mean (ms)" "min (ms)" "max (ms)" 1>&2
for i in $days; do
  d=$(printf "$i" | cut -d , -f 1)
  bin=$(printf "$i" | cut -d , -f 2)
  printf 'Group %s\n' $bin
  (cd "$d" && hyperfine $benchargs -L bin "./bin/$bin,./target/release/$bin" '{bin}' 2> /dev/null)
done | sed -u \
  -e 's/^[ ]*Time .* \([0-9.]\+\) ms .* \([0-9.]\+\) ms .*\[.*$/Time \1 \2/' \
  -e 's/^[ ]*Range .* \([0-9.]\+\) ms .* \([0-9.]\+\) ms .*$/Range \1 \2/' \
  | awk "$awkprg1" \
  | awk -F '##' "$awkprg2" \
  | awk "$awkprg3" \
  | tee /dev/stderr \
  | awk "$awkprg4" \
  | sort
