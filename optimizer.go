package colormix

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
)

// mixColors 计算混合颜色
func mixColors(weights *mat.VecDense, colorTable *mat.Dense) *mat.VecDense {
	_, cols := colorTable.Dims()
	mixed := mat.NewVecDense(cols, nil)
	mixed.MulVec(colorTable.T(), weights)
	return mixed
}

// colorDistance 计算颜色误差（欧几里得距离）
func colorDistance(weights *mat.VecDense, colorTable *mat.Dense, targetColor *mat.VecDense) float64 {
	mixed := mixColors(weights, colorTable)
	_, cols := colorTable.Dims()
	diff := mat.NewVecDense(cols, nil)
	diff.SubVec(mixed, targetColor)

	// 欧几里得距离
	dist := mat.Dot(diff, diff)

	// L1 正则化项，防止某些颜色被忽略
	// l1Penalty := 0.01 * mat.Sum(weights) // 0.01 是调节参数，可以调整
	// 约束惩罚项，确保 sum(weights) ≈ 1
	sum := 0.0
	for i := 0; i < weights.Len(); i++ {
		sum += weights.AtVec(i)
	}
	penalty := 100 * (sum - 1) * (sum - 1) // 约束惩罚项，确保 sum(weights) ≈ 1

	// 防止某些颜色权重为0的惩罚
	zeroPenalty := 0.0
	for i := 0; i < weights.Len(); i++ {
		if weights.AtVec(i) < 0.01 {
			zeroPenalty += 1000.0 // 强制避免某些颜色权重过低
		}
	}
	return dist + penalty + zeroPenalty
}

// 计算目标函数的梯度（colorDistance 的梯度）
func colorGradient(grad []float64, weights *mat.VecDense, colorTable *mat.Dense, targetColor *mat.VecDense) {
	mixed := mixColors(weights, colorTable)
	diff := mat.NewVecDense(3, nil)
	diff.SubVec(mixed, targetColor)

	// 计算梯度：梯度是误差的导数，即颜色误差的导数
	gradVec := mat.NewVecDense(weights.Len(), grad)
	gradVec.MulVec(colorTable, diff) // 这里是乘以颜色表

	// 加上约束的梯度惩罚项
	for i := 0; i < weights.Len(); i++ {
		gradVec.SetVec(i, gradVec.AtVec(i)+200*(weights.AtVec(i)-1)) // 惩罚项对梯度的影响
	}

	// 防止零惩罚项对梯度的影响
	for i := 0; i < weights.Len(); i++ {
		if weights.AtVec(i) < 0.01 {
			gradVec.SetVec(i, gradVec.AtVec(i)+2000) // 强制避免梯度过大
		}
	}
	// 将计算出来的梯度更新回传入的 grad 参数
	copy(grad, gradVec.RawVector().Data)
}

// computeColorRatios 计算颜色表的最佳混合比例
func computeColorRatios(colorTable *mat.Dense, targetColor *mat.VecDense) (*mat.VecDense, error) {
	numColors, _ := colorTable.Dims()

	// 初始猜测 (均分权重)
	initialWeights := mat.NewVecDense(numColors, nil)
	for i := 0; i < numColors; i++ {
		initialWeights.SetVec(i, 1.0/float64(numColors)) // 初始均匀分配
	}

	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			weights := mat.NewVecDense(len(x), x)
			return colorDistance(weights, colorTable, targetColor)
		},
		Grad: func(grad, x []float64) {
			weights := mat.NewVecDense(len(x), x)
			colorGradient(grad, weights, colorTable, targetColor)
		},
	}

	settings := &optimize.Settings{}
	// method := &optimize.LBFGS{}
	// 使用 Nelder-Mead 优化方法替代 LBFGS
	method := &optimize.NelderMead{}
	result, err := optimize.Minimize(problem, initialWeights.RawVector().Data, settings, method)
	if err != nil {
		return nil, fmt.Errorf("optimization failed: %w", err)
	}

	// 归一化（确保总和为 1）
	finalWeights := mat.NewVecDense(numColors, result.X)
	sum := 0.0
	for i := 0; i < numColors; i++ {
		sum += finalWeights.AtVec(i)
	}
	for i := 0; i < numColors; i++ {
		finalWeights.SetVec(i, finalWeights.AtVec(i)/sum)
	}

	return finalWeights, nil
}
