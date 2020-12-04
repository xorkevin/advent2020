use regex::Regex;
use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

#[derive(Debug, Default)]
struct Passport {
    byr: i32,
    iyr: i32,
    eyr: i32,
    hgt: String,
    hcl: String,
    ecl: String,
    pid: String,
}

struct Re {
    height: Regex,
    color: Regex,
    pid: Regex,
    eye_colors: HashSet<String>,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let re = Re {
        height: Regex::new(r"^([0-9]+)(cm|in)$")?,
        color: Regex::new(r"^#[0-9a-f]{6}$")?,
        pid: Regex::new(r"^[0-9]{9}$")?,
        eye_colors: vec![
            "amb".to_string(),
            "blu".to_string(),
            "brn".to_string(),
            "gry".to_string(),
            "grn".to_string(),
            "hzl".to_string(),
            "oth".to_string(),
        ]
        .into_iter()
        .collect::<HashSet<_>>(),
    };

    let passports = {
        let mut passports = Vec::new();
        let mut p = Passport::default();
        for line in reader.lines() {
            let line = line?;
            if line == "" {
                passports.push(p);
                p = Passport::default();
            }
            for kv in line.split_whitespace() {
                match kv.split(':').collect::<Vec<_>>()[..] {
                    [k, v] => match k {
                        "byr" => p.byr = v.parse::<i32>()?,
                        "iyr" => p.iyr = v.parse::<i32>()?,
                        "eyr" => p.eyr = v.parse::<i32>()?,
                        "hgt" => p.hgt = String::from(v),
                        "hcl" => p.hcl = String::from(v),
                        "ecl" => p.ecl = String::from(v),
                        "pid" => p.pid = String::from(v),
                        _ => {}
                    },
                    _ => return Err("invalid kv pair".into()),
                }
            }
        }
        passports.push(p);
        passports
    };

    let count = passports
        .iter()
        .map(|i| if pass_is_valid(i) { 1 } else { 0 })
        .sum::<usize>();

    let count2 = passports
        .iter()
        .map(|i| {
            if pass_is_valid(i) && pass_is_valid2(i, &re) {
                1
            } else {
                0
            }
        })
        .sum::<usize>();

    println!("Part 1: {}", count);
    println!("Part 2: {}", count2);

    Ok(())
}

fn pass_is_valid(p: &Passport) -> bool {
    p.byr != 0
        && p.iyr != 0
        && p.eyr != 0
        && p.hgt != ""
        && p.hcl != ""
        && p.ecl != ""
        && p.pid != ""
}

fn pass_is_valid2(p: &Passport, re: &Re) -> bool {
    if !in_range(p.byr, 1920, 2002) {
        return false;
    }
    if !in_range(p.iyr, 2010, 2020) {
        return false;
    }
    if !in_range(p.eyr, 2020, 2030) {
        return false;
    }
    let hm = if let Some(hm) = re.height.captures(&p.hgt) {
        hm
    } else {
        return false;
    };
    let h = if let Some(h) = hm.get(1) {
        if let Ok(h) = h.as_str().parse::<i32>() {
            h
        } else {
            return false;
        }
    } else {
        return false;
    };
    let m = if let Some(m) = hm.get(2) {
        m.as_str()
    } else {
        return false;
    };
    match m {
        "cm" => {
            if !in_range(h, 150, 193) {
                return false;
            }
        }
        "in" => {
            if !in_range(h, 59, 76) {
                return false;
            }
        }
        _ => return false,
    }
    if !re.color.is_match(&p.hcl) {
        return false;
    }
    if !re.eye_colors.contains(&p.ecl) {
        return false;
    }
    if !re.pid.is_match(&p.pid) {
        return false;
    }
    return true;
}

fn in_range(i: i32, a: i32, b: i32) -> bool {
    i >= a && i <= b
}
