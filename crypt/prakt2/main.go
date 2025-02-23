package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"os"
)

type EllipticCurve struct {
	P *big.Int 
	A *big.Int
	B *big.Int
}



type Point struct {
	X        *big.Int
	Y        *big.Int
	Infinity bool
}



func extendedEuclid(a, b *big.Int) (g, x, y *big.Int) {
	zero := big.NewInt(0)
	if b.Cmp(zero) == 0 {
		return new(big.Int).Set(a), big.NewInt(1), big.NewInt(0)
	}
	
	mod := new(big.Int).Mod(a, b)
	g, x1, y1 := extendedEuclid(b, mod)
	
	q := new(big.Int).Div(a, b)
	x = new(big.Int).Set(y1)
	y = new(big.Int).Sub(x1, new(big.Int).Mul(q, y1))
	return g, x, y
}



func modInverse(a, mod *big.Int) *big.Int {
	g, x, _ := extendedEuclid(a, mod)
	if g.Cmp(big.NewInt(1)) != 0 {
		return nil 
	}
	return new(big.Int).Mod(x, mod)
}



func modExp(base, exponent, mod *big.Int) *big.Int {
	result := big.NewInt(1)
	baseMod := new(big.Int).Mod(base, mod)
	exp := new(big.Int).Set(exponent)
	zero := big.NewInt(0)
	two := big.NewInt(2)
	for exp.Cmp(zero) > 0 {
		
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




func isPrimeFermat(n *big.Int, iterations int) bool {
	one := big.NewInt(1)
	two := big.NewInt(2)
	if n.Cmp(two) < 0 {
		return false
	}
	if n.Cmp(two) == 0 {
		return true
	}
	
	if new(big.Int).Mod(n, two).Cmp(big.NewInt(0)) == 0 {
		return false
	}
	nMinusOne := new(big.Int).Sub(n, one)
	
	for i := 0; i < iterations; i++ {
		
		
		a := big.NewInt(int64(2 + i))
		
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


func (curve *EllipticCurve) IsOnCurve(P Point) bool {
	
	if P.Infinity {
		return true
	}
	
	y2 := new(big.Int).Mul(P.Y, P.Y)
	y2.Mod(y2, curve.P)

	x3 := new(big.Int).Exp(P.X, big.NewInt(3), curve.P)
	ax := new(big.Int).Mul(curve.A, P.X)
	rhs := new(big.Int).Add(x3, ax)
	rhs.Add(rhs, curve.B)
	rhs.Mod(rhs, curve.P)

	return y2.Cmp(rhs) == 0
}


func (curve *EllipticCurve) Add(P, Q Point) Point {
	
	if P.Infinity {
		return Q
	}
	if Q.Infinity {
		return P
	}
	p := curve.P
	zero := big.NewInt(0)

	
	if P.X.Cmp(Q.X) == 0 {
		
		sumY := new(big.Int).Add(P.Y, Q.Y)
		sumY.Mod(sumY, p)
		if sumY.Cmp(zero) == 0 {
			return Point{Infinity: true}
		} else {
			
			return curve.Double(P)
		}
	}
	
	numerator := new(big.Int).Sub(Q.Y, P.Y)
	denominator := new(big.Int).Sub(Q.X, P.X)
	denomInv := modInverse(denominator, p)
	if denomInv == nil {
		return Point{Infinity: true}
	}
	lambda := new(big.Int).Mul(numerator, denomInv)
	lambda.Mod(lambda, p)

	
	x3 := new(big.Int).Mul(lambda, lambda)
	x3.Sub(x3, P.X)
	x3.Sub(x3, Q.X)
	x3.Mod(x3, p)

	
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


func (curve *EllipticCurve) Double(P Point) Point {
	if P.Infinity {
		return P
	}
	p := curve.P
	zero := big.NewInt(0)
	if P.Y.Cmp(zero) == 0 {
		return Point{Infinity: true}
	}
	
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

	
	x3 := new(big.Int).Mul(lambda, lambda)
	twoX := new(big.Int).Mul(big.NewInt(2), P.X)
	x3.Sub(x3, twoX)
	x3.Mod(x3, p)

	
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


func (curve *EllipticCurve) ScalarMult(P Point, k *big.Int) Point {
	result := Point{Infinity: true} 
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


func (curve *EllipticCurve) Sub(P, Q Point) Point {
	return curve.Add(P, curve.Neg(Q))
}


func pointToString(P Point) string {
	if P.Infinity {
		return "inf"
	}
	return P.X.String() + "," + P.Y.String()
}


func (curve *EllipticCurve) Points() []Point {
	points := []Point{}
	zero := big.NewInt(0)
	one := big.NewInt(1)

	
	pInt64 := curve.P.Int64()
	for i := int64(0); i < pInt64; i++ {
		x := big.NewInt(i)
		
		x3 := new(big.Int).Exp(x, big.NewInt(3), curve.P)
		ax := new(big.Int).Mul(curve.A, x)
		fx := new(big.Int).Add(x3, ax)
		fx.Add(fx, curve.B)
		fx.Mod(fx, curve.P)
		
		
		
		if fx.Cmp(zero) == 0 {
			points = append(points, Point{X: new(big.Int).Set(x), Y: big.NewInt(0), Infinity: false})
		} else {
			exp := new(big.Int).Sub(curve.P, one)
			exp.Div(exp, big.NewInt(2))
			legendre := modExp(fx, exp, curve.P)
			if legendre.Cmp(one) == 0 {
				
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
	
	points = append(points, Point{Infinity: true})
	return points
}



func (curve *EllipticCurve) DiscreteLog(P, Q Point, groupOrder *big.Int) (*big.Int, error) {
	
	m := new(big.Int).Sqrt(groupOrder)
	m.Add(m, big.NewInt(1))

	
	babySteps := make(map[string]*big.Int)
	current := Point{Infinity: true} 
	j := big.NewInt(0)
	one := big.NewInt(1)
	for j.Cmp(m) < 0 {
		babySteps[pointToString(current)] = new(big.Int).Set(j)
		current = curve.Add(current, P)
		j.Add(j, one)
	}
	
	mP := curve.ScalarMult(P, m)
	
	current = Q
	i := big.NewInt(0)
	for i.Cmp(m) < 0 {
		if jVal, ok := babySteps[pointToString(current)]; ok {
			
			k := new(big.Int).Mul(i, m)
			k.Add(k, jVal)
			return k, nil
		}
		current = curve.Sub(current, mP)
		i.Add(i, one)
	}
	return nil, errors.New("discrete log not found")
}


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



func (curve *EllipticCurve) PointOrder(P Point, groupOrder *big.Int) *big.Int {
	order := new(big.Int).Set(groupOrder)
	
	groupOrderInt64 := order.Int64()
	factors := primeFactors(groupOrderInt64)
	
	factorCounts := make(map[int64]int)
	for _, factor := range factors {
		factorCounts[factor]++
	}
	for q, count := range factorCounts {
		qBig := big.NewInt(q)
		for i := 0; i < count; i++ {
			temp := new(big.Int).Div(order, qBig)
			
			if curve.ScalarMult(P, temp).Infinity {
				order = temp
			} else {
				break
			}
		}
	}
	return order
}


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

	
	pts := curve.Points()
	fmt.Printf("Найдено точек: %d\n", len(pts))
	for _, pt := range pts {
		if pt.Infinity {
			fmt.Println("Infinity")
		} else {
			fmt.Printf("(%s, %s)\n", pt.X.String(), pt.Y.String())
		}
	}

	
	groupOrder := big.NewInt(int64(len(pts)))
	fmt.Printf("Порядок группы: %s\n", groupOrder.String())

	
	if len(pts) >= 2 {
		P := pts[2]
		Q := pts[1]
		R := curve.Add(P, Q)
		fmt.Printf("P = %s\nQ = %s\nP+Q = %s\n", pointToString(P), pointToString(Q), pointToString(R))
	}

	
	if len(pts) > 0 {
		P := pts[0]
		k := big.NewInt(3)
		R := curve.ScalarMult(P, k)
		fmt.Printf("3 * P = %s\n", pointToString(R))
	}

	
	if isPrimeFermat(p, 5) {
		fmt.Printf("%s является простым числом (по тесту Ферма).\n", p.String())
	} else {
		fmt.Printf("%s не является простым числом.\n", p.String())
	}

	if len(pts) > 2 {
		P := pts[1]
		kExpected := big.NewInt(3)
		Q := curve.ScalarMult(P, kExpected)
		kFound, err := curve.DiscreteLog(P, Q, groupOrder)
		if err != nil {
			fmt.Println("Не удалось провести скалярное умножение:", err)
		} else {
			fmt.Printf("удалось провесит скалярное умножение: k = %s, P={%v %v}, Q={%v %v}\n", kFound.String(), P.X, P.Y, Q.X, Q.Y)
		}
	}

	
	
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
