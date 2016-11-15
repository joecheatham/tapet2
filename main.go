package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"runtime"

	"github.com/fogleman/primitive/primitive"
)

const (
	url           string = "http://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1"
	tmpImg        string = "temp.jpg"
	backgroundImg string = "background.png"
)

type shapeConfig struct {
	Count  int
	Mode   int
	Alpha  int
	Repeat int
}

func main() {
	_, height := getScreenResolution()
	res, err := http.Get(url)

	defer res.Body.Close()
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	imageUrl, err := getInputImage([]byte(body))

	response, err := http.Get(imageUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	usr, err := user.Current()
	err = os.Mkdir(fmt.Sprint(usr.HomeDir, "/.tapet2"), 0777)
	if err != nil {
		if !os.IsExist(err) {
			println("dir NOPE")
			log.Fatal(err)
		}
	}

	file, err := os.Create(fmt.Sprint(usr.HomeDir, "/.tapet2/", tmpImg))
	if err != nil {
		println("file NOPE")
		log.Fatal(err)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	input, err := primitive.LoadImage(fmt.Sprint(usr.HomeDir, "/.tapet2/", tmpImg))

	Configs := make([]shapeConfig, 0)
	//configCount := randMinMax(50, 200)
	mode := randMinMax(0, 5)

	for i := 0; i < 20; i++ {
		Configs = append(Configs, shapeConfig{i, mode, 128, 0})
	}

	bg := primitive.MakeColor(primitive.AverageImageColor(input))

	// run primitive algorithm
	model := primitive.NewModel(input, bg, height, runtime.NumCPU())
	for j, config := range Configs {
		for i := 0; i < config.Count; i++ {
			model.Step(primitive.ShapeType(config.Mode), config.Alpha, config.Repeat)
			last := j == len(Configs)-1 && i == config.Count-1
			if last {
				primitive.SavePNG(fmt.Sprint(usr.HomeDir, "/.tapet2/", backgroundImg), model.Context.Image())
			}
		}
	}

	changeDesktopBackground(fmt.Sprint(usr.HomeDir, "/.tapet2/", backgroundImg))
}
