package ftpb

import (
	"bytes"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gosuri/uilive"
)

// FileTransferProgressBar represents a single file transfer
type FileTransferProgressBar struct {
	doneBytes         int64
	totalBytes        int64
	writer            *uilive.Writer
	started           bool
	byteBuffer        *bytes.Buffer
	ShowTransferSpeed bool
	stop              chan bool
}

// New returns a new file transfer progress bar
func New(totalBytes int64) *FileTransferProgressBar {
	writer := uilive.New()
	writer.Start()
	return &FileTransferProgressBar{
		ShowTransferSpeed: true,
		doneBytes:         0,
		totalBytes:        totalBytes,
		writer:            writer,
		stop:              make(chan bool),
	}
}

// SetWriter changes the io.Writer that the progress bar gets outputted to.
// The default writer is an instance of uilive.Writer but it can be changed for testing
func (f *FileTransferProgressBar) SetWriter(w *uilive.Writer) {
	f.writer = w
}

// SetByteBuffer makes puts the rendered output into the referenced byte buffer
// insted of outputting it to the default *uilive.Writer. Can be useful for testing,
// and more.
func (f *FileTransferProgressBar) SetByteBuffer(b *bytes.Buffer) {
	f.byteBuffer = b
}

//Start the rendering loop
func (f *FileTransferProgressBar) Start() {
	f.started = true
	time.Sleep(time.Millisecond * 100)
	go func() {
		for {
			select {
			case <-f.stop:
				f.render()
				return
			default:
				f.render()
			}
		}
	}()
}

func (f *FileTransferProgressBar) render() {
	if f.byteBuffer != nil {
		fmt.Fprintln(f.byteBuffer, f.generate())
	} else {
		fmt.Fprintln(f.writer, f.generate())
		time.Sleep(time.Millisecond * 200)
	}
}

func (f *FileTransferProgressBar) generate() string {
	if f.doneBytes > f.totalBytes {
		f.doneBytes = f.totalBytes
	}

	return fmt.Sprintf(
		"%s / %s",
		humanize.Bytes(uint64(f.doneBytes)),
		humanize.Bytes(uint64(f.totalBytes)),
	)
}

// Increase the progressbar by v amount of bytes
func (f *FileTransferProgressBar) Increase(v int64) {
	f.updateDoneBytes(f.doneBytes + v)
}

// Set the progress bar to v amount of bytes
func (f *FileTransferProgressBar) Set(v int64) {
	f.updateDoneBytes(v)
}

func (f *FileTransferProgressBar) updateDoneBytes(v int64) {
	f.doneBytes = v
	if f.byteBuffer != nil {
		f.render()
	}
}

// Done sets the progressbar to 100% and close the progressbar
func (f *FileTransferProgressBar) Done() {
	f.doneBytes = f.totalBytes
	if f.started {
		f.stop <- true
		f.writer.Stop()
	} else {
		f.render()
	}
	time.Sleep(100 * time.Millisecond)
}
