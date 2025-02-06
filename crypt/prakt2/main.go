package main

import (
	"fmt"
)

// EllipticCurve задаёт кривую y^2 = x^3 + a*x + b над полем F_p.
type EllipticCurve struct {
	P int // модуль (простое число)
	A int
	B int
}

// Point задаёт точку на эллиптической кривой.
type Point struct {
	X        int
	Y        int
	Infinity bool
}

// modInverse вычисляет обратный элемент a по модулю mod.
func modInverse(a, mod int) int {
	a = (a % mod + mod) % mod
	for x := 1; x < mod; x++ {
		if (a*x)%mod == 1 {
			return x
		}
	}
	return -1 // обратного элемента не существует
}

// modExp выполняет быстрое возведение в степень по модулю.
func modExp(base, exponent, mod int) int {
	result := 1
	base = base % mod
	for exponent > 0 {
		if exponent%2 == 1 {
			result = (result * base) % mod
		}
		exponent /= 2
		base = (base * base) % mod
	}
	return result
}

// isPrimeFermat проверяет, является ли число n простым.
func isPrimeFermat(n, iterations int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i < iterations+2; i++ {
		if modExp(i, n-1, n) != 1 {
			return false
		}
	}
	return true
}

// Add выполняет сложение точек P и Q.
func (curve *EllipticCurve) Add(P, Q Point) Point {
	if P.Infinity {
		return Q
	}
	if Q.Infinity {
		return P
	}

	if P.X == Q.X {
		if (P.Y+Q.Y)%curve.P == 0 {
			return Point{Infinity: true}
		}
		return curve.Double(P)
	}

	numerator := (Q.Y - P.Y + curve.P) % curve.P
	denominator := (Q.X - P.X + curve.P) % curve.P
	lambda := (numerator * modInverse(denominator, curve.P)) % curve.P

	x3 := (lambda*lambda - P.X - Q.X + curve.P) % curve.P
	y3 := (lambda*(P.X-x3) - P.Y + curve.P) % curve.P

	return Point{X: x3, Y: y3, Infinity: false}
}

// Double удваивает точку P на кривой.
func (curve *EllipticCurve) Double(P Point) Point {
	if P.Infinity || P.Y == 0 {
		return Point{Infinity: true}
	}

	numerator := (3*P.X*P.X + curve.A) % curve.P
	denominator := (2 * P.Y) % curve.P
	lambda := (numerator * modInverse(denominator, curve.P)) % curve.P

	x3 := (lambda*lambda - 2*P.X + curve.P) % curve.P
	y3 := (lambda*(P.X-x3) - P.Y + curve.P) % curve.P

	return Point{X: x3, Y: y3, Infinity: false}
}

// ScalarMult выполняет умножение точки на скаляр k.
func (curve *EllipticCurve) ScalarMult(P Point, k int) Point {
	result := Point{Infinity: true}
	addend := P

	for k > 0 {
		if k%2 == 1 {
			result = curve.Add(result, addend)
		}
		addend = curve.Double(addend)
		k /= 2
	}
	return result
}

// Points вычисляет все точки эллиптической кривой (наивный перебор).
func (curve *EllipticCurve) Points() []Point {
	points := []Point{}
	for x := 0; x < curve.P; x++ {
		fx := (x*x*x + curve.A*x + curve.B) % curve.P
		for y := 0; y < curve.P; y++ {
			if (y*y)%curve.P == fx {
				points = append(points, Point{X: x, Y: y, Infinity: false})
			}
		}
	}
	points = append(points, Point{Infinity: true})
	return points
}

func main() {
	p := 17
	a := 2
	b := 2
	curve := EllipticCurve{P: p, A: a, B: b}

	fmt.Printf("Эллиптическая кривая: y^2 = x^3 + %d*x + %d над F_%d\n", a, b, p)
	pts := curve.Points()
	fmt.Printf("Найдено точек: %d\n", len(pts))
	for _, pt := range pts {
		if pt.Infinity {
			fmt.Println("Infinity")
		} else {
			fmt.Printf("(%d, %d)\n", pt.X, pt.Y)
		}
	}

	if len(pts) > 1 {
		P := pts[0]
		Q := pts[1]
		R := curve.Add(P, Q)
		fmt.Printf("P = (%d, %d)\nQ = (%d, %d)\nP+Q = (%d, %d)\n", P.X, P.Y, Q.X, Q.Y, R.X, R.Y)
	}

	if len(pts) > 0 {
		P := pts[0]
		k := 3
		R := curve.ScalarMult(P, k)
		fmt.Printf("%d * P = (%d, %d)\n", k, R.X, R.Y)
	}

	if isPrimeFermat(p, 5) {
		fmt.Printf("%d является простым числом (по тесту Ферма).\n", p)
	} else {
		fmt.Printf("%d не является простым числом.\n", p)
	}
}
