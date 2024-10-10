// You can edit this code!
// Click here and start typing.
package main

import "fmt"

type MediaPlayer interface {
	PlayAudio(filename string)
}

type AdvancedMediaPlayer interface {
	PlayVideo(filename string)
}

type VideoPlayer struct{}

func (p *VideoPlayer) PlayVideo(filename string) {
	fmt.Printf("Playing video file: %s\n", filename)
}

type AudioPlayer struct{}

func (p *AudioPlayer) PlayAudio(filename string) {
	fmt.Printf("Playing audio file: %s\n", filename)
}

type MediaAdapter struct {
	player AdvancedMediaPlayer
}

func (adapter *MediaAdapter) PlayAudio(filename string) {
	// 实际上调用的是 VideoPlayer 的 PlayVideo 方法
	adapter.player.PlayVideo(filename)
}

func main() {
	audioPlayer := &AudioPlayer{}
	audioPlayer.PlayAudio("song.mp3")

	adapter := &MediaAdapter{player: &VideoPlayer{}}
	adapter.PlayAudio("movie.mp4")
}
