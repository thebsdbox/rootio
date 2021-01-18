package storage

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/thebsdbox/rootio/pkg/types.go"
)

// FileSystemCreate handles the creation of filesystems
func FileSystemCreate(f types.Filesystem) error {

	// Add filesystem flags
	f.Mount.Create.Options = append(f.Mount.Create.Options, "-t")
	f.Mount.Create.Options = append(f.Mount.Create.Options, f.Mount.Format)

	// Add force
	f.Mount.Create.Options = append(f.Mount.Create.Options, "-F")

	// Add Device to formate
	f.Mount.Create.Options = append(f.Mount.Create.Options, f.Mount.Device)

	// Format disk
	cmd := exec.Command("/sbin/mke2fs", f.Mount.Create.Options...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	var debugCMD string
	for i := range f.Mount.Create.Options {
		debugCMD = fmt.Sprintf("%s %s", debugCMD, f.Mount.Create.Options[i])
	}
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Command [%s] Filesystem [%v]", debugCMD, err)
	}
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("Command [%s] Filesystem [%v]", debugCMD, err)
	}

	return nil
}
