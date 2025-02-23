package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"os"
)

// EllipticCurve задаёт кривую вида y^2 = x^3 + a*x + b над полем F_p.
type EllipticCurve struct {
	P *big.Int // модуль (простое число)
	A *big.Int
	B *big.Int
}

// Point задаёт точку на эллиптической кривой.
// Если Infinity == true, то это точка в бесконечности.
type Point struct {
	X        *big.Int
	Y        *big.Int
	Infinity bool
}

// extendedEuclid вычисляет расширенный алгоритм Евклида для a и b.
// Возвращает (g, x, y), где g = gcd(a, b) и x*a + y*b = g.
func extendedEuclid(a, b *big.Int) (g, x, y *big.Int) {
	zero := big.NewInt(0)
	if b.Cmp(zero) == 0 {
		return new(big.Int).Set(a), big.NewInt(1), big.NewInt(0)
	}
	// рекурсивно: g, x1, y1 = extendedEuclid(b, a mod b)
	mod := new(big.Int).Mod(a, b)
	g, x1, y1 := extendedEuclid(b, mod)
	// x = y1, y = x1 - (a/b)*y1
	q := new(big.Int).Div(a, b)
	x = new(big.Int).Set(y1)
	y = new(big.Int).Sub(x1, new(big.Int).Mul(q, y1))
	return g, x, y
}

// modInverse вычисляет обратный элемент к a по модулю mod.
// Если обратного элемента не существует, возвращается nil.
func modInverse(a, mod *big.Int) *big.Int {
	g, x, _ := extendedEuclid(a, mod)
	if g.Cmp(big.NewInt(1)) != 0 {
		return nil // обратного элемента не существует
	}
	return new(big.Int).Mod(x, mod)
}

// modExp – быстрое возведение в степень по модулю.
// Вычисляет (base^exponent) mod mod.
func modExp(base, exponent, mod *big.Int) *big.Int {
	result := big.NewInt(1)
	baseMod := new(big.Int).Mod(base, mod)
	exp := new(big.Int).Set(exponent)
	zero := big.NewInt(0)
	two := big.NewInt(2)
	for exp.Cmp(zero) > 0 {
		// если exp нечётное
		if new(big.Int).And(exp, big.NewInt(1)).Cmp(big.NewInt(1)) == 0 {
			result.Mul(result, baseMod)
			result.Mod(result, mod)
		}
		exp.Div(exp, two)
		baseMod.Mul(baseMod, baseMod)
		baseMod.Mod(baseMod, mod)
	}
	return result
}

// isPrimeFermat проводит тест Ферма для n с заданным числом итераций.
// Обратите внимание: для генерации случайных чисел можно использовать crypto/rand,
// но для простоты здесь выбираются фиксированные основания.
func isPrimeFermat(n *big.Int, iterations int) bool {
	one := big.NewInt(1)
	two := big.NewInt(2)
	if n.Cmp(two) < 0 {
		return false
	}
	if n.Cmp(two) == 0 {
		return true
	}
	// Чётное число (кроме 2) – составное.
	if new(big.Int).Mod(n, two).Cmp(big.NewInt(0)) == 0 {
		return false
	}
	nMinusOne := new(big.Int).Sub(n, one)
	// Для простоты выбираем основания 2, 3, 4, … (не самый надёжный метод для больших чисел)
	for i := 0; i < iterations; i++ {
		// Если n слишком велико, генерация случайного числа усложняется.
		// Здесь просто выбираем a = 2 + i (при условии, что a < n-2).
		a := big.NewInt(int64(2 + i))
		// Если a >= n-2, можно остановиться.
		if a.Cmp(new(big.Int).Sub(n, two)) > 0 {
			a = big.NewInt(2)
		}
		res := modExp(a, nMinusOne, n)
		if res.Cmp(one) != 0 {
			return false
		}
	}
	return true
}

// IsOnCurve проверяет, принадлежит ли точка P эллиптической кривой.
func (curve *EllipticCurve) IsOnCurve(P Point) bool {
	// Точка в бесконечности всегда принадлежит группе.
	if P.Infinity {
		return true
	}
	// Вычисляем y^2 и x^3 + a*x + b по модулю p.
	y2 := new(big.Int).Mul(P.Y, P.Y)
	y2.Mod(y2, curve.P)

	x3 := new(big.Int).Exp(P.X, big.NewInt(3), curve.P)
	ax := new(big.Int).Mul(curve.A, P.X)
	rhs := new(big.Int).Add(x3, ax)
	rhs.Add(rhs, curve.B)
	rhs.Mod(rhs, curve.P)

	return y2.Cmp(rhs) == 0
}

