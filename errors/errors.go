package errors

import (
	"bufio"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
)

var (
	gopath       string
	templatePath string
)

var out = colorable.NewColorableStdout()

func IfErr(msg string, err error) {
	if err != nil {
		log.Fatal(out, "%s: %s \n",
			color.RedString(msg),
			zap.Error(err),
		)
	}
}

func IfNoErr(msg string, err error) {
	if err == nil {
		log.Print(out, "%s: %s \n",
			color.GreenString(msg),
			zap.Error(err),
		)
	}
}

func OK() bool {
	log.Println(out, "Is this OK? %ses/%so\n",
		color.YellowString("[y]"),
		color.CyanString("[n]"),
	)
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	return strings.Contains(strings.ToLower(scan.Text()), "y")
}

func Exit(msg interface{}) {
	log.Println(out, "Error:", msg)
	color.Red("%s", msg)
	os.Exit(1)
}
