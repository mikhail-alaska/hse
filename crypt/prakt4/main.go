package main

import (
    "bufio"
    "fmt"
    "math"
    "math/rand"
    "os"
    "strconv"
    "strings"
    "time"
)

// modPow выполняет a^x mod n
func modPow(a, x, n int) int {
    result := 1
    base := a % n
    exp := x

    for exp > 0 {
        if exp&1 == 1 {
            result = (result * base) % n
        }
        base = (base * base) % n
        exp >>= 1
    }
    return result
}

// checkPrimeFerma проверяет простоту числа n с помощью малой теоремы Ферма (основание 2)
func checkPrimeFerma(n int) bool {
    tmp := modPow(2, n-1, n)
    return tmp == 1
}

// checkPrime делает дополнительную проверку простоты
func checkPrime(n int) bool {
    if n == 2 {
        return true
    }
    if !checkPrimeFerma(n) {
        return false
    }
    d := 3
    limit := int(math.Sqrt(float64(n)))
    for d <= limit {
        if n%d == 0 {
            return false
        }
        d += 2
    }
    return true
}

// testFerma просто вызывает checkPrimeFerma
func testFerma(n int) bool {
    return checkPrimeFerma(n)
}

// extendedEuclid возвращает (x, y, d), такие что a*x + b*y = d = gcd(a, b)
func extendedEuclid(a, b int) (int, int, int) {
    if b == 0 {
        return 1, 0, a
    }
    x1, y1, d := extendedEuclid(b, a%b)
    x := y1
    y := x1 - (a/b)*y1
    return x, y, d
}

// generateOpenKey рассчитывает g^k mod p
func generateOpenKey(p, g, k int) int {
    return modPow(g, k, p)
}

// encodeBlock шифрует один блок (байт) m
func encodeBlock(p, g, k, m int) (int, int) {
    c1 := modPow(g, k, p)
    c2 := (m * modPow(p, k, p)) % p
    return c1, c2
}

// decodeBlock расшифровывает один блок
func decodeBlock(p, c1, c2, key int) int {
    c1Inv := modPow(c1, p-1-key, p)
    return (c2 * c1Inv) % p
}

// encrypt шифрует массив байтов (каждый элемент byte_string — это int от 0 до 255)
func encrypt(byteString []int, key, p, g int) ([]int, int) {
    lenMsg := len(byteString)
    rand.Seed(time.Now().UnixNano())
    // Случайный k
    k := rand.Intn(p-3) + 2 // randint(2, p-2) эквивалент

    // Добавим 255 и заполним нулями
    byteString = append(byteString, 255)
    for len(byteString) < 2*lenMsg {
        byteString = append(byteString, 0)
    }

    var encrypted []int
    for i := 0; i < lenMsg; i++ {
        block := byteString[i]
        c1, c2 := encodeBlock(p, g, k, block)
        encrypted = append(encrypted, c1, c2)
    }
    return encrypted, lenMsg
}

// decrypt расшифровывает массив int (каждые два числа — это c1, c2)
func decrypt(byteString []int, key, p, g int) string {
    lenMsg := len(byteString) / 2
    var result []int

    for i := 0; i < lenMsg; i++ {
        c1 := byteString[2*i]
        c2 := byteString[2*i+1]
        decoded := decodeBlock(p, c1, c2, key)
        result = append(result, decoded)
    }

    // Убираем завершающие нули
    for len(result) > 0 && result[len(result)-1] == 0 {
        result = result[:len(result)-1]
    }
    // Убираем 255, если есть
    if len(result) > 0 && result[len(result)-1] == 255 {
        result = result[:len(result)-1]
    }

    // Превращаем в байты и декодируем как utf-8
    bytesRes := make([]byte, len(result))
    for i := 0; i < len(result); i++ {
        bytesRes[i] = byte(result[i])
    }
    return string(bytesRes)
}

// readFile читает весь файл как []int (каждый байт — int)
func readFile(filename string) []int {
    data, err := os.ReadFile(filename)
    if err != nil {
        fmt.Println("Ошибка чтения файла:", err)
        return nil
    }
    result := make([]int, len(data))
    for i, b := range data {
        result[i] = int(b)
    }
    return result
}

// writeFile записывает первые 2*lenMsg элементов из result в файл
func writeFile(result []int, filename string, lenMsg int) {
    f, err := os.Create(filename)
    if err != nil {
        fmt.Println("Ошибка создания файла:", err)
        return
    }
    defer f.Close()

    writer := bufio.NewWriter(f)
    countToWrite := 2 * lenMsg
    if countToWrite > len(result) {
        countToWrite = len(result)
    }

    for i := 0; i < countToWrite; i++ {
        b := make([]byte, 1)
        b[0] = byte(result[i])
        writer.Write(b)
    }
    writer.Flush()
}

// main аналогичен python-коду: предлагает режим работы и т.д.
func main() {
    reader := bufio.NewReader(os.Stdin)
    p := 293
    g := 4
    closedKey := 0
    openKey := 0

    fmt.Println("Выберите режим работы:")
    fmt.Println("1. шифрование/дешифрование")
    fmt.Println("2. генерация открытого ключа")
    modeLine, _ := reader.ReadString('\n')
    modeLine = strings.TrimSpace(modeLine)

    switch modeLine {
    case "1":
        fmt.Print("Введите имя входного файла: ")
        inFilename, _ := reader.ReadString('\n')
        inFilename = strings.TrimSpace(inFilename)

        fmt.Print("Введите имя выходного файла: ")
        outFilename, _ := reader.ReadString('\n')
        outFilename = strings.TrimSpace(outFilename)

        fmt.Print("Введите ключ: ")
        keyLine, _ := reader.ReadString('\n')
        keyLine = strings.TrimSpace(keyLine)
        key, _ := strconv.Atoi(keyLine)

        // Шифруем
        byteString := readFile(inFilename)
        encrypted, lenMsg := encrypt(byteString, key, p, g)
        writeFile(encrypted, outFilename, lenMsg)

        // Для демонстрации сразу расшифровываем
        fmt.Print("Введите имя файла для расшифрования: ")
        decFilename, _ := reader.ReadString('\n')
        decFilename = strings.TrimSpace(decFilename)

        fmt.Print("Введите ключ для расшифрования: ")
        decKeyLine, _ := reader.ReadString('\n')
        decKeyLine = strings.TrimSpace(decKeyLine)
        decKey, _ := strconv.Atoi(decKeyLine)

        encData := readFile(outFilename)
        decryptedText := decrypt(encData, decKey, p, g)

        // Запишем расшифрованный текст в decFilename (как UTF-8)
        if err := os.WriteFile(decFilename, []byte(decryptedText), 0644); err != nil {
            fmt.Println("Ошибка при записи расшифрованного текста:", err)
        }

    case "2":
        // Пример генерации открытого ключа
        // (подробности реализации зависят от ваших потребностей)
        fmt.Println("Генерация открытого ключа (пример)")

        fmt.Print("Введите закрытый ключ: ")
        closedKeyLine, _ := reader.ReadString('\n')
        closedKeyLine = strings.TrimSpace(closedKeyLine)
        cKey, err := strconv.Atoi(closedKeyLine)
        if err != nil {
            fmt.Println("Ошибка ввода:", err)
            return
        }
        closedKey = cKey

        openKey = generateOpenKey(p, g, closedKey)
        fmt.Println("Открытый ключ:", openKey)

    default:
        fmt.Println("Неверный режим")
    }
}
