package state

type State interface {
	Consistent() (bool, error) // Check to see if the state is consistent with the operating system's state
	Dump() ([]byte, error)     // Dump the state to a JSON byte array
	Execute() error            // Execute the state if it is not already Consistent
	Initialize() error         // Initialize the state validating loaded fields
	Meta() Metadata            // Return the state's metadata ("Name", "Type", and "state")
	Requirements() []string    // Return the state's requirements // TODO: use in ordering of the state's execution and do not execute on failure
}

type Result struct {
	Consistent bool      // The state is consistent with the operating system
	Metadata   *Metadata // The metadata of the state which returned this result
	Message    string    // A message returned by the state // TODO
}

type Metadata struct {
	Name  string // Unique name to associate with a state
	Type  string // The type of state "package", "file", etc.
	State string // The desired state "installed", "rendered", etc.
}
