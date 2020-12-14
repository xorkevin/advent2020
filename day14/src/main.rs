use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BErr = Box<dyn std::error::Error>;

fn main() -> Result<(), BErr> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut mem = HashMap::new();
    let mut mem2 = HashMap::new();

    let mut mask = (0, 0, 0);
    let mut mask_pos = Vec::new();

    for line in reader.lines() {
        let line = line?;
        let s = line.split(" = ").collect::<Vec<_>>();
        if s.len() != 2 {
            return Err("Invalid line".into());
        }
        if s[0] == "mask" {
            mask_pos.clear();
            let chars = s[1].chars().collect::<Vec<_>>();
            mask = process_mask(&chars, &mut mask_pos);
            continue;
        }
        let addr = s[0][4..s[0].len() - 1].parse::<usize>()?;
        let num = s[1].parse::<usize>()?;

        mem.insert(addr, (num & mask.0) | mask.1);

        let base = (addr & mask.2) | mask.1;
        for i in &mask_pos {
            mem2.insert(base | i, num);
        }
    }

    let sum = mem.values().sum::<usize>();
    println!("Part 1: {}", sum);

    let sum2 = mem2.values().sum::<usize>();
    println!("Part 2: {}", sum2);

    Ok(())
}

fn process_mask(b: &[char], pos_mask: &mut Vec<usize>) -> (usize, usize, usize) {
    let mut pos = Vec::new();
    let mut zeros = 0;
    let mut ones = 0;
    let mut fls = 0;
    let l = b.len();
    for (n, i) in b.iter().enumerate() {
        let p = l - n - 1;
        let k = 1 << p;
        match i {
            '0' => zeros |= k,
            '1' => ones |= k,
            'X' => {
                fls |= k;
                pos.push(k);
            }
            _ => (),
        }
    }
    mask_dfs(0, &pos, pos_mask);
    (!zeros, ones, !fls)
}

fn mask_dfs(base: usize, pos: &[usize], pos_mask: &mut Vec<usize>) {
    let (head, rest) = if let Some((head, rest)) = pos.split_first() {
        (head, rest)
    } else {
        pos_mask.push(base);
        return;
    };
    mask_dfs(base, rest, pos_mask);
    mask_dfs(base | head, rest, pos_mask);
}
