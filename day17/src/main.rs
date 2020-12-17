use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BErr = Box<dyn std::error::Error>;

#[derive(Clone, Copy, PartialEq, Eq, Hash)]
struct Point(i32, i32, i32, i32);

fn main() -> Result<(), BErr> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut active = HashSet::new();
    for (i, line) in reader.lines().enumerate() {
        let line = line?;
        for (j, c) in line.chars().enumerate() {
            if c == '#' {
                active.insert(Point(j as i32, i as i32, 0, 0));
            }
        }
    }
    let mut active2 = active.clone();

    for _ in 0..6 {
        active = step(&active, false);
    }
    println!("Part 1: {}", active.len());
    for _ in 0..6 {
        active2 = step(&active2, true);
    }
    println!("Part 2: {}", active2.len());
    Ok(())
}

fn step(active: &HashSet<Point>, dim4: bool) -> HashSet<Point> {
    let mut next = HashSet::new();
    let mut nbc = HashMap::new();
    for &p in active {
        let nb = if dim4 { neighbors_w(p) } else { neighbors(p) };
        let c = count_neighbors(&nb, active);
        if c == 2 || c == 3 {
            next.insert(p);
        }
        for n in nb {
            let v = nbc.entry(n).or_insert(0);
            *v += 1;
        }
    }
    for (p, c) in nbc {
        if c == 3 {
            next.insert(p);
        }
    }
    next
}

fn count_neighbors(points: &[Point], active: &HashSet<Point>) -> usize {
    let mut count = 0;
    for p in points {
        if active.contains(p) {
            count += 1;
        }
    }
    count
}

fn neighbors(p: Point) -> Vec<Point> {
    let mut points = Vec::with_capacity(26);
    for i in p.0 - 1..=p.0 + 1 {
        for j in p.1 - 1..=p.1 + 1 {
            for k in p.2 - 1..=p.2 + 1 {
                let x = Point(i, j, k, 0);
                if x != p {
                    points.push(x);
                }
            }
        }
    }
    points
}

fn neighbors_w(p: Point) -> Vec<Point> {
    let mut points = Vec::with_capacity(26);
    for i in p.0 - 1..=p.0 + 1 {
        for j in p.1 - 1..=p.1 + 1 {
            for k in p.2 - 1..=p.2 + 1 {
                for l in p.3 - 1..=p.3 + 1 {
                    let x = Point(i, j, k, l);
                    if x != p {
                        points.push(x);
                    }
                }
            }
        }
    }
    points
}
