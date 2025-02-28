//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"io"
	"os"
	"sort"
	"strings"
)

func NewReader(input io.Reader) LineReader {
	return &customLineReader{
		bufferedReader: bufio.NewReader(input),
	}
}

func NewWriter(output io.Writer) LineWriter {
	return &customLineWriter{
		baseWriter: output,
	}
}

func Sort(output io.Writer, files ...string) error {
	var tempReaders []LineReader
	for _, fileName := range files {
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		lines, err := readAndSortLines(file)
		if err != nil {
			return err
		}

		tempFile, err := writeLinesToTempFile(lines)
		if err != nil {
			return err
		}
		defer func(tempFile *os.File) {
			err := tempFile.Close()
			if err != nil {

			}
		}(tempFile)

		tempReaders = append(tempReaders, NewReader(tempFile))
	}
	return Merge(NewWriter(output), tempReaders...)
}

func Merge(writer LineWriter, readers ...LineReader) error {
	priorityQueue := &mergePriorityQueue{}
	for _, reader := range readers {
		if line, err := reader.ReadLine(); err == nil {
			heap.Push(priorityQueue, mergeItem{line: line, source: reader})
		} else if err != io.EOF {
			return err
		}
	}
	heap.Init(priorityQueue)

	for priorityQueue.Len() > 0 {
		item := heap.Pop(priorityQueue).(mergeItem)
		if err := writer.Write(item.line); err != nil {
			return err
		}
		if nextLine, err := item.source.ReadLine(); err == nil {
			heap.Push(priorityQueue, mergeItem{line: nextLine, source: item.source})
		} else if err != io.EOF {
			return err
		}
	}
	return nil
}

func readAndSortLines(input io.Reader) ([]string, error) {
	reader := NewReader(input)
	var lines []string
	for {
		line, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	sort.Strings(lines)
	return lines, nil
}

func writeLinesToTempFile(lines []string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "sorted_")
	if err != nil {
		return nil, err
	}
	writer := NewWriter(tempFile)
	for _, line := range lines {
		if err := writer.Write(line); err != nil {
			return nil, err
		}
	}
	_, err = tempFile.Seek(0, io.SeekStart)
	return tempFile, err
}

type customLineReader struct {
	bufferedReader *bufio.Reader
}

func (cLineReader *customLineReader) ReadLine() (string, error) {
	var lineBuilder strings.Builder
	for {
		char, err := cLineReader.bufferedReader.ReadByte()
		if err != nil {
			if lineBuilder.Len() > 0 {
				return lineBuilder.String(), nil
			}
			return "", err
		}
		if char == '\n' {
			break
		}
		lineBuilder.WriteByte(char)
	}
	return lineBuilder.String(), nil
}

type customLineWriter struct {
	baseWriter io.Writer
}

func (cLineWriter *customLineWriter) Write(line string) error {
	_, err := cLineWriter.baseWriter.Write([]byte(line + "\n"))
	return err
}

type mergeItem struct {
	line   string
	source LineReader
}

type mergePriorityQueue []mergeItem

func (priorityQueue mergePriorityQueue) Len() int { return len(priorityQueue) }
func (priorityQueue mergePriorityQueue) Less(i, j int) bool {
	return priorityQueue[i].line < priorityQueue[j].line
}
func (priorityQueue mergePriorityQueue) Swap(i, j int) {
	priorityQueue[i], priorityQueue[j] = priorityQueue[j], priorityQueue[i]
}

func (priorityQueue *mergePriorityQueue) Push(item interface{}) {
	*priorityQueue = append(*priorityQueue, item.(mergeItem))
}

func (priorityQueue *mergePriorityQueue) Pop() interface{} {
	old := *priorityQueue
	n := len(old)
	element := old[n-1]
	*priorityQueue = old[:n-1]
	return element
}
