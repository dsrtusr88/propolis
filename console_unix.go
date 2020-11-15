// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package propolis

import "gitlab.com/catastrophic/assistance/ui"

var (
	titleHeader = ui.BlueBold("▻ ")
)

const (
	ArrowHeader   = "⮕ "
	internalRule  = ` -- `
	OKString      = " 🗹 "
	KOString      = " 🗷 "
	WarningString = " 🞎 "
	NeutralString = " 🛈 "

	integrityCheckOK = "Integrity checks successful for all FLACs, no ID3 tags detected."
)

func simplifyBullet(bullet string) string {
	switch bullet {
	case OKString:
		return "OK"
	case KOString:
		return "KO"
	case WarningString:
		return "!!"
	case NeutralString:
		return "--"
	default:
		return "  "
	}
}
