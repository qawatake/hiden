# hiden

A CLI tool to search and access personal memo/script directories (hiden directories) across repositories managed by ghq.

## Installation

```bash
go install github.com/qawatake/hiden@latest
```

## Usage

### Search files

```bash
# Select a file with incremental search
hiden ls

# Open in editor
vim $(hiden ls)
```

### Create date directory

```bash
# Create today's date directory in current repository
hiden mkdir
# => .hiden/2025-12-04

# Change to the created directory
cd $(hiden mkdir)
```

## Configuration

Config file: `~/.config/hiden/config.json`

```json
{
  "dirname": ".hiden"
}
```

| Field | Default | Description |
|-------|---------|-------------|
| `dirname` | `.hiden` | Name of the hiden directory |

## Directory structure example

```
~/src/github.com/
├── org1/repo1/
│   ├── .gitignore   # exclude .hiden/
│   └── .hiden/
│       ├── memo.md
│       └── scripts/test.sh
└── org2/repo2/
    ├── .gitignore   # exclude .hiden/
    └── .hiden/
        └── notes.txt
```

## License

MIT
