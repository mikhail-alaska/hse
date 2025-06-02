package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
)

const (
	N = 16
	R = 1.0 / 2.0
)

var (
	k = int(float64(N) * R)
	n = int(math.Log2(N))
	frozenPositions = []int{0, 1, 2, 3, 4, 5, 8, 9}
)

type Path struct {
	metric float64
	u      []int
}

func main() {
	fmt.Println("Что будет происходить:")
	fmt.Println("1) Программа сгенерирует случайные информационные биты длины k =", k)
	fmt.Println("2) Построит полярный код (матрица Арикана), закодирует биты, выполнит BPSK-модуляцию")
	fmt.Println("3) Добавит шум и искусственные ошибки в заданных позициях")
	fmt.Println("4) Выполнит SCL-декодирование и покажет, сколько ошибок осталось в информационных битах")
	fmt.Println()
	fmt.Println("Нажмите Enter, чтобы запустить симуляцию.")
	bufio.NewReader(os.Stdin).ReadString('\n')


	F := make(map[int]struct{}, len(frozenPositions))
	for _, idx := range frozenPositions {
		F[idx] = struct{}{}
	}
	A := make([]int, 0, N-len(F))
	for i := 0; i < N; i++ {
		if _, isFrozen := F[i]; !isFrozen {
			A = append(A, i)
		}
	}
	sort.Ints(A) 

	infoBits := make([]int, k)
	for i := range infoBits {
		infoBits[i] = rand.Intn(2)
	}
	fmt.Printf("Информационные биты: %v\n", infoBits)
	fmt.Printf("Информационные позиции (A): %v\n", A)

	uFull := make([]int, N)
	for j, pos := range A {
		uFull[pos] = infoBits[j]
	}
	fmt.Println("Вектор u (с замороженными битами):")
	fmt.Printf("%v\n", uFull)

	G := genArikanMatrix(n)

	codeword := polarEnc(uFull, G)
	fmt.Printf("Закодированное слово: %v\n", codeword)

	sig := bpsk(codeword)
	fmt.Printf("Модулированный сигнал: %v\n", roundSlice(sig, 4))

	sigma := 0.1
	errorCounts := []int{0, 1, 2, 3}

	for _, errCount := range errorCounts {
		fmt.Printf("\nОбработка случая с %d ошибками:\n", errCount)

		allowed := frozenPositions
		errorPositions := make(map[int]struct{})
		for i := 0; i < errCount; i++ {
			choice := allowed[rand.Intn(len(allowed))]
			errorPositions[choice] = struct{}{}
		}
		posList := make([]int, 0, len(errorPositions))
		for pos := range errorPositions {
			posList = append(posList, pos)
		}
		sort.Ints(posList)
		fmt.Printf("Позиции ошибок: %v\n", posList)

		yNoisy, yErrors := addErrors(sig, sigma, errorPositions)
		fmt.Printf("Принятый сигнал (только шум): %v\n", roundSlice(yNoisy, 4))
		fmt.Printf("Принятый сигнал (шум + ошибки): %v\n", roundSlice(yErrors, 4))

		receivedLLR := make([]float64, N)
		for i := 0; i < N; i++ {
			receivedLLR[i] = 2.0 * yErrors[i] / (sigma * sigma)
		}
		fmt.Println("LLR для декодера:")
		fmt.Printf("%v\n", roundSlice(receivedLLR, 4))

		decodedU := sclDecode(receivedLLR, F, 4, n)
		fmt.Printf("Декодированный вектор u: %v\n", decodedU)

		decodedInfo := make([]int, k)
		for j, pos := range A {
			decodedInfo[j] = decodedU[pos]
		}
		fmt.Printf("Декодированные информационные биты: %v\n", decodedInfo)

		bitErrors := 0
		for j := 0; j < k; j++ {
			if infoBits[j] != decodedInfo[j] {
				bitErrors++
			}
		}
		fmt.Printf("Ошибок в информационных битах: %d\n", bitErrors)
	}
}

func genArikanMatrix(n int) [][]int {
	G := [][]int{
		{1, 0},
		{1, 1},
	}
	result := G
	for level := 1; level < n; level++ {
		result = kronecker(result, G)
	}
	return result
}

func kronecker(A, B [][]int) [][]int {
	p := len(A)
	q := len(A[0])
	r := len(B)
	s := len(B[0])
	C := make([][]int, p*r)
	for i := 0; i < p*r; i++ {
		C[i] = make([]int, q*s)
	}
	for i := 0; i < p; i++ {
		for j := 0; j < q; j++ {
			for u := 0; u < r; u++ {
				for v := 0; v < s; v++ {
					C[i*r+u][j*s+v] = A[i][j] * B[u][v]
				}
			}
		}
	}
	return C
}

