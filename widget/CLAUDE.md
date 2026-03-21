# widget

Package `widget`. Core `Widget` interface and `BaseWidget` implementation. All concrete widgets embed `BaseWidget` and implement `Compose()` (declare children), `Update(Msg) Cmd` (handle messages), and `Render(region) []strip.Strip` (produce output). Also contains focus management logic.
