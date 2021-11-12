/**
 * Template golang project
 *
 * This sample will show how WaitGroups can be used to ensure child threads have completed.
 * We will also see how channels can be used to communicate between threads.
 */
package main

// We are going to require some builtin features - namely, text formatting, and thread synchronization.
// Oh, let's also provide an example of using plugins!
import (
	"fmt"
	"plugin"
	"sync"
)

// Demonstrating a template dependency requires importing it first
import "github.com/BrianHannay/golang-template-example/dependency"

// the main function is the entrypoint to the compiled go program
func main() {

	var funcName = "Print"
	fmt.Printf("Running a plugin's %s method\n", funcName)
	pluginExample(funcName)

	// Print 5 times using goroutines and a wait group:
	fmt.Println("Threading 5 Print calls in goroutines")
	goroutineExample()

	// Our sample dependency defines a print method. Let's defer printing 5 hello worlds:
	fmt.Println("Deferring 5 Print calls")
	deferredExample()

}

// This shows how to load plugins and access data and methods contained within
func pluginExample(functionName string) {

	// Load the plugin by file path
	plugged, err := plugin.Open("./plugins/example_plugin.so")
	if err != nil {
		fmt.Printf("Failed to load plugin plugins/example_plugin.so: %+v\n", err)
		panic(err)
	}

	// Find a public symbol
	loaded, err := plugged.Lookup(functionName)
	if err != nil {
		fmt.Printf("Failed to lookup symbol with name %s\n: %+v\n", functionName, err)
		panic(err)
	}

	// Cast loaded symbol to a function, and call it.
	loaded.(func())()
}

// This is an example of using WaitGroups to synchronize goroutine calls with their
// parent thread, ensuring completion before main thread sends an exit signal
func goroutineExample() {
	print_group := &sync.WaitGroup{}

	fmt.Println("Creating 5 goroutine threads (numbered 1-5), then waiting for them to complete in nondeterministic order:")
	for i := 1; i <= 5; i++ {

		//Register a new entry in the WaitGroup for each goroutine we run below.
		print_group.Add(1)

		// Let's use goroutines to insert these numbers in nondeterministic positions
		go dependency.Print(
			fmt.Sprintf("Hello World (Goroutine plus waitgroup %d)!", i),
			print_group,
		)
	}
	print_group.Wait()
	fmt.Println("All goroutine threads have completed")
}

// This is an example of using the defer keyword to stack up function calls
func deferredExample() {

	// Produce 5 deferred Print calls
	for i := 5; i <= 10; i++ {
		defer dependency.Print(
			fmt.Sprintf(
				"Hello World (Deferred Thread %d)!",
				i,
			),
			// Since these are deferred, no need to pass a WaitGroup and synchronize:
			nil,
		)
	}
	fmt.Println("Deferred functions (numbered 5-10) should now (upon parent function completion) execute in reverse order:")
}
