/*
Executors accept State objects to validate that the operating system's state is consistent with what is
described in the loaded object as well as executing that state against the OS.
 */

package executors

type Executor interface {
	Consistent() (bool, error) // The operating system state is consistent with the loaded state object
	Execute() (Result, error)  // Apply the loaded state object to the operating system
	Load([]byte) error         // Load a state object into the executor
}

type Result struct{}
