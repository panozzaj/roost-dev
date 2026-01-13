package icons

import (
	_ "embed"
	"strings"
)

//go:embed gear.svg
var Gear string

//go:embed copy.svg
var Copy string

//go:embed external-link.svg
var ExternalLink string

//go:embed clipboard.svg
var Clipboard string

//go:embed clipboard-agent.svg
var ClipboardAgent string

//go:embed claude.svg
var Claude string

//go:embed check.svg
var Check string

//go:embed x.svg
var X string

//go:embed trash.svg
var Trash string

func init() {
	// Trim whitespace from embedded SVGs
	Gear = strings.TrimSpace(Gear)
	Copy = strings.TrimSpace(Copy)
	ExternalLink = strings.TrimSpace(ExternalLink)
	Clipboard = strings.TrimSpace(Clipboard)
	ClipboardAgent = strings.TrimSpace(ClipboardAgent)
	Claude = strings.TrimSpace(Claude)
	Check = strings.TrimSpace(Check)
	X = strings.TrimSpace(X)
	Trash = strings.TrimSpace(Trash)
}

// CheckGreen returns a check icon with green stroke
func CheckGreen() string {
	return strings.Replace(Check, `stroke="currentColor"`, `stroke="#22c55e"`, 1)
}

// XRed returns an X icon with red stroke
func XRed() string {
	return strings.Replace(X, `stroke="currentColor"`, `stroke="#ef4444"`, 1)
}

// JSObject returns a JavaScript object containing all icons for client-side use
func JSObject() string {
	return `var ICONS = {
    gear: '` + escapeJS(Gear) + `',
    copy: '` + escapeJS(Copy) + `',
    externalLink: '` + escapeJS(ExternalLink) + `',
    clipboard: '` + escapeJS(Clipboard) + `',
    clipboardAgent: '` + escapeJS(ClipboardAgent) + `',
    claude: '` + escapeJS(Claude) + `',
    check: '` + escapeJS(Check) + `',
    checkGreen: '` + escapeJS(CheckGreen()) + `',
    x: '` + escapeJS(X) + `',
    xRed: '` + escapeJS(XRed()) + `',
    trash: '` + escapeJS(Trash) + `'
};`
}

func escapeJS(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}
