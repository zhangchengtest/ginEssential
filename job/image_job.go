package job

import (
	"fmt"
	"ginEssential/config"
	"ginEssential/model"
	"ginEssential/util"
	"github.com/hbagdi/go-unsplash/unsplash"
	strftime "github.com/itchyny/timefmt-go"
	"github.com/zhangchengtest/simple/sqls"
	"golang.org/x/oauth2"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

func getImage(filename string, dir string) (image.Image, error) {
	f, err := os.Open(dir + "/" + filename)
	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func randomImage(filename string, dir string) {

	//pqIb2UEnFKOc3p-G51tm9JrznDqrTDmdWun7CuM-lFk

	ts := oauth2.StaticTokenSource(
		// note Client-ID in front of the access token
		&oauth2.Token{AccessToken: "pqIb2UEnFKOc3p-G51tm9JrznDqrTDmdWun7CuM-lFk"},
	)
	client := oauth2.NewClient(oauth2.NoContext, ts)
	//use the http.Client to instantiate unsplash
	unsplash2 := unsplash.New(client)

	var opt = &unsplash.PhotoOpt{Height: 100, Width: 100}
	//opt.Crop = true
	tt, resp, err := unsplash2.Photos.Photo("random", opt)
	if err != nil {
		panic(err)
	}
	log.Printf("photos : %v", tt)
	log.Printf("resp : %v", resp)
	log.Printf("tt : %v", tt.Urls.Raw.String())
	log.Printf("tt : %v", tt.Urls.Full.String())
	log.Printf("tt : %v", tt.Urls.Full.String()+"&w=100&h=100")

	err2 := downLoad(dir, filename, tt.Urls.Full.String()+"&fit=crop&w=900&h=900")
	if err2 != nil {
		fmt.Println("Download pic file failed!", err2)
	} else {
		fmt.Println("Download file success.")
	}
}

// 下载图片信息
func downLoad(base string, filename string, url string) error {
	pic := base
	pic += "/" + filename
	v, err := http.Get(url)
	if err != nil {
		fmt.Printf("Http get [%v] failed! %v", url, err)
		return err
	}
	defer v.Body.Close()
	content, err := ioutil.ReadAll(v.Body)
	if err != nil {
		fmt.Printf("Read http response failed! %v", err)
		return err
	}
	err = ioutil.WriteFile(pic, content, 0666)
	if err != nil {
		fmt.Printf("Save to file failed! %v", err)
		return err
	}
	return nil
}

func SplitImageWithIterator() {

	filename := "full.jpg"
	t := time.Now()
	dir2 := strftime.Format(t, "%Y%m%d%H%M%S")

	_, err := os.Stat(config.Instance.Uploader.Local.Path + "/" + dir2)
	if os.IsNotExist(err) {
		os.Mkdir(config.Instance.Uploader.Local.Path+"/"+dir2, os.ModePerm)
	}

	path := config.Instance.Uploader.Local.Path + "/" + dir2

	randomImage(filename, path)
	in, err := getImage(filename, path)
	if err != nil {
		log.Fatal("failed to load test image: %s", err)
	}

	// set cfg
	cfg := util.Config{X: 3, Y: 3}

	// split
	it, err := util.SplitImageWithIterator(in, cfg)
	if err != nil {
		log.Fatal("unexpected failure: %s", err)
	}

	// drain the parts to check expected count
	outs := util.ConsumeIterator(it)

	for i, tt := range outs {
		savePuzzle(tt, filename, i, dir2)
	}

}

func savePuzzle(img image.Image, ret string, sort int, dir string) {

	//bounds := m.Bounds()
	//fmt.Println(bounds, formatString)

	//Encode from image format to writer
	pngFilename := config.Instance.Uploader.Local.Path + "/" + dir + "/" + "output" + strconv.Itoa(sort) + ".png"
	f, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Png file", pngFilename, "created")

	DB := sqls.DB()
	var s = util.Worker1{}
	// 创建图
	newUser := model.PuzzlePiece{
		Id:       s.GetId(),
		Content:  config.Instance.Uploader.Local.Host + "images/" + dir + "/" + "output" + strconv.Itoa(sort) + ".png",
		Title:    ret,
		Url:      config.Instance.Uploader.Local.Host + "images/" + dir + "/" + ret,
		Sort:     sort,
		CreateDt: time.Now(),
		CreateBy: "",
	}

	sysType := runtime.GOOS
	fmt.Printf(sysType)

	if sysType != "windows" {
		DB.Create(&newUser)
	}

}
