package scorekeeper

// Stat does some math on a set of scores and returns the result as a string
// We will only have one kind of Statistic for this project (Average),
// but through this interface we could extend to other kinds easily.
type Stat interface {
	// Compute generates the statistic on a set of scores immediately.
	Compute(ss []Score) (string, error)
	// Step takes a score and, if possible, includes it in the running computation
	// For example, Average.Step() will add a score to the running average
	Step(s Score) error
	// Report returns the result of the running computation.
	Report() (string, error)
}
