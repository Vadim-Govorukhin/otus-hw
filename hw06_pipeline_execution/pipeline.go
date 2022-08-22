package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)

	if len(stages) == 0 {
		close(out)
		return out
	}

	outStage := in
	for _, stage := range stages {
		outStage = stage(outStage)
	}

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-outStage:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()

	return out
}
