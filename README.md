# Golang Queue Manager

Golang Queue Manager package limits the number of goroutines that are allowed to run concurrently.

### Installation

Run the following command to install this package:

```
$ go get -u github.com/MSK998/goq
```

### Example

```go
package main

import (
    "fmt"
    "github.com/MSK998/goq"
    "time"
)

func main()  {
    // Limit 3 goroutines to run concurrently.
    c := goq.New(3)
    
    for i := 1; i <= 10; i++ {
    	
        // This function has to call before any goroutine
        c.Wait()
        
        go func(i int) {
            fmt.Printf("Job %d is running\n", i)
            time.Sleep(2 * time.Second)
            
            // This function has to when a goroutine has finished
            // Or you can use `defer c.Done()` at the top of goroutine.
            c.Done()
        }(i)
    }
    
    // This function has to call to ensure all goroutines have finished 
    // after close the main program.
    c.WaitAllDone()
}
```

### List of supported functions

```go
package main

import "github.com/MSK998/goq"

func main() {
    // Create the concurrency manager
    // The first argument is the maximum number of goroutines to run concurrently.
    c := goq.New(10)
    
    // Wait until a slot is available for the new goroutine.
    c.Wait()
    
    // Mark a goroutine as finished
    c.Done()
    
    // Wait for all goroutines are done
    c.WaitAllDone()
    
    // Close the manager manually
    c.Close()
    
    // Returns the number of goroutines which are running
    c.RunningCount()
}
```