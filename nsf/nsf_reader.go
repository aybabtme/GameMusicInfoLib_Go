package nsf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
)

type NSFSong struct {
	headerMagic      string
	versionNumber    byte
	totalSongs       byte
	startingSong     byte
	loadAddress      uint16
	initAddress      uint16
	playAddress      uint16
	songName         string
	artistName       string
	copyright        string
	songTicks        uint16
	usesNTSC         bool
	isDualSupportive bool

	// Audio chip flags
	usingVRC6    bool
	usingVRC7    bool
	usingFDS     bool
	usingMMC5    bool
	usingNamco   bool
	usingSunsoft bool
}

// NewSongFromFile tries to open the NSF file at filepath and extract information from it.
func NewSongFromFile(filepath string) (*NSFSong, error) {

	readOrPanic := func(n int, err error) {
		if err != nil {
			panic(err)
		}
	}

	// Attempt to open the given file
	nsf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening Nintendo NSF file, %v", err)
	}

	bufNsf := bytes.NewReader(nsf)

	magicBuf := make([]byte, 5)
	readOrPanic(bufNsf.Read(magicBuf))
	headerMagic := string(magicBuf)

	versionNumber, err := bufNsf.ReadByte()
	if err != nil {
		panic(err)
	}

	totalSongs, err := bufNsf.ReadByte()
	if err != nil {
		panic(err)
	}

	startingSong, err := bufNsf.ReadByte()
	if err != nil {
		panic(err)
	}

	// Read the data addresses
	var loadAddress uint16
	var initAddress uint16
	var playAddress uint16
	binary.Read(bufNsf, binary.LittleEndian, &loadAddress)
	binary.Read(bufNsf, binary.LittleEndian, &initAddress)
	binary.Read(bufNsf, binary.LittleEndian, &playAddress)

	songNameBuf := make([]byte, 32)
	readOrPanic(bufNsf.Read(songNameBuf))
	songName := string(songNameBuf)

	artistNameBuf := make([]byte, 32)
	readOrPanic(bufNsf.Read(artistNameBuf))
	artistName := string(artistNameBuf)

	copyrightBuf := make([]byte, 32)
	readOrPanic(bufNsf.Read(copyrightBuf))
	copyright := string(copyrightBuf)

	usesNTSC := isNTSC(bufNsf)
	isDualSupportive := isDualSupportive(bufNsf)
	songTicks := speedTicks(bufNsf)

	// Evaluate sound chip mapper support.
	bufNsf.Seek(0x7B, os.SEEK_SET)
	chipByte, err := bufNsf.ReadByte()
	if err != nil {
		panic(err)
	}
	usingVRC6 := isUsingVRC6(chipByte)
	usingVRC7 := isUsingVRC7(chipByte)
	usingFDS := isUsingFDS(chipByte)
	usingMMC5 := isUsingMMC5(chipByte)
	usingNamco := isUsingNamco(chipByte)
	usingSunsoft := isUsingSunsoft(chipByte)

	return &NSFSong{
		headerMagic:      headerMagic,
		versionNumber:    versionNumber,
		totalSongs:       totalSongs,
		startingSong:     startingSong,
		loadAddress:      loadAddress,
		initAddress:      initAddress,
		playAddress:      playAddress,
		songName:         songName,
		artistName:       artistName,
		copyright:        copyright,
		songTicks:        songTicks,
		usesNTSC:         usesNTSC,
		isDualSupportive: isDualSupportive,
		usingVRC6:        usingVRC6,
		usingVRC7:        usingVRC7,
		usingFDS:         usingFDS,
		usingMMC5:        usingMMC5,
		usingNamco:       usingNamco,
		usingSunsoft:     usingSunsoft,
	}, nil
}

// HeaderMagic the initial header magic for the NSF file.
func (n *NSFSong) HeaderMagic() string {
	return n.headerMagic
}

// VersionNumber The version number of the NSF specification used within this NSF file.
func (n *NSFSong) VersionNumber() byte {
	return n.versionNumber
}

// TotalSongs the total number of songs in the NSF file.
// NOTE: This is 1-based. A result of '1' means exactly one song, etc.
func (n *NSFSong) TotalSongs() byte {
	return n.totalSongs
}

