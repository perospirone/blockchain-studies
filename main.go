package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"

	"golang.org/x/crypto/sha3"
)

type Block struct {
	Index     int64
	Nonce     int64
	Timestamp int64
	Data      string
	Hash      []byte
}

func (b *Block) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write Index
	if err := binary.Write(buf, binary.LittleEndian, b.Index); err != nil {
		return nil, err
	}
	// Write Nonce
	if err := binary.Write(buf, binary.LittleEndian, b.Nonce); err != nil {
		return nil, err
	}
	// Write Timestamp
	if err := binary.Write(buf, binary.LittleEndian, b.Timestamp); err != nil {
		return nil, err
	}
	// Write Data string
	dataLen := int32(len(b.Data)) // store the length of the string
	if err := binary.Write(buf, binary.LittleEndian, dataLen); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, []byte(b.Data)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (b *Block) CalculateHash() error {
	serializedData, err := b.Serialize()

	if err != nil {
		return err
	}

	hash := sha256.Sum256(serializedData)

	fmt.Println("hash: ", hash)

	b.Hash = hash[:]

	return nil
}

// Deserialize deserializes binary data into a Block
func Deserialize(data []byte) (*Block, error) {
	buf := bytes.NewReader(data)
	block := &Block{}

	// Read Index
	if err := binary.Read(buf, binary.LittleEndian, &block.Index); err != nil {
		return nil, err
	}
	// Read Nonce
	if err := binary.Read(buf, binary.LittleEndian, &block.Nonce); err != nil {
		return nil, err
	}
	// Read Timestamp
	if err := binary.Read(buf, binary.LittleEndian, &block.Timestamp); err != nil {
		return nil, err
	}
	// Read Data length
	var dataLen int32
	if err := binary.Read(buf, binary.LittleEndian, &dataLen); err != nil {
		return nil, err
	}
	// Read Data string
	dataBytes := make([]byte, dataLen)
	if err := binary.Read(buf, binary.LittleEndian, &dataBytes); err != nil {
		return nil, err
	}
	block.Data = string(dataBytes)

	return block, nil
}

func main() {
	genesisBlock := Block{
    	Index:     0,
    	Nonce:     16383,
    	Timestamp: time.Now().Unix(),
    	Data:      "pudim",
    	Hash:      []byte{},
	}

	sucrilhos, err := genesisBlock.Serialize()

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(sucrilhos)
	fmt.Println(Deserialize(sucrilhos))
	fmt.Println(genesisBlock.CalculateHash())
	if err := genesisBlock.CalculateHash(); err != nil {
		panic(err.Error())
	}
	fmt.Println(genesisBlock)







	return
}

func oldMain() {
	bytes := writeNonce("pudim", 16383)

	fmt.Println(bytes)

	fmt.Println(getNonce(bytes))

	keccak := sha3.NewLegacyKeccak256()

	keccak.Write(bytes)

	hash := keccak.Sum(nil)

	hashString := fmt.Sprintf("%x", hash)

	fmt.Println("Full Hash:", hashString)
	fmt.Println("Hash Size (bytes):", len(hash))
}

func writeNonce(str string, nonce int) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(nonce))

	finalBytes := append([]byte(str), bytes...)

	return finalBytes
}

func getNonce(bytes []byte) int {
	lastBytes := bytes[len(bytes)-4:]
	reconvertedValue := binary.BigEndian.Uint32(lastBytes)

	return int(reconvertedValue)
}
