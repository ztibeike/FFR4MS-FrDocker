package algo

import (
	"math"

	"gitee.com/zengtao321/frdocker/config"
	"gonum.org/v1/gonum/floats"
)

func calculatewithHistoryByTEDA(data []float64, variance []float64, sigma float64, n int64) (float64, float64, []float64, float64) {
	k := float64(n)
	// 更新variance
	floats.Scale(k-1, variance)
	floats.Add(variance, data)
	floats.Scale(1.0/k, variance)
	// 更新sigma
	sub := make([]float64, len(data))
	floats.SubTo(sub, data, variance)
	sigma = sigma*(k-1)/k + 1.0/(k-1)*floats.Dot(sub, sub)
	// 计算离心率
	ecc := (floats.Norm(sub, 2)/sigma + 1.0) / (2 * k)
	// 计算阈值
	threshold := (math.Pow(config.TEDA_N_SIGMA, 2) + 1.0) / (2 * k)
	return ecc, threshold, variance, sigma
}
