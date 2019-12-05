package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

func wget(urlsList, log string) {
	wg := &sync.WaitGroup{}

	var urlsFile *os.File
	logFilename := log

	urlsFile, err := os.Open(urlsList)

	if err != nil {
		panic("Can't open the URL list file.")
	}

	var lines []string
	s := bufio.NewScanner(urlsFile)

	for s.Scan() {
		lines = append(lines, s.Text())
	}

	for _, url := range lines {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// wgetする。log.datにログが残る。
			cmd := exec.Command("cmd.exe", "/c", "wget -p -k -nH -np -E --append-output="+logFilename, url) // -r で再帰、-ndでディレクトリを作成しない
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}(url)
	}
	wg.Wait()
}

func main() {
	var saveDir string
	var urlsList string
	var log = "log.dat"
	var logFile *os.File

	// wgetしてくるファイルの保存場所
	fmt.Print("Save folder: ")
	fmt.Scan(&saveDir)

	os.Chdir(saveDir) // 保存場所に移動する

	// URLリスト
	fmt.Print("URL list: ")
	fmt.Scan(&urlsList)

	// プログレスバー
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Color("green")
	s.Start()

	wget(urlsList, log)

	logFile, err := os.Open(log)

	if err != nil {
		panic("Can't open the log.dat file.")
	}

	logScanner := bufio.NewScanner(logFile)

	// ログにfailedがあったら出力
	for logScanner.Scan() {
		logText := logScanner.Text()

		// failedが見つかったら-1以外になる
		if strings.Index(logText, "failed") != -1 {
			fmt.Printf("\x1b[31m%s\x1b[0m", "\nError:\n")
			fmt.Println(logText)
		}
	}

	s.Stop() // プログレスバー終了

	fmt.Print("\nDone!\n")
}
