package main

import (
	"flag"
	"fmt"
	"io"
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
	defer m.Close()

	memory, err := codes.OpenMemory(m)
	if err != nil {
		log.Fatalf("failed to open file: %s", err.Error())
	}
	defer memory.Close()

	//native.Ccodes_handle_new_from_message_copy(native.DefaultContext, content)

	n := 0
	for {
		err = process(memory, n)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("failed to get message (#%d) from index: %s", n, err.Error())
		}
		n++
	}
}

func process(memory codes.Memory, n int) error {
	start := time.Now()

	msg, err := memory.Next()
	if err != nil {
		return err
	}
	defer msg.Close()

	log.Printf("============= BEGIN MESSAGE N%d ==========\n", n)

	shortName, err := msg.GetString("shortName")
	if err != nil {
		return errors.Wrap(err, "failed to get 'shortName' value")
	}
	name, err := msg.GetString("name")
	if err != nil {
		return errors.Wrap(err, "failed to get 'name' value")
	}

	log.Printf("Variable = [%s](%s)\n", shortName, name)

	// just to measure timing
	_, _, _, err = msg.Data()
	if err != nil {
		return errors.Wrap(err, "failed to get data (latitudes, longitudes, values)")
	}

	log.Printf("elapsed=%.0f ms", time.Since(start).Seconds()*1000)
	log.Printf("============= END MESSAGE N%d ============\n\n", n)

	debug.FreeOSMemory()

	return nil
}
