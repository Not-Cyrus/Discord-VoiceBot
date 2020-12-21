// https://raw.githubusercontent.com/bwmarrin/dgvoice/master/dgvoice.go

package audio

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

const (
	channels  int = 2
	frameRate int = 48000
	frameSize int = 960
	maxBytes  int = 3840
)

var (
	speakers    map[uint32]*gopus.Decoder
	opusEncoder *gopus.Encoder
	mu          sync.Mutex
	run         *exec.Cmd
)

// why this function was exposed to the global scope in the actual version? Who knows

var onError = func(str string, err error) {
	prefix := "dgVoice: " + str

	if err != nil {
		os.Stderr.WriteString(prefix + ": " + err.Error())
	} else {
		os.Stderr.WriteString(prefix)
	}
}

func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	var err error

	opusEncoder, err = gopus.NewEncoder(frameRate, channels, gopus.Audio)

	if err != nil {
		onError("NewEncoder Error", err)
		return
	}

	for {

		recv, ok := <-pcm
		if !ok {
			//onError("PCM Channel closed", nil) THIS IS SO FUCKING ANNOYING LIKE STOP I GET IT THE CHANNELS CLOSED!!!
			return
		}

		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			onError("Encoding Error", err)
			return
		}

		if v.Ready == false || v.OpusSend == nil {
			return
		}
		v.OpusSend <- opus
	}
}

func Skip() {
	if run != nil {
		run.Process.Kill()
	}
}
func PlayAudioFile(v *discordgo.VoiceConnection, filename string, stop <-chan bool) {

	run = exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	ffmpegout, err := run.StdoutPipe()
	if err != nil {
		onError("StdoutPipe Error", err)
		return
	}

	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	err = run.Start()
	if err != nil {
		onError("RunStart Error", err)
		return
	}

	go func() {
		<-stop
		err = run.Process.Kill()
	}()

	err = v.Speaking(true)
	if err != nil {
		onError("Couldn't set speaking", err)
	}

	defer func() {
		err := v.Speaking(false)
		if err != nil {
			onError("Couldn't stop speaking", err)
		}
	}()

	send := make(chan []int16, 2)
	defer close(send)

	close := make(chan bool)
	go func() {
		SendPCM(v, send)
		close <- true
	}()

	for {

		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		if err != nil {
			onError("error reading from ffmpeg stdout", err)
			return
		}

		select {
		case send <- audiobuf:
		case <-close:
			return
		}
	}
}
