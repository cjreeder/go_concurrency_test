package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	//	"github.com/cjreeder/via_networking_script/via"
	"github.com/spf13/pflag"
)

type ViaList struct {
	gateway_id string
	vianame    string
}

const xthreads = 750

func ReadCsv(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// read file into a variable to be able to usue later
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

func doSomething(a int, wg *sync.WaitGroup, requests <-chan ViaList) {
	defer wg.Done()
	for req := range requests {
		fmt.Printf("Doing Stuff - Job number: %v Job name: %s \n", a, req.vianame)
		time.Sleep(time.Second * 60)
		fmt.Printf("Finishing job for: %s \n", req.vianame)
	}
}

func main() {
	var (
		ifile string
		ofile string
		wg    sync.WaitGroup
	)

	pflag.StringVarP(&ifile, "input", "i", "", "Input file containing a list of VIAs to Migrate")
	pflag.StringVarP(&ofile, "output", "o", "", "file to log all output to")
	pflag.Parse()

	lines, err := ReadCsv(ifile)
	if err != nil {
		fmt.Printf("File cannot be found or read: %v", err.Error())
	}

	var requests = make(chan ViaList, 1000) // This number 50 can be anything as long as it's larger than xthreads
	// loop through the lines and turn it into an object
	a := 0
	for _, line := range lines {
		data := ViaList{
			gateway_id: line[0],
			vianame:    line[1],
		}
		fmt.Printf("Moving Data into Channel: %s \n", data.vianame)
		requests <- data
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go doSomething(a, &wg, requests)
	}
	close(requests)
	wg.Wait()
	/*
		fmt.Printf("Changing over %v\n", data.vianame)
		a++
		go func(a int, vianame string) {
			defer wg.Done()
			time.Sleep(time.Second * 60)

			fmt.Printf("Doing Stuff - Job number: %v Job name: %s \n", a, vianame)
		}(a, data.vianame)

		//err := SetNetwork(data.vianame, data.oldaddress, data.ipaddress, data.subnetmask, data.gateway, data.dns)dd(xthreads)
		/*wg.Add(xthreads)
		for i := 0; i < xthreads; i++ {
			go func() {
				for {
					a, ok := <-ch
					if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
						wg.Done()
						return
					}
					doSomething(a, data.vianame) // do the thing
				}
			}()
		}

		// Now the jobs can be added to the channel, which is used as a queue
		for i := 0; i < 50; i++ {
			ch <- i // add i to the queue
		}

		close(ch) // This tells the goroutines there's nothing else to do
	*/
	//wg.Wait() // Wait for the threads to finish}

	//if err != nil {
	/*	fmt.Printf("%v returned an error: %v\n", data.vianame, err)
		} else {
			fmt.Printf("Change over script sent to %v\n", data.vianame)
		}
		time.Sleep(120 * time.Second)
		GetNetwork(data.vianame, data.ipaddress)
		if err != nil {
			fmt.Printf("Error: %v", err)
		}*/

}
