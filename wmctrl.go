package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Calls external apps (wmctrl, xprop) and processes their output.
// One day, this might talk directly to the windows manager.

// WMCTRLWindow is the internal representation of windows returned by wmctrl
type WMCTRLWindow struct {
	WindowID  int32
	DesktopID int32
	ProcessID int32
	Machine   string
	Title     string
}

// UID gets a (hopefully) unique id for the window, built from the window ID and process ID.
// Even if both are recycled when a window/application closes, the probability of collisions should be small?!
func (ww WMCTRLWindow) UID() WUID {
	return fmt.Sprintf("%d/0x%x", ww.ProcessID, ww.WindowID)
}

func (ww WMCTRLWindow) String() string {
	return "[" + ww.UID() + "]"
}

func getWindowListFromWMCTRLOutput(out string) map[WUID]WMCTRLWindow {
	lines := strings.Split(out, "\n")
	rex := regexp.MustCompile(`0x([0-9a-f]+)\s+([\-0-9]+)\s+([0-9]+)\s+([^\s]+)\s+(.+)`)
	wws := make(map[WUID]WMCTRLWindow)
	for _, line := range lines {
		matches := rex.FindStringSubmatch(line)
		if len(matches) < 6 {
			continue
		}

		windowID, windowIDOK := strconv.ParseInt(matches[1], 16, 32)
		desktopID, desktopIDOK := strconv.Atoi(matches[2])
		processID, processIDOK := strconv.Atoi(matches[3])
		if windowIDOK != nil || desktopIDOK != nil || processIDOK != nil {
			continue
		}
		ww := WMCTRLWindow{
			WindowID:  int32(windowID),
			DesktopID: int32(desktopID),
			ProcessID: int32(processID),
			Machine:   matches[4],
			Title:     matches[5],
		}
		wws[ww.UID()] = ww
	}

	return wws
}

func getWMCTRLWindows() map[WUID]WMCTRLWindow {
	out := execCmd(`wmctrl`, `-lp`) // window list with process ids
	return getWindowListFromWMCTRLOutput(out)
}

func getDesktopForWindowID(UID WUID) (desktopID int, found bool) {
	wws := getWMCTRLWindows()
	ww, found := wws[UID]
	if found {
		desktopID = int(ww.DesktopID)
	}
	return
}

// TODO currently we only look at desktop Ids, so if the user reorders their desktops,
// we will use the wrong ones. However, if we use names and the desktops are renamed, that could be just as bad...
/*func getDesktops() {
	out := execCmd(`wmctrl`, `-d`)
	getDesktopListFromOutput(out)
}*/

func moveWindowToDesktop(windowID int32, desktopID int32) {
	/*out :=*/
	if desktopID == -1 {
		desktopID = -2 // peculiarity of wmctrl: to send to all desktops, we must use -2, not -1 as listed
		// TODO: some sources also recommend "wmctrl -i -r 0x03600011 -b add,sticky", maybe use both and check success via getDesktopForWindowID
	}
	execCmd(
		`wmctrl`,
		`-i`,                              // Option: interpret window arguments (<WIN>) as a numeric value rather than a string name for the window.
		`-r`, strconv.Itoa(int(windowID)), // Action: specify target window
		`-t`, strconv.Itoa(int(desktopID)), // Action: move a window that has been specified with the -r action to the desktop <DESK>.
	)
}
