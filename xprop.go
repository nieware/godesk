package main

import (
	"regexp"
	"strconv"
	"strings"
)

// XProp contains an X window property returned by xprop
type XProp struct {
	// e.g. "_NET_WM_STATE(ATOM) = _NET_WM_STATE_MAXIMIZED_VERT, _NET_WM_STATE_MAXIMIZED_HORZ"
	Name   string   // property name (e.g. "_NET_WM_STATE")
	Type   string   // property type (e.g. "ATOM")
	Value  string   // property value as single string (e.g. "_NET_WM_STATE_MAXIMIZED_VERT, _NET_WM_STATE_MAXIMIZED_HORZ")
	Values []string // property value split into multiple strings (e.g. {"_NET_WM_STATE_MAXIMIZED_VERT", "_NET_WM_STATE_MAXIMIZED_HORZ"})
}

// get all (if key == "") or one (if key given) XProps as a slice of XProp records
func getXProps(windowID int32, key string) map[string]XProp {
	// xprop -notype _NET_WM_STATE -id 0x6a00001
	out := execCmd(
		`xprop`,
		key,
		`-id`, `0x`+strconv.FormatInt(int64(windowID), 16),
	)
	rex := regexp.MustCompile(`^([_A-Z0-9]+)(\(*)([_A-Z0-9]*)(\)*)([\s=:]+)(.*)$`)
	lines := strings.Split(out, "\n")
	props := make(map[string]XProp)
	for _, line := range lines {
		matches := rex.FindStringSubmatch(line)
		if len(matches) < 7 {
			continue
		}

		xp := XProp{
			Name:   matches[1],
			Type:   matches[3],
			Value:  matches[6],
			Values: strings.Split(matches[6], ", "),
		}

		props[xp.Name] = xp
	}
	return props
}

func getXProp(windowID int32, key string) (prop XProp, found bool) {
	props := getXProps(windowID, key)
	prop, found = props[key]
	return
}

func getWindowState(windowID int32) string {
	// xprop _NET_WM_STATE -id 0xa00194
	prop, found := getXProp(windowID, `_NET_WM_STATE`)

	if found {
		return prop.Value
	}
	return ""
}
