# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build ./...          # Build all packages
go test ./...           # Run all tests
go run ./examples/counter  # Run an example app
```

No Makefile — standard Go tooling only.

## Architecture

**go-textual** is a Go port of Python's Textual TUI framework. It uses a single-goroutine event loop with message-passing to build interactive terminal applications.

### Core Design

**Event Loop** (`app/loop.go`): The single goroutine that owns all widget state. Each cycle:
1. Drains all queued messages from `a.events` channel
2. Dispatches each message to the screen/widget tree, collecting returned `Cmd`s
3. Spawns each `Cmd` as a goroutine (Cmds do async work and return a Msg)
4. Re-renders if any widget is dirty
5. Blocks on `a.events` until next message

No mutexes are needed because only the event loop goroutine touches widget state.

**Message System** (`msg/`):
- `Msg` is a sealed interface — all messages embed `BaseMsg`
- `Cmd` is `func(context.Context) Msg` — async work that produces a message
- `Update(Msg) Cmd` is called on each widget to handle messages; return nil for no-op
- Built-in: `KeyMsg`, `MouseMsg`, `QuitMsg`, `TickMsg`, `ResizeMsg`, `PanicMsg`
- Helpers: `msg.Batch(...)` fans out multiple Cmds; `msg.Sequence(...)` chains them; `msg.Tick(d)` produces recurring ticks

**Widget Tree** (`widget/`, `screen/`):
- `Widget` interface: `Compose()`, `Update(Msg) Cmd`, `Render(region) []strip.Strip`
- `BaseWidget`: default implementation embedded in all concrete widgets
- `Screen` is the root widget (application entry point); implement `OnMount() Cmd` for initialization
- Widgets declare children in `Compose()` — the framework builds the tree

**Rendering** (`strip/`, `layout/`):
- `strip.Strip`: immutable horizontal line of Rich-styled text segments
- `Render()` returns one `Strip` per screen row within its region
- Layout algorithms (`layout/vertical.go`, `horizontal.go`, `grid.go`, `stream.go`) partition regions among children

**DOM & Styling** (`dom/`, `css/`):
- `DOMNode` tracks the parent-child tree and CSS selector matching
- CSS-inspired property system for colors, padding, borders, dimensions

**Platform Drivers** (`driver/`, `app/driver_unix.go`, `app/driver_windows.go`):
- `driver.Driver` interface abstracts terminal I/O
- Platform implementations handle raw mode, alternate screen, SIGWINCH, SIGTSTP
- Headless driver available for testing

### Key Packages

| Package | Role |
|---------|------|
| `app` | Application lifecycle, event loop, message dispatch |
| `widget` | Widget interface and BaseWidget |
| `screen` | Screen (root widget) interface |
| `msg` | All message types and Cmd utilities |
| `widgets` | Concrete widgets: Button, Label, Input, TextArea, ListView, DataTable, TabbedContent, Header, Footer, RichLog, DirectoryTree, ProgressBar, etc. |
| `layout` | Vertical, Horizontal, Grid, Stream layout algorithms |
| `strip` | Immutable styled text line primitives |
| `dom` | DOM tree and CSS selector matching |
| `css` | CSS parsing and styling model |
| `geometry` | Size, Region, Offset types |
| `binding` | Key binding management |
| `internal` | ANSI, caching, layout resolution, spatial utilities |

### Examples

`examples/` contains: `hello`, `counter`, `calculator`, `stopwatch`, `todo`, `dashboard`, `file_browser` — all runnable with `go run ./examples/<name>`.
