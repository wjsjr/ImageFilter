package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"proj2/png"
	"sync"
)

//Worker context
type bspWorkerContext struct {
	inPath               string
	reader               *json.Decoder
	imageTaskPipeline    [][]png.ImageTask
	imageResultsPipeline [][]png.ImageTask
	superStep            string
	barriers             []barrier
	numThreads           int
}

//Barrier struct
type barrier struct {
	capacity int
	size     int
	l        *sync.Mutex
	con      *sync.Cond
}

//Helper function which initalizes 2d array, with each thread having its own array to hold its tasks
func initPipe(numThreads int) [][]png.ImageTask {
	taskPipe := [][]png.ImageTask{}
	for i := 0; i < numThreads; i++ {
		taskPipe = append(taskPipe, []png.ImageTask{})
	}
	return taskPipe
}

//Create array of 3 single use barriers for 3 super steps
func createBarriers(capacity int) []barrier {
	return []barrier{newBarrier(capacity), newBarrier(capacity), newBarrier(capacity)}
}

//Await function which causes thread to wait at barrier until all threads have reached barrier
func (b *barrier) await() {
	b.l.Lock()
	b.size++
	//Wait until size of barrier = capacity of barrier
	if b.size == b.capacity {
		b.con.Broadcast()
	} else {
		for b.size < b.capacity {
			b.con.Wait()
		}
	}
	b.l.Unlock()

}

//Create new barrier with given capacity
func newBarrier(capacity int) barrier {
	l := new(sync.Mutex)
	return barrier{capacity, 0, l, sync.NewCond(l)}
}

//Initialize BSP context
func NewBSPContext(config Config) *bspWorkerContext {
	basePath := "../data/"
	inPath := basePath + "in/" + config.DataDirs + "/"
	effectsPathFile := fmt.Sprintf(basePath + "effects.txt")
	effectsFile, _ := os.Open(effectsPathFile)
	reader := json.NewDecoder(effectsFile)

	//Make sure image task pipeline should be buffered
	return &bspWorkerContext{inPath, reader, initPipe(config.ThreadCount), initPipe(config.ThreadCount), "reading", createBarriers(config.ThreadCount), config.ThreadCount}
}

//Edit Super Step. While imageResultsPipeline is not empty, pop a task off the task pipeline, and apply the first effect
//If no more effects left, push to results pipeline, else push task back onto task pipeline
func edit(id int, ctx *bspWorkerContext) {
	task := popTask(id, ctx)
	if task != nil {
		if len(task.Effects) == 0 {
			task.Image.SwapOut()
			appendResult(id, *task, ctx)

		} else {
			effect := task.Effects[0]
			png.ApplyEffect(task.Image, effect)
			if len(task.Effects) > 1 {
				task.Effects = task.Effects[1:]
				task.Image.SwapIn()
				appendTask(id, *task, ctx)
			} else {
				appendResult(id, *task, ctx)
			}

		}
	} else {
		ctx.barriers[1].await()
		ctx.superStep = "writing"

	}

}

//Write super step, Try to pop a result off results pipeline. If succesful, write result and increment numWritten
//Else, wait to exit
func write(id int, ctx *bspWorkerContext) {
	result := popResult(id, ctx)
	if result != nil {
		result.Save()
	} else {
		ctx.barriers[2].await()
		ctx.superStep = "done"
	}

}

//Helper function to read in effects JSON, build image tasks and populate image task pipeline
func read(ctx *bspWorkerContext) {
	id := 0
	for ctx.reader.More() {
		var rawTask png.RawTask
		if err := ctx.reader.Decode(&rawTask); err != nil {
			os.Exit(0)
		}
		image, err := png.Load(ctx.inPath + rawTask.InPath)
		if err != nil {
			fmt.Println("ERROR READING IMAGE")
			os.Exit(1)
		}
		imageTask := png.BuildImageTask(rawTask, image)
		appendTask(id, imageTask, ctx)
		id++
		if id == ctx.numThreads {
			id = 0
		}
	}
	ctx.superStep = "editing"
}

//Loop until all super steps have been completed, then break out
func RunBSPWorker(id int, ctx *bspWorkerContext) {
	for {
		if ctx.superStep == "reading" {
			if id == 0 {
				read(ctx)
			}
			ctx.barriers[0].await()
		} else if ctx.superStep == "editing" {
			edit(id, ctx)
		} else if ctx.superStep == "writing" {
			write(id, ctx)
		} else {
			break
		}

	}

}

//Add result to results pipe
func appendResult(id int, task png.ImageTask, ctx *bspWorkerContext) {
	ctx.imageResultsPipeline[id] = append(ctx.imageResultsPipeline[id], task)
}

//Add task to task pipe
func appendTask(id int, task png.ImageTask, ctx *bspWorkerContext) {
	ctx.imageTaskPipeline[id] = append(ctx.imageTaskPipeline[id], task)
}

//Pop task off task pipe
func popTask(id int, ctx *bspWorkerContext) *png.ImageTask {
	var task png.ImageTask
	if len(ctx.imageTaskPipeline[id]) > 0 {
		task, ctx.imageTaskPipeline[id] = ctx.imageTaskPipeline[id][0], ctx.imageTaskPipeline[id][1:]
		return &task
	} else {
		return nil
	}
}

//Pop result off results pipe
func popResult(id int, ctx *bspWorkerContext) *png.ImageTask {
	var result png.ImageTask
	if len(ctx.imageResultsPipeline[id]) > 0 {
		result, ctx.imageResultsPipeline[id] = ctx.imageResultsPipeline[id][0], ctx.imageResultsPipeline[id][1:]
		return &result
	} else {
		return nil
	}
}
