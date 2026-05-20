package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Ant struct {
	ID, Pos int
	Path    []string
}

type Path struct {
	Rooms []string
	Ants  int
}

func addEdge(g map[string][]string, a, b string) {
	g[a] = append(g[a], b)
	g[b] = append(g[b], a)
}

func contains(s []string, t string) bool {
	for _, v := range s {
		if v == t {
			return true
		}
	}
	return false
}

func dfs(g map[string][]string, start, end string, vis map[string]bool, path []string, paths *[][]string) {
	vis[start] = true
	path = append(path, start)

	if start == end {
		*paths = append(*paths, append([]string{}, path...))
		vis[start] = false
		return
	}

	for _, n := range g[start] {
		if !vis[n] {
			dfs(g, n, end, vis, path, paths)
		}
	}

	vis[start] = false
}

func interceptPath(a, b []string) bool {
	for i := 1; i < len(a)-1; i++ {
		if contains(b[1:len(b)-1], a[i]) {
			return true
		}
	}
	return false
}

func turns(p []Path, ants int) int {
	t := make([]Path, len(p))

	for i := 0; i < ants; i++ {
		b := 0

		for j := range t {
			if len(t[j].Rooms)+t[j].Ants < len(t[b].Rooms)+t[b].Ants {
				b = j
			}
		}

		t[b].Ants++
	}

	m := 0

	for _, v := range t {
		if v.Ants == 0 {
			continue
		}

		x := len(v.Rooms) + v.Ants - 2

		if x > m {
			m = x
		}
	}
	return m
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Correct format: go run (filename) (textfile)")
		return
	}
	filePath := os.Args[1]
	file, err := os.ReadFile(filePath + ".txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(file) == 0 {
		fmt.Println("Error: Empty File Content")
		return
	}

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")

	antsNum, err := strconv.Atoi(lines[0])
	if err != nil || antsNum <= 0 {
		fmt.Println("Error: invalid ant count")
		return
	}

	g := map[string][]string{}
	re := regexp.MustCompile(`^\d+-\d+$`)

	var start, end string

	for i, l := range lines {
		if l == "##start" {
			if i+1 >= len(lines) || len(strings.Fields(lines[i+1])) == 0 {
				fmt.Println("Error: invalid start")
				return
			}
			start = strings.Fields(lines[i+1])[0]
		}

		if l == "##end" {
			if i+1 >= len(lines) || len(strings.Fields(lines[i+1])) == 0 {
				fmt.Println("Error: invalid end")
				return
			}
			end = strings.Fields(lines[i+1])[0]
		}

		if re.MatchString(l) {
			p := strings.Split(l, "-")

			if len(p) != 2 || p[0] == "" || p[1] == "" {
				fmt.Println("Error: invalid coordinate")
				return
			}

			addEdge(g, p[0], p[1])
		}
	}

	if start == "" || end == "" {
		fmt.Println("Error: missing start point/ending point")
		return
	}

	if start == end {
		fmt.Println("Error: starting point is same as ending point")
		return
	}

	if _, ok := g[start]; !ok {
		fmt.Println("Error: invalid start")
		return
	}

	if _, ok := g[end]; !ok {
		fmt.Println("Error: invalid end")
		return
	}

	var allPathWays [][]string

	dfs(g, start, end, map[string]bool{}, []string{}, &allPathWays)

	if len(allPathWays) == 0 {
		fmt.Println("Error: no path found")
		return
	}

	fmt.Println(allPathWays)

	sort.Slice(allPathWays, func(i, j int) bool {
		return len(allPathWays[i]) < len(allPathWays[j])
	})

	use := []Path{{Rooms: allPathWays[0]}}

	best := turns(use, antsNum)

	for _, p := range allPathWays[1:] {
		ok := true

		for _, u := range use {
			if interceptPath(u.Rooms, p) {
				ok = false
				break
			}
		}

		if !ok {
			continue
		}

		tmp := append(append([]Path{}, use...), Path{Rooms: p})

		if t := turns(tmp, antsNum); t < best {
			best = t
			use = tmp
		}
	}

	assign := [][]string{}
	tmp := make([]Path, len(use))
	copy(tmp, use)

	for range antsNum {
		b := 0

		for j := range tmp {
			if len(tmp[j].Rooms)+tmp[j].Ants < len(tmp[b].Rooms)+tmp[b].Ants {
				b = j
			}
		}

		tmp[b].Ants++
		assign = append(assign, tmp[b].Rooms)
	}

	ants := make([]Ant, antsNum)

	for i := range ants {
		ants[i] = Ant{i + 1, 0, assign[i]}
	}

	finished := 0

	for finished < antsNum {
		occ := map[string]bool{}
		var out []string

		for i := range ants {
			a := &ants[i]

			if a.Pos >= len(a.Path)-1 {
				continue
			}

			next := a.Path[a.Pos+1]

			if a.Pos == 0 && occ[next] {
				continue
			}

			if next != end && occ[next] {
				continue
			}

			occ[next] = true
			a.Pos++

			out = append(out, fmt.Sprintf("L%d-%s", a.ID, next))

			if next == end {
				finished++
			}
		}

		if len(out) > 0 {
			fmt.Println(strings.Join(out, " "))
		}
	}
}
