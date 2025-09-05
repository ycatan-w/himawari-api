package colors

import "fmt"

var (
	RED                 = "\033[0;31m"
	RED_FLASH           = "\033[5;31m"
	BACKGROUND_RED      = "\033[0;41;97m"
	GREEN               = "\033[0;32m"
	GREEN_FLASH         = "\033[5;32m"
	YELLOW              = "\033[0;33m"
	YELLOW_UNDERLINE    = "\033[4;33m"
	YELLOW_FLASH        = "\033[5;33m"
	YELLOW_BOLD         = "\033[1;33m"
	BRIGHT_MAGENTA_BOLD = "\033[1;95m"
	BRIGHT_YELLOW       = "\033[2;93m"
	BRIGHT_YELLOW_FLASH = "\033[5;2;93m"
	CYAN                = "\033[0;36m"
	BRIGHT_CYAN         = "\033[0;96m"
	BLUE                = "\033[0;34m"
	NC                  = "\033[0m"
)

func wrapColor(color, s string) string {
	return fmt.Sprintf("%s%s%s", color, s, NC)
}

func Red(s string) string {
	return wrapColor(RED, s)
}
func RedFlash(s string) string {
	return wrapColor(RED_FLASH, s)
}
func BackgroundRed(s string) string {
	return wrapColor(BACKGROUND_RED, s)
}
func Green(s string) string {
	return wrapColor(GREEN, s)
}
func GreenFlash(s string) string {
	return wrapColor(GREEN_FLASH, s)
}
func Yellow(s string) string {
	return wrapColor(YELLOW, s)
}
func YellowUnderline(s string) string {
	return wrapColor(YELLOW_UNDERLINE, s)
}
func YellowFlash(s string) string {
	return wrapColor(YELLOW_FLASH, s)
}
func YellowBold(s string) string {
	return wrapColor(YELLOW_BOLD, s)
}
func BrightMagentaBold(s string) string {
	return wrapColor(BRIGHT_MAGENTA_BOLD, s)
}
func BrightYellow(s string) string {
	return wrapColor(BRIGHT_YELLOW, s)
}
func BrightYellowFlash(s string) string {
	return wrapColor(BRIGHT_YELLOW_FLASH, s)
}
func Cyan(s string) string {
	return wrapColor(CYAN, s)
}
func BrightCyan(s string) string {
	return wrapColor(BRIGHT_CYAN, s)
}
func Blue(s string) string {
	return wrapColor(BLUE, s)
}
