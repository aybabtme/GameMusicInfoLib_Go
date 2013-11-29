package nsf

import (
	"fmt"
	"path/filepath"
	"testing"
)

var tt = []struct {
	testfile string
	want     *NSFSong
}{
	{"Super Mario Bros 2.nsf", &NSFSong{
		headerMagic:      "NESM ",
		versionNumber:    1,
		totalSongs:       34,
		startingSong:     1,
		loadAddress:      32768,
		initAddress:      34816,
		playAddress:      32768,
		songName:         "Super Mario Bros. 2",
		artistName:       "<?>",
		copyright:        "1988 Nintendo",
		songTicks:        16666,
		usesNTSC:         true,
		isDualSupportive: true,
		usingVRC6:        false,
		usingVRC7:        false,
		usingFDS:         false,
		usingMMC5:        false,
		usingNamco:       false,
		usingSunsoft:     false,
	}},
}

func TestMatchTestData(t *testing.T) {

	for _, testCase := range tt {
		want := testCase.want
		filename := fmt.Sprintf("../testdata/%s", testCase.testfile)
		// FromSlash cleans up / separated path to be OS specific
		got, err := NewSongFromFile(filepath.FromSlash(filename))
		if err != nil {
			t.Errorf("For file %s, %v", filename, err)
		}

		assertEqual(t, want, got)
	}
}

func assertEqual(t *testing.T, want, got *NSFSong) {
	assertStringEqual(t, "HeaderMagic", want.HeaderMagic(), got.HeaderMagic())

	if want.VersionNumber() != got.VersionNumber() {
		t.Errorf("VersionNumber, want '%v' got '%v'", want.VersionNumber(), got.VersionNumber())
	}
	if want.TotalSongs() != got.TotalSongs() {
		t.Errorf("TotalSongs, want '%v' got '%v'", want.TotalSongs(), got.TotalSongs())
	}
	if want.StartingSong() != got.StartingSong() {
		t.Errorf("StartingSong, want '%v' got '%v'", want.StartingSong(), got.StartingSong())
	}
	if want.LoadAddress() != got.LoadAddress() {
		t.Errorf("LoadAddress, want '%v' got '%v'", want.LoadAddress(), got.LoadAddress())
	}
	if want.InitAddress() != got.InitAddress() {
		t.Errorf("InitAddress, want '%v' got '%v'", want.InitAddress(), got.InitAddress())
	}
	if want.PlayAddress() != got.PlayAddress() {
		t.Errorf("PlayAddress, want '%v' got '%v'", want.PlayAddress(), got.PlayAddress())
	}

	assertStringEqual(t, "SongName", want.SongName(), got.SongName())
	assertStringEqual(t, "ArtistName", want.ArtistName(), got.ArtistName())
	assertStringEqual(t, "Copyright", want.Copyright(), got.Copyright())

	if want.SongTicks() != got.SongTicks() {
		t.Errorf("SongTicks, want '%v' got '%v'", want.SongTicks(), got.SongTicks())
	}
	if want.IsNTSC() != got.IsNTSC() {
		t.Errorf("IsNTSC, want '%v' got '%v'", want.IsNTSC(), got.IsNTSC())
	}
	if want.IsDualSupportive() != got.IsDualSupportive() {
		t.Errorf("IsDualSupportive, want '%v' got '%v'", want.IsDualSupportive(), got.IsDualSupportive())
	}
	if want.UsingVRC6() != got.UsingVRC6() {
		t.Errorf("UsingVRC6, want '%v' got '%v'", want.UsingVRC6(), got.UsingVRC6())
	}
	if want.UsingVRC7() != got.UsingVRC7() {
		t.Errorf("UsingVRC7, want '%v' got '%v'", want.UsingVRC7(), got.UsingVRC7())
	}
	if want.UsingFDS() != got.UsingFDS() {
		t.Errorf("UsingFDS, want '%v' got '%v'", want.UsingFDS(), got.UsingFDS())
	}
	if want.UsingMMC5() != got.UsingMMC5() {
		t.Errorf("UsingMMC5, want '%v' got '%v'", want.UsingMMC5(), got.UsingMMC5())
	}
	if want.UsingNamco() != got.UsingNamco() {
		t.Errorf("UsingNamco, want '%v' got '%v'", want.UsingNamco(), got.UsingNamco())
	}
	if want.UsingSunsoft() != got.UsingSunsoft() {
		t.Errorf("UsingSunsoft, want '%v' got '%v'", want.UsingSunsoft(), got.UsingSunsoft())
	}
}

func assertStringEqual(t *testing.T, testname, want, got string) {
	if want == got {
		return
	}
	t.Errorf("%s, want '%v' got '%v'", testname, want, got)

	if len(want) != len(got) {
		t.Errorf("\tWant length of %d, got %d\n", len(want), len(got))

	}

	var minLen int
	if len(want) < len(got) {
		minLen = len(want)
	} else {
		minLen = len(got)
	}

	for i := 0; i < minLen; i++ {
		if want[i] == got[i] {
			continue
		}
		t.Errorf("\tByte %d of string, want %.8b got %.8b\n", i, want[i], got[i])
	}
}
