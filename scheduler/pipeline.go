package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"proj2/png"
)

func RunPipeline(config Config) {
	//Initialize Channels
	imgTaskChannel := make(chan png.ImageTask, 1000)
	imgResultsChannel := make(chan png.ImageTask, 1000)
	allDone := make(chan bool, 2)
	workersFinished := make(chan bool, config.ThreadCount)

	//Create slice to hold worker channels
	workerChans := []chan bool{}

	//Launch all workers, not each worker has its own channel, as well as access to global channels
	for i := 0; i < config.ThreadCount; i++ {
		ch := make(chan bool)
		workerChans = append(workerChans, ch)
		go worker(ch, workersFinished, imgTaskChannel, imgResultsChannel, config.ThreadCount)
	}

	//Launch Thread to aggregate image results
	go ImageResultsAggregator(allDone, imgResultsChannel, workersFinished, config.ThreadCount)

	//Launch Thread to generate tasks
	go ImageTaskGenerator(workerChans, imgTaskChannel, config)

	//Wait for allDone signal to exit
	for {
		select {
		case _ = <-allDone:
			return
		}
	}

}

/*
NOTE: WAS UNABLE TO COME UP WITH WORKING SOLUTION USING MINI-WORKERS. EACH INDIVIDUAL WORKER HANDLES AN ENTIRE SINGLE IMAGE
*/
func worker(taskChannelClosed <-chan bool, workersFinished chan<- bool, imgTaskChannel <-chan png.ImageTask, imgResultsChannel chan<- png.ImageTask, threads int) {
	var task png.ImageTask
	done := false
	for {
		select {
		//Look for tasks coming in on task channel or signal saying task channel is closed
		case task = <-imgTaskChannel:
			//Pop task off task channel, apply effects to it and then add it to results channel
			task.ApplyEffects()
			imgResultsChannel <- task
		case _ = <-taskChannelClosed:
			//If signal comes from task generator indicating that task channel is closed, mark done flag as true
			done = true

		default:
			//If done flag marked as true and nothing in image task channel, tell results aggregator that a worker finished and return
			if done {
				workersFinished <- true
				return
			}
		}

	}

}

//Launches thread that aggregates image results
func ImageResultsAggregator(allDone chan<- bool, imgResultsChannel <-chan png.ImageTask, workersFinished <-chan bool, numThreads int) {
	var imageResult png.ImageTask
	numWorkersDone := 0
	done := false

	for {
		select {

		//Save all image results that come in
		case imageResult = <-imgResultsChannel:
			imageResult.Save()

		//Keep track of the number of workers which have finished
		case _ = <-workersFinished:
			numWorkersDone++

			//Once all workers finished, mark done flag as true
			if numWorkersDone == numThreads {
				done = true
			}

		//If no more image results in channel, and done flag marked exit, and signal to main thread that system is all done
		default:
			if done {
				allDone <- true
				return
			}

		}

	}

}

//Reads through effects json, builds image tasks and passes them to image task channel
func ImageTaskGenerator(workerChans []chan bool, imgTaskChannel chan<- png.ImageTask, config Config) {
	basePath := "../data/"
	inPath := basePath + "in/" + config.DataDirs + "/"
	effectsPathFile := fmt.Sprintf(basePath + "effects.txt")
	effectsFile, epfErr := os.Open(effectsPathFile)
	if epfErr != nil {
		fmt.Println("FILE ERROR")
		os.Exit(1)
	}
	reader := json.NewDecoder(effectsFile)

	for reader.More() {
		var rawTask png.RawTask
		if err := reader.Decode(&rawTask); err != nil {
			os.Exit(0)
		}

		image, err := png.Load(inPath + rawTask.InPath)

		if err != nil {
			fmt.Println("ERROR READING IMAGE")
			os.Exit(0)
		}

		imgTask := png.BuildImageTask(rawTask, image)
		imgTaskChannel <- imgTask

	}

	//Tell each worker that tasks have been added to channel
	for i := 0; i < config.ThreadCount; i++ {
		workerChans[i] <- true
	}

}
