use std::collections::VecDeque;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BErr = Box<dyn std::error::Error>;

struct TargetReader {
    buf: Box<dyn Iterator<Item = i32>>,
    iter: Box<dyn Iterator<Item = std::io::Result<String>>>,
}

impl TargetReader {
    fn new(iter: Box<dyn Iterator<Item = std::io::Result<String>>>) -> TargetReader {
        TargetReader {
            buf: Box::new(Vec::new().into_iter()),
            iter,
        }
    }

    fn read_line(&mut self) -> Result<i32, BErr> {
        let line = self.iter.next().ok_or("EOF")??;
        let num = line.parse::<i32>()?;
        Ok(num)
    }

    fn scan(&mut self) -> Result<i32, BErr> {
        let mut b = Vec::new();
        for _ in 0..25 {
            let k = self.read_line()?;
            b.push(k);
        }

        let mut idx = 0;
        loop {
            let k = self.read_line()?;
            b.push(k);
            if !has_parts(k, &b[idx..idx + 25]) {
                break;
            }
            idx += 1;
        }
        let last = b[b.len() - 1];
        self.buf = Box::new(b.into_iter());
        Ok(last)
    }

    fn next_int(&mut self) -> Result<i32, BErr> {
        if let Some(k) = self.buf.next() {
            return Ok(k);
        }
        self.read_line()
    }
}

fn has_parts(a: i32, b: &[i32]) -> bool {
    for i in 0..b.len() {
        for j in i + 1..b.len() {
            if b[i] + b[j] == a {
                return true;
            }
        }
    }
    false
}

fn main() -> Result<(), BErr> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut r = TargetReader::new(Box::new(reader.lines()));

    let target = r.scan()?;
    println!("Part 1: {}", target);

    let mut sum = 0;
    let mut buf = VecDeque::new();
    while sum != target {
        if sum < target {
            let k = r.next_int()?;
            sum += k;
            buf.push_back(k);
        } else {
            if buf.len() == 0 {
                break;
            }
            let first = buf[0];
            sum -= first;
            buf.pop_front();
        }
    }
    if sum == target {
        buf.make_contiguous();
        let (slice, _) = buf.as_slices();
        let (min, max) = min_max(slice)?;
        println!("Part 2: {}", min + max);
    }
    Ok(())
}

fn min_max(b: &[i32]) -> Result<(i32, i32), BErr> {
    if b.len() == 0 {
        return Err("No elements".into());
    }
    let mut min = *b.first().ok_or("No elements")?;
    let mut max = min;
    for &i in b {
        if i < min {
            min = i;
        }
        if i > max {
            max = i;
        }
    }
    Ok((min, max))
}
