package gui

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type CreationCallback = func(*fsnotify.Event)

// Depending on the system, a single "write" can generate many Write events; for
// example compiling a large Go program can generate hundreds of Write events on
// the binary.
//
// The general strategy to deal with this is to wait a short time for more write
// events, resetting the wait period for every new event.
func dedup(callback CreationCallback, paths ...string) {
	if len(paths) < 1 {
		panic("must specify at least one path to watch")
	}

	// Create a new watcher.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprintf("creating a new watcher: %s", err))
	}
	defer w.Close()

	// Start listening for events.
	go watchForCreated(w, callback)

	// Add all paths from the commandline.
	for _, p := range paths {
		err = w.Add(p)
		if err != nil {
			panic(fmt.Sprintf("%q: %s", p, err))
		}
	}
}

func watchForCreated(w *fsnotify.Watcher, creatcb CreationCallback) {
	var (
		// Wait 100ms for new events; each new event resets the timer.
		waitFor = 100 * time.Millisecond
		// why is this not a const?

		// Keep track of the timers, as path â†’ timer.
		mu     sync.Mutex
		timers = make(map[string]*time.Timer)

		// Callback we run.
		cb = func(e fsnotify.Event) {
			creatcb(&e)

			// Don't need to remove the timer if you don't have a lot of files.
			mu.Lock()
			delete(timers, e.Name)
			mu.Unlock()
		}
	)

	for {
		select {
		// Read from Errors.
		case _, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
		// Read from Events.
		case e, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}

			// We just want to watch for file creation, so ignore everything
			// outside of Create and Write.
			if !e.Has(fsnotify.Create) && !e.Has(fsnotify.Write) {
				continue
			}

			// Get timer.
			mu.Lock()
			t, ok := timers[e.Name]
			mu.Unlock()

			// No timer yet, so create one.
			if !ok {
				t = time.AfterFunc(math.MaxInt64, func() { cb(e) })
				t.Stop()

				mu.Lock()
				timers[e.Name] = t
				mu.Unlock()
			}

			// Reset the timer for this path, so it will start from 100ms again.
			t.Reset(waitFor)
		}
	}
}

type NewSaveCallback = func(string)

func watchForNewSave(cb NewSaveCallback) {
	watcher, err := fsnotify.NewWatcher()
	w := watcher
	if err != nil {
		panic("Could not create watcher")
	}
	defer watcher.Close()
	err = watcher.Add(guessGameDir())
	if err != nil {
		panic(fmt.Sprintf("Could not create add dir to watcher %v", err))
	}
	watchForCreated(w, func(e *fsnotify.Event) {
		base := filepath.Base(e.Name)
		if base != "save_loading" && strings.HasPrefix(base, "save") {
			cb(e.Name)
		}
	})
}
