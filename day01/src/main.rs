use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let nums = {
        let mut nums = HashSet::new();
        for line in reader.lines() {
            nums.insert(line?.parse::<i32>()?);
        }
        nums
    };

    for i in &nums {
        let k = 2020 - i;
        if nums.contains(&k) {
            println!("Part 1: {}", i * k);
            break;
        }
    }

    'pt2loop: for i in &nums {
        for j in &nums {
            let k = 2020 - i - j;
            if nums.contains(&k) {
                println!("Part 2: {}", i * j * k);
                break 'pt2loop;
            }
        }
    }

    Ok(())
}
