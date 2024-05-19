package procedure

import (
	"gonum.org/v1/gonum/mat"
	"math/rand"
	"testing"
)

func Test_BuildTree(t *testing.T) {
	rand.Seed(42) // Seed the random number generator for reproducibility

	// Create a small dataset for testing
	data := mat.NewDense(6, 2, []float64{
		1, 2,
		3, 4,
		5, 6,
		7, 8,
		9, 10,
		11, 12,
	})

	height := 3
	tree := BuildTree(data, height)

	if tree == nil {
		t.Errorf("expected non-nil tree, got nil")
	}
	if tree.Size != 6 {
		t.Errorf("expected tree size to be 6, got %d", tree.Size)
	}
	if tree.IsLeaf {
		t.Errorf("expected tree to not be a leaf, but it is")
	}
}

func Test_MinMax(t *testing.T) {
	array := []float64{1, 2, 3, 4, 5, 6}
	min, max := minMax(array)

	if min != 1 {
		t.Errorf("expected min to be 1, got %f", min)
	}
	if max != 6 {
		t.Errorf("expected max to be 6, got %f", max)
	}
}

func Test_PartitionData(t *testing.T) {
	data := mat.NewDense(6, 2, []float64{
		1, 2,
		3, 4,
		5, 6,
		7, 8,
		9, 10,
		11, 12,
	})

	splitAttr := 0
	splitValue := 5.0
	leftData, rightData := partitionData(data, splitAttr, splitValue)

	leftRows, leftCols := leftData.Dims()
	rightRows, rightCols := rightData.Dims()

	if leftRows != 2 || leftCols != 2 {
		t.Errorf("expected left partition to have dimensions 2x2, got %dx%d", leftRows, leftCols)
	}
	if rightRows != 4 || rightCols != 2 {
		t.Errorf("expected right partition to have dimensions 4x2, got %dx%d", rightRows, rightCols)
	}

	// Check the values in the partitions
	expectedLeft := mat.NewDense(2, 2, []float64{
		1, 2,
		3, 4,
	})
	expectedRight := mat.NewDense(4, 2, []float64{
		5, 6,
		7, 8,
		9, 10,
		11, 12,
	})

	if !mat.Equal(leftData, expectedLeft) {
		t.Errorf("expected left partition to be \n%v\nbut got \n%v", mat.Formatted(expectedLeft), mat.Formatted(leftData))
	}
	if !mat.Equal(rightData, expectedRight) {
		t.Errorf("expected right partition to be \n%v\nbut got \n%v", mat.Formatted(expectedRight), mat.Formatted(rightData))
	}
}
