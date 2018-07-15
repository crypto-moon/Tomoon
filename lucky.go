package main

import (
	"fmt"
	"hash"
	"crypto/md5"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
	"crypto/sha1"
	"crypto/sha512"
	"crypto/sha256"
	"encoding/hex"
	"encoding/binary"
	"strconv"
	"sort"
)

const REPEAT = 10000
var hashFunc []hash.Hash

func init() {
	hashFunc = make([]hash.Hash, 0)
	hashFunc = append(hashFunc, md5.New())
	hashFunc = append(hashFunc, md4.New())
	hashFunc = append(hashFunc, ripemd160.New())
	hashFunc = append(hashFunc, sha256.New())
	hashFunc = append(hashFunc, sha1.New())
	hashFunc = append(hashFunc, sha512.New())
}

type Score struct {
	id uint64
	score uint64
}

type ScoreSlice []Score

func(s ScoreSlice) Len() int {
	return len(s)
}

func(s ScoreSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func(s ScoreSlice) Less(i, j int) bool {
	return s[i].score > s[j].score
}

func getLuckyNum(hash string) []byte{
	lucky, err := hex.DecodeString(hash)
	if err != nil {
		panic(fmt.Sprintf("block hash error: %s", err))
	}
	fmt.Println("Rolling lucky num:")
	for i:=0; i < REPEAT; i++{
		for _, f := range hashFunc{
			f.Write(lucky)
			lucky = f.Sum(nil)
			f.Reset()
		}
	}
	fmt.Printf("lucky num: %x\n", lucky)
	return lucky
}

func getScore(lucky []byte, id uint64) uint64{
	encoder := sha256.New()
	encoder.Write(lucky)
	bytes := make([]byte, 8, 8)
	binary.LittleEndian.PutUint64(bytes, id)
	encoder.Write(bytes)
	scoreBytes := encoder.Sum(nil)
	var score uint64
	for i := 0; i < 8; i++ {
		score = score * 16 * 16 + uint64(scoreBytes[i])
	}
	return score
}

func main() {
	var blockHash, lotteryNumStr string
	fmt.Printf("Please enter the block hash: ")
	fmt.Scanln(&blockHash)
	fmt.Printf("Please enter lottery num: ")
	fmt.Scanln(&lotteryNumStr)
	luckyNum := getLuckyNum(blockHash)
	lotteryNum, err := strconv.ParseUint(lotteryNumStr, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("strconv error: %s", err))
	}
	scoreList := make([]Score, lotteryNum)
	for i := uint64(0); i < lotteryNum; i ++ {
		scoreList[i] = Score{i + 1, getScore(luckyNum, i + 1)}
	}
	sort.Sort(ScoreSlice(scoreList))
	for k, v := range scoreList {
		fmt.Printf("%d----score:%d, id:%d\n", k+1, v.score, v.id)
	}
	fmt.Printf("Press enter to exit...\n")
	fmt.Scanln()
}