package detect

import (
	"github.com/mdshahjahanmiah/fraud-detection/pkg/procedure"
	"gonum.org/v1/gonum/mat"
	"math"
	"testing"
)

// Test_ExpectedPathLengthCoefficient tests the calculation of the expected path length coefficient.
func Test_ExpectedPathLengthCoefficient(t *testing.T) {
	size := 10
	expected := 2*math.Log(float64(size-1)) + 0.5772156649 - (2 * (float64(size-1) / float64(size)))
	result := expectedPathLengthCoefficient(size)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// Test_CalculatePathLength tests the calculation of the path length within an isolation tree.
func Test_CalculatePathLength(t *testing.T) {
	data := mat.NewDense(4, 1, []float64{1.0, 2.0, 3.0, 4.0})
	height := int(math.Ceil(math.Log2(float64(4))))
	tree := procedure.BuildTree(data, height)

	point := []float64{2.5}
	currentLength := 0.0
	result := calculatePathLength(tree, point, currentLength)
	expected := 1.0 // Initial path length increment
	for {
		if point[tree.SplitAttr] < tree.SplitValue {
			tree = tree.Left
		} else {
			tree = tree.Right
		}
		if tree.IsLeaf {
			expected += expectedPathLengthCoefficient(tree.Size)
			break
		} else {
			expected++
		}
	}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// Test_CalculateAnomalyScore tests the calculation of the anomaly score.
func Test_CalculateAnomalyScore(t *testing.T) {
	pathLength := 3.0
	dataSize := 10
	expectedPathLength := expectedPathLengthCoefficient(dataSize)
	expected := math.Pow(2, -pathLength/expectedPathLength)
	result := calculateAnomalyScore(pathLength, dataSize)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// Test_DetectAnomalies tests the detection of anomalies in the dataset using isolation trees.
func Test_DetectAnomalies(t *testing.T) {
	data := mat.NewDense(4, 1, []float64{1.0, 2.0, 3.0, 100.0})
	height := int(math.Ceil(math.Log2(float64(4))))
	numTrees := 100
	trees := make([]*procedure.IsolationTree, numTrees)
	for i := 0; i < numTrees; i++ {
		trees[i] = procedure.BuildTree(data, height)
	}

	scores := DetectAnomalies(data, trees)

	if len(scores) != 4 {
		t.Errorf("Expected %v scores, got %v", 4, len(scores))
	}

	for _, score := range scores {
		if score < 0 || score > 1 {
			t.Errorf("Anomaly score out of range: %v", score)
		}
	}
}
