package utils

import (
	"math"
)

const (
	epsilon = 1e-5
)

func NormalizeVector(v [768]float32) [768]float32 {
	var normalizedVec [768]float32
	var norm float32 = .0
	var sum float32 = .0

	for _, v := range v {
		sum += (v * v)
	}

	norm = float32(math.Sqrt(float64(sum)))
	for i, _ := range normalizedVec {
		normalizedVec[i] = v[i] / norm
	}

	return normalizedVec
}

func DotProduct(a, b [768]float32) float32 {
	var result float32 = 0.0
	for i, _ := range a {
		result += a[i] * b[i]
	}

	return result
}

func CosineProduct(a, b [768]float32) float32 {
	aNormalized := NormalizeVector(a)
	bNormalized := NormalizeVector(b)
	cossineValue := DotProduct(aNormalized, bNormalized)

	return cossineValue
}

// Gemma embeddings já é normalizado
func CosineProductPreNormalized(a, b [768]float32) float32 {
	cosineValue := DotProduct(a, b)
	return cosineValue
}

func EuclideanDistance(a, b [768]float32) float32 {
	var dist float32 = 0.0
	var coor float32 = 0.0

	for i, _ := range a {
		coor = a[i] - b[i]
		dist += (coor * coor)
	}

	return float32(math.Sqrt(float64(dist)))
}
