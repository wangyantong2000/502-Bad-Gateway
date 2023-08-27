package main

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
)

// 从视频流中截取一帧并返回
func GetImageBuffer(videoPath string) (*bytes.Buffer, error) {

	imgBuffer := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(imgBuffer, os.Stdout).
		Run()

	if err != nil {
		log.Printf("截取帧失败：", err)
		return nil, err
	}
	var imgByte []byte
	imgBuffer.Write(imgByte)
	return imgBuffer, nil
}
