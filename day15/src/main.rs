use std::collections::HashMap;

fn main() {
    let puzzle_input = vec![19, 20, 14, 0, 9, 1];
    let boundary = 1 << 24;
    let mut hist = vec![0; boundary];
    let mut hist2: HashMap<usize, usize> = HashMap::new();
    for (n, &i) in puzzle_input[..puzzle_input.len() - 1].iter().enumerate() {
        if i < boundary {
            hist[i] = n + 1;
        } else {
            hist2.insert(i, n + 1);
        }
    }
    let mut idx = puzzle_input.len();
    let mut prev = puzzle_input[puzzle_input.len() - 1];
    while idx < 2020 {
        let v = if prev < boundary {
            let v = hist[prev];
            hist[prev] = idx;
            v
        } else {
            let v = *hist2.get(&prev).unwrap_or(&0);
            hist2.insert(prev, idx);
            v
        };
        if v == 0 {
            prev = 0;
        } else {
            prev = idx - v;
        }
        idx += 1;
    }
    println!("Part 1: {}", prev);
    while idx < 30_000_000 {
        let v = if prev < boundary {
            let v = hist[prev];
            hist[prev] = idx;
            v
        } else {
            let v = *hist2.get(&prev).unwrap_or(&0);
            hist2.insert(prev, idx);
            v
        };
        if v == 0 {
            prev = 0;
        } else {
            prev = idx - v;
        }
        idx += 1;
    }
    println!("Part 2: {}", prev);
}
