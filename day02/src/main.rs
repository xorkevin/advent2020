use regex::Regex;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let re = Regex::new(r"^(\d+)-(\d+) (a-z): ([a-z]+)$")?;

    let mut valid = 0;
    let mut valid2 = 0;

    for line in reader.lines() {
        let caps = re.captures(&line?)?;
        let a = caps.get(1)?.parse::<i32>()?;
        let b = caps.get(2)?.parse::<i32>()?;
        let c = caps.get(3)?.parse::<char>()?;
        let pass = caps.get(4)?;
        let count = pass.chars().map(|i| if i == c { 1 } else { 0 }).sum();
        if count >= a && count <= b {
            valid += 1;
        }
    }

    println!("Part 1: {}", valid);
    Ok(())
}