func polarEnc(u []int, G [][]int) []int {
	N := len(u)
	code := make([]int, N)
	for i := 0; i < N; i++ {
		sum := 0
		for j := 0; j < N; j++ {
			sum += u[j] * G[j][i]
		}
		code[i] = sum % 2
	}
	return code
}

func bpsk(codeword []int) []float64 {
	N := len(codeword)
	s := make([]float64, N)
	for i := 0; i < N; i++ {
		if codeword[i] == 0 {
			s[i] = 1.0
		} else {
			s[i] = -1.0
		}
	}
	return s
}

func addErrors(s []float64, sigma float64, errorPositions map[int]struct{}) ([]float64, []float64) {
	N := len(s)
	yNoisy := make([]float64, N)
	yErrors := make([]float64, N)
	for i := 0; i < N; i++ {
		noise := rand.NormFloat64() * sigma
		yNoisy[i] = s[i] + noise
		yErrors[i] = yNoisy[i]
	}
	for pos := range errorPositions {
		if pos >= 0 && pos < N {
			yErrors[pos] = -yErrors[pos]
		}
	}
	return yNoisy, yErrors
}

func f(a, b float64) float64 {
	signA := 1.0
	if a < 0 {
		signA = -1.0
	}
	signB := 1.0
	if b < 0 {
		signB = -1.0
	}
	absA := math.Abs(a)
	absB := math.Abs(b)
	minAB := absA
	if absB < absA {
		minAB = absB
	}
	return signA * signB * minAB
}

func g(a, b float64, uBit int) float64 {
	factor := 1.0 - 2.0*float64(uBit)
	return b + factor*a
}

func calculateLLR(i int, uPrev []int, channelLLR []float64, depth int) float64 {
	if depth == 0 {
		return channelLLR[0]
	}
	Nhalf := 1 << (depth - 1)
	firstHalf := channelLLR[:Nhalf]
	secondHalf := channelLLR[Nhalf:]
	llrLeft := make([]float64, Nhalf)
	for j := 0; j < Nhalf; j++ {
		llrLeft[j] = f(firstHalf[j], secondHalf[j])
	}
	if i < Nhalf {
		return calculateLLR(i, uPrev[:min(len(uPrev), Nhalf)], llrLeft, depth-1)
	}
	uLeft := uPrev[:min(len(uPrev), Nhalf)]
	llrRight := make([]float64, Nhalf)
	for j := 0; j < Nhalf; j++ {
		uBit := 0
		if j < len(uLeft) {
			uBit = uLeft[j]
		}
		llrRight[j] = g(firstHalf[j], secondHalf[j], uBit)
	}
	return calculateLLR(i-Nhalf, suffix(uPrev, Nhalf), llrRight, depth-1)
}

func sclDecode(channelLLR []float64, F map[int]struct{}, L, depth int) []int {
	paths := []Path{{metric: 0.0, u: []int{}}}

	for i := 0; i < N; i++ {
		newPaths := make([]Path, 0, len(paths)*2)
		for _, path := range paths {
			var llrI float64
			if len(path.u) > 0 {
				llrI = calculateLLR(i, path.u, channelLLR, depth)
			} else {
				llrI = calculateLLR(i, []int{}, channelLLR, depth)
			}

			if _, isFrozen := F[i]; isFrozen {
				uBit := 0
				metricUpdate := metricIncrement(llrI, uBit)
				newMetric := path.metric + metricUpdate
				newU := append(path.u, uBit)
				newPaths = append(newPaths, Path{metric: newMetric, u: newU})
			} else {
				for _, bit := range []int{0, 1} {
					metricUpdate := metricIncrement(llrI, bit)
					newMetric := path.metric + metricUpdate
					newU := append(path.u, bit)
					newPaths = append(newPaths, Path{metric: newMetric, u: newU})
				}
			}
		}
		sort.Slice(newPaths, func(a, b int) bool {
			return newPaths[a].metric < newPaths[b].metric
		})
		if len(newPaths) > L {
			newPaths = newPaths[:L]
		}
		paths = newPaths
	}

	best := paths[0]
	for _, p := range paths[1:] {
		if p.metric < best.metric {
			best = p
		}
	}
	if len(best.u) < N {
		padding := make([]int, N-len(best.u))
		best.u = append(best.u, padding...)
	}
	return best.u
}

func metricIncrement(llr float64, uBit int) float64 {
	factor := 1.0 - 2.0*float64(uBit)
	if factor*llr >= 0 {
		return 0.0
	}
	return math.Abs(llr)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func suffix(uPrev []int, k int) []int {
	if len(uPrev) > k {
		return uPrev[k:]
	}
	return []int{}
}

func roundSlice(input []float64, decimals int) []float64 {
	factor := math.Pow(10, float64(decimals))
	out := make([]float64, len(input))
	for i, val := range input {
		out[i] = math.Round(val*factor) / factor
	}
	return out
}
