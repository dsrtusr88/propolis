// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package main

import "gitlab.com/catastrophic/assistance/ui"

var (
	titleHeader = ui.BlueBold("▻ ")
)

const (
	arrowHeader   = "⮕ "
	internalRule  = ` -- `
	OKString      = " 🗹 "
	KOString      = " 🗷 "
	WarningString = " 🞎 "
	NeutralString = " 🛈 "

	integrityCheckOK = "Integrity checks successful for all FLACs, no ID3 tags detected."
)
