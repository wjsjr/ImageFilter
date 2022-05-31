package png

import (
	"fmt"
	"os"
)

//RawTask struct used for reading in from JSON
type RawTask struct {
	InPath  string
	OutPath string
	Effects []string
}

//Image Task adds PNG.Image object to Raw Task and is used to apply effects and write image
type ImageTask struct {
	InPath  string
	OutPath string
	Effects []string
	Image   *Image
}

//Build Image Task from raw task
func BuildImageTask(r RawTask, img *Image) ImageTask {
	return ImageTask{r.InPath, r.OutPath, r.Effects, img}
}

//Applies all effects for given image task.
func (task *ImageTask) ApplyEffects() {
	//Effects is empty, pass img.In to img.Out and exit
	if len(task.Effects) == 0 {
		task.Image.SwapOut()
	} else {
		ApplyEffect(task.Image, task.Effects[0])
		var effect string
		for i := 1; i < len(task.Effects); i++ {
			task.Image.SwapIn()
			effect = task.Effects[i]
			ApplyEffect(task.Image, effect)
		}
	}
}

func (task *ImageTask) Save() {
	basePath := "../data/out/"
	err := task.Image.Save(basePath + task.OutPath)
	if err != nil {
		fmt.Println("ERROR WRITING IMAGE")
		os.Exit(0)
	}
}
