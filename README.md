# Raspberry Pi Camera: Youtube Streamer

## 0. What is this?

It's a go application for streaming video through [Youtube Live Stream](https://www.youtube.com/live_dashboard), using [raspivid](https://www.raspberrypi.org/documentation/usage/camera/raspicam/raspivid.md) and [ffmpeg](https://ffmpeg.org/) on [Raspberry Pi](https://www.raspberrypi.org/).

## 1. What do I need before running it?

You need:

* Raspberry Pi
* Raspberry Pi Camera Module enabled, and its cable correctly connected
* [golang installed on Raspberry Pi](https://github.com/meinside/rpi-configs/blob/master/bin/prep_go.sh)
* [ffmpeg installed on Raspberry Pi](https://github.com/meinside/rpi-configs/blob/master/bin/prep_ffmpeg.sh)
* and this README.md.

## 2. How can I build it?

```bash
$ go get -u github.com/meinside/rpi-camera-youtube-streamer
$ cd $GOPATH/src/github.com/meinside/rpi-camera-youtube-streamer
$ go build
```

or

```bash
$ git clone https://github.com/meinside/rpi-camera-youtube-streamer.git
$ cd rpi-camera-youtube-streamer
$ go build
```

## 3. How can I run it?

You need to create your own config file.

Sample file is included, so feel free to copy it and change values as you want.

```bash
$ cp config.json.sample config.json
$ vi config.json
```

At least you have to change the **youtube_stream_key** value for running it.

You can get your key in your [Live Stream Dashboard](https://www.youtube.com/live_dashboard), which is labeled as 'Stream name/key'.

After the configuration is finished, just execute the binary:

```bash
$ ./rpi-cameera-youtube-streamer
```

If nothing goes wrong, you'll see your live streaming in your dashboard in several seconds.

## 4. How can I run it as a service?

### systemd

```bash
$ sudo cp systemd/rpi-camera-youtube-streamer.service /lib/systemd/system/
$ sudo vi /lib/systemd/system/rpi-camera-youtube-streamer.service
```

and edit **User**, **Group**, **WorkingDirectory**, and **ExecStart** values.

You can simply start/stop it with:

```
$ sudo systemctl start rpi-camera-youtube-streamer.service
$ sudo systemctl stop rpi-camera-youtube-streamer.service
```

If you want to launch it automatically on boot:

```bash
$ sudo systemctl enable rpi-camera-youtube-streamer.service
```

## 998. Any trouble?

Please open an issue.

## 999. License?

MIT
