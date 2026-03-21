# app

Package `app`. Application lifecycle and event loop. This is the heart of the framework — it owns the single-goroutine event loop (`loop.go`), dispatches messages to the widget tree, spawns `Cmd` goroutines, drives rendering, and initializes the platform driver. Also contains the headless driver for testing and platform-specific signal/driver wiring for Unix and Windows.
