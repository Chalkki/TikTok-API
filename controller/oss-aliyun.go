package controller

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
	"os"
)

//see details in the documentations of aliyun oss api
// "https://help.aliyun.com/document_detail/32145.html"
func UploadObject(objectName string, filePath string) error {
	// create a new OSS client instance
	client, err := oss.New(Endpoint, OSSAccessKeyID, OSSAccessSecret)
	if err != nil {
		return err
	}
	// get the bucket
	bucket, err := client.Bucket(BucketName)
	if err != nil {
		return err
	}

	// file uploading
	err = bucket.PutObjectFromFile(objectName, filePath)
	if err != nil {
		return err
	}
	return nil
}

//Get the snapshot from the video using open-source ffmpeg tool
//Upload the video and snapshot stored in the public file
func UploadVideo(videoName string, snapShotName string, videoFilePath string) error {

	err := getSnapShot(snapShotName, videoFilePath)
	if err != nil {
		return err
	}
	err = UploadObject(videoName, videoFilePath)
	if err != nil {
		return err
	}
	err = UploadObject(snapShotName, "./public/"+snapShotName)
	if err != nil {
		return err
	}
	return nil
}

//see details at "https://github.com/u2takey/ffmpeg-go"
func ExampleReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal(err)
	}
	return buf
}

// get the first frame of the uploaded video
func getSnapShot(snapShotName string, videoFilePath string) error {
	//because we use "id-video.mp4" e.g. "1-video.mp4" as the video name, we need to split the videoName and
	//get the id of the video

	reader := ExampleReadFrameAsJpeg(videoFilePath, 1)
	img, err := imaging.Decode(reader)
	if err != nil {
		return err
	}
	err = imaging.Save(img, "./public/"+snapShotName)
	if err != nil {
		return err
	}
	return nil
}
