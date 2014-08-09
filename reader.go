package AdtExtractor

import (
	"os"
	"log"
	"fmt"
	"math"
	"unicode/utf8"
	"bytes"
	"encoding/binary"
	"encoding/json"
)

type adtChunk struct {
	t string
	length int32
	offset int64
}
type adtAreaChunk struct {
	MapName string
	ZoneId int32
	Xp int32
	Yp int32
}

func Extract(fileName string, x int32, y int32, mapName string) {
	file, err := os.Open(fileName) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	var chunks = make(map [int64]adtChunk)
	var offset int64
	offset = 0
	for {
		chunk, offsetNew, err := readChunk(file, offset)
		if err != nil {
			break
		}
		offset = offsetNew
		chunks[offset] = chunk
	}

	var n = 0
	for _, chunk := range chunks {
		if chunk.t == "MCNK" {
			file.Seek(int64(chunk.offset) + 0x3C, 0)
			data := make([]byte, 4)
			_, err = file.Read(data)
			if err != nil {
				fmt.Print("%v", err)
				break
			}
			var zoneId = read_int32(data)

			var i = n
			sx := int32(i % 16)
			var syFloor = float64(i / 16)
			sy := int32(math.Floor(syFloor))

			xp := (x * 256) + (sx * 16)
			yp := (y * 256) + (sy * 16)

			b, err := json.Marshal(adtAreaChunk{MapName: mapName, ZoneId: zoneId, Xp: xp, Yp: yp})
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Println(string(b))

			n++
		}
	}
	file.Close()
}

func readChunk(file *os.File, offsetOriginal int64) (chunk adtChunk, offset int64, err error) {
	chunkType, err := readType(file)
	if err == nil {
		chunkLength, err2 := readLength(file)
		offset = offsetOriginal + 8 + int64(chunkLength)
		if err2 == nil {
			chunk = adtChunk{t: chunkType, length: chunkLength, offset: offset}
			file.Seek(int64(chunkLength), 1)
		}
	}
	return
}

func readType(file *os.File) (chunkType string, err error) {
	data := make([]byte, 4)
	count, err := file.Read(data)
	count = count
	if err == nil {
		chunkType = reverse(string(data[:count]))
	}
	return
}

func readLength(file *os.File) (length int32, err error) {
	data := make([]byte, 4)
	_, err = file.Read(data)
	if err == nil {
		length = read_int32(data)
	}
	return
}

func readContents(file *os.File, length int32) {
	data := make([]byte, length)
	_, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
}

/**
 * Found this on: http://codereview.stackexchange.com/a/18169
 */
func read_int32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

/**
 * Found this on: http://www.snip2code.com/Snippet/13468/golang--reverse-string-
 */
func reverse(s string) string {
	cs := make([]rune, utf8.RuneCountInString(s))
	i := len(cs)
	for _, c := range s {
		i--
		cs[i] = c
	}
	return string(cs)
}
