package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"runtime/debug"
	"time"

	"github.com/zachaller/go-errors"

	"github.com/zachaller/go-eccodes"
	cio "github.com/zachaller/go-eccodes/io"
)

func main() {
	filename := flag.String("file", "", "io path, e.g. /tmp/ARPEGE_0.1_SP1_00H12H_201709290000.grib2")

	flag.Parse()

	content, err := ioutil.ReadFile(*filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(content[0])

	m, err := cio.OpenMemory(content, len(content))
	if err != nil {
		log.Fatalf("failed to open file on file system: %s", err.Error())
	}

	memory, err := codes.OpenMemory(m)
	if err != nil {
		log.Fatalf("failed to open file: %s", err.Error())
	}

	start := time.Now()

	msg, err := memory.GetSingleMessage()
	if err != nil {
		log.Println(err)
	}

	log.Printf("============= BEGIN MESSAGE N%d ==========\n", 1)

	shortName, err := msg.GetString("shortName")
	if err != nil {
		log.Println(errors.Wrap(err, "failed to get 'shortName' value"))
	}
	name, err := msg.GetString("name")
	if err != nil {
		log.Println(errors.Wrap(err, "failed to get 'name' value"))
	}

	log.Printf("Variable = [%s](%s)\n", shortName, name)

	// just to measure timing
	lat, long, value, err := msg.Data()

	fmt.Println(lat[0], long[0], value[0])
	if err != nil {
		log.Println(errors.Wrap(err, "failed to get data (latitudes, longitudes, values)"))
	}

	log.Printf("elapsed=%.0f ms", time.Since(start).Seconds()*1000)
	log.Printf("============= END MESSAGE N%d ============\n\n", 1)

	//Close to keep memory from leaking
	m.Close()
	memory.Close()
	msg.Close()

	debug.FreeOSMemory()
}
