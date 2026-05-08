package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Ant struct {
	ID  int
	Pos int
}

func addEdge(struc map[string][]string, x, y string) {
	struc[x] = append(struc[x], y)
	struc[y] = append(struc[y], x)
}

func pathWays(struc map[string][]string, start string, end string) string {
	bfs := start
	for string(bfs[len(bfs)-1]) != end {
		for _, n := range struc[string(bfs[len(bfs)-1])] {
			if strings.Contains(bfs, n) {
				continue
			} else {
				bfs += n
			}
		}
	}
	return bfs
}

func main() {
	file, err := os.ReadFile("test1.txt")
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	antsNum, err := strconv.Atoi(string(file[0]))
	if err != nil {
		fmt.Println(err)
		return
	}

	re := regexp.MustCompile(`^\d+-\d+$`)
	slicedFile := strings.Split(string(file), "\n")
	struc := make(map[string][]string)
	var start string
	var end string

	for i, line := range slicedFile {
		if string(line) == "##start" {
			start = string(slicedFile[i+1][0])
		} else if string(line) == "##end" {
			end = string(slicedFile[i+1][0])
		}
		if re.MatchString(line) {
			addEdge(struc, string(line[0]), string(line[2]))
		}
	}
	path := pathWays(struc, start, end)
	ants := make([]Ant, antsNum)

	for i := 0; i < antsNum; i++ {
		ants[i] = Ant{
			ID:  i + 1,
			Pos: 0,
		}
	}

	finished := 0

	for finished < len(ants) {
		occupied := make(map[int]bool)

		for i := range ants {
			ant := &ants[i]
			if ant.Pos == len(path)-1 {
				continue
			}
			next := ant.Pos + 1
			if next != len(path)-1 && occupied[next] {
				continue
			}

			occupied[next] = true
			ant.Pos = next
			fmt.Printf("L%d-%s ", ant.ID, string(path[next]))

			if ant.Pos == len(path)-1 {
				finished++
			}
		}
		fmt.Println()
	}
}
