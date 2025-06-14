package loggerfx

import (
	"bufio"
	"context"
	"io"
	"time"

	"go.uber.org/multierr"
)

type BufferIOConfig struct {
	ChannelSize int
	BufferSize  int
}

type nonBlockingBufferedWriter struct {
	cancel context.CancelFunc
	ch     chan []byte
	buffer *bufio.Writer
	closer io.Closer
}

func newNonBlockingBufferedWriter(sink Sink) (*nonBlockingBufferedWriter, error) {
	f, err := extractFile(sink)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	cfg := BufferIOConfig{
		ChannelSize: 1024,
		BufferSize:  32 * 1024,
	}

	if len(sink.Args) > 1 {
		config, ok := sink.Args[1].(BufferIOConfig)

		if ok {
			cfg = config
		}
	}

	ch := make(chan []byte, cfg.ChannelSize)

	buffer := bufio.NewWriterSize(f, cfg.BufferSize)
	go func() {
		timer := time.NewTicker(1 * time.Second)

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				_ = buffer.Flush()
			case data := <-ch:
				_, _ = buffer.Write(data)
			}
		}
	}()

	return &nonBlockingBufferedWriter{
		cancel: cancel,
		ch:     ch,
		closer: f,
		buffer: buffer,
	}, nil
}

func (n *nonBlockingBufferedWriter) Write(data []byte) (int, error) {
	n.ch <- data

	return len(data), nil
}

func (n *nonBlockingBufferedWriter) Close() error {
	n.cancel()
	close(n.ch)
	err := multierr.Append(nil, n.buffer.Flush())

	return multierr.Append(err, n.closer.Close())
}
