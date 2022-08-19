package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/podhmo/commentof"
)

var options struct {
	IncludeTestFile bool
}

func main() {
	flag.BoolVar(&options.IncludeTestFile, "include-test-file", false, "inglude *_test.go")
	flag.Parse()

	fset := token.NewFileSet()
	for _, filename := range os.Args[1:] {
		if filename == "-" {
			continue
		}

		stat, err := os.Stat(filename)
		if err != nil {
			stdSrcFilename := filepath.Join(runtime.GOROOT(), "src", filename)
			stat, err = os.Stat(stdSrcFilename)
			if err != nil {
				if err := runModulePackage(fset, filename); err != nil {
					log.Printf("!! %+v", err)
				}
				continue
			}
			filename = stdSrcFilename
		}

		if stat.IsDir() {
			if err := runDir(fset, filename); err != nil {
				log.Printf("!! %+v", err)
			}
		} else {
			if err := runFile(fset, filename); err != nil {
				log.Printf("!! %+v", err)
			}
		}
	}
}

func runDir(fset *token.FileSet, dirname string) error {
	filter := func(info fs.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}
	if options.IncludeTestFile {
		filter = nil
	}
	tree, err := parser.ParseDir(fset, dirname, filter, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse dir: %w", err)
	}

	names := make([]string, 0, len(tree))
	for name := range tree {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		p := tree[name]
		result, err := commentof.Package(fset, p)
		if err != nil {
			return fmt.Errorf("collect: dir=%s, name=%s, %w", dirname, name, err)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "	")
		if err := enc.Encode(result); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
	}
	return nil
}

func runFile(fset *token.FileSet, filename string) error {
	tree, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	result, err := commentof.File(fset, tree)
	if err != nil {
		return fmt.Errorf("collect: file=%s, %w", filename, err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "	")
	if err := enc.Encode(result); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func runModulePackage(fset *token.FileSet, pkgpath string) error {
	goModFile, err := goMod()
	if err != nil {
		return err
	}

	// Try to download dependencies that are not in the module cache in order to
	// to show their documentation.
	// This may fail if module downloading is disallowed (GOPROXY=off) or due to
	// limited connectivity, in which case we print errors to stderr and show
	// documentation only for packages that are available.
	fillModuleCache(os.Stderr, goModFile)

	mods, err := buildList(goModFile)
	if err != nil {
		return err
	}
	for _, mod := range mods {
		if mod.Path == pkgpath {
			return runDir(fset, mod.Dir)
		}
	}
	return fmt.Errorf("package %s is not found", pkgpath)
}

// from: golang.org/x/tools/cmd/godoc/main.go

// goMod returns the go env GOMOD value in the current directory
// by invoking the go command.
//
// GOMOD is documented at https://golang.org/cmd/go/#hdr-Environment_variables:
//
//	The absolute path to the go.mod of the main module,
//	or the empty string if not using modules.
func goMod() (string, error) {
	out, err := exec.Command("go", "env", "-json", "GOMOD").Output()
	if ee := (*exec.ExitError)(nil); errors.As(err, &ee) {
		return "", fmt.Errorf("go command exited unsuccessfully: %v\n%s", ee.ProcessState.String(), ee.Stderr)
	} else if err != nil {
		return "", err
	}
	var env struct {
		GoMod string
	}
	err = json.Unmarshal(out, &env)
	if err != nil {
		return "", err
	}
	return env.GoMod, nil
}

// fillModuleCache does a best-effort attempt to fill the module cache
// with all dependencies of the main module in the current directory
// by invoking the go command. Module download logs are streamed to w.
// If there are any problems encountered, they are also written to w.
// It should only be used in module mode, when vendor mode isn't on.
//
// See https://golang.org/cmd/go/#hdr-Download_modules_to_local_cache.
func fillModuleCache(w io.Writer, goMod string) {
	if goMod == os.DevNull {
		// No module requirements, nothing to do.
		return
	}

	cmd := exec.Command("go", "mod", "download", "-json")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = w
	err := cmd.Run()
	if ee := (*exec.ExitError)(nil); errors.As(err, &ee) && ee.ExitCode() == 1 {
		// Exit code 1 from this command means there were some
		// non-empty Error values in the output. Print them to w.
		fmt.Fprintf(w, "documentation for some packages is not shown:\n")
		for dec := json.NewDecoder(&out); ; {
			var m struct {
				Path    string // Module path.
				Version string // Module version.
				Error   string // Error loading module.
			}
			err := dec.Decode(&m)
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Fprintf(w, "error decoding JSON object from go mod download -json: %v\n", err)
				continue
			}
			if m.Error == "" {
				continue
			}
			fmt.Fprintf(w, "\tmodule %s@%s is not in the module cache and there was a problem downloading it: %s\n", m.Path, m.Version, m.Error)
		}
	} else if err != nil {
		fmt.Fprintf(w, "there was a problem filling module cache: %v\n", err)
	}
}

type mod struct {
	Path string // Module path.
	Dir  string // Directory holding files for this module, if any.
}

// buildList determines the build list in the current directory
// by invoking the go command. It should only be used in module mode,
// when vendor mode isn't on.
//
// See https://golang.org/cmd/go/#hdr-The_main_module_and_the_build_list.
func buildList(goMod string) ([]mod, error) {
	if goMod == os.DevNull {
		// Empty build list.
		return nil, nil
	}

	out, err := exec.Command("go", "list", "-m", "-json", "all").Output()
	if ee := (*exec.ExitError)(nil); errors.As(err, &ee) {
		return nil, fmt.Errorf("go command exited unsuccessfully: %v\n%s", ee.ProcessState.String(), ee.Stderr)
	} else if err != nil {
		return nil, err
	}
	var mods []mod
	for dec := json.NewDecoder(bytes.NewReader(out)); ; {
		var m mod
		err := dec.Decode(&m)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		mods = append(mods, m)
	}
	return mods, nil
}
