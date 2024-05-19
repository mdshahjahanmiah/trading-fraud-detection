package detect

import (
	"github.com/mdshahjahanmiah/fraud-detection/pkg/procedure"
	"gonum.org/v1/gonum/mat"
	"math"
)

// expectedPathLengthCoefficient calculates the expected path length for a data point
// given the size of the dataset. This value is used as part of the anomaly score calculation.
func expectedPathLengthCoefficient(size int) float64 {
	if size <= 1 {
		return 0
	}
	return 2*math.Log(float64(size-1)) + 0.5772156649 - (2 * (float64(size-1) / float64(size)))
}

// calculatePathLength recursively computes the path length of a given data point
// within the isolation tree. The path length is used to determine the anomaly score.
func calculatePathLength(node *procedure.IsolationTree, point []float64, currentLength float64) float64 {
	if node.IsLeaf {
		return currentLength + expectedPathLengthCoefficient(node.Size)
	}
	if point[node.SplitAttr] < node.SplitValue {
		return calculatePathLength(node.Left, point, currentLength+1)
	}
	return calculatePathLength(node.Right, point, currentLength+1)
}

// calculateAnomalyScore computes the anomaly score for a data point based on its path length
// and the size of the dataset. The score indicates how anomalous the data point is.
func calculateAnomalyScore(pathLength float64, dataSize int) float64 {
	expectedPathLength := expectedPathLengthCoefficient(dataSize)
	return math.Pow(2, -pathLength/expectedPathLength)
}

// DetectAnomalies detects anomalies in the given data using an ensemble of isolation trees.
// It returns a slice of anomaly scores for each data point in the dataset.
func DetectAnomalies(data *mat.Dense, trees []*procedure.IsolationTree) []float64 {
	rows, _ := data.Dims()
	scores := make([]float64, rows)

	for i := 0; i < rows; i++ {
		point := mat.Row(nil, i, data)
		totalPathLength := 0.0
		for _, tree := range trees {
			totalPathLength += calculatePathLength(tree, point, 0)
		}
		avgPathLength := totalPathLength / float64(len(trees))
		scores[i] = calculateAnomalyScore(avgPathLength, rows)
	}
	return scores
}
