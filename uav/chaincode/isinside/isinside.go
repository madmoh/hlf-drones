package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"gonum.org/v1/gonum/spatial/kdtree"
	"gonum.org/v1/gonum/spatial/r3"
)

func main() {
	v, f := ParseOBJ("armadillo.obj")
	q := make([][3]float64, 100*100*100)
	for i := range q {
		q[i] = [3]float64{-1.377701, -1.285421, -1.947002}
	}
	fmt.Println(InOrOut(v, f, q))
}

func ParseOBJ(filepath string) ([][3]float64, [][3]uint64) {
	file, _ := os.Open(filepath)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	vCount, fCount := 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == 'v' {
			vCount = vCount + 1
		} else if line[0] == 'f' {
			fCount = fCount + 1
		}
	}

	file.Seek(0, io.SeekStart)
	scanner = bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	vertices := make([][3]float64, vCount)
	facets := make([][3]uint64, fCount)
	vIndex, fIndex := 0, 0
	for scanner.Scan() {
		vorf := scanner.Text()
		if vorf == "v" {
			scanner.Scan()
			p0, _ := strconv.ParseFloat(scanner.Text(), 64)
			scanner.Scan()
			p1, _ := strconv.ParseFloat(scanner.Text(), 64)
			scanner.Scan()
			p2, _ := strconv.ParseFloat(scanner.Text(), 64)
			vertices[vIndex] = [3]float64{p0, p1, p2}
			vIndex = vIndex + 1
		}
		if vorf == "f" {
			scanner.Scan()
			v0, _ := strconv.ParseUint(scanner.Text(), 10, 64)
			scanner.Scan()
			v1, _ := strconv.ParseUint(scanner.Text(), 10, 64)
			scanner.Scan()
			v2, _ := strconv.ParseUint(scanner.Text(), 10, 64)
			facets[fIndex] = [3]uint64{v0 - 1, v1 - 1, v2 - 1}
			fIndex = fIndex + 1
		}
	}

	return vertices, facets
}

func GetIncenterNormal(vertices [3][3]float64) ([3]float64, [3]float64) {
	v0 := r3.Vec{X: vertices[0][0], Y: vertices[0][1], Z: vertices[0][2]}
	v1 := r3.Vec{X: vertices[1][0], Y: vertices[1][1], Z: vertices[1][2]}
	v2 := r3.Vec{X: vertices[2][0], Y: vertices[2][1], Z: vertices[2][2]}

	// Calculating incenter
	a := r3.Norm(r3.Sub(v1, v2))
	b := r3.Norm(r3.Sub(v2, v0))
	c := r3.Norm(r3.Sub(v0, v1))
	abc := a + b + c
	iV0 := r3.Scale(a/abc, v0)
	iV1 := r3.Scale(b/abc, v1)
	iV2 := r3.Scale(c/abc, v2)
	iRes := r3.Add(iV0, r3.Add(iV1, iV2))

	// Calculating normal
	nV0 := r3.Sub(v1, v0)
	nV1 := r3.Sub(v2, v1)
	nV2 := r3.Cross(nV0, nV1)
	nRes := r3.Scale(1/r3.Norm(nV2), nV2)

	return [3]float64{iRes.X, iRes.Y, iRes.Z}, [3]float64{nRes.X, nRes.Y, nRes.Z}
}

func GetIncentersNormals(vertices [][3]float64, facets [][3]uint64) ([][3]float64, [][3]float64) {
	incenters := make([][3]float64, len(facets))
	normals := make([][3]float64, len(facets))
	for i := 0; i < len(facets); i++ {
		fVertices := [3][3]float64{
			vertices[facets[i][0]],
			vertices[facets[i][1]],
			vertices[facets[i][2]],
		}
		incenters[i], normals[i] = GetIncenterNormal(fVertices)
	}
	return incenters, normals
}

func IndexOf(value [3]float64, array [][3]float64) int {
	for i, v := range array {
		if v[0] == value[0] && v[1] == value[1] && v[2] == value[2] {
			return i
		}
	}
	return -1
}

func GetDistance(vertex [3]float64, normal [3]float64, query [3]float64) float64 {
	vVec := r3.Vec{X: vertex[0], Y: vertex[1], Z: vertex[2]}
	nVec := r3.Vec{X: normal[0], Y: normal[1], Z: normal[2]}
	qVec := r3.Vec{X: query[0], Y: query[1], Z: query[2]}
	dist := r3.Sub(qVec, vVec)
	dist = r3.Vec{X: dist.X * nVec.X, Y: dist.Y * nVec.Y, Z: dist.Z * nVec.Z}
	return dist.X + dist.Y + dist.Z
}

func InOrOut(vertices [][3]float64, facets [][3]uint64, query [][3]float64) []float64 {
	incenters, normals := GetIncentersNormals(vertices, facets)
	incentersPoints := make([]kdtree.Point, len(facets))
	for i := 0; i < len(facets); i++ {
		incentersPoints[i] = kdtree.Point(incenters[i][:])
	}
	tree := kdtree.New(kdtree.Points(incentersPoints), false)
	distances := make([]float64, len(query))
	for i := 0; i < len(query); i++ {
		fmt.Println(i)
		closestPoint, _ := tree.Nearest(kdtree.Point(query[i][:]))
		closest := [3]float64(closestPoint.(kdtree.Point))
		closestIndex := IndexOf(closest, incenters)
		distances[i] = GetDistance(closest, normals[closestIndex], query[i])
	}
	return distances
}
