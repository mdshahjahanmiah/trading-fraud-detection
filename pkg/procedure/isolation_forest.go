package procedure

import (
	"gonum.org/v1/gonum/mat"
	"log/slog"
	"math/rand"
)

type IsolationTree struct {
	Left, Right *IsolationTree
	SplitAttr   int
	SplitValue  float64
	Size        int
	IsLeaf      bool
}

func BuildTree(data *mat.Dense, height int) *IsolationTree {
	rows, cols := data.Dims()
	if rows <= 1 || height <= 0 {
		return &IsolationTree{Size: rows, IsLeaf: true}
	}

	splitAttr := rand.Intn(cols)
	col := mat.Col(nil, splitAttr, data)
	minVal, maxVal := minMax(col)

	if minVal == maxVal {
		return &IsolationTree{Size: rows, IsLeaf: true}
	}

	splitValue := minVal + rand.Float64()*(maxVal-minVal)
	leftData, rightData := partitionData(data, splitAttr, splitValue)

	leftRows, leftCols := leftData.Dims()
	rightRows, rightCols := rightData.Dims()

	slog.Debug("partitioned data dimensions", "leftRows", leftRows, "leftCols", leftCols, "rightRows", rightRows, "rightCols", rightCols)

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

func partitionData(data *mat.Dense, splitAttr int, splitValue float64) (*mat.Dense, *mat.Dense) {
	rows, cols := data.Dims()
	leftRows := [][]float64{}
	rightRows := [][]float64{}

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