// Add выполняет сложение двух точек P и Q на эллиптической кривой.
func (curve *EllipticCurve) Add(P, Q Point) Point {
	// Если одна из точек – точка в бесконечности, возвращаем другую.
	if P.Infinity {
		return Q
	}
	if Q.Infinity {
		return P
	}
	p := curve.P
	zero := big.NewInt(0)

	// Если x координаты равны:
	if P.X.Cmp(Q.X) == 0 {
		// Если y1 + y2 = 0 mod p, то P + Q = бесконечность.
		sumY := new(big.Int).Add(P.Y, Q.Y)
		sumY.Mod(sumY, p)
		if sumY.Cmp(zero) == 0 {
			return Point{Infinity: true}
		} else {
			// В противном случае – удвоение точки.
			return curve.Double(P)
		}
	}
	// Вычисляем λ = (y2 - y1)/(x2 - x1) mod p.
	numerator := new(big.Int).Sub(Q.Y, P.Y)
	denominator := new(big.Int).Sub(Q.X, P.X)
	denomInv := modInverse(denominator, p)
	if denomInv == nil {
		return Point{Infinity: true}
	}
	lambda := new(big.Int).Mul(numerator, denomInv)
	lambda.Mod(lambda, p)

	// Вычисляем x3 = λ² - x1 - x2 mod p.
	x3 := new(big.Int).Mul(lambda, lambda)
	x3.Sub(x3, P.X)
	x3.Sub(x3, Q.X)
	x3.Mod(x3, p)

	// Вычисляем y3 = λ*(x1 - x3) - y1 mod p.
	y3 := new(big.Int).Sub(P.X, x3)
	y3.Mul(lambda, y3)
	y3.Sub(y3, P.Y)
	y3.Mod(y3, p)

	return Point{
		X:        x3,
		Y:        y3,
		Infinity: false,
	}
}

// Double удваивает точку P на кривой.
func (curve *EllipticCurve) Double(P Point) Point {
	if P.Infinity {
		return P
	}
	p := curve.P
	zero := big.NewInt(0)
	if P.Y.Cmp(zero) == 0 {
		return Point{Infinity: true}
	}
	// Вычисляем λ = (3*x^2 + a)/(2*y) mod p.
	three := big.NewInt(3)
	numerator := new(big.Int).Mul(three, new(big.Int).Mul(P.X, P.X))
	numerator.Add(numerator, curve.A)
	denom := new(big.Int).Mul(big.NewInt(2), P.Y)
	denomInv := modInverse(denom, p)
	if denomInv == nil {
		return Point{Infinity: true}
	}
	lambda := new(big.Int).Mul(numerator, denomInv)
	lambda.Mod(lambda, p)

	// x3 = λ² - 2*x
	x3 := new(big.Int).Mul(lambda, lambda)
	twoX := new(big.Int).Mul(big.NewInt(2), P.X)
	x3.Sub(x3, twoX)
	x3.Mod(x3, p)

	// y3 = λ*(x - x3) - y
	y3 := new(big.Int).Sub(P.X, x3)
	y3.Mul(lambda, y3)
	y3.Sub(y3, P.Y)
	y3.Mod(y3, p)

	return Point{
		X:        x3,
		Y:        y3,
		Infinity: false,
	}
}

// ScalarMult вычисляет k * P методом «двоичного умножения».
func (curve *EllipticCurve) ScalarMult(P Point, k *big.Int) Point {
	result := Point{Infinity: true} // нейтральный элемент
	addend := P

	kCopy := new(big.Int).Set(k)
	zero := big.NewInt(0)
	two := big.NewInt(2)
	for kCopy.Cmp(zero) > 0 {
		if new(big.Int).And(kCopy, big.NewInt(1)).Cmp(big.NewInt(1)) == 0 {
			result = curve.Add(result, addend)
		}
		addend = curve.Double(addend)
		kCopy.Div(kCopy, two)
	}
	return result
}

// Neg возвращает обратную точку для P: -P = (x, -y mod p).
func (curve *EllipticCurve) Neg(P Point) Point {
	if P.Infinity {
		return P
	}
	negY := new(big.Int).Neg(P.Y)
	negY.Mod(negY, curve.P)
	return Point{
		X:        new(big.Int).Set(P.X),
		Y:        negY,
		Infinity: false,
	}
}

