package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	fcpu, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(fcpu); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	Execute()

	fmem, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(fmem); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	fmem.Close()
}
