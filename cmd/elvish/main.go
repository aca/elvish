// Elvish is a cross-platform shell, supporting Linux, BSDs and Windows. It
// features an expressive programming language, with features like namespacing
// and anonymous functions, and a fully programmable user interface with
// friendly defaults. It is suitable for both interactive use and scripting.
package main

import (
	"fmt"
	"os"

	"src.elv.sh/pkg/buildinfo"
	"src.elv.sh/pkg/daemon"
	"src.elv.sh/pkg/lsp"
	"src.elv.sh/pkg/prog"
	"src.elv.sh/pkg/shell"
)

func main() {
	fmt.Println("flaked2")
	os.Exit(prog.Run(
		[3]*os.File{os.Stdin, os.Stdout, os.Stderr}, os.Args,
		prog.Composite(
			&buildinfo.Program{}, &daemon.Program{}, &lsp.Program{},
			&shell.Program{ActivateDaemon: daemon.Activate})))
}
