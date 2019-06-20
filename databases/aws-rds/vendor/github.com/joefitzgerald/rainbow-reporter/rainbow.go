package reporter

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"testing"

	"github.com/sclevine/spec"
)

var denoter = "•"

func init() {
	if runtime.GOOS == "windows" {
		denoter = "+"
	}
}

const defaultStyle = "\x1b[0m"
const boldStyle = "\x1b[1m"
const redColor = "\x1b[91m"
const greenColor = "\x1b[32m"
const yellowColor = "\x1b[33m"
const cyanColor = "\x1b[36m"
const grayColor = "\x1b[90m"
const lightGrayColor = "\x1b[37m"

func colorize(colorCode string, format string, args ...interface{}) string {
	var out string

	if len(args) > 0 {
		out = fmt.Sprintf(format, args...)
	} else {
		out = format
	}

	return fmt.Sprintf("%s%s%s", colorCode, out, defaultStyle)
}

// Rainbow reports specs via stdout with colorization.
type Rainbow struct{}

// Start is called when the plan is executed.
func (Rainbow) Start(_ *testing.T, plan spec.Plan) {
	fmt.Println("Suite:", plan.Text)

	fmt.Printf("Total: %d | Focused: %d | Pending: %d\n", plan.Total, plan.Focused, plan.Pending)
	if plan.HasRandom {
		fmt.Println("Random seed:", plan.Seed)
	}
	if plan.HasFocus {
		fmt.Println("Focus is active.")
	}
}

// Specs is called while specs are being run.
func (Rainbow) Specs(_ *testing.T, specs <-chan spec.Spec) {
	var passed, failed, skipped int
	for s := range specs {
		switch {
		case s.Failed:
			failed++
			if !testing.Verbose() {
				fmt.Print(colorize(redColor+boldStyle, "x"))
			} else {
				if out, err := ioutil.ReadAll(s.Out); err == nil {
					fmt.Println(colorize(redColor+boldStyle, "%s", out))
				}
			}
		case s.Skipped:
			skipped++
			if !testing.Verbose() {
				fmt.Print(colorize(cyanColor+boldStyle, "s"))
			}
		default:
			passed++
			if !testing.Verbose() {
				fmt.Print(colorize(greenColor+boldStyle, denoter))
			}
		}
	}
	fmt.Printf("\n%s | %s | %s\n", colorize(greenColor+boldStyle, "Passed: %d", passed), colorize(redColor+boldStyle, "Failed: %d", failed), colorize(cyanColor+boldStyle, "Skipped: %d", skipped))
}
