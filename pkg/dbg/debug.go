package dbg

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Debugger struct {
	*zap.Logger
	ModeAST, ModeDebug bool
}

var debugger *Debugger

// Startup logs the name, version, commit, runtime, and log level of binary and
// returns a logger or exits if the logger cannot be constructed.
// <name>: <version> and pid: <process id> are set to the global logger
func NewDebugger(debug bool) *Debugger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05:")

	l := zap.Must(cfg.Build())
	zap.ReplaceGlobals(l)

	l.With(zap.String(Name(), Version()), zap.Int("pid", os.Getpid())).
		Info(l.Level().CapitalString(), zap.String("commit", Commit()),
			zap.String("runtime", runtime.Version()))

	debugger = &Debugger{
		Logger:    l,
		ModeDebug: debug,
	}
	return debugger
}

func (d *Debugger) PrintAst(args ...interface{}) {
	if !d.ModeAST {
		return
	}
	l := d.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	pp.Println(args...)

	l.Info("\n<enter> to continue...")
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(input) == "" {
			break
		}
		l.Info("\n<enter> to continue...")
	}
}

func (d *Debugger) Glamour(args ...interface{}) {
	if !d.ModeDebug {
		return
	}
	content := "```bash\n"
	for _, arg := range args {
		content += fmt.Sprintf("%s", arg)
	}
	content += "\n```"
	d.Render(content)
}

func (d *Debugger) Render(content string) {
	out, err := glamour.Render(content, "dark")
	if err != nil {
		d.Fatal("Glamour", zap.Error(err))
	}
	fmt.Print(out)
}
