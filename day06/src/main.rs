use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut count = 0;
    let mut set = HashSet::new();

    let mut count2 = 0;
    let mut set2 = HashSet::new();
    let mut is_init = false;

    for line in reader.lines() {
        let line = line?;
        if line == "" {
            count += set.len();
            set = HashSet::new();

            count2 += set2.len();
            is_init = false;
            continue;
        }
        let chars = line.chars().collect::<HashSet<char>>();
        set = set.union(&chars).cloned().collect();

        if is_init {
            set2 = set2.intersection(&chars).cloned().collect();
        } else {
            is_init = true;
            set2 = chars;
        }
    }

    count += set.len();
    count2 += set2.len();

    println!("Part 1: {}", count);
    println!("Part 2: {}", count2);

    Ok(())
}
