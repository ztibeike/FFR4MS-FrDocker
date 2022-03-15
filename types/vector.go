package types

import (
	"errors"
	"math"
)

type Vector struct {
	Data []float64
}

func NewVector(data []float64) *Vector {
	v := &Vector{}
	copy(v.Data, data)
	return v
}

func (v *Vector) Len() int {
	return len(v.Data)
}

func (v *Vector) AddVec(b *Vector) *Vector {
	if v.Len() != b.Len() {
		panic(errors.New("error length"))
	}
	var res = &Vector{
		Data: make([]float64, v.Len()),
	}
	for idx, val := range v.Data {
		res.Data[idx] = b.Data[idx] + val
	}
	return res
}

func (v *Vector) SubVec(b *Vector) *Vector {
	if v.Len() != b.Len() {
		panic(errors.New("error length"))
	}
	var res = &Vector{
		Data: make([]float64, v.Len()),
	}
	for idx, val := range v.Data {
		res.Data[idx] = val - b.Data[idx]
	}
	return res
}

func (v *Vector) ScaleVec(a float64) *Vector {
	var res = &Vector{
		Data: make([]float64, v.Len()),
	}
	for idx, val := range v.Data {
		res.Data[idx] = val * a
	}
	return res
}

func (v *Vector) MulVec(b *Vector) float64 {
	var res float64
	for idx, val := range v.Data {
		res += val * b.Data[idx]
	}
	return res
}

func (v *Vector) Copy() *Vector {
	var res = &Vector{
		Data: make([]float64, v.Len()),
	}
	copy(res.Data, v.Data)
	return res
}

func (v *Vector) T() *Vector {
	return v.Copy()
}

// 二范数
func (v *Vector) Norm() float64 {
	n := v.Len()
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += math.Pow(v.Data[i], 2)
	}
	return math.Sqrt(sum)
}
