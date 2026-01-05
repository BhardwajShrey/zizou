# Documentation Organization

This document explains the structure and organization of Zizou's documentation.

## Directory Structure

### Root Directory (`/`)

Only essential user-facing files:

- **README.md** - Project overview, quick start, basic usage
- **RUN_ZIZOU.md** - Complete running instructions with all options

### Documentation Directory (`/docs`)

Detailed technical documentation:

- **README.md** - Documentation index and navigation
- **DIFF_GUIDE.md** - Complete guide to understanding git diff format
- **DIFF_FORMAT_EXPLAINED.md** - Line-by-line explanation of diff format
- **QUICK_REFERENCE.md** - Quick reference card for diff format
- **DIFF_PARSING.md** - Parser architecture and usage guide

### Examples Directory (`/examples`)

- **README.md** - Overview of available examples
- **line_by_line_demo.go** - Interactive diff explanation
- **simple_usage.go** - Common use cases
- **parse_demo.go** - Comprehensive parser demo
- **test_parser.go** - Test utility

## Navigation Path

### For Users

1. Start with [README.md](../README.md) in root
2. Read [RUN_ZIZOU.md](../RUN_ZIZOU.md) for usage details
3. Explore [examples/](../examples/) for code samples

### For Developers

1. Read [docs/README.md](README.md) for documentation index
2. Study [DIFF_PARSING.md](DIFF_PARSING.md) for parser architecture
3. Review [examples/](../examples/) for implementation patterns

### For Learning Diff Format

1. [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Quick cheat sheet
2. [DIFF_GUIDE.md](DIFF_GUIDE.md) - Complete guide
3. [DIFF_FORMAT_EXPLAINED.md](DIFF_FORMAT_EXPLAINED.md) - Deep dive
4. Run `go run examples/line_by_line_demo.go` for interactive learning

## File Purposes

### README.md (root)
- **Audience:** All users
- **Purpose:** Quick start, installation, basic usage
- **Depth:** Minimal, focuses on getting started

### RUN_ZIZOU.md (root)
- **Audience:** Users who want to run Zizou
- **Purpose:** Complete usage guide with all flags and options
- **Depth:** Comprehensive usage documentation

### docs/DIFF_GUIDE.md
- **Audience:** Users learning diff format
- **Purpose:** Line-by-line explanation of git diff
- **Depth:** Educational, detailed examples

### docs/DIFF_FORMAT_EXPLAINED.md
- **Audience:** Advanced users, contributors
- **Purpose:** Complete technical breakdown of diff format
- **Depth:** Very detailed, covers edge cases

### docs/QUICK_REFERENCE.md
- **Audience:** Users who need quick lookups
- **Purpose:** Cheat sheet for diff format
- **Depth:** Quick reference, minimal explanation

### docs/DIFF_PARSING.md
- **Audience:** Developers, contributors
- **Purpose:** Parser architecture and API documentation
- **Depth:** Technical, includes code examples

### docs/README.md
- **Audience:** Anyone seeking documentation
- **Purpose:** Navigation hub for all documentation
- **Depth:** Index with links and descriptions

## Maintenance Guidelines

### Root Directory
- Keep minimal (only README.md and RUN_ZIZOU.md)
- Focus on user-facing content
- No technical deep dives

### docs/ Directory
- All detailed documentation goes here
- Technical content, architecture, guides
- Educational materials

### examples/ Directory
- Runnable code examples only
- Include README.md to explain examples
- Keep examples up to date with API changes

## Contributing

When adding new documentation:

1. **User guides** → Keep in root only if absolutely essential, otherwise use `docs/`
2. **Technical docs** → Always goes in `docs/`
3. **Code examples** → Always goes in `examples/`
4. **API reference** → Goes in `docs/`

## Quick Links

**Root:**
- [Main README](../README.md)
- [Running Zizou](../RUN_ZIZOU.md)

**Documentation:**
- [Docs Index](README.md)
- [Diff Guide](DIFF_GUIDE.md)
- [Quick Reference](QUICK_REFERENCE.md)
- [Parser Guide](DIFF_PARSING.md)

**Examples:**
- [Examples Index](../examples/README.md)
- [Line by Line Demo](../examples/line_by_line_demo.go)
