package algo

// 通过历史数据计算, 返回计算值、阈值、新均值、新方差
func CalculateWithHistory(data []float64, mean []float64, sigma float64, n int64) (float64, float64, []float64, float64) {
	return calculatewithHistoryByTEDA(data, mean, sigma, n)
}

func CalculateWithSample() {
}
