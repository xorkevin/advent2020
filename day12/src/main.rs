use regex::Regex;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BErr = Box<dyn std::error::Error>;

#[derive(Clone, Copy, PartialEq, Eq)]
enum Action {
    N,
    S,
    E,
    W,
    L,
    R,
    F,
}

struct Instr(Action, i32);

#[derive(Clone, Copy, PartialEq, Eq)]
enum Dir {
    N,
    E,
    S,
    W,
}

impl Dir {
    fn left(&self) -> Dir {
        match self {
            Dir::N => Dir::W,
            Dir::E => Dir::N,
            Dir::S => Dir::E,
            Dir::W => Dir::S,
        }
    }

    fn right(&self) -> Dir {
        match self {
            Dir::N => Dir::E,
            Dir::E => Dir::S,
            Dir::S => Dir::W,
            Dir::W => Dir::N,
        }
    }
}

struct Ship {
    dir: Dir,
    x: i32,
    y: i32,
    wx: i32,
    wy: i32,
}

impl Ship {
    fn new() -> Ship {
        Ship {
            dir: Dir::E,
            x: 0,
            y: 0,
            wx: 10,
            wy: -1,
        }
    }

    fn step(&mut self, a: &Instr) {
        let Instr(act, v) = a;
        match act {
            Action::N => self.y -= v,
            Action::S => self.y += v,
            Action::E => self.x += v,
            Action::W => self.x -= v,
            Action::L => {
                for _ in 0..(v / 90) {
                    self.dir = self.dir.left();
                }
            }
            Action::R => {
                for _ in 0..(v / 90) {
                    self.dir = self.dir.right();
                }
            }
            Action::F => match self.dir {
                Dir::N => self.y -= v,
                Dir::E => self.x += v,
                Dir::S => self.y += v,
                Dir::W => self.x -= v,
            },
        }
    }

    fn turn_waypoint_left(&mut self) {
        let a = self.wx;
        let b = self.wy;
        self.wx = b;
        self.wy = -a;
    }

    fn turn_waypoint_right(&mut self) {
        let a = self.wx;
        let b = self.wy;
        self.wx = -b;
        self.wy = a;
    }

    fn step_waypoint(&mut self, a: &Instr) {
        let Instr(act, v) = a;
        match act {
            Action::N => self.wy -= v,
            Action::S => self.wy += v,
            Action::E => self.wx += v,
            Action::W => self.wx -= v,
            Action::L => {
                for _ in 0..(v / 90) {
                    self.turn_waypoint_left();
                }
            }
            Action::R => {
                for _ in 0..(v / 90) {
                    self.turn_waypoint_right();
                }
            }
            Action::F => {
                self.x += self.wx * v;
                self.y += self.wy * v;
            }
        }
    }
}

fn main() -> Result<(), BErr> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let re = Regex::new(r"^([A-Z])([0-9]+)$")?;

    let mut s = Ship::new();
    let mut s2 = Ship::new();

    for line in reader.lines() {
        let line = line?;
        let m = re.captures(&line).ok_or("Invalid line format")?;
        let act = match m.get(1).ok_or("Invalid line format")?.as_str() {
            "N" => Action::N,
            "S" => Action::S,
            "E" => Action::E,
            "W" => Action::W,
            "L" => Action::L,
            "R" => Action::R,
            "F" => Action::F,
            _ => return Err("Invalid action".into()),
        };
        let num = m
            .get(2)
            .ok_or("Invalid line format")?
            .as_str()
            .parse::<i32>()?;
        let a = Instr(act, num);
        s.step(&a);
        s2.step_waypoint(&a);
    }

    println!("Part 1: {}", s.x.abs() + s.y.abs());
    println!("Part 2: {}", s2.x.abs() + s2.y.abs());

    Ok(())
}
