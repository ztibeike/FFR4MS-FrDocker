package algo

// 通过历史数据计算, 返回计算值、阈值、新方差、新均值
func CalculateWithHistory(data []float64, variance []float64, sigma float64, n int64) (float64, float64, []float64, float64) {
	return calculatewithHistoryByTEDA(data, variance, sigma, n)
}

func CalculateWithSample() {
}
