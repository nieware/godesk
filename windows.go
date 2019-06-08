package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"os/user"
	"strings"
)

// Windows is a list that stores all windows currently known to our application,
// including ones that are no longer "current" (have been closed).
type Windows map[WUID]Window

// NewWindows creates a new Windows object
func NewWindows() Windows {
	return make(Windows)
}

// Refresh refreshes the window list
func (ws Windows) Refresh() {
	wws := getWMCTRLWindows()
	// set windows which are not in the current window list to "non-current"
	for UID, w := range ws {
		_, found := wws[UID]
		if w.Current && !found {
			fmt.Println("window", w, "closed")
			w.Current = false
			ws[UID] = w
		}
	}
	for UID, ww := range wws {
		w, ok := ws[UID]
		if !ok {
			// not present in list
			ws.HandleNewWindow(ww)
			continue
		}
		// present in list -> refresh info
		if w.DesktopID != ww.DesktopID {
			fmt.Println("Window", w, "moved from desktop", w.DesktopID, "to desktop", ww.DesktopID)
		}
		w.DesktopID = ww.DesktopID
		w.Machine = ww.Machine
		w.AddTitle(ww.Title)
		w.Current = true
		ws[w.UID()] = w
	}
}

func (ws Windows) WindowShouldBeIgnored(w Window) bool {
	return w.HasState("_NET_WM_STATE_MODAL")
}

// HandleNewWindow adds a newly created window to the list and
// places it on the correct desktop.
func (ws Windows) HandleNewWindow(ww WMCTRLWindow) {
	state := getWindowState(ww.WindowID)
	addW := WindowFromWMCTRLWindow(ww, state)

	// some windows should not be moved - e.g. modal dialogs, because they belong to their parent window
	if ws.WindowShouldBeIgnored(addW) {
		fmt.Println("Ignoring", addW, "state is", addW.State)
		return
	}

	// if we have a matching non-current window, move it to the correct desktop!
	w, found := ws.FindWindowByTitle(ww.Title, true)
	if found {
		fmt.Println("found prev. window", w, "for new window", ww, "using title", ww.Title, ". MOVING to desktop", w.DesktopID)
		moveWindowToDesktop(ww.WindowID, w.DesktopID)
		// "appropriate" the found entry so we don't have several entries for the "same" window
		addW = w
		delete(ws, w.UID())
		addW.WindowID = ww.WindowID
		addW.ProcessID = ww.ProcessID
		fmt.Println("deleting prev. window", w)
	}

	// not present in list -> add it
	fmt.Println("adding new window", ww, "with title", ww.Title)
	ws[ww.UID()] = addW
}

// FindWindowByTitle returns the first window that has the given title in its slice of Titles.
// If excludeCurrent is true, windows that are currently open (w.Current == true) are not considered.
func (ws Windows) FindWindowByTitle(title string, excludeCurrent bool) (w Window, found bool) {
	if strings.TrimSpace(title) == "" {
		// ignore empty titles
		return
	}

	for _, w = range ws {
		if (!w.Current || w.Current && !excludeCurrent) && w.HasTitle(title) {
			found = true
			return
		}
	}
	return
}

func (ws Windows) getConfigFileName() (fn string, err error) {
	// from https://stackoverflow.com/questions/7922270/obtain-users-home-directory
	usr, err := user.Current()
	if err != nil {
		return
	}
	fn = usr.HomeDir + "/godeskcfg"
	return
}

// SaveToFile saves the Windows list to the given filename.
func (ws Windows) SaveToFile() error {
	fn, err := ws.getConfigFileName()
	if err != nil {
		return err
	}

	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	err = encoder.Encode(ws)
	if err != nil {
		return err
	}

	return nil
}

// ReadFromFile reads the Windows list from a given filename.
// If reset is true, reset the window and process IDs.
func (ws Windows) ReadFromFile(reset bool) error {
	fn, err := ws.getConfigFileName()
	if err != nil {
		fmt.Println(1, err)
		return err
	}

	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(2, err)
		return err
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&ws)
	if err != nil {
		fmt.Println(3, err)
		return err
	}

	if reset {
		for _, w := range ws {
			// reset the window and process ids when reading from file,
			// as a match would be a mere coincidence after reboot.
			w.WindowID = 0
			w.ProcessID = 0
		}
	}

	return nil
}
