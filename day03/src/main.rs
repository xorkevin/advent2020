use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let m = {
        let mut m = Vec::new();
        for line in reader.lines() {
            m.push(line?.chars().collect::<Vec<char>>());
        }
        m
    };

    println!("Part 1: {}", get_count(3, 1, &m));
    println!(
        "Part 2: {}",
        get_count(1, 1, &m)
            * get_count(3, 1, &m)
            * get_count(5, 1, &m)
            * get_count(7, 1, &m)
            * get_count(1, 2, &m)
    );

    Ok(())
}

fn get_count(x: usize, y: usize, m: &Vec<Vec<char>>) -> usize {
    let mut count = 0;
    let mut i = 0;
    let mut j = 0;
    let w = m[0].len();
    while i < m.len() {
        if m[i][j] == '#' {
            count += 1;
        }
        i += y;
        j = (j + x) % w;
    }
    count
}
