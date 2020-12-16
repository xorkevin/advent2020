use regex::Regex;
use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BErr = Box<dyn std::error::Error>;

struct Rule(usize, usize, usize, usize);

fn main() -> Result<(), BErr> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);
    let mut lines = reader.lines();

    let re = Regex::new(r"^([a-z ]+): ([0-9]+)-([0-9]+) or ([0-9]+)-([0-9]+)$")?;

    let (rules, rule_names) = {
        let mut rules = HashMap::new();
        let mut rule_names = HashSet::new();

        while let Some(line) = lines.next() {
            let line = line?;
            if line == "" {
                break;
            }
            let m = re.captures(&line).ok_or("Invalid rule line")?;
            let text = m[1].to_string();
            let text2 = m[1].to_string();
            let n1 = m[2].parse::<usize>()?;
            let n2 = m[3].parse::<usize>()?;
            let n3 = m[4].parse::<usize>()?;
            let n4 = m[5].parse::<usize>()?;
            rules.insert(text, Rule(n1, n2, n3, n4));
            rule_names.insert(text2);
        }
        (rules, rule_names)
    };

    let own_ticket = {
        let mut own_ticket = Vec::new();
        while let Some(line) = lines.next() {
            let line = line?;
            if line == "" {
                break;
            }
            if line == "your ticket:" {
                continue;
            }
            for i in line.split(",") {
                own_ticket.push(i.parse::<usize>()?);
            }
        }
        own_ticket
    };

    let (other_tickets, part1) = {
        let mut other_tickets = Vec::new();
        let mut part1 = 0;
        while let Some(line) = lines.next() {
            let line = line?;
            if line == "nearby tickets:" {
                continue;
            }
            let mut ticket = Vec::new();
            for i in line.split(",") {
                ticket.push(i.parse::<usize>()?);
            }
            if let Some(i) = is_invalid(&ticket, &rules) {
                part1 += i;
            } else {
                other_tickets.push(ticket);
            }
        }
        (other_tickets, part1)
    };
    println!("Part 1: {}", part1);

    let mut possible = (0..own_ticket.len())
        .map(|_| rule_names.clone())
        .collect::<Vec<_>>();
    let mut determined: HashMap<String, usize> = HashMap::new();
    while determined.len() < own_ticket.len() {
        loop {
            let mut changed = false;
            possible = possible
                .into_iter()
                .enumerate()
                .map(|(n, mut i)| {
                    if i.len() < 2 {
                        return i;
                    }
                    for j in &other_tickets {
                        let l = i.len();
                        i = i
                            .into_iter()
                            .filter(|k| in_range(j[n], &rules[k]))
                            .collect::<HashSet<_>>();
                        if i.len() != l {
                            changed = true;
                        }
                    }
                    i
                })
                .collect();
            if !changed {
                break;
            }
        }
        for (n, i) in possible.iter().enumerate() {
            if i.len() == 0 {
                return Err("Invalid constraints".into());
            }
            if i.len() == 1 {
                for k in i {
                    determined.insert(k.clone(), n);
                }
            }
        }
        for (k, &v) in &determined {
            possible = possible
                .into_iter()
                .enumerate()
                .map(|(n, mut i)| {
                    if n == v {
                        return i;
                    }
                    i.remove(k);
                    i
                })
                .collect();
        }
    }

    let part2 = determined
        .iter()
        .filter(|(k, _)| k.starts_with("departure"))
        .fold(1, |acc, (_, &v)| acc * own_ticket[v]);
    println!("Part 2: {}", part2);
    Ok(())
}

fn is_invalid(ticket: &[usize], rules: &HashMap<String, Rule>) -> Option<usize> {
    'outer: for &i in ticket {
        for v in rules.values() {
            if in_range(i, v) {
                continue 'outer;
            }
        }
        return Some(i);
    }
    None
}

fn in_range(i: usize, r: &Rule) -> bool {
    i >= r.0 && i <= r.3 && (i <= r.1 || i >= r.2)
}
