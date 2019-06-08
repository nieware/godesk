package main

import (
	"fmt"
	"strings"
	"time"
)

// WUID is the (hopefully) unique ID of a window
type WUID = string

// Window is a window known to our application. This is also used to
// remember windows that are no longer "Current" (have been closed)
type Window struct {
	Current   bool     // does this window currently exist?
	WindowID  int32    // window ID
	DesktopID int32    // desktop ID (-1 = all desktops)
	ProcessID int32    // process ID
	Machine   string   // Machine Name
	Titles    []string // All "stable" titles we have seen for this window
	State     string   // e.g. "_NET_WM_STATE_MODAL", "_NET_WM_STATE_MAXIMIZED_VERT, _NET_WM_STATE_MAXIMIZED_HORZ"

	cTitle   string    // current title
	cTitleTS time.Time // timestamp where current title was first seen
}

// WindowFromWMCTRLWindow creates a new Window object based on a given WMCTRLWindow.
// The state string must be provided separately, since it's not part of wmctrl's output.
func WindowFromWMCTRLWindow(ww WMCTRLWindow, state string) Window {
	return Window{
		WindowID:  ww.WindowID,
		DesktopID: ww.DesktopID,
		ProcessID: ww.ProcessID,
		Machine:   ww.Machine,
		Titles:    []string{ww.Title},
		State:     state,
	}
}

// HasTitle returns true is w has the given title in its Titles
func (w *Window) HasTitle(title string) bool {
	for _, t := range w.Titles {
		if t == title {
			return true
		}
	}
	return false
}

// HasState checks if the State string contains the specified state as a substring
func (w *Window) HasState(state string) bool {
	return strings.Contains(w.State, state)
}

// AddTitle adds the given title to a Window's Titles
func (w *Window) AddTitle(title string) {
	if w.HasTitle(title) {
		return
	}
	fmt.Println("adding title", title, "to window", w)
	w.Titles = append(w.Titles, title)
}

// String implements the Stringer interface
func (w Window) String() string {
	return "[" + w.UID() + "]"
}

// UID gets a (hopefully) unique id for the window, built from the window ID and process ID.
// Even if both are recycled when a window/application closes, the probability of collisions should be small?!
func (w Window) UID() WUID {
	return fmt.Sprintf("%d/0x%x", w.ProcessID, w.WindowID)
}
