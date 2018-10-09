package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/blackjack/webcam"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func camera() []byte {
	cam, err := webcam.Open("/dev/video0") // Open webcam
	if err != nil {
		panic(err.Error())
	}

	defer cam.Close()

	//	format_desc := cam.GetSupportedFormats()
	//
	//	fmt.Println("Available formats:")
	//	for a, s := range format_desc {
	//		fmt.Fprintln(os.Stderr, s)
	//		fmt.Fprintln(os.Stderr, a)
	//
	//
	//
	//		for _, b := range cam.GetSupportedFrameSizes(a) {
	//			fmt.Fprintln(os.Stderr, b.GetString())
	//		}
	//
	//	}

	cam.SetImageFormat(0x56595559, 640, 480)

	err = cam.StartStreaming()
	if err != nil {
		panic(err.Error())
	}

	for {

		err = cam.WaitForFrame(10000)

		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			fmt.Fprint(os.Stderr, err.Error())
			continue
		default:
			panic(err.Error())
		}

		frame, err := cam.ReadFrame()

		if len(frame) != 0 {

			cpBuf := make([]byte, len(frame))
			copy(cpBuf, frame)

			yuyv := image.NewYCbCr(image.Rect(0, 0, 640, 480), image.YCbCrSubsampleRatio422)
			for i := range yuyv.Cb {
				ii := i * 4
				yuyv.Y[i*2] = cpBuf[ii]
				yuyv.Y[i*2+1] = cpBuf[ii+2]
				yuyv.Cb[i] = cpBuf[ii+1]
				yuyv.Cr[i] = cpBuf[ii+3]

			}

			buf := &bytes.Buffer{}
			if err := jpeg.Encode(buf, yuyv, nil); err != nil {
				panic(err)
			}

			return buf.Bytes()
		} else if err != nil {
			panic(err.Error())
		}
	}

}

func file() []byte {
	imgFile, err := os.Open(os.Args[1]) // a QR code image

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	fmt.Println(len(buf))
	return buf

}

func main() {

	for {
		var buf []byte = nil
		if len(os.Args) > 1 {
			buf = file()
		} else {
			buf = camera()
		}
		fmt.Println(len(buf))

		// convert the buffer bytes to base64 string
		imgBase64Str := base64.StdEncoding.EncodeToString(buf)

		//   fmt.Println(imgBase64Str)

		jsonStr := []byte(fmt.Sprintf(`{"device":"countcamera1", "readings":[{"name":"cameraeiamge", "value":"%s"}]}`, imgBase64Str))

		//   fmt.Println(len(jsonStr))
		response, err := http.Post("http://localhost:48080/api/v1/event", "application/json", bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Printf("The HTTP request failedith error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}

		if len(os.Args) == 1 {
			time.Sleep(time.Second)
		}
	}
	fmt.Println("Terminating the application...")
}
