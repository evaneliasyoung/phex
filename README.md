# phex

**Phex — Phaser Texture Manager.**
A Go CLI that packs **Phaser 3** sprite atlases from individual images _and_ unpacks existing atlases back into frames.

---

## Features

- **Pack** a directory of sprites into a single atlas image + Phaser 3 JSON.
- **Unpack** a Phaser 3 atlas (`.json`) into individual frame images.
- ✂️ **Trim** transparent borders when packing (optional).
- 🧰 **Dedupe** identical frames when packing (optional).
- 📐 **Efficient bin-packing** to minimize atlas size.
- ⚙️ **Fast + concurrent** for big batches (hundreds to thousands of images).
- 🎯 Focused on **Phaser 3 (Matter)** pipelines.

---

## Installation

Build locally:

```bash
git clone https://github.com/evaneliasyoung/phex.git
cd phex
go build -o phex ./cmd/phex
```

Or install directly:

```bash
go install github.com/evaneliasyoung/phex@latest
```

> On first use, ensure `$GOPATH/bin` (or your Go bin dir) is on your `PATH`.

---

## Usage

### Pack

Create an atlas (image + JSON) from a folder of sprites:

```bash
phex pack <input-dir> -o <out-dir>
```

### Unpack

Extract frames from a Phaser atlas JSON:

```bash
phex unpack <atlas.json> -o <out-dir>
```

---

## Flags

Common flags

| Flag                 | Description      | Default                                   |
| -------------------- | ---------------- | ----------------------------------------- |
| `-o, --output <dir>` | Output directory | for `pack`: `./output`, for `unpack`: `.` |

Unpack flags

| Flag                  | Description                  | Default                  |
| --------------------- | ---------------------------- | ------------------------ |
| `-w, --workers <num>` | Number of concurrent workers | 2×Thread Count, up to 32 |
| `--no-progress`       | Disable progress bars        | disabled if non-TTY      |

Pack-specific:

| Flag                   | Description                                                                         | Default |
| ---------------------- | ----------------------------------------------------------------------------------- | ------- |
| `-m, --maxsize <size>` | Maximum texture sheet size (as a square)                                            | `2048`  |
| `-p, --padding <px>`   | Padding pixels between sprites in the sheet                                         | `0`     |
| `-n, --name <base>`    | The name of the sprite sheets and the atlas file (e.g., `atlas.json`, `atlas.webp`) | `atlas` |

> If a flag isn’t implemented yet in your codebase, it will land as part of the Roadmap below.

---

## Examples

```bash
# Unpack textures to ./sprites_output/*.png
phex unpack assets/sprites.json -o sprites_output

# Use 8 workers regardless of CPU count
phex unpack assets/sprites.json --workers 8

# Pack textures to ./atlas/atlas.json
phex pack sprites -o .

# Pack textures to ./items/atlas.json
phex pack sprites -o ./items

# Pack textures to ./rooms/coffee.json with a larger maximum sheet size
phex pack sprites -o ./rooms -n coffee -m 4096
```

---

## Compatibility

- **Input (pack):** `.png` (subfolders supported).
- **Output (pack):** atlas image (`.webp`) + Phaser 3 JSON.
- **Input (unpack):** Phaser-compatible atlas JSON (Phex/TexturePacker style).
- **Output (unpack):** individual `.png` frames.

---

## Dependencies

- [`gen2brain/webp`](https://github.com/gen2brain/webp) — WebP encoder
- [`spf13/cobra`](https://github.com/spf13/cobra) — CLI framework
- [`vbauerster/mpb`](https://github.com/vbauerster/mpb) — Progress bars
- [`golang.org/x/image/webp`](https://pkg.go.dev/golang.org/x/image/webp) — WEBP decoder
- [`golang.org/x/term`](https://pkg.go.dev/golang.org/x/term) — Determine if TTY
