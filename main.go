package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"time"
)

type Result struct {
	Name      string
	Available bool
	Worker    int
}

// Generates powerset from provided string
func generatePowerset(s string) []string {
	res := []string{""}

	for _, letter := range s {
		var subset []string

		for _, substr := range res {
			subset = append(subset, substr+string(letter))
		}

		res = append(res, subset...)
	}

	return res
}

func checkNames(names []string, workers int, sleep int) chan Result {
	splitted_names := split(names, workers)
	results := make(chan Result)

	for workerId, names := range splitted_names {
		go checkWorker(results, names, workerId, sleep)
	}

	return results
}

// Goroutine responsible for checking slice of names
func checkWorker(results chan Result, names []string, workerId int, sleep int) {
	for _, name := range names {
		results <- Result{
			Name:      name,
			Available: checkAvailability(name),
			Worker:    workerId,
		}
		time.Sleep(time.Duration(sleep))
	}
}

func checkAvailability(name string) bool {
	if len(name) < 4 { // automatically reject names shorter than 4
		return false
	}
	resp, err := http.Get("https://twitter.com/" + name)
	panicIf(err)
	return resp.StatusCode == http.StatusNotFound
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// Splits array of names into x equal slices.
func split(array []string, x int) [][]string {
	floatX := float64(x)
	length := len(array)
	single := int(math.Ceil(float64(length) / floatX))
	res := make([][]string, x, single)

	for i := 0; i < x; i++ {
		from := i * single
		to := int(math.Min(float64((i+1)*single), float64(length)))
		res[i] = array[from:to]
	}

	return res
}

// Listens on the `results` channel for `size` results, and displays them
func displayResults(results chan Result, size int, verbose bool) {
	for i := 0; i < size; i++ {
		result := <-results
		if verbose {
			printResult(result)
		} else {
			if result.Available {
				fmt.Printf("%s\n", result.Name)
			}
		}
	}
}

// Verbose print of a result
func printResult(result Result) {
	if result.Available {
		fmt.Printf("[%s] %s is available\n", result.Worker, result.Name)
	} else {
		fmt.Printf("[%s] %s is unavailable\n", result.Worker, result.Name)
	}
}

func main() {
	name := flag.String("name", "", "base name used for generation")
	verbose := flag.Bool("verbose", false, "output all names, even taken ones")
	workers := flag.Int("workers", 4, "number of workers querying Twitter")
	sleep := flag.Int("sleep", 500, "sleep in ms between requests")

	flag.Parse()

	powerset := generatePowerset(*name)
	results := checkNames(powerset, *workers, *sleep)
	displayResults(results, len(powerset), *verbose)
}
