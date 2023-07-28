package graphos

import (
	"fmt"
	"os"

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
)

var xxx = true

func playWavLoop(filename string) error {

	var xs *audio.InfiniteLoop
	if xxx {
		xxx = false
		// Abra o arquivo WAV
		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("não foi possível abrir o arquivo WAV: %v", err)
		}
		//defer f.Close() // TODO: close file after playing (test if loop)

		// Decodifique o arquivo WAV
		wavStream, err := wav.Decode(audioContext, f)
		if err != nil {
			return fmt.Errorf("não foi possível decodificar o arquivo WAV: %v", err)
		}

		xs = audio.NewInfiniteLoop(wavStream, wavStream.Length())

		// Crie um player de áudio para reproduzir o arquivo em loop
		//audioPlayer, err = audio.NewPlayer(audioContext, wavStream)
		audioPlayer, err = audioContext.NewPlayer(xs)
		if err != nil {
			return fmt.Errorf("não foi possível criar o player de áudio: %v", err)
		}
	}

	audioPlayer.Play()

	// Defina o player de áudio para repetir em loop
	// audioPlayer.SetLooping(true) // see NewInfiniteLoop

	// Inicie a reprodução

	isPlaying = true

	return nil
}

func (p *Instance) StopWav() {
	if audioPlayer != nil {
		audioPlayer.Pause()
		//audioPlayer.Seek(0)
		isPlaying = false
	}
}

func (p *Instance) UpdateSound() error {
	err := playWavLoop("fixture/tik.wav")
	return err
}

func (p *Instance) PlaySound() {
	if audioPlayer != nil {
		audioPlayer.Play()
	}
}

func (p *Instance) InitSound() {
	audioContext = audio.NewContext(sampleRate)
}
