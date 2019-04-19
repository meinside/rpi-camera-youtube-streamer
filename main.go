package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	configFilename = "config.json"

	raspividBinPath = "/usr/bin/raspivid"
	ffmpegBinPath   = "/usr/local/bin/ffmpeg"
)

type config struct {
	YoutubeStreamKey string `json:"youtube_stream_key"`

	VideoWidth    int    `json:"video_width"`
	VideoHeight   int    `json:"video_height"`
	VideoRotation int    `json:"video_rotation"`
	VideoExposure string `json:"video_exposure"`
	VideoAWB      string `json:"video_awb"`

	IsVerbose bool `json:"is_verbose"`
}

var youtubeStreamKey string
var videoWidth int
var videoHeight int
var videoRotation int
var videoExposure string
var videoAWB string
var isVerbose bool

// Read config
func getConfig() (conf config, err error) {
	var execFilepath string
	if execFilepath, err = os.Executable(); err == nil {
		var file []byte
		if file, err = ioutil.ReadFile(filepath.Join(filepath.Dir(execFilepath), configFilename)); err == nil {
			if err = json.Unmarshal(file, &conf); err == nil {
				return conf, nil
			}
		}
	}
	return config{}, err
}

func init() {
	// read config
	if conf, err := getConfig(); err != nil {
		panic(err)
	} else {
		youtubeStreamKey = conf.YoutubeStreamKey
		videoWidth = conf.VideoWidth
		videoHeight = conf.VideoHeight
		videoRotation = conf.VideoRotation
		videoExposure = conf.VideoExposure
		videoAWB = conf.VideoAWB
		isVerbose = conf.IsVerbose
	}
}

func main() {
	// http://maxogden.com/hd-live-streaming-cats.html
	raspividArgs := []string{
		"-o", "-", // to STDIO
		"-t", "0", // no timeout
		"-fps", "30", // 30 frames/sec
		"-b", "6000000", // 6 MBits

		// read from config
		"-w", strconv.Itoa(videoWidth), // width
		"-h", strconv.Itoa(videoHeight), // height
		"-rot", strconv.Itoa(videoRotation), // rotation
		"-ex", videoExposure, // exposure
		"-awb", videoAWB, // white balance
	}
	ffmpegArgs := []string{
		"-re",
		"-ar", "44100",
		"-ac", "2",
		"-acodec", "pcm_s16le",
		"-f", "s16le",
		"-ac", "2",
		"-i", "/dev/zero",
		"-f", "h264",
		"-i", "-",
		"-vcodec", "copy",
		"-acodec", "aac",
		"-ab", "128k",
		"-g", "60", // number of frames per keyframe
		"-strict", "experimental",
		"-f", "flv",
		"rtmp://a.rtmp.youtube.com/live2/" + youtubeStreamKey,
	}

	// commands
	raspivid := exec.Command(raspividBinPath, raspividArgs...)
	ffmpeg := exec.Command(ffmpegBinPath, ffmpegArgs...)

	// pipe raspivid's STDOUT to ffmpeg's STDIN
	reader, writer := io.Pipe()
	raspivid.Stdout = writer
	ffmpeg.Stdin = reader

	// run raspivid
	go func() {
		defer writer.Close()

		if err := raspivid.Start(); err != nil {
			panic(err)
		}
		if isVerbose {
			log.Printf("started raspivid with args: %s", strings.Join(raspividArgs, " "))
		} else {
			log.Printf("started raspivid")
		}
		if err := raspivid.Wait(); err != nil {
			panic(err)
		}
	}()

	// run ffmpeg
	if err := ffmpeg.Start(); err != nil {
		panic(err)
	}
	if isVerbose {
		log.Printf("started ffmpeg with args: %s", strings.Join(ffmpegArgs, " "))
	} else {
		log.Printf("started ffmpeg")
	}
	if err := ffmpeg.Wait(); err != nil {
		log.Printf("ffmpeg finished with error: %s", err)
	}
}
