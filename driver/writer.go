package driver

import (
	"os"
)

// maxQueuedWrites is the maximum number of pending write strings buffered in
// the Writer queue before blocking. When the queue is full Write blocks until
// the writer goroutine drains it.
const maxQueuedWrites = 30

// Writer serialises writes to a file (typically stdout/stderr) through a
// background goroutine, batching and flushing when the queue drains.
type Writer struct {
	queue chan string
	file  *os.File
	done  chan struct{}
}

// NewWriter creates a Writer targeting file. Call [Writer.Start] to start the
// background goroutine before calling [Writer.Write].
func NewWriter(file *os.File) *Writer {
	return &Writer{
		queue: make(chan string, maxQueuedWrites),
		file:  file,
		done:  make(chan struct{}),
	}
}

// Start launches the background writer goroutine. It must be called exactly
// once before any calls to Write.
func (w *Writer) Start() {
	go w.run()
}

// Write enqueues text for writing. It blocks if the queue is full.
// Passing an empty string is a no-op.
func (w *Writer) Write(text string) {
	if text == "" {
		return
	}
	w.queue <- text
}

// Stop signals the writer goroutine to flush remaining items and exit. It
// blocks until the goroutine has finished.
func (w *Writer) Stop() {
	// A nil sentinel signals the goroutine to exit after draining.
	w.queue <- ""
	<-w.done
}

// Flush is a no-op; the writer goroutine flushes automatically when the queue
// drains. Provided so that Writer satisfies the same informal interface as
// other flushing writers.
func (w *Writer) Flush() {}

// run is the writer goroutine body.
func (w *Writer) run() {
	defer close(w.done)
	var buf []string
	for {
		// Block until there is at least one item.
		text := <-w.queue
		if text == "" {
			// Sentinel: flush remaining buffer and exit.
			w.flush(buf)
			return
		}
		buf = append(buf, text)

		// Drain any further items without blocking.
	drain:
		for {
			select {
			case text := <-w.queue:
				if text == "" {
					w.flush(buf)
					return
				}
				buf = append(buf, text)
			default:
				break drain
			}
		}
		// Write everything accumulated in this batch.
		w.flush(buf)
		buf = buf[:0]
	}
}

func (w *Writer) flush(items []string) {
	for _, s := range items {
		_, _ = w.file.WriteString(s)
	}
}
