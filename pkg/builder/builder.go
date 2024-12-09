package builder

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/tools/imports"
)

// SSAIn ...
type SSAIn struct {
	FuncName string `json:"funcname"`
	GcFlags  string `json:"gcflags"`
	Code     string `json:"code"`
}

// SSAOut ...
type SSAOut struct {
	ID  string `json:"build_id"`
	Msg string `json:"msg"`
}

// Build SSA IR
func Build(code string) error {
	out := SSAOut{
		ID: uuid.Must(uuid.NewUUID()).String(),
	}
	path := filepath.Join("build", "/"+out.ID)

	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create buildbox, err: %v", err)
	}

	// 2. write code
	var in = SSAIn{
		FuncName: "main",
		GcFlags:  "",
		Code:     code,
	}

	if !findSSAFunc(in.Code, in.FuncName) {
		os.Remove(path)
		return fmt.Errorf("cannot find FUNC=%s", in.FuncName)
	}

	var buildFile string
	isTest := isPackageTest(in.Code)
	if !isTest {
		buildFile = filepath.Join(path, "/main.go")
	} else {
		buildFile = filepath.Join(path, "/main_test.go")
	}

	imports, err := autoimports([]byte(in.Code))
	if err != nil {
		os.Remove(path)
		return fmt.Errorf("autoimports failed: \n%v", err)
	}

	err = os.WriteFile(buildFile, imports, os.ModePerm)
	if err != nil {
		os.Remove(path)
		return fmt.Errorf("failed to write output: \n%v", err)
	}

	// 3.2 go mod init gossa && go mod tidy
	err = initModules(path)
	if err != nil {
		os.Remove(path)
		return fmt.Errorf("go modules failed: \n%v", err)
	}

	// 3.3 GOSSAFUNC=foo go build
	outFile := filepath.Join(path, "/main.out")
	err = buildSSA(in.FuncName, in.GcFlags, outFile, buildFile, isTest)
	if err != nil {
		os.Remove(path)
		return fmt.Errorf("build ssa failed: \n%v", err)
	}

	// After successful build, parse the SSA output
	ssaFile := filepath.Join(path, "/ssa.html")
	result, err := ParseSSAOutput(ssaFile)
	if err != nil {
		return fmt.Errorf("failed to parse SSA output: %v", err)
	}
	fmt.Printf("Function: %s\n", result.FunctionName)
	for _, phase := range result.Phases {
		fmt.Printf("\nPhase: %s\n", phase.Name)
		for _, line := range phase.Code {
			if strings.TrimSpace(line) == "" {
				continue
			}
			instruction := ParseSSAInstruction(line)
			fmt.Printf("  Value: %s, Op: %s, Args: %s\n",
				instruction["value"],
				instruction["operation"],
				instruction["args"])
		}
	}
	return nil
}

/*
According to spec there are two cases we need to handle:

  - Function declaration, in the form of 'func f() {...}' , see:
    https://golang.org/ref/spec#Function_declarations

  - Function literal/anonymous function, in the form of
    'myfunc := func() {...}' or 'go func() {...}', see:
    https://golang.org/ref/spec#Function_literals

As users can use some tricks like raw string to bypass our check, we
only do check conservatively, which means it is mainly used for
preventing misspell and wrong format.

All cases:
// <i>,<j>,<k> means unique function index in the scope of outer function, see:
https://github.com/golang/go/blob/84162b88324aa7993fe4a8580a2b65c6a7055f88/src/cmd/compile/internal/typecheck/func.go#L182

- func foo()	// most common case
- glob..func<i>	// global function literal
  - glob..func<i>.<j>.<k>...		// inner anonymous function

- foo.func<i>	// anonymous function inside function 'foo'
  - foo.func<i>.<j>.<k>...

- (*T).foo()	// method expression with explicit receiver, see
https://golang.org/ref/spec#Method_expressions

Note that non-ascii letters are unsupported, as our intention is to dig
into go ssa IR.
*/
func findSSAFunc(code, funcname string) bool {
	// The dot character is not allowed to appear in function name.
	// See https://golang.org/ref/spec#Identifiers
	if strings.IndexByte(funcname, '.') != -1 {
		if funcname[0] == '(' {
			methodReg := regexp.MustCompile(`^\([\w\*]+\)\.\w+$`)
			return methodReg.MatchString(funcname)
		} else if strings.HasPrefix(funcname, "glob") {
			globReg := regexp.MustCompile(`^glob\.\.func\d+(\.\d)*$`)
			return globReg.MatchString(funcname)
		} else {
			anonyReg := regexp.MustCompile(`^\w+\.func\d+(\.\d)*$`)
			return anonyReg.MatchString(funcname)
		}
	}
	// func Foo (
	re := regexp.MustCompile(fmt.Sprintf(`func[ \t]+%s[ \t]*\(`, funcname))
	return re.FindString(code) != ""
}

func isPackageTest(code string) bool {
	// package *_test
	re := regexp.MustCompile(`package .*\_test`)
	return re.FindString(code) != ""
}

func autoimports(code []byte) ([]byte, error) {
	out, err := imports.Process("", code, &imports.Options{
		Fragment:  true,
		AllErrors: true,
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func initModules(path string) error {
	// 1. go mod init
	cmd := exec.Command("go", "mod", "init", "test")
	cmd.Dir = path
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	if err != nil {
		return errors.New(cmd.Stderr.(*bytes.Buffer).String())
	}

	// 2. go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = path
	cmd.Stderr = &bytes.Buffer{}
	err = cmd.Run()
	if err != nil {
		return errors.New(cmd.Stderr.(*bytes.Buffer).String())
	}

	return nil
}

func buildSSA(funcname, gcflags, outf, buildf string, isTest bool) error {
	var (
		cmd      *exec.Cmd
		buildDir string
	)

	// Restrict the ssa.html target to the target ssa build folder.
	// See https://github.com/golang-design/ssaplayground/issues/9
	buildDir = filepath.Dir(buildf)
	outf = filepath.Base(outf)
	buildf = filepath.Base(buildf)

	if !isTest {
		cmd = exec.Command("go", "build", "-mod=readonly", fmt.Sprintf(`-gcflags=%s`, gcflags), "-o", outf, buildf)
	} else {
		cmd = exec.Command("go", "test", "-mod=readonly", fmt.Sprintf(`-gcflags=%s`, gcflags), buildf)
	}
	cmd.Dir = buildDir
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOSSAFUNC=%s", funcname))
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	if err != nil {
		return errors.New(cmd.Stderr.(*bytes.Buffer).String())
	}
	if err != nil {
		return fmt.Errorf("failed to build SSA: %v", err)
	}
	return nil
}