// Sub вычисляет разность точек: P - Q = P + (-Q).
func (curve *EllipticCurve) Sub(P, Q Point) Point {
	return curve.Add(P, curve.Neg(Q))
}

// pointToString возвращает строковое представление точки (для использования в качестве ключа).
func pointToString(P Point) string {
	if P.Infinity {
		return "inf"
	}
	return P.X.String() + "," + P.Y.String()
}

// Points возвращает все точки эллиптической кривой над полем F_p (наивный перебор).
func (curve *EllipticCurve) Points() []Point {
	points := []Point{}
	zero := big.NewInt(0)
	one := big.NewInt(1)

	// Предполагаем, что p достаточно мало, чтобы перебрать все x.
	pInt64 := curve.P.Int64()
	for i := int64(0); i < pInt64; i++ {
		x := big.NewInt(i)
		// Вычисляем f(x) = x^3 + a*x + b mod p.
		x3 := new(big.Int).Exp(x, big.NewInt(3), curve.P)
		ax := new(big.Int).Mul(curve.A, x)
		fx := new(big.Int).Add(x3, ax)
		fx.Add(fx, curve.B)
		fx.Mod(fx, curve.P)
		// Определяем число решений уравнения y^2 = f(x).
		// Если f(x) == 0, то y = 0 – одно решение.
		// Иначе, по критерию Эйлера: если f(x)^((p-1)/2) mod p == 1, то решений два.
		if fx.Cmp(zero) == 0 {
			points = append(points, Point{X: new(big.Int).Set(x), Y: big.NewInt(0), Infinity: false})
		} else {
			exp := new(big.Int).Sub(curve.P, one)
			exp.Div(exp, big.NewInt(2))
			legendre := modExp(fx, exp, curve.P)
			if legendre.Cmp(one) == 0 {
				// Для поиска корней наивно перебираем возможные y.
				for j := int64(0); j < pInt64; j++ {
					y := big.NewInt(j)
					y2 := new(big.Int).Mul(y, y)
					y2.Mod(y2, curve.P)
					if y2.Cmp(fx) == 0 {
						points = append(points, Point{X: new(big.Int).Set(x), Y: new(big.Int).Set(y), Infinity: false})
					}
				}
			}
		}
	}
	// Добавляем точку в бесконечности.
	points = append(points, Point{Infinity: true})
	return points
}

// DiscreteLog реализует метод больших и малых шагов для поиска k,
// такого что Q = k * P, при условии, что порядок группы известен.
func (curve *EllipticCurve) DiscreteLog(P, Q Point, groupOrder *big.Int) (*big.Int, error) {
	// m = ceil(sqrt(n))
	m := new(big.Int).Sqrt(groupOrder)
	m.Add(m, big.NewInt(1))

	// Построение таблицы baby steps: jP для j = 0 ... m-1.
	babySteps := make(map[string]*big.Int)
	current := Point{Infinity: true} // 0 * P
	j := big.NewInt(0)
	one := big.NewInt(1)
	for j.Cmp(m) < 0 {
		babySteps[pointToString(current)] = new(big.Int).Set(j)
		current = curve.Add(current, P)
		j.Add(j, one)
	}
	// Вычисляем mP.
	mP := curve.ScalarMult(P, m)
	// Ищем i такое, что Q - i*(mP) совпадает с каким-либо jP.
	current = Q
	i := big.NewInt(0)
	for i.Cmp(m) < 0 {
		if jVal, ok := babySteps[pointToString(current)]; ok {
			// k = i*m + jVal.
			k := new(big.Int).Mul(i, m)
			k.Add(k, jVal)
			return k, nil
		}
		current = curve.Sub(current, mP)
		i.Add(i, one)
	}
	return nil, errors.New("discrete log not found")
}

// primeFactors возвращает срез простых делителей числа n (для int64).
func primeFactors(n int64) []int64 {
	factors := []int64{}
	for i := int64(2); i*i <= n; i++ {
		for n%i == 0 {
			factors = append(factors, i)
			n /= i
		}
	}
	if n > 1 {
		factors = append(factors, n)
	}
	return factors
}

