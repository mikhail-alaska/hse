package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

const (
	N = 16
	R = 1.0 / 2.0
)

var (
	k = int(float64(N) * R)
	n = int(math.Log2(N))
	// Замороженные позиции битов
	frozenPositions = []int{0, 1, 2, 3, 4, 5, 8, 9}
)

// Path представляет кандидатную траекторию в SCL-декодировании:
// накопленная метрика и принятые биты u.
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
	sort.Ints(A) // Убедимся, что A отсортирован

	// 1) Генерируем случайные информационные биты длины k
	infoBits := make([]int, k)
	for i := range infoBits {
		infoBits[i] = rand.Intn(2)
	}
	fmt.Printf("Информационные биты: %v\n", infoBits)
	fmt.Printf("Информационные позиции (A): %v\n", A)

	// 2) Построим u_full длины N: замороженные биты = 0, информационные вставляем по A
	uFull := make([]int, N)
	for j, pos := range A {
		uFull[pos] = infoBits[j]
	}
	fmt.Println("Вектор u (с замороженными битами):")
	fmt.Printf("%v\n", uFull)

	// 3) Построим матрицу Арикана размера N×N
	G := genArikanMatrix(n)

	// 4) Полярное кодирование: codeword = uFull × G mod 2
	codeword := polarEnc(uFull, G)
	fmt.Printf("Закодированное слово: %v\n", codeword)

	// 5) BPSK-модуляция: s[i] = 1 - 2*codeword[i]
	sig := bpsk(codeword)
	fmt.Printf("Модулированный сигнал: %v\n", roundSlice(sig, 4))

	// 6) Смоделировать шум и инвертировать биты; затем выполнить SCL-декодирование
	sigma := 0.1
	errorCounts := []int{0, 1, 2, 3}

	for _, errCount := range errorCounts {
		fmt.Printf("\nОбработка случая с %d ошибками:\n", errCount)

		// a) Случайным образом выбрать errCount позиций ошибок из frozenPositions
		allowed := frozenPositions
		errorPositions := make(map[int]struct{})
		for i := 0; i < errCount; i++ {
			choice := allowed[rand.Intn(len(allowed))]
			errorPositions[choice] = struct{}{}
		}
		// Соберём и отсортируем позиции для вывода
		posList := make([]int, 0, len(errorPositions))
		for pos := range errorPositions {
			posList = append(posList, pos)
		}
		sort.Ints(posList)
		fmt.Printf("Позиции ошибок: %v\n", posList)

		// b) Добавить шум и инвертировать указанные биты
		yNoisy, yErrors := addErrors(sig, sigma, errorPositions)
		fmt.Printf("Принятый сигнал (только шум): %v\n", roundSlice(yNoisy, 4))
		fmt.Printf("Принятый сигнал (шум + ошибки): %v\n", roundSlice(yErrors, 4))

		// c) Рассчитать LLR: LLR[i] = 2*yErrors[i]/(sigma^2)
		receivedLLR := make([]float64, N)
		for i := 0; i < N; i++ {
			receivedLLR[i] = 2.0 * yErrors[i] / (sigma * sigma)
		}
		fmt.Println("LLR для декодера:")
		fmt.Printf("%v\n", roundSlice(receivedLLR, 4))

		// d) Выполнить SCL-декодирование с размером списка L=4
		decodedU := sclDecode(receivedLLR, F, 4, n)
		fmt.Printf("Декодированный вектор u: %v\n", decodedU)

		// e) Извлечь декодированные информационные биты по индексам A
		decodedInfo := make([]int, k)
		for j, pos := range A {
			decodedInfo[j] = decodedU[pos]
		}
		fmt.Printf("Декодированные информационные биты: %v\n", decodedInfo)

		// f) Посчитать количество ошибок в информационных битах
		bitErrors := 0
		for j := 0; j < k; j++ {
			if infoBits[j] != decodedInfo[j] {
				bitErrors++
			}
		}
		fmt.Printf("Ошибок в информационных битах: %d\n", bitErrors)
	}
}

