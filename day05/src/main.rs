use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let v = {
        let mut v = Vec::new();
        for line in reader.lines() {
            let line = line?;
            let chars = line.chars().collect::<Vec<_>>();
            v.push(calc_id(&chars));
        }
        v.sort();
        v
    };

    println!("Part 1: {}", v.last().ok_or("No ids")?);

    let mut prev = v.first().ok_or("No ids")?;
    for i in &v {
        if i - prev > 1 {
            println!("Part 2: {}", i - 1)
        }
        prev = i
    }

    Ok(())
}

fn calc_id(b: &Vec<char>) -> usize {
    let mut n = 0;
    for &i in b {
        n *= 2;
        if i == 'B' || i == 'R' {
            n += 1;
        }
    }
    n
}
