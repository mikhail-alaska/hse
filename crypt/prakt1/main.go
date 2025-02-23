package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
)

var digits string = "0123456789abcdefghijklmnopqrstuvwxyz"

func SliceEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func FindSlice(a []int, b [][]int) int {
	for i, j := range b {
		if SliceEq(a, j) {
			return i
		}
	}
	return -1
}

func AsIntegerRatio(num float64) (int, int) {
	a, b := num, 1.0
	for a != math.Floor(a) {
		a *= 10
		b *= 10
	}
	return int(a), int(b)
}

func AllZeroes(arr []int) bool {
	for _, v := range arr {
		if v != 0 {
			return false
		}
	}
	return true
}

func Power(a, b int) int {
	res := 1
	for i := 0; i < b; i++ {
		res = res * a
	}
	return res
}

func TrimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func RemoveZeroes(a []int) []int {
	if AllZeroes(a) {
		return nil
	}
	for {
		if a[len(a)-1] != 0 {
			return a
		}
		a = a[:len(a)-1]
	}
}

func Convert(number, base int) string {
	if base > len(digits) {
		return ""
	}
	var result []string
	for {
		if number == 0 {
			break
		}
		result = append([]string{string(digits[number%base])}, result...)
		number = number / base
	}
	return strings.Join(result, "")
}

func ZeroFill(s string, n int) string {
    r := strings.Split(s, "")
    for i := 0; i < n-len(s); i++ {
        r = append([]string{"0"}, r...)
    }
    return strings.Join(r, "")
}

func CreateGalua(p, n int) [][]int {

	field := [][]int{}
	for i := 0; i < Power(p, n); i++ {
		a := Convert(i, p)
		poly := []int{}
		zeroa := ZeroFill(a, n)
		zeroa = ZeroFill(zeroa, n)
		zeroa = ZeroFill(zeroa, n)
		for _, z := range zeroa {
			poly = append(poly, int(z-'0'))
		}
		field = append(field, poly)
	}
	return field
}

func PolySum(a, b []int, n int) []int {
	m := min(len(a), len(b))
	res := []int{}
	for i := 0; i < m; i++ {
		x := (a[i] + b[i]) % n
		res = append(res, x)
	}
	return res
}

func Neprevodim(p, n int) ([][]int, error) {
	b := [][]int{}
	for i := Power(p, n); i < Power(p, n+1); i++ {
		a := Convert(i, p)
		flag := 0
		poly := []int{}
		for _, j := range a {

			poly = append(poly, int(j-'0'))
		}
		for j := 0; j < p; j++ {
			sum := 0
			for l := 0; l < len(poly); l++ {

				sum = sum + poly[l]*(Power(j, n+1-l-1))
			}
			if sum%p == 0 {
				flag = 1
			}
		}

		if n%2 == 0 {
			galua := CreateGalua(p, n)

			galua1 := galua[1:]
			for _, j := range galua1 {

				if j[1:] == nil {
					continue
				}
				res, err := PolyDiv(poly, j, p)

				if res == nil {
					flag = 1
				}
				if err != nil {
					return b, errors.New("Нет неприводимых значений")
				}
			}
		}
		if flag == 0 {
			b = append(b, poly)
		}
	}
	if len(b) != 0 {
		return b, nil
	}
	return b, errors.New("Нет неприводимых значений")
}

func PolyMult(a, b []int, p int, nepr []int) []int {
	result := make([]int, len(a)+len(b)-1)
	for i := range a {
		for j := range b {
			result[i+j] += a[i] * b[j]
		}
	}
	for i := range result {
		result[i] %= p
	}
	result, err := PolyDiv(result, nepr, p)
	if err != nil {
		panic(err)
	}
	result = RemoveZeroes(result)
	for {
		if len(result) >= len(a) {
			break
		}
		result = append(result, 0)
	}
	for i := range result {
		result[i] %= p
	}
	return result
}

func FindMultip(galua [][]int, p, n int, neprevod []int) [][]interface{} {
	b := [][]interface{}{}

	for i := 1; i < len(galua); i++ {
		if len(galua[i]) > 1 && AllZeroes(galua[i][1:]) && galua[i][0] == 1 {
			continue
		}

		flag := false
		proizved := append([]int(nil), galua[i]...)
		st := 1
		steps := [][]interface{}{{proizved, st}}

		for !flag {
			proizved = PolyMult(proizved, galua[i], p, neprevod)
			st++
			steps = append(steps, []interface{}{proizved, st})

			if len(proizved) > 1 && AllZeroes(proizved[1:]) && proizved[0] == 1 {
				flag = true
			}
		}

		if st == (Power(p, n) - 1) {
			b = append(b, []interface{}{galua[i], steps})
		}
	}

	return b
}

func FracModule(p int, num float64) float64 {
	a, b := AsIntegerRatio(num)
	a %= p
	if a < 0 {
		a += p
	}
	for i := 1; i < p; i++ {

		if (b*i)%p == 1 {
			b = i
			break
		}
	}
	result := float64((a * b) % p)
	if result < 0 {
		result += float64(p)
	}
	return result

}