// PointOrder вычисляет порядок точки P, используя делимость порядка группы.
// Предполагается, что groupOrder – малое число (int64), чтобы можно было выполнить факторизацию.
func (curve *EllipticCurve) PointOrder(P Point, groupOrder *big.Int) *big.Int {
	order := new(big.Int).Set(groupOrder)
	// Преобразуем groupOrder в int64 для факторизации (подходит для малых чисел).
	groupOrderInt64 := order.Int64()
	factors := primeFactors(groupOrderInt64)
	// Подсчёт кратностей
	factorCounts := make(map[int64]int)
	for _, factor := range factors {
		factorCounts[factor]++
	}
	for q, count := range factorCounts {
		qBig := big.NewInt(q)
		for i := 0; i < count; i++ {
			temp := new(big.Int).Div(order, qBig)
			// Если (order/q)*P = бесконечность, то порядок делится на q.
			if curve.ScalarMult(P, temp).Infinity {
				order = temp
			} else {
				break
			}
		}
	}
	return order
}

// FindPointsOfPrimeOrder ищет точки, порядок которых равен заданному простому числу.
func (curve *EllipticCurve) FindPointsOfPrimeOrder(prime int64, groupOrder *big.Int) []Point {
	candidates := []Point{}
	pts := curve.Points()
	target := big.NewInt(prime)
	for _, P := range pts {
		if P.Infinity {
			continue
		}
		ord := curve.PointOrder(P, groupOrder)
		if ord != nil && ord.Cmp(target) == 0 {
			candidates = append(candidates, P)
		}
	}
	return candidates
}

func main() {
	// Пример: выбираем малое простое p и коэффициенты a, b.
	// Для исследования кривых с |E(F_p)| от 2^10 до 2^512 следует подбирать соответствующие параметры.
	fmt.Println("введите через пробел числа p, a, b")
	in := bufio.NewReader(os.Stdin)
    var pR, aR, bR int64
    fmt.Fscan(in, &pR, &aR, &bR)
    p := big.NewInt(pR)
	a := big.NewInt(aR)
	b := big.NewInt(bR)
	curve := EllipticCurve{
		P: p,
		A: a,
		B: b,
	}
	fmt.Printf("Эллиптическая кривая: y^2 = x^3 + %s*x + %s над F_%s\n", a.String(), b.String(), p.String())

	// Вычисляем все точки кривой (наивный перебор).
	pts := curve.Points()
	fmt.Printf("Найдено точек: %d\n", len(pts))
	for _, pt := range pts {
		if pt.Infinity {
			fmt.Println("Infinity")
		} else {
			fmt.Printf("(%s, %s)\n", pt.X.String(), pt.Y.String())
		}
	}

	// Группа точек имеет порядок равный числу найденных точек.
	groupOrder := big.NewInt(int64(len(pts)))
	fmt.Printf("Порядок группы: %s\n", groupOrder.String())

	// Пример: сложение двух точек.
	if len(pts) >= 2 {
		P := pts[2]
		Q := pts[1]
		R := curve.Add(P, Q)
		fmt.Printf("P = %s\nQ = %s\nP+Q = %s\n", pointToString(P), pointToString(Q), pointToString(R))
	}

	// Пример: умножение точки на число (например, 3*P).
	if len(pts) > 0 {
		P := pts[0]
		k := big.NewInt(3)
		R := curve.ScalarMult(P, k)
		fmt.Printf("3 * P = %s\n", pointToString(R))
	}

	// Тест Ферма для проверки простоты p.
	if isPrimeFermat(p, 5) {
		fmt.Printf("%s является простым числом (по тесту Ферма).\n", p.String())
	} else {
		fmt.Printf("%s не является простым числом.\n", p.String())
	}

	if len(pts) > 2 {
		P := pts[0]
		kExpected := big.NewInt(3)
		Q := curve.ScalarMult(P, kExpected)
		kFound, err := curve.DiscreteLog(P, Q, groupOrder)
		if err != nil {
			fmt.Println("Не удалось найти дискретный логарифм:", err)
		} else {
			fmt.Printf("Найден дискретный логарифм: k = %s, P={%v}, Q={%v}\n", kFound.String(), P, Q)
		}
	}

	// Поиск подгрупп простого порядка.
	// Факторизуем порядок группы (подразумевается, что он малый).
	groupOrderInt64 := groupOrder.Int64()
	pFactors := primeFactors(groupOrderInt64)
	uniquePrimes := make(map[int64]bool)
	for _, q := range pFactors {
		uniquePrimes[q] = true
	}
	for q := range uniquePrimes {
		ptsPrime := curve.FindPointsOfPrimeOrder(q, groupOrder)
		fmt.Printf("Точки простого порядка %d:\n", q)
		for _, pt := range ptsPrime {
			fmt.Println(pointToString(pt))
		}
	}
}
