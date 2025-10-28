#!/usr/bin/env nu

def main [new_ver: string] {
  let cur_ver = git tag | lines | last | str replace -r "^v" ""

  print $"v($cur_ver) -> v($new_ver)"

  # Checkout main
  git switch main
  git fetch
  git pull

  # Merge
  git merge dev --no-ff --no-commit

  # Replace in files
  let path = "./internal/root.go"
  open -r $path
    | str replace $"\"($cur_ver)\"" $"\"($new_ver)\""
    | save -f $path

  # Commit
  let commit_message = $"chore: :bookmark: v($new_ver)"
  let tag = $"v($new_ver)"
  git add .
  git commit -m $commit_message
  git tag $tag -m $tag

  # Changelog
  nu generate_changelog.nu
  git add CHANGELOG.md
  git commit --amend -m $commit_message

  # Update tag
  git tag -d $tag
  git tag $tag -m $tag
}
