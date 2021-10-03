package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"

	arg "github.com/alexflint/go-arg"
)

var patterns = map[string][]byte{
	"coins": {
		'\x43', '\x6F', '\x69', '\x6E', '\x00', '\x29', '\x00', '\x00',
		'\x00', '\x43', '\x6F', '\x75', '\x6E', '\x74', '\x5F', '\x35',
		'\x5F', '\x43', '\x35', '\x33', '\x41', '\x36', '\x39', '\x30',
		'\x30', '\x34', '\x42', '\x42', '\x30', '\x31', '\x45', '\x45',
		'\x42', '\x42', '\x41', '\x30', '\x32', '\x43', '\x44', '\x38',
		'\x36', '\x44', '\x34', '\x43', '\x42', '\x41', '\x36', '\x43',
		'\x36',
	},
	// "survivalPoints": {
	// 	'\x53', '\x5F', '\x50', '\x74', '\x5F', '\x37', '\x36', '\x5F',
	// 	'\x31', '\x44', '\x42', '\x35', '\x44', '\x31', '\x46', '\x37',
	// 	'\x34', '\x46', '\x30', '\x38', '\x44', '\x37', '\x46', '\x39',
	// 	'\x41', '\x46', '\x42', '\x41', '\x35', '\x36', '\x38', '\x42',
	// 	'\x35', '\x36', '\x32', '\x41', '\x30', '\x41', '\x35', '\x44',
	// },
	// "extractionPoints": {
	// 	'\x45', '\x5F', '\x50', '\x74', '\x5F', '\x36', '\x34', '\x5F',
	// 	'\x36', '\x46', '\x38', '\x34', '\x30', '\x44', '\x35', '\x44',
	// 	'\x34', '\x31', '\x42', '\x46', '\x42', '\x42', '\x46', '\x46',
	// 	'\x35', '\x45', '\x34', '\x39', '\x43', '\x42', '\x38', '\x46',
	// 	'\x37', '\x38', '\x42', '\x43', '\x32', '\x33', '\x34', '\x36',
	// },
	"reputation": {
		'\x44', '\x79', '\x6E', '\x61', '\x73', '\x74', '\x79', '\x52',
		'\x65', '\x70', '\x75', '\x74', '\x61', '\x74', '\x69', '\x6F',
		'\x6E', '\x00',
	},
	"age": {
		'\x41', '\x67', '\x65', '\x5F', '\x31', '\x37', '\x5F', '\x39',
		'\x35', '\x44', '\x34', '\x43', '\x32', '\x35', '\x42', '\x34',
		'\x30', '\x34', '\x37', '\x39', '\x36', '\x42', '\x30', '\x33',
		'\x39', '\x46', '\x35', '\x30', '\x45', '\x41', '\x42', '\x36',
		'\x30', '\x33', '\x30', '\x34', '\x44', '\x38', '\x32',
	},
}

func find(toFind, saveBytes []byte, offset int) int64 {
	index := bytes.Index(saveBytes, toFind)
	if index == -1 {
		return -1
	}
	index += len(toFind) + offset
	return int64(index)
}

func readValue(f *os.File, saveBytes, toFind []byte, offset int) (int64, uint32, error) {
	index := find(toFind, saveBytes, offset)
	if index == -1 {
		return index, 0, errors.New("Couldn't find bytes.")
	}
	_, err := f.Seek(index, io.SeekStart)
	if err != nil {
		return index, 0, err
	}
	buffer := make([]byte, 4)
	_, err = f.Read(buffer)
	if err != nil {
		return index, 0, err
	}
	value := binary.LittleEndian.Uint32(buffer)
	return index, value, nil
}

func write(f *os.File, buffer []byte, index int64) error {
	_, err := f.Seek(index, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = f.Write(buffer)
	return err
}

func writeValues(f *os.File, offsets map[string]int64, args *Args) error {
	buffer := make([]byte, 4)
	if args.Coins != -1 {
		binary.LittleEndian.PutUint32(buffer, uint32(args.Coins*10))
		err := write(f, buffer, offsets["coins"])
		if err != nil {
			return err
		}
	}
	if args.Age != -1 {
		binary.LittleEndian.PutUint32(buffer, math.Float32bits(args.Age))
		err := write(f, buffer, offsets["age"])
		if err != nil {
			return err
		}
	}
	if args.Reputation != -1 {
		binary.LittleEndian.PutUint32(buffer, uint32(args.Reputation))
		err := write(f, buffer, offsets["reputation"])
		if err != nil {
			return err
		}
	}
	return nil
}

func readValues(f *os.File, saveBytes []byte) (map[string]int64, error) {
	offsets := map[string]int64{}
	// Coins
	coinsIndex, coins, err := readValue(f, saveBytes, patterns["coins"], 26)
	if err != nil {
		return nil, err
	}
	coinsInt32 := int32(coins)
	coinsFloat := float32(coinsInt32) / 10
	offsets["coins"] = coinsIndex

	// Dynasty reputation
	dynRepIndex, dynastyRep, err := readValue(f, saveBytes, patterns["reputation"], 25)
	if err != nil {
		return nil, err
	}
	offsets["reputation"] = dynRepIndex

	// Age
	ageIndex, age, err := readValue(f, saveBytes, patterns["age"], 28)
	if err != nil {
		return nil, err
	}
	offsets["age"] = ageIndex

	// spIndex, survivalPoints, err := readValue(f, saveBytes, patterns["survivalPoints"], 26)
	// if err != nil && spIndex != -1 {
	// 	return nil, err
	// }
	// epIndex, extractionPoints, err := readValue(f, saveBytes, patterns["extractionPoints"], 26)
	// if err != nil && epIndex != -1 {
	// 	return nil, err
	// }
	fmt.Printf("Age:                %f -> 0x%X\n", math.Float32frombits(age), ageIndex)
	fmt.Printf("Coins:              %.1f -> 0x%X\n", coinsFloat, coinsIndex)
	fmt.Printf("Dynasty reputation: %d -> 0x%X\n", int32(dynastyRep), dynRepIndex)
	// if spIndex != -1 {
	// 	fmt.Printf("Survival points:    %d -> 0x%X\n", int32(survivalPoints), spIndex)
	// }
	// if epIndex != -1 {
	// 	fmt.Printf("Extraction points:  %d -> 0x%X\n", int32(extractionPoints), epIndex)
	// }
	return offsets, nil
}

func checkHeader(header []byte) bool {
	return bytes.Equal(header, []byte{71, 86, 65, 83})
}

func parseArgs() *Args {
	var args Args
	arg.MustParse(&args)
	return &args
}

func main() {
	args := parseArgs()
	savePath := args.Path
	saveBytes, err := ioutil.ReadFile(savePath)
	if err != nil {
		panic(err)
	}
	ok := checkHeader(saveBytes[:4])
	if !ok {
		panic("Invalid save file header.")
	}
	f, err := os.OpenFile(savePath, os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	offsets, err := readValues(f, saveBytes)
	if err != nil {
		panic(err)
	}
	if len(os.Args) == 2 {
		fmt.Println("\nPress enter to exit.")
		fmt.Scanln()
	} else {
		err = writeValues(f, offsets, args)
		if err != nil {
			panic(err)
		}
		fmt.Println("\nOK. Press enter to exit.")
		fmt.Scanln()
	}
}
