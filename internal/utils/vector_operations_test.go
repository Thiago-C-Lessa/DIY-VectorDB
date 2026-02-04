package utils

import (
	"math"
	"math/rand/v2"
	"testing"
)

func almostEqual(a, b float32) bool {
	return math.Abs(float64(a-b)) < epsilon
}

func TestDotProduct(t *testing.T) {
	var a, b [768]float32

	a[0] = 1
	a[1] = 2
	a[2] = 3

	b[0] = 4
	b[1] = 5
	b[2] = 6

	var expected float32 = 32.0
	result := DotProduct(a, b)

	if result != expected {
		t.Errorf("DotProduct failed: expected %v, got %v", expected, result)
	}
}

func TestCosineProductPreNormalized(t *testing.T) {
	var a, b [768]float32
	var expected float32 = 1.0
	var f float32 = .0

	for i, _ := range a {
		f = rand.Float32()
		a[i] = f
		b[i] = f
	}

	a = NormalizeVector(a)
	b = NormalizeVector(b)

	result := CosineProductPreNormalized(a, b)

	if !almostEqual(result, expected) {
		t.Errorf("CosineProductPreNormalized failed: expected %v, got %v", expected, result)
	}
}

func TestNormalizeVector(t *testing.T) {
	var v [768]float32
	v[0] = 3
	v[1] = 4

	normalized := NormalizeVector(v)

	var norm float32 = 0.0
	for _, val := range normalized {
		norm += val * val
	}
	norm = float32(math.Sqrt(float64(norm)))

	if !almostEqual(norm, 1.0) {
		t.Errorf("NormalizeVector failed: expected norm 1, got %v", norm)
	}
}

func TestCosineProduct_IdenticalVectors(t *testing.T) {
	var v [768]float32
	v[0] = 1
	v[1] = 2
	v[2] = 3

	result := CosineProduct(v, v)

	if !almostEqual(result, 1.0) {
		t.Errorf("CosineProduct identical vectors: expected 1, got %v", result)
	}
}

func TestCosineProduct_OrthogonalVectors(t *testing.T) {
	var a, b [768]float32

	a[0] = 1
	a[1] = 0

	b[0] = 0
	b[1] = 1

	result := CosineProduct(a, b)

	if !almostEqual(result, 0.0) {
		t.Errorf("CosineProduct orthogonal vectors: expected 0, got %v", result)
	}
}

func TestCosineProduct_KnownValue(t *testing.T) {
	var a, b [768]float32

	a[0] = 1
	a[1] = 0

	b[0] = 0.8
	b[1] = 0.6

	result := CosineProduct(a, b)

	var expected float32 = 0.8
	if !almostEqual(result, expected) {
		t.Errorf("CosineProduct known value: expected %v, got %v", expected, result)
	}
}

func TestMinUint64(t *testing.T) {
	var a uint64 = math.MaxUint64
	var b uint64 = 1
	var c uint64 = MinUint64(a, b)
	if c != b {
		t.Errorf("MinUint64 failed: expected %v, got %v", b, c)
	}
}

func TestEuclideanDistance(t *testing.T) {
	var a, b [768]float32

	// Caso 1: vetores iguais → distância = 0
	for i := 0; i < 768; i++ {
		a[i] = 1.0
		b[i] = 1.0
	}

	dist := EuclideanDistance(a, b)
	if dist != 0 {
		t.Errorf("expected distance 0, got %f", dist)
	}

	// Caso 2: um vetor zero, outro unitário em uma dimensão
	a = [768]float32{}
	b = [768]float32{}
	b[0] = 1.0

	dist = EuclideanDistance(a, b)
	if math.Abs(float64(dist-1.0)) > 1e-6 {
		t.Errorf("expected distance 1, got %f", dist)
	}

	// Caso 3: vetores opostos normalizados
	// distância esperada = sqrt( (1 - (-1))^2 ) = 2
	a = [768]float32{}
	b = [768]float32{}
	a[0] = 1.0
	b[0] = -1.0

	dist = EuclideanDistance(a, b)
	if math.Abs(float64(dist-2.0)) > 1e-6 {
		t.Errorf("expected distance 2, got %f", dist)
	}

	// Caso 4: simetria
	a = [768]float32{}
	b = [768]float32{}
	a[10] = 0.3
	b[10] = -0.7

	d1 := EuclideanDistance(a, b)
	d2 := EuclideanDistance(b, a)

	if math.Abs(float64(d1-d2)) > 1e-6 {
		t.Errorf("distance not symmetric: %f vs %f", d1, d2)
	}

	// Caso 5: consistência com soma dos quadrados
	a = [768]float32{}
	b = [768]float32{}
	a[0] = 0.6
	b[0] = 0.8

	expected := float32(math.Sqrt(float64((0.6 - 0.8) * (0.6 - 0.8))))
	dist = EuclideanDistance(a, b)

	if math.Abs(float64(dist-expected)) > 1e-6 {
		t.Errorf("expected %f, got %f", expected, dist)
	}
}

func BenchmarkNormalizeVector(b *testing.B) {
	var v [768]float32

	for i := range v {
		v[i] = rand.Float32()
	}

	b.ResetTimer()
	for range b.N {
		result := NormalizeVector(v)
		result[2] = 0.0
		result[3] = 0.3
		result[4] -= 0.1
	}
}

func BenchmarkCosineProduct(b *testing.B) {
	var v [768]float32

	for i := range v {
		v[i] = rand.Float32()
	}
	b.ResetTimer()

	for range b.N {
		result := CosineProduct(v, v)
		result += .1
		result -= .2
		result += .1
	}
}

func BenchmarkCosineProductPreNormalized(b *testing.B) {
	var v [768]float32

	for i := range v {
		v[i] = rand.Float32()
	}
	v = NormalizeVector(v)

	b.ResetTimer()
	for range b.N {
		result := CosineProductPreNormalized(v, v)
		result += .1
		result -= .2
		result += .4
	}
}

func BenchmarkEuclideanDistance(b *testing.B) {
	var v [768]float32
	for i := range v {
		v[i] = rand.Float32()
	}

	b.ResetTimer()

	for range b.N {
		result := EuclideanDistance(v, v)
		result += .1
		result -= .2
		result += .4
	}
}
