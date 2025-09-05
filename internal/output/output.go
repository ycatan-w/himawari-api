package output

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ycatan-w/himawari-api/internal/output/colors"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func visibleLen(s string) int {
	return len(ansiRegexp.ReplaceAllString(s, ""))
}
func AppNameGreen() string {
	return colors.Green("himawari-server")
}
func PrintBox(s string) {
	line := strings.Repeat("─", visibleLen(s)+2)

	fmt.Printf("\n┌%s┐\n", line)
	fmt.Printf("│ %s │\n", s)
	fmt.Printf("└%s┘\n", line)
}

func PrintHeader(s string) {
	now := time.Now().UTC().Format(time.RFC3339)
	header := []string{
		colors.BrightCyan("\n==========["),
		colors.BrightYellow(now),
		fmt.Sprintf(colors.BrightCyan(" • %s ]==========\n"), s),
	}
	fmt.Print(strings.Join(header, ""))
}

func PrintSubHeader(s string) {
	fmt.Printf(colors.Cyan("\n--- %s\n"), s)
}

func PrintSuccess(s string) {
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Printf("[%s] %s %s\n", colors.BrightYellow(now), colors.Green("✔"), s)
}

func PrintFail(s string) {
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Printf("[%s] %s %s\n", colors.BrightYellow(now), colors.Red("✖"), s)
}

func PrintWarn(s string) {
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Printf("[%s] %s %s\n", colors.BrightYellow(now), colors.Yellow("‼"), s)
}

func PrintInfo(s string) {
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Printf("[%s] %s %s\n", colors.BrightYellow(now), colors.Blue("ⓘ"), s)
}
