package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"proj2/png"
	//"proj2/png"
)

func RunSequential(config Config) {
	//Set up reader
	basePath := "../data/"
	inPath := basePath + "in/" + config.DataDirs + "/"
	effectsPathFile := fmt.Sprintf(basePath + "effects.txt")
	effectsFile, _ := os.Open(effectsPathFile)
	reader := json.NewDecoder(effectsFile)

	var rawTask png.RawTask

	//Read in imageTasks from reader one by one
	for reader.More() {

		//Read in "Raw Task" from effects
		if err := reader.Decode(&rawTask); err != nil {
			os.Exit(0)
		}

		//Load image
		image, err := png.Load(inPath + rawTask.InPath)

		if err != nil {
			fmt.Println("ERROR READING IMAGE")
			os.Exit(0)
		}

		//Build imagetask with loaded Image object, effects and write path
		task := png.BuildImageTask(rawTask, image)

		task.ApplyEffects()
		task.Save()

	}

}
