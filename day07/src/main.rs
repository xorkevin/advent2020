use regex::Regex;
use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let re = Regex::new(r"^([a-z ]+) bags contain")?;
    let re2 = Regex::new(r"([0-9]+) ([a-z ]+) bags?[,.]")?;
    let re3 = Regex::new(r"no other bags.$")?;

    let graph = {
        let mut graph = HashMap::new();
        for line in reader.lines() {
            let line = line?;
            let m = re.captures(&line).ok_or("Invalid line format")?;
            let m1 = m.get(1).ok_or("Invalid line format")?.as_str().to_owned();
            let m3 = re3.is_match(&line);
            let mut subgraph = HashMap::new();
            if m3 {
                graph.insert(m1, subgraph);
                continue;
            }
            if !re2.is_match(&line) {
                return Err("Invalid line format".into());
            }
            for i in re2.captures_iter(&line) {
                let num = i
                    .get(1)
                    .ok_or("Invalid line format")?
                    .as_str()
                    .parse::<usize>()?;
                let k = i.get(2).ok_or("Invalid line format")?.as_str().to_owned();
                subgraph.insert(k, num);
            }
            graph.insert(m1, subgraph);
        }
        graph
    };

    let mut count = 0;
    let mut cache = HashMap::new();
    for i in graph.keys() {
        if has_path_to(i, "shiny gold", &graph, &mut cache) {
            count += 1;
        }
    }
    println!("Part 1: {}", count);

    let mut cache2 = HashMap::new();
    let count2 = count_children("shiny gold", &graph, &mut cache2);
    println!("Part 2: {}", count2);

    Ok(())
}

fn has_path_to<'a>(
    a: &'a str,
    b: &str,
    graph: &'a HashMap<String, HashMap<String, usize>>,
    cache: &mut HashMap<&'a str, bool>,
) -> bool {
    if let Some(&v) = cache.get(a) {
        return v;
    }
    let edges = match graph.get(a) {
        Some(edges) => edges,
        None => return false,
    };
    for i in edges.keys() {
        if b == i {
            cache.insert(a, true);
            return true;
        }
    }
    for i in edges.keys() {
        if has_path_to(i, b, graph, cache) {
            cache.insert(a, true);
            return true;
        }
    }
    cache.insert(a, false);
    return false;
}

fn count_children<'a>(
    a: &'a str,
    graph: &'a HashMap<String, HashMap<String, usize>>,
    cache: &mut HashMap<&'a str, usize>,
) -> usize {
    if let Some(&v) = cache.get(a) {
        return v;
    }
    let mut count = 1;
    let edges = match graph.get(a) {
        Some(edges) => edges,
        None => return 1,
    };
    for (k, &v) in edges {
        count += v * count_children(k, graph, cache);
    }
    cache.insert(a, count);
    return count;
}
