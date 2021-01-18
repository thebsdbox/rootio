package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/rootio/pkg/storage"
	"github.com/thebsdbox/rootio/pkg/types.go"
)

var metadata *types.Metadata

// Release - this struct contains the release information populated when building rootio
var Release struct {
	Version string
	Build   string
}

var rootioCmd = &cobra.Command{
	Use:   "rootio",
	Short: "This is a tool for managing storage for bare-metal servers",
}

func init() {

	rootioCmd.AddCommand(rootioFormat)
	rootioCmd.AddCommand(rootioPartition)
	rootioCmd.AddCommand(rootioVersion)

	// Find configuration
	var err error
	if os.Getenv("TEST") != "" {
		// TEST MODE
		metadata, err = test()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		metadata, err = types.RetreieveData()
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Succesfully parsed the MetaData, Found [%d] Disks\n", len(metadata.Storage.Disks))

}

// Execute - starts the command parsing process
func Execute() {
	if err := rootioCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootioFormat = &cobra.Command{
	Use:   "format",
	Short: "Use rootio to format disks based upon metadata",
	Run: func(cmd *cobra.Command, args []string) {

		for fileSystem := range metadata.Storage.Filesystems {
			err := storage.FileSystemCreate(metadata.Storage.Filesystems[fileSystem])
			if err != nil {
				log.Error(err)
			}
		}
	},
}

var rootioPartition = &cobra.Command{
	Use:   "partition",
	Short: "Use rootio to partition disks based upon metadata",
	Run: func(cmd *cobra.Command, args []string) {

		for disk := range metadata.Storage.Disks {
			err := storage.VerifyBlockDevice(metadata.Storage.Disks[disk].Device)
			if err != nil {
				log.Error(err)
			}
			err = storage.ExamineDisk(metadata.Storage.Disks[disk])
			if err != nil {
				log.Error(err)
			}

			if metadata.Storage.Disks[disk].WipeTable {
				err = storage.Wipe(metadata.Storage.Disks[disk])
				log.Infoln("Wiping")
				if err != nil {
					log.Error(err)
				}
			}
			log.Infoln("Partitioning")
			err = storage.Partition(metadata.Storage.Disks[disk])
			if err != nil {
				log.Error(err)
			}
		}
	},
}

// var bootyPush = &cobra.Command{
// 	Use:   "push",
// 	Short: "This is will direct BOOTy to push the contents of a disk to a remote server",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		push.Image()
// 	},
// }

// var bootyServer = &cobra.Command{
// 	Use:   "server",
// 	Short: "This is for starting BOOTy as a simple (test) web server",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		server.Serve()
// 	},
// }

var rootioVersion = &cobra.Command{
	Use:   "version",
	Short: "Version and Release information about the rootio storage manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("rootio Release Information\n")
		fmt.Printf("Version:  %s\n", Release.Version)
		fmt.Printf("Build:    %s\n", Release.Build)
	},
}

func test() (*types.Metadata, error) {
	// Open our jsonFile
	jsonFile, err := os.Open("./test/test.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully Opened test data")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var mdata types.Metadata

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &mdata)

	return &mdata, nil
}
