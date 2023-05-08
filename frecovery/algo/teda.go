package algo

import (
	"math"

	"gitee.com/zengtao321/frdocker/config"
	"gonum.org/v1/gonum/floats"
)

func calculatewithHistoryByTEDA(data []float64, mean []float64, sigma float64, n int64) (float64, float64, []float64, float64) {
	k := float64(n)
	// 更新mean
	floats.Scale(k-1, mean)
	floats.Add(mean, data)
	floats.Scale(1.0/k, mean)
	// 更新sigma
	sub := make([]float64, len(data))
	floats.SubTo(sub, data, mean)
	sigma = sigma*(k-1)/k + 1.0/(k-1)*floats.Dot(sub, sub)
	// 计算离心率
	ecc := (floats.Dot(sub, sub)/sigma + 1.0) / (2.0 * k)
	// 计算阈值
	threshold := (math.Pow(config.TEDA_N_SIGMA, 2) + 1.0) / (2.0 * k)
	return ecc, threshold, mean, sigma
}
