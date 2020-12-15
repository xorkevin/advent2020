fn main() {
    let puzzle_input = vec![19, 20, 14, 0, 9, 1];
    let mut hist = vec![0; 30_000_000];
    for (n, &i) in puzzle_input[..puzzle_input.len() - 1].iter().enumerate() {
        hist[i] = n + 1;
    }
    let mut idx = puzzle_input.len();
    let mut prev = puzzle_input[puzzle_input.len() - 1];
    while idx < 2020 {
        let v = hist[prev];
        hist[prev] = idx;
        if v == 0 {
            prev = 0;
        } else {
            prev = idx - v;
        }
        idx += 1;
    }
    println!("Part 1: {}", prev);
    while idx < 30_000_000 {
        let v = hist[prev];
        hist[prev] = idx;
        if v == 0 {
            prev = 0;
        } else {
            prev = idx - v;
        }
        idx += 1;
    }
    println!("Part 2: {}", prev);
}
