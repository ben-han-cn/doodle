package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	inputDir  string
	outputDir string
)

func init() {
	flag.StringVar(&inputDir, "i", "", "input dir")
	flag.StringVar(&outputDir, "o", "", "output dir")
}

func transform(input string, outDir string) {
	in, err := os.Open(input)
	if err != nil {
		log.Printf("open %s failed:%s", input, err.Error())
		return
	}
	defer in.Close()

	output := filepath.Join(outDir, strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))+".enb")
	out, err := os.Create(output)
	if err != nil {
		log.Printf("create output file %s failed:%s", output, err.Error())
		return
	}
	defer out.Close()

	transformedLen := doTransform(bufio.NewReader(in), bufio.NewWriter(out))
	log.Printf("transform %s[%d]\n", filepath.Base(input), transformedLen)
}

func doTransform(r *bufio.Reader, w *bufio.Writer) int {
	totolLen := 0
	buf := make([]byte, 1024)
	swapBuf := make([]byte, 512)
	for {
		n, _ := io.ReadFull(r, buf)
		if n == 0 {
			break
		}
		totolLen += n

		if n == 1024 {
			copy(swapBuf, buf[:512])
			copy(buf, buf[512:])
			copy(buf[512:], swapBuf)
		}

		wn, err := w.Write(buf[:n])
		if err != nil || n != wn {
			log.Fatalf("write file failed %v", err)
		}
	}

	w.Flush()
	return totolLen
}

func main() {
	flag.Parse()

	inputDir, err := filepath.Abs(inputDir)
	if err != nil {
		log.Fatalf("input dir %s isn't valid:%s", inputDir, err.Error())
	}
	outputDir, err := filepath.Abs(outputDir)
	if err != nil {
		log.Fatalf("output dir %s isn't valid:%s", outputDir, err.Error())
	}

	fs, err := filepath.Glob(filepath.Join(inputDir, "*"))
	if err != nil {
		log.Fatalf("get file in input dir %s failed %s", inputDir, err.Error())
	}

	var wg sync.WaitGroup
	for _, f := range fs {
		if strings.HasPrefix(filepath.Base(f), ".") {
			continue
		}
		wg.Add(1)
		go func(input, outputDir string) {
			log.Printf("%s -> %s\n", input, outputDir)
			defer wg.Done()
			//transform(input, outputDir)
		}(f, outputDir)
	}
	wg.Wait()
}
