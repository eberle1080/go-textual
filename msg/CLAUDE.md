# msg

Package `msg`. All framework message types and `Cmd` utilities. Defines `Msg` (sealed interface), `Cmd` (`func(context.Context) Msg`), and built-in messages: `KeyMsg`, `MouseMsg`, `QuitMsg`, `TickMsg`, `ResizeMsg`, `PanicMsg`. Also provides helpers: `msg.Batch(...)` fans out multiple Cmds, `msg.Sequence(...)` chains them, and `msg.Tick(d)` produces recurring ticks.
