package procedure

import (
	"gonum.org/v1/gonum/mat"
	"log/slog"
	"math/rand"
)

// IsolationTree represents a node in the isolation tree.
type IsolationTree struct {
	Left, Right *IsolationTree
	SplitAttr   int
	SplitValue  float64
	Size        int
	IsLeaf      bool
}

// BuildTree recursively builds an isolation tree given data and a height limit.
// The data matrix is partitioned based on randomly selected attributes and split values.
func BuildTree(data *mat.Dense, height int) *IsolationTree {
	rows, cols := data.Dims()
	if rows <= 1 || height <= 0 {
		return &IsolationTree{Size: rows, IsLeaf: true}
	}

	splitAttr := rand.Intn(cols)
	col := mat.Col(nil, splitAttr, data)
	minVal, maxVal := minMax(col)

	// If all values are the same, we can't split further
	if minVal == maxVal {
		return &IsolationTree{Size: rows, IsLeaf: true}
	}

	splitValue := minVal + rand.Float64()*(maxVal-minVal)
	leftData, rightData := partitionData(data, splitAttr, splitValue)

	leftRows, leftCols := leftData.Dims()
	rightRows, rightCols := rightData.Dims()

	slog.Debug("partitioned data dimensions", "leftRows", leftRows, "leftCols", leftCols, "rightRows", rightRows, "rightCols", rightCols)

	// Handle an edge case where one partition is empty
	if leftRows == 0 || rightRows == 0 {
		slog.Error("partition resulted in zero length matrix dimension.")
		return &IsolationTree{Size: rows, IsLeaf: true}
	}

	return &IsolationTree{
		Left:       BuildTree(leftData, height-1),
		Right:      BuildTree(rightData, height-1),
		SplitAttr:  splitAttr,
		SplitValue: splitValue,
		Size:       rows,
	}
}

// minMax finds the minimum and maximum values in a slice of float64.
func minMax(array []float64) (float64, float64) {
	minVal, maxVal := array[0], array[0]
	for _, v := range array {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	return minVal, maxVal
}

// partitionData splits the data matrix into two based on the specified attribute and value.
func partitionData(data *mat.Dense, splitAttr int, splitValue float64) (*mat.Dense, *mat.Dense) {
	rows, cols := data.Dims()
	var leftRows [][]float64
	var rightRows [][]float64

	for i := 0; i < rows; i++ {
		row := mat.Row(nil, i, data)
		if row[splitAttr] < splitValue {
			leftRows = append(leftRows, row)
		} else {
			rightRows = append(rightRows, row)
		}
	}

	if len(leftRows) == 0 || len(rightRows) == 0 {
		slog.Warn("partition resulted in zero rows for one of the partitions.")
	}

	leftData := mat.NewDense(len(leftRows), cols, nil)
	rightData := mat.NewDense(len(rightRows), cols, nil)

	for i, row := range leftRows {
		leftData.SetRow(i, row)
	}
	for i, row := range rightRows {
		rightData.SetRow(i, row)
	}

	return leftData, rightData
}
