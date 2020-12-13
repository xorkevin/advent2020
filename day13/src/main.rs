const PUZZLEINPUT: i64 = 1006605;
const PUZZLEINPUT2: &str = "19,x,x,x,x,x,x,x,x,x,x,x,x,37,x,x,x,x,x,883,x,x,x,x,x,x,x,23,x,x,x,x,13,x,x,x,17,x,x,x,x,x,x,x,x,x,x,x,x,x,797,x,x,x,x,x,x,x,x,x,41,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,29";

type BErr = Box<dyn std::error::Error>;

fn main() -> Result<(), BErr> {
    let mut nums = Vec::new();
    let mut rem = Vec::new();

    for (n, i) in (0..).zip(PUZZLEINPUT2.split(",")) {
        if i == "x" {
            continue;
        }
        let num = i.parse::<i64>()?;
        nums.push(num);
        rem.push(num - n);
    }

    let mut i = PUZZLEINPUT;
    loop {
        if let Some(k) = can_take(i, &nums) {
            println!("Part 1: {}", k * (i - PUZZLEINPUT));
            break;
        }
        i += 1;
    }

    let x = crt(&nums, &rem)?;
    println!("Part 2: {}", x);

    Ok(())
}

fn crt(nums: &[i64], rem: &[i64]) -> Result<i64, BErr> {
    let p = nums.iter().fold(1, |acc, i| acc * i);
    let mut k = 0;
    for (&n, &i) in nums.iter().zip(rem) {
        let h = p / n;
        let t = mul_inv(h, n)?;
        k += i * t * h;
        k %= p;
    }
    Ok((k + p) % p)
}

/// Returns t, where `a*t = 1 (mod n)`
fn mul_inv(a: i64, n: i64) -> Result<i64, BErr> {
    let mut t = (0, 1);
    let mut r = (n, a);
    while r.1 != 0 {
        let q = r.0 / r.1;
        r = (r.1, r.0 - q * r.1);
        t = (t.1, t.0 - q * t.1);
    }
    if r.0 != 1 {
        return Err(format!("No mul inverse for {} mod {}", a, n).into());
    }
    if t.0 < 0 {
        return Ok(t.0 + n);
    }
    Ok(t.0)
}

fn can_take(t: i64, nums: &[i64]) -> Option<i64> {
    for &i in nums {
        if t % i == 0 {
            return Some(i);
        }
    }
    None
}
