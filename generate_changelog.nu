#!/usr/bin/env nu

const repo_url = "https://github.com/evaneliasyoung/phex"

def pairwise [empty_head: bool]: list -> list<list> {
  let list = $in
  mut pairs = []

  mut pair_start = if ($empty_head) {""} else {$list.0}
  let itr = if ($empty_head) {$list} else {$list | skip}
  for $entry in $itr {
    $pairs = $pairs | append [[$pair_start, $entry]]
    $pair_start = $entry
  }

  return $pairs
}

def get_tags []: nothing -> list<string> {
  git show-ref --tags
    | str replace -arm '[0-9a-f]+\s+refs/tags/(?<tag>.*)' '$tag'
    | split row "\n"
}

def get_commits_between_tags [start: string, end: string] {
  let rev_range = if ($start | is-empty) { $end } else { $start + ".." + $end}
  git rev-list $rev_range --pretty="%h»|«%s»|«%as"
    | lines
    | every 2 --skip
    | split column "»|«" sha1 desc date
}

def "format commit" []: record -> string {
  let commit = $in
  let hash_title = $"[`($commit.sha1)`]"
  let hash_url = $"($repo_url)/commit/($commit.sha1)"
  $"- [($hash_title)\(($hash_url)\)] ($commit.desc)"
}

def main [] {
  let tag_pairs = get_tags | pairwise true | reverse
  mut line_buffer = ["# Changelog" ""]

  for $pair in $tag_pairs {
    let previous_tag = $pair | first
    let tag = $pair | last
    let commits = get_commits_between_tags $previous_tag $tag
      | where $it.desc !~ ":twisted_rightwards_arrows:"
      | where $it.desc !~ ":bookmark:"

    $line_buffer = $line_buffer
      | append $"## ($tag | str replace -r '^v' '')"
      | append ""
      | append $"> (($commits | first).date)"
      | append ""
      | append ($commits | each { $in | format commit })
      | append ""
  }

  $line_buffer | str join "\n" | save -f "CHANGELOG.md"
}
