use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BErr = Box<dyn std::error::Error>;

fn main() -> Result<(), BErr> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let nums = {
        let mut nums = Vec::new();
        nums.push(0);
        for line in reader.lines() {
            nums.push(line?.parse::<u64>()?);
        }
        nums.sort();
        nums.push(nums.last().ok_or("Invalid vec")? + 3);
        nums
    };

    let mut start = *nums.first().ok_or("Invalid vec")?;
    let mut diff1 = 0;
    let mut diff3 = 0;
    for &i in &nums {
        let k = i - start;
        if k == 1 {
            diff1 += 1;
        } else if k == 3 {
            diff3 += 1;
        }
        start = i;
    }
    println!("Part 1: {}", diff1 * diff3);
    println!("Part 2: {}", count_paths(&nums));
    Ok(())
}

fn count_paths(nums: &[u64]) -> u64 {
    let mut cache = vec![0; nums.len()];
    for i in 0..nums.len() {
        if i == 0 {
            cache[0] = 1;
            continue;
        }
        let k = nums[i];
        let a = get_val(k, i - 1, nums, &cache);
        let b = get_val(k, i - 2, nums, &cache);
        let c = get_val(k, i - 3, nums, &cache);
        cache[i] = a + b + c;
    }
    cache[cache.len() - 1]
}

fn get_val(start: u64, i: usize, nums: &[u64], cache: &[u64]) -> u64 {
    if i >= cache.len() {
        return 0;
    }
    if start - nums[i] > 3 {
        return 0;
    }
    cache[i]
}
