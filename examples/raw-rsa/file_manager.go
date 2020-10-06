package main

import (
	"bufio"
	cryptoRand "crypto/rand"
	"math/big"
	"math/rand"
	"os"

	"github.com/DiscreteTom/rawrsa"
)

// FileManagerSession stores all session resources.
type FileManagerSession struct {
	dataFile   *os.File
	tagFile    *os.File
	scanner    *bufio.Scanner
	dataBuffer []byte
}

// FileManager will manage data files & tag files.
type FileManager struct {
	rr           *rawrsa.RawRsa
	dataFilename string
	tagFilename  string
	chunkSize    int
	blockCount   int

	ss *FileManagerSession
}

// NewFiles will use the given params to create files and return the FileManager
func NewFiles(keyFilename string, rsaKeyBits int, dataFilename, tagFilename string, chunkSize, blockCount int) *FileManager {
	// create keyfile and save
	var rr, err = rawrsa.NewRawRsa(cryptoRand.Reader, rsaKeyBits)
	if err != nil {
		panic(err)
	}
	err = rr.Save(keyFilename)
	if err != nil {
		panic(err)
	}

	// create data file
	dataFile, err := os.OpenFile(dataFilename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()
	// create tag file
	tagFile, err := os.OpenFile(tagFilename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer tagFile.Close()
	dataFile.Truncate(0) // clear file
	tagFile.Truncate(0)  // clear file

	// write files
	data := make([]byte, chunkSize)
	for i := 0; i < blockCount; i++ {
		// data file consists of random bytes
		rand.Read(data)
		if _, err = dataFile.Write(data); err != nil {
			panic(err)
		}
		// each line in tag file is a tag
		if _, err = tagFile.WriteString(rr.RawEncrypt(new(big.Int).SetBytes(data)).String() + "\n"); err != nil {
			panic(err)
		}
	}
	return &FileManager{
		rr:           rr,
		dataFilename: dataFilename,
		tagFilename:  tagFilename,
		chunkSize:    chunkSize,
		blockCount:   blockCount,
		ss:           nil,
	}
}

// LoadFiles will use existing files and return a FileManager.
func LoadFiles(keyFilename, dataFilename, tagFilename string, chunkSize int, blockCount int) *FileManager {
	// load key file
	rr, err := rawrsa.Load(keyFilename)
	if err != nil {
		panic(err)
	}
	return &FileManager{
		rr:           rr,
		dataFilename: dataFilename,
		tagFilename:  tagFilename,
		chunkSize:    chunkSize,
		blockCount:   blockCount,
		ss:           nil,
	}
}

// StartSession will open files for reading operations.
func (fm *FileManager) StartSession() {
	// open files
	dataFile, err := os.Open(fm.dataFilename)
	if err != nil {
		panic(err)
	}
	tagFile, err := os.Open(fm.tagFilename)
	if err != nil {
		panic(err)
	}

	fm.ss = &FileManagerSession{
		dataFile:   dataFile,
		tagFile:    tagFile,
		scanner:    bufio.NewScanner(tagFile),
		dataBuffer: make([]byte, fm.chunkSize),
	}
}

// EndSession will close files.
func (fm *FileManager) EndSession() {
	fm.ss.dataFile.Close()
	fm.ss.tagFile.Close()
	fm.ss = nil
}

// NextBlockData will read from data file and return bytes of a block size.
func (fm *FileManager) NextBlockData() *big.Int {
	// for data file, read a chunk
	_, err := fm.ss.dataFile.Read(fm.ss.dataBuffer)
	if err != nil {
		panic(err)
	}
	return new(big.Int).SetBytes(fm.ss.dataBuffer)
}

// NextBlockTag will read from tag file and return the tag.
func (fm *FileManager) NextBlockTag() *big.Int {
	fm.ss.scanner.Scan()
	result, _ := new(big.Int).SetString(fm.ss.scanner.Text(), 10)
	return result
}
