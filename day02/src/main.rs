use regex::Regex;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let re = Regex::new(r"^(\d+)-(\d+) ([a-z]): ([a-z]+)$")?;

    let mut valid = 0;
    let mut valid2 = 0;

    for line in reader.lines() {
        let line = line?;
        let caps = re.captures(&line).ok_or("Failed to match line")?;
        let a = caps
            .get(1)
            .ok_or("Failed to get 1")?
            .as_str()
            .parse::<usize>()?;
        let b = caps
            .get(2)
            .ok_or("Failed to get 2")?
            .as_str()
            .parse::<usize>()?;
        let ch = caps.get(3).ok_or("Failed to get 3")?.as_str();
        let c = ch.parse::<char>()?;
        let pass = caps.get(4).ok_or("Failed to get 4")?.as_str();
        let count = pass
            .chars()
            .map(|i| if i == c { 1 } else { 0 })
            .sum::<usize>();
        if count >= a && count <= b {
            valid += 1;
        }
        match (pass.get((a - 1)..a), pass.get((b - 1)..b)) {
            (Some(c1), Some(c2)) => {
                if c1 != c2 && (c1 == ch || c2 == ch) {
                    valid2 += 1;
                }
            }
            (_, _) => (),
        }
    }

    println!("Part 1: {}", valid);
    println!("Part 2: {}", valid2);
    Ok(())
}
