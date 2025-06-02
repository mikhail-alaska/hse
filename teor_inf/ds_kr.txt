package main

import (
	"fmt"
	"math"
	"math/rand"
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
	// Frozen bit indices
	frozenPositions = []int{0, 1, 2, 3, 4, 5, 8, 9}
)

// Path represents a candidate path in SCL decoding, with its accumulated metric and decided bits.
type Path struct {
	metric float64
	u      []int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Build frozen set F and information set A
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
	sort.Ints(A) // Ensure A is sorted

	// 1) Generate random information bits of length k
	infoBits := make([]int, k)
	for i := range infoBits {
		infoBits[i] = rand.Intn(2)
	}
	fmt.Printf("Информационные биты: %v\n", infoBits)
	fmt.Printf("Информационные позиции (A): %v\n", A)

	// 2) Build u_full of length N, with frozen bits = 0 and information bits placed at indices in A
	uFull := make([]int, N)
	for j, pos := range A {
		uFull[pos] = infoBits[j]
	}
	fmt.Println("Вектор u (с замороженными битами):")
	fmt.Printf("%v\n", uFull)

	// 3) Generate Arıkan matrix of size N×N
	G := genArikanMatrix(n)

	// 4) Polar encoding: codeword = uFull × G mod 2
	codeword := polarEnc(uFull, G)
	fmt.Printf("Закодированное слово: %v\n", codeword)

	// 5) BPSK modulation: s[i] = 1 - 2*codeword[i]
	sig := bpsk(codeword)
	fmt.Printf("Модулированный сигнал: %v\n", roundSlice(sig, 4))

	// 6) Simulate noise and bit-flips; then decode via SCL
	sigma := 0.1
	errorCounts := []int{0, 1, 2, 3}

	for _, errCount := range errorCounts {
		fmt.Printf("\nОбработка случая с %d ошибками:\n", errCount)

		// a) Randomly choose error positions (from the frozen set indices) errCount times
		allowed := frozenPositions
		errorPositions := make(map[int]struct{})
		for i := 0; i < errCount; i++ {
			choice := allowed[rand.Intn(len(allowed))]
			errorPositions[choice] = struct{}{}
		}
		// Collect and sort positions for printing
		posList := make([]int, 0, len(errorPositions))
		for pos := range errorPositions {
			posList = append(posList, pos)
		}
		sort.Ints(posList)
		fmt.Printf("Позиции ошибок: %v\n", posList)

		// b) Add AWGN noise and flip specified bits
		yNoisy, yErrors := addErrors(sig, sigma, errorPositions)
		fmt.Printf("Принятый сигнал (только шум): %v\n", roundSlice(yNoisy, 4))
		fmt.Printf("Принятый сигнал (шум + ошибки): %v\n", roundSlice(yErrors, 4))

		// c) Compute channel LLRs: LLR[i] = 2*yErrors[i]/(sigma^2)
		receivedLLR := make([]float64, N)
		for i := 0; i < N; i++ {
			receivedLLR[i] = 2.0 * yErrors[i] / (sigma * sigma)
		}
		fmt.Println("LLR для декодера:")
		fmt.Printf("%v\n", roundSlice(receivedLLR, 4))

		// d) SCL decode with list size L=4
		decodedU := sclDecode(receivedLLR, F, 4, n)
		fmt.Printf("Декодированный вектор u: %v\n", decodedU)

		// e) Extract decoded information bits at indices in A
		decodedInfo := make([]int, k)
		for j, pos := range A {
			decodedInfo[j] = decodedU[pos]
		}
		fmt.Printf("Декодированные информационные биты: %v\n", decodedInfo)

		// f) Count bit errors between original infoBits and decodedInfo
		bitErrors := 0
		for j := 0; j < k; j++ {
			if infoBits[j] != decodedInfo[j] {
				bitErrors++
			}
		}
		fmt.Printf("Ошибок в информационных битах: %d\n", bitErrors)
	}
}

