package algo

// 通过历史数据计算, 返回计算值、阈值、新均值、新方差
func CalculateWithHistory(data []float64, mean []float64, sigma float64, n int64) (float64, float64, []float64, float64) {
	return calculatewithHistoryByTEDA(data, mean, sigma, n)
}

// 根据总体样本计算, 返回每个样本的离心率和总体阈值
func CalculateWithSample(data [][]float64) ([]float64, float64) {
	return calculateWithSampleByTEDA(data)
}