func PolyDiv(a, b []int, p int) ([]int, error) {

	c1 := []int{}
	c2 := []int{}
	for _, i := range a {
		c1 = append(c1, i)
	}
	for _, i := range b {
		c2 = append(c2, i)
	}
	c1 = RemoveZeroes(c1)
	c2 = RemoveZeroes(c2)
	if b == nil {
		return nil, errors.New("Division by zero")
	}
	if len(c1) < len(c2) {
		return a, nil
	}

	if len(c1) == len(c2) && c1[len(c1)-1] < c2[len(c2)-1] {
		c1, c2 = c2, c1
	}

	for len(c1) >= len(c2) {

            fmt.Println(11)
		i := len(c1) - 1
		j := len(c2) - 1
		k := c1[i] / c2[j]
		for i >= 0 && j >= 0 {
			c1[i] -= c2[j] * k
			i--
			j--
		}
		c1 = RemoveZeroes(c1)
	}
	for len(c1) != len(c2) {
		c1 = append(c1, 0)
	}

	for i := range c1 {
		for c1[i] < 0 {
			c1[i] += p
		}
	}
	return c1, nil
}

func Opposite(galua [][]int, k int, nepr []int, p int) []int {
	for i := 1; i < len(galua); i++ {
		temp := PolyMult(galua[k], galua[i], p, nepr)

		one := make([]int, len(temp))
		one[0] = 1
		if SliceEq(temp, one) {
			return galua[i]
		}
	}
	return nil
}

func Alpha(stepen int) []int {
	fmt.Println("Ввод альфа: введите", stepen, "коэфицентов, начиная с свободного члена")
	key_alpha := []int{}
	for i := 0; i < stepen; i++ {
		s := -1
		for s < 0 || s >= stepen {
			r := bufio.NewReader(os.Stdin)
			char1, _, err := r.ReadRune()
			if err != nil {
				fmt.Println("i have a problem with your input")
			}
			s = int(char1 - '0')
			if s < 0 || s >= stepen {
				fmt.Println("try to enter again", s)
			}
		}
		if len(key_alpha) > 0 {
			key_alpha = slices.Insert(key_alpha, 0, s)
		} else {
			key_alpha = append(key_alpha, s)
		}
	}
	if key_alpha == nil {
		panic("Error ERROR ОШИБКА ОШИБКА")
	}
	return key_alpha
}

func Beta(stepen int) []int {
	fmt.Println("Ввод бета: введите", stepen, "коэфицентов, начиная с свободного члена")
	key_beta := []int{}
	for i := 0; i < stepen; i++ {
		s := -1
		for s < 0 || s >= stepen {
			r := bufio.NewReader(os.Stdin)
			char1, _, err := r.ReadRune()
			if err != nil {
				fmt.Println("i have a problem with your input")
			}
			s = int(char1 - '0')
			if s < 0 || s >= stepen {
				fmt.Println("try to enter again", s)
			}
		}
		if len(key_beta) > 0 {
			key_beta = slices.Insert(key_beta, 0, s)
		} else {
			key_beta = append(key_beta, s)
		}
	}
	if key_beta == nil {
		panic("Error ERROR ОШИБКА ОШИБКА")
	}
	return key_beta
}

func AffineEncode(message,
	alphabet string,
	stepen int,
	key_alpha, key_beta []int,
	nepr []int) string {
	shifr := []string{}
	alph := strings.Split(alphabet, "")

	for i := 0; i < len(message); i++ {
		if !strings.Contains(alphabet, string(message[i])) {
			shifr = append(shifr, string(message[i]))
		} else {
			for k, j := range alphabet {
				if string(message[i]) == string(j) {
					bukva := Convert(k, 2)
					bukva_galua := []int{}
					temp := ZeroFill(bukva, stepen)
					for _, z := range temp {
						if len(bukva_galua) > 0 {
							bukva_galua = slices.Insert(bukva_galua, 0, int(z-'0'))
						} else {
							bukva_galua = append(bukva_galua, int(z-'0'))
						}
					}
					multiplication := PolyMult(bukva_galua, key_alpha, 2, nepr)
					sum := PolySum(multiplication, key_beta, 2)
					index := 0
					for i1 := range sum {
						index = index + Power(2, i1)*sum[i1]
					}
					shifr = append(shifr, alph[index%len(alphabet)])
					break
				}
			}
		}

	}

	return strings.Join(shifr, "")
}