// StartingSong is the very first song number
// NOTE: This is 1-based. ie. If this returns '1' then it
//       means song 1.
func (n *NSFSong) StartingSong() byte {
	return n.startingSong
}

// Can be from 0x8000 to 0xFFFF
func (n *NSFSong) LoadAddress() uint16 {
	return n.loadAddress
}

// Can be from 0x8000 to 0xFFFF
func (n *NSFSong) InitAddress() uint16 {
	return n.initAddress
}

// Can be from 0x8000 to 0xFFFF
func (n *NSFSong) PlayAddress() uint16 {
	return n.playAddress
}

// SongName the internal song name of this NSF file.
// NOTE: If no actual name is present in the NSF file, then "<?>" will be returned
func (n *NSFSong) SongName() string {
	return n.songName
}

// ArtistName the name of the artist/composer of the NSF files.
// NOTE: If no actual name is present in the NSF file, then "<?>" will be returned
func (n *NSFSong) ArtistName() string {
	return n.artistName
}

// Copyright the copyright string embedded within the NSF file.
// NOTE: If an actual copyright string is not present, this will return "<?>"
func (n *NSFSong) Copyright() string {
	return n.copyright
}

// Song speed in 1/1000000th sec ticks
func (n *NSFSong) SongTicks() uint16 {
	return n.songTicks
}

// IsNTSC Checks whether or not the NSF file is using the NTSC NES clock rate.
func (n *NSFSong) IsNTSC() bool {
	return n.usesNTSC
}

// IsDualSupportive Checks if the NSF file supports both NTSC and PAL clock rates.
func (n *NSFSong) IsDualSupportive() bool {
	return n.isDualSupportive
}

// Sound chip support
func (n *NSFSong) UsingVRC6() bool {
	return n.usingVRC6
}
func (n *NSFSong) UsingVRC7() bool {
	return n.usingVRC7
}
func (n *NSFSong) UsingFDS() bool {
	return n.usingFDS
}
func (n *NSFSong) UsingMMC5() bool {
	return n.usingMMC5
}
func (n *NSFSong) UsingNamco() bool {
	return n.usingNamco
}
func (n *NSFSong) UsingSunsoft() bool {
	return n.usingSunsoft
}

// isNTSC Checks whether or not the NSF file is using the NTSC NES clock rate
// NOTE:  If this returns false, then the NSF file is using the PAL NES clock rate.
func isNTSC(nsf *bytes.Reader) bool {
	nsf.Seek(0x7A, os.SEEK_SET)
	playbackByte, err := nsf.ReadByte()
	if err != nil {
		panic(err)
	}

	if playbackByte&1 == 0 {
		return true
	}

	return false
}

// isDualSupportive True if the NSF file supports both NTSC and PAL clock rates.
func isDualSupportive(nsf *bytes.Reader) bool {
	nsf.Seek(0x7A, os.SEEK_SET)
	playbackByte, err := nsf.ReadByte()
	if err != nil {
		panic(err)
	}

	if playbackByte&2 == 0 {
		return true
	}

	return false
}

// Song speed in 1/1000000th sec ticks
func speedTicks(nsf *bytes.Reader) uint16 {
	var isNtsc bool = isNTSC(nsf)
	var ticks uint16 = 0

	if isNtsc {
		nsf.Seek(0x6E, os.SEEK_SET)
		binary.Read(nsf, binary.LittleEndian, &ticks)
	} else {
		nsf.Seek(0x78, os.SEEK_SET)
		binary.Read(nsf, binary.LittleEndian, &ticks)
	}

	return ticks
}

func isUsingVRC6(chipByte byte) bool {
	if chipByte&1 != 0 {
		return true
	}

	return false
}
func isUsingVRC7(chipByte byte) bool {
	if chipByte&2 != 0 {
		return true
	}

	return false
}
func isUsingFDS(chipByte byte) bool {
	if chipByte&4 != 0 {
		return true
	}

	return false
}
func isUsingMMC5(chipByte byte) bool {
	if chipByte&8 != 0 {
		return true
	}

	return false
}
func isUsingNamco(chipByte byte) bool {
	if chipByte&16 != 0 {
		return true
	}

	return false
}
func isUsingSunsoft(chipByte byte) bool {
	if chipByte&32 != 0 {
		return true
	}

	return false
}
