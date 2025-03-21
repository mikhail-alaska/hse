package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func modPow(a, x, n int) int {
	result := 1
	for x > 0 {
		if x%2 == 1 {
			result = (result * a) % n
		}
		a = (a * a) % n
		x /= 2
	}
	return result
}

func checkPrimeFerma(n int) bool {
	return modPow(2, n-1, n) == 1
}

func checkPrime(n int) bool {
	if n == 2 {
		return true
	}
	if n == 1 || n%2 == 0 {
		return false
	}
	if !checkPrimeFerma(n) {
		return false
	}
	for d := 3; d*d <= n; d += 2 {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func extendedEuclid(a, b int) (int, int, int) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, x1, y1 := extendedEuclid(b, a%b)
	x := y1
	y := x1 - (a/b)*y1
	return gcd, x, y
}

func generateOpenKey(p, g, k int) int {
	return modPow(g, k, p)
}

func encodeBlock(h, m, p, g, k int) (int, int) {
	c1 := modPow(g, k, p)
	c2 := (m * modPow(h, k, p)) % p
	return c1, c2
}

func decodeBlock(x, c1, c2, p int) int {
	c1Exp := modPow(c1, x, p)
	_, _, inv := extendedEuclid(p, c1Exp)
	inv = ((inv % p) + p) % p
	result := (c2 * inv) % p
	return result
}

func readFileBytes(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func writeFileBytes(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

func encrypt(byteString []byte, key, p, g int) ([]int, int) {
	bitLen := 0
	temp := p - 2
	for temp > 0 {
		bitLen++
		temp /= 2
	}
	lenMsg := (bitLen - 1) / 8
	if lenMsg < 1 {
		lenMsg = 1
	}

	byteString = append(byteString, 255)
	padding := lenMsg - (len(byteString) % lenMsg)
	for i := 0; i < padding; i++ {
		byteString = append(byteString, 0)
	}

	result := []int{}
	for i := 0; i < len(byteString); i += lenMsg {
		blockInt := 0
		for j := 0; j < lenMsg; j++ {
			blockInt = (blockInt << 8) | int(byteString[i+j])
		}
		k := rand.Intn(p-3) + 2
		c1, c2 := encodeBlock(key, blockInt, p, g, k)
		fmt.Printf("block=%d\nk=%d\n", blockInt, k)
		result = append(result, c1, c2)
	}
	return result, lenMsg
}

func decrypt(byteString []byte, key, p, g int) []int {
	bitLen := 0
	temp := p - 2
	for temp > 0 {
		bitLen++
		temp /= 2
	}
	lenMsg := (bitLen - 1) / 8
	if lenMsg < 1 {
		lenMsg = 1
	}
	result := []int{}
	blockSize := 2 * (lenMsg + 1)
	for i := 0; i < len(byteString); i += blockSize {
		c1 := 0
		for j := 0; j < lenMsg+1; j++ {
			c1 = (c1 << 8) | int(byteString[i+j])
		}
		c2 := 0
		for j := 0; j < lenMsg+1; j++ {
			c2 = (c2 << 8) | int(byteString[i+lenMsg+1+j])
		}
		decodedBlock := decodeBlock(key, c1, c2, p)
		fmt.Printf("decodedBlock=%d\nc1=%d\nc2=%d\n", decodedBlock, c1, c2)
		result = append(result, decodedBlock)
	}
	return result
}

func writeIntsToFile(result []int, filename string, lenMsg int) error {
	data := []byte{}
	for _, num := range result {
		b := make([]byte, lenMsg)
		for i := lenMsg - 1; i >= 0; i-- {
			b[i] = byte(num & 0xFF)
			num >>= 8
		}
		data = append(data, b...)
	}
	return writeFileBytes(filename, data)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	p := 293
	g := 4
	closedKey := 5
	openKey := 0

	fmt.Println("Выберите режим работы:")
	fmt.Println("1. Шифрование/дешифрование")
	fmt.Println("2. Генерация открытого ключа")

	var mode int
	_, err := fmt.Scan(&mode)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		os.Exit(1)
	}

	switch mode {
	case 1:
		inBytes, err := readFileBytes("in.txt")
		if err != nil {
			fmt.Println("Ошибка чтения in.txt:", err)
			return
		}
		openKey = generateOpenKey(p, g, closedKey)
		encrypted, lenMsg := encrypt(inBytes, openKey, p, g)
		err = writeIntsToFile(encrypted, "out_en.txt", lenMsg+1)
		if err != nil {
			fmt.Println("Ошибка записи out_en.txt:", err)
			return
		}
		encBytes, err := readFileBytes("out_en.txt")
		if err != nil {
			fmt.Println("Ошибка чтения out_en.txt:", err)
			return
		}
		decrypted := decrypt(encBytes, closedKey, p, g)
		err = writeIntsToFile(decrypted, "out_dec.txt", 1)
		if err != nil {
			fmt.Println("Ошибка записи out_dec.txt:", err)
			return
		}
	case 2:
		fmt.Println("Открытый ключ:", generateOpenKey(p, g, closedKey))
	default:
		fmt.Println("Неверный режим")
	}
}
