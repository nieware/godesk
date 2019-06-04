package main

import "fmt"

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
	Titles    []string // All titles we have seen for this window
}

// WindowFromWMCTRLWindow creates a new Window object based on a given WMCTRLWindow
func WindowFromWMCTRLWindow(ww WMCTRLWindow) Window {
	return Window{
		WindowID:  ww.WindowID,
		DesktopID: ww.DesktopID,
		ProcessID: ww.ProcessID,
		Machine:   ww.Machine,
		Titles:    []string{ww.Title},
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
