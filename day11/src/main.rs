use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BErr = Box<dyn std::error::Error>;

fn main() -> Result<(), BErr> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut grid = Vec::new();
    let mut tmp_grid = Vec::new();
    for line in reader.lines() {
        let k = line?.chars().collect::<Vec<_>>();
        tmp_grid.push(vec!['\0'; k.len()]);
        grid.push(k);
    }
    let mut grid2 = grid.clone();
    let mut tmp_grid2 = tmp_grid.clone();

    until_stable(false, &mut grid, &mut tmp_grid);
    println!("Part 1: {}", count_seats(&grid));

    until_stable(true, &mut grid2, &mut tmp_grid2);
    println!("Part 2: {}", count_seats(&grid2));

    Ok(())
}

fn count_seats(grid: &Vec<Vec<char>>) -> usize {
    let mut count = 0;
    for i in grid {
        for &j in i {
            if j == '#' {
                count += 1;
            }
        }
    }
    count
}

fn until_stable<'a>(
    mode: bool,
    mut grid: &'a mut Vec<Vec<char>>,
    mut tmp_grid: &'a mut Vec<Vec<char>>,
) {
    loop {
        if !next_grid(mode, grid, tmp_grid) {
            break;
        }
        let tmp = grid;
        grid = tmp_grid;
        tmp_grid = tmp;
    }
}

fn next_grid(mode: bool, grid: &mut Vec<Vec<char>>, next: &mut Vec<Vec<char>>) -> bool {
    let mut change = false;
    for i in 0..grid.len() {
        for j in 0..grid[i].len() {
            let c = grid[i][j];
            let k = if mode {
                next_seat2(i, j, grid)
            } else {
                next_seat(i, j, grid)
            };
            if k != c {
                change = true;
            }
            next[i][j] = k;
        }
    }
    change
}

fn next_seat2(i: usize, j: usize, grid: &Vec<Vec<char>>) -> char {
    let c = grid[i][j];
    let s1 = seat_val2(i, j, 0, 1, 0, 0, grid);
    let s2 = seat_val2(i, j, 0, 1, 1, 0, grid);
    let s3 = seat_val2(i, j, 0, 0, 1, 0, grid);
    let s4 = seat_val2(i, j, 1, 0, 1, 0, grid);
    let s5 = seat_val2(i, j, 1, 0, 0, 0, grid);
    let s6 = seat_val2(i, j, 1, 0, 0, 1, grid);
    let s7 = seat_val2(i, j, 0, 0, 0, 1, grid);
    let s8 = seat_val2(i, j, 0, 1, 0, 1, grid);
    let k = s1 + s2 + s3 + s4 + s5 + s6 + s7 + s8;
    if c == 'L' && k == 0 {
        return '#';
    }
    if c == '#' && k >= 5 {
        return 'L';
    }
    return c;
}

fn seat_val2(
    mut i: usize,
    mut j: usize,
    ui: usize,
    di: usize,
    uj: usize,
    dj: usize,
    grid: &Vec<Vec<char>>,
) -> usize {
    i = i + ui - di;
    j = j + uj - dj;
    while i < grid.len() && j < grid[0].len() {
        if grid[i][j] == 'L' {
            return 0;
        }
        if grid[i][j] == '#' {
            return 1;
        }
        i = i + ui - di;
        j = j + uj - dj;
    }
    0
}

fn next_seat(i: usize, j: usize, grid: &Vec<Vec<char>>) -> char {
    let c = grid[i][j];
    let s1 = seat_val(i - 1, j, grid);
    let s2 = seat_val(i - 1, j + 1, grid);
    let s3 = seat_val(i, j + 1, grid);
    let s4 = seat_val(i + 1, j + 1, grid);
    let s5 = seat_val(i + 1, j, grid);
    let s6 = seat_val(i + 1, j - 1, grid);
    let s7 = seat_val(i, j - 1, grid);
    let s8 = seat_val(i - 1, j - 1, grid);
    let k = s1 + s2 + s3 + s4 + s5 + s6 + s7 + s8;
    if c == 'L' && k == 0 {
        return '#';
    }
    if c == '#' && k >= 4 {
        return 'L';
    }
    return c;
}

fn seat_val(i: usize, j: usize, grid: &Vec<Vec<char>>) -> usize {
    if i >= grid.len() {
        return 0;
    }
    if j >= grid[0].len() {
        return 0;
    }
    if grid[i][j] == '#' {
        1
    } else {
        0
    }
}
