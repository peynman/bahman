package main

import (
	"fmt"
	"os"
	"path/filepath"
	"path"
	"strings"
	"strconv"
	"time"
	"os/exec"
)

var (
	GOPath = ""
	Log = true
	DebugMode = false
	Version = "0.1-alpha"
	Code = 1
	Platform = "OSX"
	Dependencies = []string {
		"",
	}
	ServerEntry = "app/server/server.go"
	ConsoleEntry = "app/cli/cli.go"
	AppModules = "app/modules/*/main/*.go"
	LoggerModules = "app/modules/services/logger/*/main/*.go"
	ConsoleModules = "app/cli/modules/*/main/*.go"
	AppPackage = "github.com/peyman-abdi/bahman/app/modules/services/app"
	LDFLags = "-ldflags \"" + packageConstants() + "\""
	DEBUGFLAGS = "-gcflags='-N -l'"
)

func main() {
	fmt.Println("Bahman Compiler")
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please choose a compile mode: dev, dist, server, cli, modules, front")
		return
	}

	fmt.Println("Found", len(scan(AppModules)), "App Module(s)")
	fmt.Println("Found", len(scan(LoggerModules)), "Logger Module(s)")
	fmt.Println("Found", len(scan(ConsoleModules)), "Console Module(s)")

	switch args[1] {
	case "dev":
		DebugMode = true
		buildServer()
	}
}
func exeName(path string) string {
	outputName := filepath.Base(ServerEntry)
	return outputName[:strings.Index(outputName,".go")]
}
func buildServer() {
	buildExec(DebugMode, ServerEntry, abs(path.Join("bin", Platform, variant(), exeName(ServerEntry))))
}
func variant() string {
	if DebugMode {
		return "DEBUG"
	}
	return "PRODUCTION"
}
func packageConstants() string {
	return " -X "+AppPackage+".Version="+Version+
		" -X "+AppPackage+".Platform="+Platform+
		" -X "+AppPackage+".Code="+strconv.Itoa(Code)+
		" -X "+AppPackage+".Variant="+variant()+
		" -X "+AppPackage+".BuildTime="+time.Now().String()
}

func root() string {
	root, err := os.Getwd()
	if err != nil { panic(err) }
	return root
}
func abs(p string) string {
	return path.Join(root(), p)
}
func scan(match string) []string {
	var founds []string
	matchConst := match[0:strings.Index(match, "*")]
	searchRoot := abs(matchConst)
	err := filepath.Walk(searchRoot, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			m := match[strings.Index(match, "*"):]
			p := path[len(searchRoot)+1:]
			matched, _ := filepath.Match(m, p)
			if matched {
				founds = append(founds, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return founds
}
func buildExec(debug bool, entry string, output string) {
	//debugFlags := ""
	//if debug { debugFlags = DEBUGFLAGS }
	cmds := []string {
		"go", "build", "-o", output, entry,
	}
	fmt.Println(cmds)
	exe := exec.Command(cmds[0], cmds[1:]...)
	out, err := exe.CombinedOutput()
	if err != nil {
		println(err.Error())
		return
	}
	print(string(out))

}