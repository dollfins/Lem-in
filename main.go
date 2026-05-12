package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Path []string
type Ant struct {
	ID, Pos int
	Path    Path
}

func add(g map[string][]string, a, b string) {
	g[a] = append(g[a], b)
	g[b] = append(g[b], a)
}

func dfs(g map[string][]string, v, end string, vis map[string]bool, p Path, out *[]Path) {
	vis[v] = true
	p = append(p, v)
	if v == end {
		*out = append(*out, append(Path{}, p...))
		vis[v] = false
		return
	}
	for _, n := range g[v] {
		if !vis[n] {
			dfs(g, n, end, vis, p, out)
		}
	}
	vis[v] = false
}

func interceptPath(a, b Path) bool {
	m := map[string]bool{}
	for i := 1; i < len(a)-1; i++ {
		m[a[i]] = true
	}
	for i := 1; i < len(b)-1; i++ {
		if m[b[i]] {
			return true
		}
	}
	return false
}

func pathSelection(ps []Path) []Path {
	sort.Slice(ps, func(i, j int) bool { return len(ps[i]) < len(ps[j]) })
	var r []Path
	for _, p := range ps {
		ok := true
		for _, x := range r {
			if interceptPath(p, x) {
				ok = false
				break
			}
		}
		if ok {
			r = append(r, p)
		}
	}
	return r
}

func main() {
	b, _ := os.ReadFile("test1.txt")
	l := strings.Split(strings.TrimSpace(string(b)), "\n")

	n, _ := strconv.Atoi(l[0])
	re := regexp.MustCompile(`^\d+-\d+$`)

	g := map[string][]string{}
	var s, e string

	for i, v := range l {
		if v == "##start" {
			s = l[i+1]
			if s != "" {
				s = strings.Fields(l[i+1])[0]
			} else {
				fmt.Println("Error, empty ##start")
				return
			}
		}
		if v == "##end" {
			e = l[i+1]
			if e != "" {
				e = strings.Fields(l[i+1])[0]
			} else {
				fmt.Println("Error, empty ##end")
				return
			}
		}
		if re.MatchString(v) {
			p := strings.Split(v, "-")
			add(g, p[0], p[1])
		}
		if re.MatchString(v) {
			p := strings.Split(v, "-")
			if len(p) != 2 || p[0] == "" || p[1] == "" {
				fmt.Println("Error: invalid coordinate line:", v)
				return
			}
			add(g, p[0], p[1])
		}
	}

	var all []Path
	dfs(g, s, e, map[string]bool{}, nil, &all)

	if _, ok := g[s]; !ok {
		fmt.Println("Error: start node not in graph")
		return
	}

	if _, ok := g[e]; !ok {
		fmt.Println("rror: end node not in graph")
		return
	}

	if len(all) < 1 {
		fmt.Println("Error: pathway is empty")
		return
	}

	min := len(all[0])
	for _, p := range all {
		if len(p) < min {
			min = len(p)
		}
	}

	var short []Path
	for _, p := range all {
		if len(p) == min {
			short = append(short, p)
		}
	}

	paths := pathSelection(short)

	ants := []Ant{}
	next := 1
	done := 0

	for done < n {
		occ := map[string]bool{}
		var out []string

		for i := range ants {
			a := &ants[i]
			if a.Pos >= len(a.Path)-1 {
				continue
			}
			nr := a.Path[a.Pos+1]
			if nr != e && occ[nr] {
				continue
			}
			occ[nr] = true
			a.Pos++
			out = append(out, fmt.Sprintf("L%d-%s", a.ID, nr))
			if nr == e {
				done++
			}
		}

		for _, p := range paths {
			if next > n {
				break
			}
			if occ[p[1]] {
				continue
			}
			occ[p[1]] = true
			ants = append(ants, Ant{next, 1, p})
			out = append(out, fmt.Sprintf("L%d-%s", next, p[1]))
			next++
		}

		if len(out) > 0 {
			fmt.Println(strings.Join(out, " "))
		}
	}
}
