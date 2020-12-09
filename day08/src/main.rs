use std::collections::HashMap;
use std::convert::TryFrom;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

#[derive(Clone, Copy, PartialEq, Eq)]
enum Code {
    Noop,
    Acc,
    Jmp,
}

#[derive(Clone, Copy)]
struct Instr(Code, i32);

impl Instr {
    fn parse(line: &str) -> Result<Instr, Box<dyn std::error::Error>> {
        let fields = line.split_whitespace().collect::<Vec<_>>();
        let code = match fields.get(0).ok_or("Invalid line")? {
            &"nop" => Code::Noop,
            &"acc" => Code::Acc,
            &"jmp" => Code::Jmp,
            bogus => return Err(format!("Invalid code: {}", bogus).into()),
        };
        let arg = fields.get(1).ok_or("Invalid line")?.parse::<i32>()?;
        Ok(Instr(code, arg))
    }
}

struct Machine<'a> {
    instrs: &'a Vec<Instr>,
    ip: i32,
    acc: i32,
}

impl<'a> Machine<'a> {
    fn new(instrs: &'a Vec<Instr>) -> Machine<'a> {
        Machine {
            instrs,
            ip: 0,
            acc: 0,
        }
    }

    fn step(&mut self) -> Result<bool, Box<dyn std::error::Error>> {
        let ip = usize::try_from(self.ip)?;
        if ip == self.instrs.len() {
            return Ok(true);
        }
        let Instr(code, arg) = self
            .instrs
            .get(ip)
            .ok_or(format!("Ip out of bounds: {}", self.ip))?;
        match code {
            Code::Noop => self.ip += 1,
            Code::Acc => {
                self.acc += arg;
                self.ip += 1;
            }
            Code::Jmp => self.ip += arg,
        }
        Ok(false)
    }

    fn run(&mut self, loop_limit: usize) -> Result<(), Box<dyn std::error::Error>> {
        let mut run_set = HashMap::new();
        loop {
            let &mut e = run_set.entry(self.ip).and_modify(|e| *e += 1).or_insert(1);
            if e > loop_limit {
                return Err(format!("Looped over {} at {}", loop_limit, self.ip).into());
            }
            if self.step()? {
                return Ok(());
            }
        }
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut instrs = Vec::new();
    for line in reader.lines() {
        let line = line?;
        instrs.push(Instr::parse(&line)?);
    }

    let mut m = Machine::new(&instrs);
    let _ = m.run(1);
    println!("Part 1: {}", m.acc);

    for i in 0..instrs.len() {
        let Instr(code, arg) = instrs[i];
        if code == Code::Acc {
            continue;
        }
        let swap_to_jmp = code == Code::Noop;
        if swap_to_jmp {
            instrs[i] = Instr(Code::Jmp, arg);
        } else {
            instrs[i] = Instr(Code::Noop, arg);
        }
        let mut m = Machine::new(&instrs);
        if let Ok(_) = m.run(1) {
            println!("Part 2: {}", m.acc);
            break;
        }
        if swap_to_jmp {
            instrs[i] = Instr(Code::Noop, arg);
        } else {
            instrs[i] = Instr(Code::Jmp, arg);
        }
    }

    Ok(())
}