func AffineDecode(message,
	alphabet string,
	stepen int,
	key_alpha, key_beta []int,
	nepr []int) string {

	shifr := []string{}
	alph := strings.Split(alphabet, "")

	for i := 0; i < len(message); i++ {
		if !strings.Contains(alphabet, string(message[i])) {
			shifr = append(shifr, string(message[i]))
		} else {
			for k, j := range alphabet {
				if string(message[i]) == string(j) {
					bukva := Convert(k, 2)
					bukva_galua := []int{}
					temp := ZeroFill(bukva, stepen)
					temp = ZeroFill(temp, stepen)
					for _, z := range temp {
						if len(bukva_galua) > 0 {
							bukva_galua = slices.Insert(bukva_galua, 0, int(z-'0'))
						} else {
							bukva_galua = append(bukva_galua, int(z-'0'))
						}
					}
					sum := PolySum(bukva_galua, key_beta, 2)
					galua := CreateGalua(2, stepen)

					galuaindex := FindSlice(key_alpha, galua)
					if galuaindex == -1 {
						panic("WE HAVE A PROBLEM")
					}
					oppositeAlpha := Opposite(galua, galuaindex, nepr, 2)

					multiplication := PolyMult(sum, oppositeAlpha, 2, nepr)
					index := 0
					for i2 := range multiplication {
						index = int(index + Power(2, i2)*multiplication[i2])
					}
					shifr = append(shifr, alph[index%len(alphabet)])
					break
				}
			}
		}

	}

	return strings.Join(shifr, "")
}

func WorkPole() {
	fmt.Println("Все многочлены представлены в виде [число, x, x^2, x^3...]")

	fmt.Println("Введите через пробел числа p, n")
	in := bufio.NewReader(os.Stdin)
    var p, n int
    fmt.Fscan(in, &p, &n)
	fmt.Println("Размеры поля: p =", p, "n =", n)

	pole := CreateGalua(p, n)
	fmt.Println("Поле Галуа")
	fmt.Println(pole)

	neprev, err := Neprevodim(p, n)
	if err != nil {
		panic(err)
	}
	fmt.Println("Неприводимые элементы:")
	for c, i := range neprev {
		fmt.Println(c+1, i)
	}

	fmt.Println("Введите номер элемента с которым будем работать")
	neprChosenNumber := 1
    fmt.Fscan(in, &neprChosenNumber)
	nepr := neprev[neprChosenNumber-1]

	fmt.Println("Вами выбран неприводимый элемент", nepr)

	obraz := FindMultip(pole, p, n, nepr)
	fmt.Println("Образующие группы")
	for c, i := range obraz {
		fmt.Println(c+1, i[0])
	}


	fmt.Println("Введите номер элемента с которым будем работать")
	obrazChosen := 0
    fmt.Fscan(in, &obrazChosen)
	fmt.Println("Выбран", obraz[obrazChosen][0])
	steps, ok := obraz[obrazChosen][1].([][]interface{})
	if !ok {
		fmt.Println("Ошибка приведения типа")
		return
	}

	length := len(steps)
	for i := 0; i < length; i++ {
		fmt.Println("Степень", steps[i][1], steps[i][0])
	}

	flag := false
	for !flag {
		fmt.Println("Если вы хотите выполнить действия с многочленами - press [y]")
		fmt.Println("Если больше задач нет - press [n]")

		inp := 'n'
		switch inp {
		case 'y':
			fmt.Println("Введите два многочлена 1+ x +2x^2 + 3 x^3 ... в виде 123...")
			fmt.Println("\nОбратите внимание длина многочлена -", n)
		case 'n':
			flag = true
		}
	}
	fmt.Println("Программа выполнилась успешно")
}

func WorkShifr() {
	fmt.Println("Работа с шифром")

	fmt.Println("Введите сообщение")
	in := bufio.NewReader(os.Stdin)
	var message string
	alphabet := "123456abcdefghijklmnopqrstuvwxyz"
	input, err := in.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}
	input = strings.TrimSpace(input)
    message = input

	stepen := 0

	for Power(2, stepen) < len(alphabet) {
		stepen += 1
	}
	key_a := Alpha(stepen)
	key_b := Beta(stepen)

	fmt.Println("Все возможные неприводимые элементы")
	nepr_all, err := Neprevodim(2, stepen)
	if err != nil {
		fmt.Println("ОШибка Ошиткбка", err)
		return
	}
	for i, j := range nepr_all {
		fmt.Println(i+1, j)
	}

	fmt.Println("Введите номер элемента с которым будем работать")
	nepr_c := 0
	fmt.Fscan(in, &nepr_c)

	nepr_chosen := nepr_all[nepr_c]

	fmt.Println("Сначала шифруем")

	shifr := AffineEncode(message, alphabet, stepen, key_a, key_b, nepr_chosen)
	fmt.Println("Зашифрованное сообщение", shifr)

	fmt.Println("Теперь расшифруем")

	deshifr := AffineDecode(shifr, alphabet, stepen, key_a, key_b, nepr_chosen)
	fmt.Println("Расшифрованное сообщение", deshifr)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Println(`Введите "y", если надо поработать с полем и "n", если с шифром`)
	var what string
	fmt.Fscan(in, &what)

	switch what {
	case "y":
		WorkPole()
	case "n":
		WorkShifr()
	}
}
