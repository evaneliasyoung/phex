# Changelog

## 0.2.1

> 2026-02-28

- [[`3b89c5d`](https://github.com/evaneliasyoung/phex/commit/3b89c5d)] fix: :bug: run go mod tidy, fix `go.sum`
- [[`6920a7d`](https://github.com/evaneliasyoung/phex/commit/6920a7d)] docs: :memo: note that most image formats are supported for packing
- [[`d5ad149`](https://github.com/evaneliasyoung/phex/commit/d5ad149)] build: :arrow_up: upgrade `github.com/spf13/cobra`
- [[`e5e549e`](https://github.com/evaneliasyoung/phex/commit/e5e549e)] build: :arrow_up: upgrade `github.com/vbauerster/mpb/v8`
- [[`f49bbf2`](https://github.com/evaneliasyoung/phex/commit/f49bbf2)] build: :arrow_up: upgrade `golang.org/x/term`
- [[`bd49f34`](https://github.com/evaneliasyoung/phex/commit/bd49f34)] fix: :pencil2: correct punctuation in the command description
- [[`56f9a55`](https://github.com/evaneliasyoung/phex/commit/56f9a55)] refactor: :recycle: use more semantic UTF-8 folding when checking atlas extension
- [[`9ea2f7a`](https://github.com/evaneliasyoung/phex/commit/9ea2f7a)] fix: :bug: derive pack name from the pack file rather than the parent directory
- [[`6058da8`](https://github.com/evaneliasyoung/phex/commit/6058da8)] fix: :bug: increment total textures unpacked when using the no-progress reporter
- [[`af9cd7f`](https://github.com/evaneliasyoung/phex/commit/af9cd7f)] fix: :adhesive_bandage: update literals for unpack when the texture sheet fails to decode
- [[`b05403c`](https://github.com/evaneliasyoung/phex/commit/b05403c)] fix: :bug: convert the extension to lowercase when checking unpack atlas validity
- [[`f1bad24`](https://github.com/evaneliasyoung/phex/commit/f1bad24)] fix: :adhesive_bandage: update literals for pack when no source is provided
- [[`b90f140`](https://github.com/evaneliasyoung/phex/commit/b90f140)] fix: :bug: close texture sheet file when unpacking
- [[`523da4d`](https://github.com/evaneliasyoung/phex/commit/523da4d)] fix: :bug: add check for non-positive worker threads on unpack

## 0.2.0

> 2026-02-16

- [[`2bd0e99`](https://github.com/evaneliasyoung/phex/commit/2bd0e99)] test: :white_check_mark: add `phaser` tests
- [[`37eccf1`](https://github.com/evaneliasyoung/phex/commit/37eccf1)] perf: :zap: improve packing algorithm w/ better heuristics
- [[`4b67e7d`](https://github.com/evaneliasyoung/phex/commit/4b67e7d)] fix: :bug: don't record empty sprite sheets
- [[`0498667`](https://github.com/evaneliasyoung/phex/commit/0498667)] build: :wrench: add Taskfile
- [[`93ee3b7`](https://github.com/evaneliasyoung/phex/commit/93ee3b7)] chore: :memo: add issue templates
- [[`f0eb6fe`](https://github.com/evaneliasyoung/phex/commit/f0eb6fe)] refactor: :recycle: add levels to `phex -v` version command
- [[`b6fb3ae`](https://github.com/evaneliasyoung/phex/commit/b6fb3ae)] ci: :construction_worker: add `ci` action to run tests, move `pipeline` to `release`
- [[`7d73cbe`](https://github.com/evaneliasyoung/phex/commit/7d73cbe)] test: :white_check_mark: add `maxrects` tests
- [[`e15328d`](https://github.com/evaneliasyoung/phex/commit/e15328d)] test: :white_check_mark: add `image` tests
- [[`03774c2`](https://github.com/evaneliasyoung/phex/commit/03774c2)] refactor: :coffin: remove latent references to sprite rotation
- [[`e7775e1`](https://github.com/evaneliasyoung/phex/commit/e7775e1)] test: :white_check_mark: add `atlas` tests
- [[`2402b93`](https://github.com/evaneliasyoung/phex/commit/2402b93)] fix: :bug: don't split a sprite to multiple sheets if it's too big
- [[`2de08ba`](https://github.com/evaneliasyoung/phex/commit/2de08ba)] docs: :memo: remove note on flags from README
- [[`8f236f5`](https://github.com/evaneliasyoung/phex/commit/8f236f5)] chore: :hammer: add nushell script to publish a new version

## 0.1.2

> 2025-10-28

- [[`70c5165`](https://github.com/evaneliasyoung/phex/commit/70c5165)] chore: :hammer: add nushell script to generate the `CHANGELOG.md`
- [[`5b23b63`](https://github.com/evaneliasyoung/phex/commit/5b23b63)] fix: :rotating_light: fix two linter warnings
- [[`333e496`](https://github.com/evaneliasyoung/phex/commit/333e496)] chore: :wrench: add VS Code recommended extensions
- [[`0abd4d4`](https://github.com/evaneliasyoung/phex/commit/0abd4d4)] ci: :construction_worker: add `golangci-lint` action
- [[`c30ab4e`](https://github.com/evaneliasyoung/phex/commit/c30ab4e)] ci: :construction_worker: update CI to `setup-go@v6` and use go `stable`
- [[`144049a`](https://github.com/evaneliasyoung/phex/commit/144049a)] docs: :memo: change dependency list to table
- [[`2ae4b8d`](https://github.com/evaneliasyoung/phex/commit/2ae4b8d)] chore: :page_facing_up: add LGPL license
- [[`34d9d06`](https://github.com/evaneliasyoung/phex/commit/34d9d06)] feat: :sparkles: add packing support for GIF, JPG, BMP, TIFF, and WebP images
- [[`16b9a70`](https://github.com/evaneliasyoung/phex/commit/16b9a70)] build: :heavy_plus_sign: add `github.com/h2non/filetype`
- [[`b328684`](https://github.com/evaneliasyoung/phex/commit/b328684)] fix: :bug: write errors if they occur
- [[`a51b706`](https://github.com/evaneliasyoung/phex/commit/a51b706)] fix: :bug: resolve #1
- [[`5b62ae9`](https://github.com/evaneliasyoung/phex/commit/5b62ae9)] refactor: :pencil2: remove optional notes from trimming and deduping

## 0.1.1

> 2025-10-24

- [[`67fbaba`](https://github.com/evaneliasyoung/phex/commit/67fbaba)] chore: :memo: update examples in `README.md`
- [[`3845855`](https://github.com/evaneliasyoung/phex/commit/3845855)] fix: :bug: prevent crash when leaving `--output` empty for `phex pack`
- [[`88aabeb`](https://github.com/evaneliasyoung/phex/commit/88aabeb)] refactor: :fire: remove rotation from `phex pack` (it's broken)
- [[`80f6123`](https://github.com/evaneliasyoung/phex/commit/80f6123)] refactor: :recycle: directly embed version
- [[`6a3c778`](https://github.com/evaneliasyoung/phex/commit/6a3c778)] chore: :memo: add `README.md`

## 0.1.0

> 2025-10-24

- [[`1bdfb68`](https://github.com/evaneliasyoung/phex/commit/1bdfb68)] chore: :tada: initial commit
