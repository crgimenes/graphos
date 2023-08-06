package graphos

import (
	"embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	sampleRate = 44100
)

var (
	audioContext *audio.Context
	audioPlayer  *audio.Player
	isPlaying    = false

	//go:embed fixture
	assets embed.FS
)

func PrepareWavLoop(filename string) error {

	var xs *audio.InfiniteLoop

	// f, err := os.Open(filename)
	f, err := assets.Open(filename)
	if err != nil {
		return fmt.Errorf("file not found: %v, error: %v", filename, err)
	}
	//defer f.Close() // TODO: close file after playing (test if loop)

	// Decodifique o arquivo WAV
	wavStream, err := wav.Decode(audioContext, f)
	if err != nil {
		return fmt.Errorf("decode wav file: %v, error: %v", filename, err)
	}

	xs = audio.NewInfiniteLoop(wavStream, wavStream.Length())

	audioPlayer, err = audioContext.NewPlayer(xs)
	if err != nil {
		return fmt.Errorf("new player: %v, error: %v", filename, err)
	}

	return nil
}

func (p *Instance) Stop() {
	if audioPlayer != nil {
		audioPlayer.Pause()
		//audioPlayer.Seek(0)
		isPlaying = false
	}
}

func (p *Instance) Play() {
	if isPlaying || audioPlayer == nil {
		return
	}
	audioPlayer.Play()
}

func (p *Instance) InitSound() {
	// TODO: reimplement sound (loops, individual files, wave forms, play frequency, etc)
	audioContext = audio.NewContext(sampleRate)
	err := PrepareWavLoop("fixture/tik.wav")
	if err != nil {
		panic(err)
	}
}
