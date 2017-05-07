package main

import (
	"time"
	"net"
	"fmt"
	"io"
)

const (
	BUFFER_SIZE = 16 * 1024
	PROXY_STATS_PUSH_INTERVAL = 1 * time.Second
)

type ReadWriteCount struct {
	CountRead  uint
	CountWrite uint
	Target     Target
}

func (this ReadWriteCount) IsZero() bool {
	return this.CountRead == 0 && this.CountWrite == 0
}

type Target struct {
	Host string
	Port string
}
/**
 * Compare to other target
 */
func (t *Target) EqualTo(other Target) bool {
	return t.Host == other.Host &&
		t.Port == other.Port
}

/**
 * Get target full address
 * host:port
 */
func (this *Target) Address() string {
	return this.Host + ":" + this.Port
}

/**
 * To String conversion
 */
func (this *Target) String() string {
	return this.Address()
}

func proxy(to net.Conn, from net.Conn, timeout time.Duration) <- chan ReadWriteCount {
	stats := make(chan ReadWriteCount)
	outStats := make(chan ReadWriteCount)

	rwcBuffer := ReadWriteCount{}
	ticker := time.NewTicker(PROXY_STATS_PUSH_INTERVAL)
	flushed := false

	go func() {
		if timeout > 0 {
			from.SetReadDeadline(time.Now().Add(timeout))
		}
		for {
			select {
			case <-ticker.C:
				if !rwcBuffer.IsZero() {
					outStats <- rwcBuffer
				}
				flushed = true
			case rwc, ok := <-stats:
				if !ok {
					ticker.Stop()
					if !flushed && !rwcBuffer.IsZero() {
						outStats <- rwcBuffer
					}
					close(outStats)
					return
				}

				if timeout > 0 && rwc.CountRead > 0 {
					from.SetReadDeadline(time.Now().Add(timeout))
				}
				if flushed {
					rwcBuffer = rwc
				} else {
					rwcBuffer.CountWrite += rwc.CountWrite
					rwcBuffer.CountRead += rwc.CountRead
				}
				flushed = false
			}

		}
	}()

	go func() {
		err := Copy(to, from, stats)
		e, ok := err.(*net.OpError)
		if err != nil && (!ok || e.Err.Error() != "use of closed network connection") {
			fmt.Printf("Error occured i dont know %v\n", err)
		}
		to.Close()
		from.Close()
		close(stats)
	}()
	return outStats
}

func Copy(to io.Writer, from io.Reader, ch chan <- ReadWriteCount) error {
	buf := make([]byte, BUFFER_SIZE)
	var err error = nil
	for {
		readN, readError := from.Read(buf)
		if readN > 0 {
			writeN, writeError := to.Write(buf[0:readN])
			if writeN > 0 {
				ch <- ReadWriteCount{CountRead:uint(readN), CountWrite:uint(writeN)}
			}
			if writeError != nil {
				err = writeError
				break
			}

			if readN != writeN {
				err = io.ErrShortWrite
				break
			}
		}

		if readError == io.EOF {
			break
		}
		if readError != nil {
			err = readError
			break
		}
	}
	return err
}
