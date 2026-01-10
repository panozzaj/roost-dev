package logo

import "strings"

// Base logo with no leading spaces on first line
const base = `__            __
   _________  ____  _____/ /_      ____/ /__ _   __
  / ___/ __ \/ __ \/ ___/ __/_____/ __  / _ \ | / /
 / /  / /_/ / /_/ (__  ) /_/_____/ /_/ /  __/ |/ /
/_/   \____/\____/____/\__/      \__,_/\___/|___/`

// Get returns the logo with specified leading spaces on the first line
func Get(indent int) string {
	return strings.Repeat(" ", indent) + base
}

// CLI returns logo formatted for terminal output (26 spaces)
func CLI() string {
	return Get(26)
}

// Web returns logo formatted for web/HTML output (17 spaces)
func Web() string {
	return Get(17)
}
