package main

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)


// Aof is a struct that represents an append-only file. It contains an underlying
// os.File and a bufio.Reader, as well as a sync.Mutex for synchronizing access.
type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

// NewAof creates a new Aof instance with the given file path. It opens the file
// for reading and writing, and starts a goroutine that syncs the file to disk
// every 1 second.
func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// start go routine to sync aof to disk every 1 second
	go func() {
		for {
			aof.mu.Lock()

			aof.file.Sync()

			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

// Close closes the underlying file for the Aof instance. This method is thread-safe
// and ensures that the file is properly closed and synced to disk before returning.
func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

// Write appends the given Value to the append-only file. It acquires a lock to
// ensure thread-safety, writes the marshaled value to the file, and then
// releases the lock. Any errors encountered during the write operation are
// returned.
func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

// Read reads all values from the append-only file and calls the provided
// function for each value. It acquires a lock to ensure thread-safety,
// seeks to the start of the file, and then reads each value, passing it
// to the provided function. Any errors encountered during the read
// operation are returned.
//
// NOTE: This is very slow when starting up when the DB has a lot of data.
func (aof *Aof) Read(fn func(value Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	aof.file.Seek(0, io.SeekStart)

	reader := NewResp(aof.file)

	for {
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		fn(value)
	}

	return nil
}