// genArikanMatrix builds the Arıkan polar transform matrix of size 2^n × 2^n.
// It starts from G = [[1,0],[1,1]] and takes n-fold Kronecker power.
func genArikanMatrix(n int) [][]int {
	// Base 2×2 kernel
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

// kronecker computes the Kronecker product of matrices A (p×q) and B (r×s),
// returning a (pr)×(qs) matrix.
func kronecker(A, B [][]int) [][]int {
	p := len(A)
	q := len(A[0])
	r := len(B)
	s := len(B[0])
	// Resulting size: (p*r) × (q*s)
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

// polarEnc multiplies the binary vector u (length N) by G (N×N) modulo 2.
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

// bpsk maps bits {0,1} → symbols {+1, -1}, returning a slice of float64.
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

// addErrors adds AWGN noise (σ) to signal s, and then flips the sign at errorPositions.
// Returns (yNoisy, yWithErrors).
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

// calculateLLR recursively computes the LLR for bit index i given previously decided bits uPrev,
// the channel LLR vector channelLLR (length 2^n), and the current recursion depth n.
func calculateLLR(i int, uPrev []int, channelLLR []float64, depth int) float64 {
	if depth == 0 {
		return channelLLR[0]
	}
	Nhalf := 1 << (depth - 1)
	firstHalf := channelLLR[:Nhalf]
	secondHalf := channelLLR[Nhalf:]
	// Compute left-branch LLRs
	llrLeft := make([]float64, Nhalf)
	for j := 0; j < Nhalf; j++ {
		llrLeft[j] = f(firstHalf[j], secondHalf[j])
	}
	if i < Nhalf {
		// Recurse on left half
		return calculateLLR(i, uPrev[:min(len(uPrev), Nhalf)], llrLeft, depth-1)
	}
	// Prepare right-branch LLRs
	uLeft := uPrev[:min(len(uPrev), Nhalf)]
	llrRight := make([]float64, Nhalf)
	for j := 0; j < Nhalf; j++ {
		// If uLeft is shorter than Nhalf, assume missing bits = 0
		uBit := 0
		if j < len(uLeft) {
			uBit = uLeft[j]
		}
		llrRight[j] = g(firstHalf[j], secondHalf[j], uBit)
	}
	return calculateLLR(i-Nhalf, suffix(uPrev, Nhalf), llrRight, depth-1)
}

// sclDecode performs Successive Cancellation List decoding on the channel LLR vector,
// with frozen set F, list size L, and depth n. Returns the best-decoded u vector (length N).
func sclDecode(channelLLR []float64, F map[int]struct{}, L, depth int) []int {
	// Initialize with a single empty path
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

			// Check if i is a frozen index
			if _, isFrozen := F[i]; isFrozen {
				// Force u_i = 0
				uBit := 0
				metricUpdate := metricIncrement(llrI, uBit)
				newMetric := path.metric + metricUpdate
				newU := append(path.u, uBit)
				newPaths = append(newPaths, Path{metric: newMetric, u: newU})
			} else {
				// Branch for bit = 0 and bit = 1
				for _, bit := range []int{0, 1} {
					metricUpdate := metricIncrement(llrI, bit)
					newMetric := path.metric + metricUpdate
					newU := append(path.u, bit)
					newPaths = append(newPaths, Path{metric: newMetric, u: newU})
				}
			}
		}
		// Keep only the best L paths (smallest metric)
		sort.Slice(newPaths, func(a, b int) bool {
			return newPaths[a].metric < newPaths[b].metric
		})
		if len(newPaths) > L {
			newPaths = newPaths[:L]
		}
		paths = newPaths
	}

	// Return the path with the smallest metric
	best := paths[0]
	for _, p := range paths[1:] {
		if p.metric < best.metric {
			best = p
		}
	}
	// Ensure the returned u has length N
	if len(best.u) < N {
		// Pad with zeros if for some reason it's shorter (shouldn't happen)
		padding := make([]int, N-len(best.u))
		best.u = append(best.u, padding...)
	}
	return best.u
}

// metricIncrement computes the branch metric increment for deciding bit uBit
// given LLR value llr. If (1 - 2*uBit)*llr >= 0, increment is 0; else it's |llr|.
func metricIncrement(llr float64, uBit int) float64 {
	factor := 1.0 - 2.0*float64(uBit)
	if factor*llr >= 0 {
		return 0.0
	}
	return math.Abs(llr)
}

// Helper: return the minimum of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// suffix returns the slice uPrev[k:] if len(uPrev) > k, else returns an empty slice.
func suffix(uPrev []int, k int) []int {
	if len(uPrev) > k {
		return uPrev[k:]
	}
	return []int{}
}

// roundSlice returns a new []float64 where each element of input
// is rounded to 'decimals' decimal places.
func roundSlice(input []float64, decimals int) []float64 {
	factor := math.Pow(10, float64(decimals))
	out := make([]float64, len(input))
	for i, val := range input {
		out[i] = math.Round(val*factor) / factor
	}
	return out
}
