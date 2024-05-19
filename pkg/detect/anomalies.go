package detect

import (
	"github.com/mdshahjahanmiah/fraud-detection/pkg/procedure"
	"gonum.org/v1/gonum/mat"
	"math"
)

func expectedPathLengthCoefficient(size int) float64 {
	if size <= 1 {
		return 0
	}
	return 2*math.Log(float64(size-1)) + 0.5772156649 - (2 * (float64(size-1) / float64(size)))
}

func calculatePathLength(node *procedure.IsolationTree, point []float64, currentLength float64) float64 {
	if node.IsLeaf {
		return currentLength + expectedPathLengthCoefficient(node.Size)
	}
	if point[node.SplitAttr] < node.SplitValue {
		return calculatePathLength(node.Left, point, currentLength+1)
	}
	return calculatePathLength(node.Right, point, currentLength+1)
}

func calculateAnomalyScore(pathLength float64, dataSize int) float64 {
	expectedPathLength := expectedPathLengthCoefficient(dataSize)
	return math.Pow(2, -pathLength/expectedPathLength)
}

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
