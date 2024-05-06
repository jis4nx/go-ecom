package rabbit

import "runtime"

func MaxParallelism() int {
	maxProcess := runtime.GOMAXPROCS(0)
	numCpu := runtime.NumCPU()
	if maxProcess < numCpu {
		return maxProcess
	}
	return numCpu
}