// genArikanMatrix строит матрицу Арикана размером 2^n × 2^n.
// Базовый ядро G = [[1,0],[1,1]], затем n-кратное тензорное произведение.
func genArikanMatrix(n int) [][]int {
	// Базовый 2×2-ядро
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

// kronecker вычисляет тензорное произведение матриц A (p×q) и B (r×s),
// возвращая матрицу размером (p*r)×(q*s).
func kronecker(A, B [][]int) [][]int {
	p := len(A)
	q := len(A[0])
	r := len(B)
	s := len(B[0])
	// Результирующий размер: (p*r)×(q*s)
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

// polarEnc умножает бинарный вектор u (длина N) на G (N×N) по модулю 2.
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

// bpsk отображает биты {0,1} → символы {+1, -1}.
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

// addErrors добавляет AWGN-шум (σ) к сигналу s и затем инвертирует знак
// в битах, перечисленных в errorPositions. Возвращает (yNoisy, yWithErrors).
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

// f(a,b) = sign(a)*sign(b)*min(|a|,|b|)
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

// g(a,b,uBit) = b + (1 - 2*uBit)*a
func g(a, b float64, uBit int) float64 {
	factor := 1.0 - 2.0*float64(uBit)
	return b + factor*a
}

// calculateLLR рекурсивно вычисляет LLR для бита с индексом i,
// учитывая ранее принятые биты uPrev, вектор LLR от канала channelLLR
// длины 2^depth и текущую глубину depth.
func calculateLLR(i int, uPrev []int, channelLLR []float64, depth int) float64 {
	if depth == 0 {
		return channelLLR[0]
	}
	Nhalf := 1 << (depth - 1)
	firstHalf := channelLLR[:Nhalf]
	secondHalf := channelLLR[Nhalf:]
	// Вычисляем LLR для левой ветви
	llrLeft := make([]float64, Nhalf)
	for j := 0; j < Nhalf; j++ {
		llrLeft[j] = f(firstHalf[j], secondHalf[j])
	}
	if i < Nhalf {
		// Рекурсия по левой половине
		return calculateLLR(i, uPrev[:min(len(uPrev), Nhalf)], llrLeft, depth-1)
	}
	// Готовим LLR для правой ветви
	uLeft := uPrev[:min(len(uPrev), Nhalf)]
	llrRight := make([]float64, Nhalf)
	for j := 0; j < Nhalf; j++ {
		// Если uLeft короче Nhalf, отсутствующие биты = 0
		uBit := 0
		if j < len(uLeft) {
			uBit = uLeft[j]
		}
		llrRight[j] = g(firstHalf[j], secondHalf[j], uBit)
	}
	return calculateLLR(i-Nhalf, suffix(uPrev, Nhalf), llrRight, depth-1)
}

// sclDecode выполняет SCL-декодирование L-списком для входного вектора LLR канала,
// заданного множества замороженных индексов F, размера списка L и глубины depth.
// Возвращает лучший декодированный вектор u длины N.
func sclDecode(channelLLR []float64, F map[int]struct{}, L, depth int) []int {
	// Инициализируем единственный пустой путь
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

			// Если i в F — принудительно u_i = 0
			if _, isFrozen := F[i]; isFrozen {
				uBit := 0
				metricUpdate := metricIncrement(llrI, uBit)
				newMetric := path.metric + metricUpdate
				newU := append(path.u, uBit)
				newPaths = append(newPaths, Path{metric: newMetric, u: newU})
			} else {
				// Ветвление по u_i = 0 и u_i = 1
				for _, bit := range []int{0, 1} {
					metricUpdate := metricIncrement(llrI, bit)
					newMetric := path.metric + metricUpdate
					newU := append(path.u, bit)
					newPaths = append(newPaths, Path{metric: newMetric, u: newU})
				}
			}
		}
		// Оставляем только L путей с наименьшими метриками
		sort.Slice(newPaths, func(a, b int) bool {
			return newPaths[a].metric < newPaths[b].metric
		})
		if len(newPaths) > L {
			newPaths = newPaths[:L]
		}
		paths = newPaths
	}

	// Выбираем путь с наименьшей метрикой
	best := paths[0]
	for _, p := range paths[1:] {
		if p.metric < best.metric {
			best = p
		}
	}
	// Убеждаемся, что длина best.u = N (дополняем нулями, если требуется)
	if len(best.u) < N {
		padding := make([]int, N-len(best.u))
		best.u = append(best.u, padding...)
	}
	return best.u
}

// metricIncrement вычисляет приращение метрики для решения uBit при LLR = llr.
// Если (1 - 2*uBit)*llr >= 0, приращение = 0, иначе = |llr|.
func metricIncrement(llr float64, uBit int) float64 {
	factor := 1.0 - 2.0*float64(uBit)
	if factor*llr >= 0 {
		return 0.0
	}
	return math.Abs(llr)
}

// min возвращает минимум из a и b.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// suffix возвращает подмассив uPrev[k:], если uPrev длиннее k, иначе пустой срез.
func suffix(uPrev []int, k int) []int {
	if len(uPrev) > k {
		return uPrev[k:]
	}
	return []int{}
}

// roundSlice округляет каждый элемент входного []float64 до decimals знаков после запятой.
func roundSlice(input []float64, decimals int) []float64 {
	factor := math.Pow(10, float64(decimals))
	out := make([]float64, len(input))
	for i, val := range input {
		out[i] = math.Round(val*factor) / factor
	}
	return out
}
