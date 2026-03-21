# driver

Package `driver`. Terminal I/O abstraction layer. Defines the `Driver` interface and platform-specific implementations for entering raw mode, switching to the alternate screen, writing styled output (`writer.go`), parsing ANSI escape sequences (`escapes.go`), and handling Unix signals like `SIGWINCH` and `SIGTSTP`.
