# go-textual

A Go port of [Textual](https://github.com/Textualize/textual), the Python TUI framework. Build interactive terminal applications with a clean, event-driven API.

## Overview

go-textual brings Textual's architecture to Go:

- **Event-driven**: A single event loop goroutine owns all widget state — no locks, no races
- **Message-passing**: All inter-component communication flows through typed messages
- **Composable widgets**: Build UIs by composing a tree of widgets with `Compose()`
- **CSS-inspired styling**: Style widgets with a CSS-like property system
- **Rich rendering**: Styled, colored output via [go-rich](https://github.com/eberle1080/go-rich)

## Install

```
go get github.com/eberle1080/go-textual
```

## Quick Start

```go
package main

import (
    "context"

    "github.com/eberle1080/go-textual/app"
    "github.com/eberle1080/go-textual/geometry"
    "github.com/eberle1080/go-textual/msg"
    "github.com/eberle1080/go-textual/screen"
    "github.com/eberle1080/go-textual/strip"
    "github.com/eberle1080/go-textual/widget"
    "github.com/eberle1080/go-textual/widgets"
)

type HelloScreen struct {
    screen.BaseScreen
}

func (s *HelloScreen) Compose() []widget.Widget {
    return []widget.Widget{
        widgets.NewLabel("Hello, world!"),
    }
}

func (s *HelloScreen) Update(ctx context.Context, m msg.Msg) msg.Cmd {
    switch m.(type) {
    case msg.KeyMsg:
        return msg.Quit()
    }
    return nil
}

func main() {
    a := app.New()
    a.Run(context.Background(), &HelloScreen{})
}
```

## Widgets

| Widget | Description |
|--------|-------------|
| `Button` | Clickable button, emits `ButtonPressedMsg` |
| `Label` | Static text display |
| `Input` | Single-line text input |
| `TextArea` | Multi-line text editor |
| `ListView` | Scrollable list with selection |
| `DataTable` | Tabular data display |
| `TabbedContent` | Tab-based content switcher |
| `Header` / `Footer` | App chrome |
| `RichLog` | Scrollable rich text log |
| `DirectoryTree` | File system tree browser |
| `ProgressBar` | Progress visualization |
| `Sparkline` | Compact data chart |
| `Digits` | 7-segment style digit display |

## Architecture

### Messages & Commands

All state changes flow through messages. An `Update()` method receives a message and optionally returns a `Cmd` — a function that runs asynchronously and sends a message back to the loop.

```go
func (s *MyScreen) Update(ctx context.Context, m msg.Msg) msg.Cmd {
    switch v := m.(type) {
    case msg.KeyMsg:
        if v.Key == "q" {
            return msg.Quit()
        }
    case MyCustomMsg:
        s.value = v.Data
    }
    return nil
}
```

### Periodic Updates

Use `msg.Tick()` to drive time-based updates:

```go
func (s *MyScreen) OnMount(ctx context.Context) msg.Cmd {
    return msg.Tick(time.Second, func(t time.Time) msg.Msg {
        return TickMsg{T: t}
    })
}
```

### Layout

Widgets are arranged using one of four layout algorithms, configurable via CSS-style properties:

- `VerticalLayout` — stack children top to bottom (default)
- `HorizontalLayout` — stack children left to right
- `GridLayout` — place children in a grid
- `StreamLayout` — flow children like inline elements

## Examples

The `examples/` directory contains:

| Example | Description |
|---------|-------------|
| `hello` | Minimal label + button |
| `counter` | Stateful counter with keyboard and button handling |
| `calculator` | 4-function calculator |
| `stopwatch` | Timer with start/stop/reset using `msg.Tick` |
| `todo` | Todo list with `ListView` and `Input` |
| `dashboard` | Multi-section layout |
| `file_browser` | File system browser using `DirectoryTree` |

Run any example:

```
go run ./examples/counter
```

## Status

Early development — core infrastructure works and the examples run. APIs may change.

## License

MIT — see [LICENSE](LICENSE).

Inspired by [Textual](https://github.com/Textualize/textual) by Will McGugan / Textualize.
