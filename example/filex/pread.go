package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/filex"
	"github.com/tal-tech/go-zero/core/fx"
	"github.com/tal-tech/go-zero/core/logx"
	"gopkg.in/cheggaaa/pb.v1"
)

var (
	file       = flag.String("f", "", "the input file")
	concurrent = flag.Int("c", runtime.NumCPU(), "concurrent goroutines")
	wordVecDic TXDictionary
)

type (
	Vector []float64

	TXDictionary struct {
		EmbeddingCount int64
		Dim            int64
		Dict           map[string]Vector
	}

	pair struct {
		key string
		vec Vector
	}
)

func FastLoad(filename string) error {
	if filename == "" {
		return errors.New("no available dictionary")
	}

	now := time.Now()
	defer func() {
		logx.Infof("article2vec init dictionary end used %v", time.Since(now))
	}()

	dicFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer dicFile.Close()

	header, err := filex.FirstLine(filename)
	if err != nil {
		return err
	}

	total := strings.Split(header, " ")
	wordVecDic.EmbeddingCount, err = strconv.ParseInt(total[0], 10, 64)
	if err != nil {
		return err
	}

	wordVecDic.Dim, err = strconv.ParseInt(total[1], 10, 64)
	if err != nil {
		return err
	}

	wordVecDic.Dict = make(map[string]Vector, wordVecDic.EmbeddingCount)

	ranges, err := filex.SplitLineChunks(filename, *concurrent)
	if err != nil {
		return err
	}

	info, err := os.Stat(filename)
	if err != nil {
		return err
	}

	bar := pb.New64(info.Size()).SetUnits(pb.U_BYTES).Start()
	fx.From(func(source chan<- interface{}) {
		for _, each := range ranges {
			source <- each
		}
	}).Walk(func(item interface{}, pipe chan<- interface{}) {
		offsetRange := item.(filex.OffsetRange)
		scanner := bufio.NewScanner(filex.NewRangeReader(dicFile, offsetRange.Start, offsetRange.Stop))
		scanner.Buffer([]byte{}, 1<<20)
		reader := filex.NewProgressScanner(scanner, bar)
		if offsetRange.Start == 0 {
			// skip header
			reader.Scan()
		}
		for reader.Scan() {
			text := reader.Text()
			elements := strings.Split(text, " ")
			vec := make(Vector, wordVecDic.Dim)
			for i, ele := range elements {
				if i == 0 {
					continue
				}

				v, err := strconv.ParseFloat(ele, 64)
				if err != nil {
					return
				}

				vec[i-1] = v
			}
			pipe <- pair{
				key: elements[0],
				vec: vec,
			}
		}
	}).ForEach(func(item interface{}) {
		p := item.(pair)
		wordVecDic.Dict[p.key] = p.vec
	})

	return nil
}

func main() {
	flag.Parse()

	start := time.Now()
	if err := FastLoad(*file); err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(wordVecDic.Dict))
	fmt.Println(time.Since(start))
}
