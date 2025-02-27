package main

import (
    "encoding/hex"
    "fmt"
    "io"
    "log"
    "os"
)

// Вариант хранения S-box, rCon и invSbox в виде глобальных массивов (или срезов)

var sBox = [256]byte{
    0x63, 0x7C, 0x77, 0x7B, 0xF2, 0x6B, 0x6F, 0xC5, 0x30, 0x01, 0x67, 0x2B, 0xFE, 0xD7, 0xAB, 0x76,
    0xCA, 0x82, 0xC9, 0x7D, 0xFA, 0x59, 0x47, 0xF0, 0xAD, 0xD4, 0xA2, 0xAF, 0x9C, 0xA4, 0x72, 0xC0,
    0xB7, 0xFD, 0x93, 0x26, 0x36, 0x3F, 0xF7, 0xCC, 0x34, 0xA5, 0xE5, 0xF1, 0x71, 0xD8, 0x31, 0x15,
    0x04, 0xC7, 0x23, 0xC3, 0x18, 0x96, 0x05, 0x9A, 0x07, 0x12, 0x80, 0xE2, 0xEB, 0x27, 0xB2, 0x75,
    0x09, 0x83, 0x2C, 0x1A, 0x1B, 0x6E, 0x5A, 0xA0, 0x52, 0x3B, 0xD6, 0xB3, 0x29, 0xE3, 0x2F, 0x84,
    0x53, 0xD1, 0x00, 0xED, 0x20, 0xFC, 0xB1, 0x5B, 0x6A, 0xCB, 0xBE, 0x39, 0x4A, 0x4C, 0x58, 0xCF,
    0xD0, 0xEF, 0xAA, 0xFB, 0x43, 0x4D, 0x33, 0x85, 0x45, 0xF9, 0x02, 0x7F, 0x50, 0x3C, 0x9F, 0xA8,
    0x51, 0xA3, 0x40, 0x8F, 0x92, 0x9D, 0x38, 0xF5, 0xBC, 0xB6, 0xDA, 0x21, 0x10, 0xFF, 0xF3, 0xD2,
    0xCD, 0x0C, 0x13, 0xEC, 0x5F, 0x97, 0x44, 0x17, 0xC4, 0xA7, 0x7E, 0x3D, 0x64, 0x5D, 0x19, 0x73,
    0x60, 0x81, 0x4F, 0xDC, 0x22, 0x2A, 0x90, 0x88, 0x46, 0xEE, 0xB8, 0x14, 0xDE, 0x5E, 0x0B, 0xDB,
    0xE0, 0x32, 0x3A, 0x0A, 0x49, 0x06, 0x24, 0x5C, 0xC2, 0xD3, 0xAC, 0x62, 0x91, 0x95, 0xE4, 0x79,
    0xE7, 0xC8, 0x37, 0x6D, 0x8D, 0xD5, 0x4E, 0xA9, 0x6C, 0x56, 0xF4, 0xEA, 0x65, 0x7A, 0xAE, 0x08,
    0xBA, 0x78, 0x25, 0x2E, 0x1C, 0xA6, 0xB4, 0xC6, 0xE8, 0xDD, 0x74, 0x1F, 0x4B, 0xBD, 0x8B, 0x8A,
    0x70, 0x3E, 0xB5, 0x66, 0x48, 0x03, 0xF6, 0x0E, 0x61, 0x35, 0x57, 0xB9, 0x86, 0xC1, 0x1D, 0x9E,
    0xE1, 0xF8, 0x98, 0x11, 0x69, 0xD9, 0x8E, 0x94, 0x9B, 0x1E, 0x87, 0xE9, 0xCE, 0x55, 0x28, 0xDF,
    0x8C, 0xA1, 0x89, 0x0D, 0xBF, 0xE6, 0x42, 0x68, 0x41, 0x99, 0x2D, 0x0F, 0xB0, 0x54, 0xBB, 0x16,
}

var invSBox = [256]byte{
    0x52, 0x09, 0x6A, 0xD5, 0x30, 0x36, 0xA5, 0x38, 0xBF, 0x40, 0xA3, 0x9E, 0x81, 0xF3, 0xD7, 0xFB,
    0x7C, 0xE3, 0x39, 0x82, 0x9B, 0x2F, 0xFF, 0x87, 0x34, 0x8E, 0x43, 0x44, 0xC4, 0xDE, 0xE9, 0xCB,
    0x54, 0x7B, 0x94, 0x32, 0xA6, 0xC2, 0x23, 0x3D, 0xEE, 0x4C, 0x95, 0x0B, 0x42, 0xFA, 0xC3, 0x4E,
    0x08, 0x2E, 0xA1, 0x66, 0x28, 0xD9, 0x24, 0xB2, 0x76, 0x5B, 0xA2, 0x49, 0x6D, 0x8B, 0xD1, 0x25,
    0x72, 0xF8, 0xF6, 0x64, 0x86, 0x68, 0x98, 0x16, 0xD4, 0xA4, 0x5C, 0xCC, 0x5D, 0x65, 0xB6, 0x92,
    0x6C, 0x70, 0x48, 0x50, 0xFD, 0xED, 0xB9, 0xDA, 0x5E, 0x15, 0x46, 0x57, 0xA7, 0x8D, 0x9D, 0x84,
    0x90, 0xD8, 0xAB, 0x00, 0x8C, 0xBC, 0xD3, 0x0A, 0xF7, 0xE4, 0x58, 0x05, 0xB8, 0xB3, 0x45, 0x06,
    0xD0, 0x2C, 0x1E, 0x8F, 0xCA, 0x3F, 0x0F, 0x02, 0xC1, 0xAF, 0xBD, 0x03, 0x01, 0x13, 0x8A, 0x6B,
    0x3A, 0x91, 0x11, 0x41, 0x4F, 0x67, 0xDC, 0xEA, 0x97, 0xF2, 0xCF, 0xCE, 0xF0, 0xB4, 0xE6, 0x73,
    0x96, 0xAC, 0x74, 0x22, 0xE7, 0xAD, 0x35, 0x85, 0xE2, 0xF9, 0x37, 0xE8, 0x1C, 0x75, 0xDF, 0x6E,
    0x47, 0xF1, 0x1A, 0x71, 0x1D, 0x29, 0xC5, 0x89, 0x6F, 0xB7, 0x62, 0x0E, 0xAA, 0x18, 0xBE, 0x1B,
    0xFC, 0x56, 0x3E, 0x4B, 0xC6, 0xD2, 0x79, 0x20, 0x9A, 0xDB, 0xC0, 0xFE, 0x78, 0xCD, 0x5A, 0xF4,
    0x1F, 0xDD, 0xA8, 0x33, 0x88, 0x07, 0xC7, 0x31, 0xB1, 0x12, 0x10, 0x59, 0x27, 0x80, 0xEC, 0x5F,
    0x60, 0x51, 0x7F, 0xA9, 0x19, 0xB5, 0x4A, 0x0D, 0x2D, 0xE5, 0x7A, 0x9F, 0x93, 0xC9, 0x9C, 0xEF,
    0xA0, 0xE0, 0x3B, 0x4D, 0xAE, 0x2A, 0xF5, 0xB0, 0xC8, 0xEB, 0xBB, 0x3C, 0x83, 0x53, 0x99, 0x61,
    0x17, 0x2B, 0x04, 0x7E, 0xBA, 0x77, 0xD6, 0x26, 0xE1, 0x69, 0x14, 0x63, 0x55, 0x21, 0x0C, 0x7D,
}

var rCon = [18]byte{
    0x00, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40,
    0x80, 0x1B, 0x36, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00,
}

// keySize задаёт количество раундов шифрования (Nr) в зависимости от длины ключа AES
var keySize = map[int]int{
    16: 10, // AES-128
    24: 12, // AES-192
    32: 14, // AES-256
}

// matrix – это просто 4x4 массив байт.
// Удобно хранить состояние AES именно так.
type matrix [4][4]byte

// ---------- Функции умножения в поле GF(2^8) ----------
func mulBy02(b byte) byte {
    // Если самый старший бит == 1, то после << 1 надо сделать XOR c 0x1B.
    var shifted = b << 1
    if (b & 0x80) != 0 {
        shifted ^= 0x1B
    }
    return shifted
}

func mulBy03(b byte) byte {
    return mulBy02(b) ^ b
}

func mulBy09(b byte) byte {
    // 0x09 = (x^3 + 1) => 2*2*2(b) ^ b
    return mulBy02(mulBy02(mulBy02(b))) ^ b
}

func mulBy0b(b byte) byte {
    // 0x0b = (x^3 + x + 1) => 2*2*2(b) ^ 2(b) ^ b
    return mulBy02(mulBy02(mulBy02(b))) ^ mulBy02(b) ^ b
}

func mulBy0d(b byte) byte {
    // 0x0d = (x^3 + x^2 + 1)
    return mulBy02(mulBy02(mulBy02(b))) ^ mulBy02(mulBy02(b)) ^ b
}

func mulBy0e(b byte) byte {
    // 0x0e = (x^3 + x^2 + x)
    return mulBy02(mulBy02(mulBy02(b))) ^ mulBy02(mulBy02(b)) ^ mulBy02(b)
}

// ---------- Преобразования для матриц/байт ----------
func bytesToMatrix(src []byte) matrix {
    // Перекладываем 16 байт (src) в 4x4 матрицу
    var m matrix
    for i := 0; i < 16; i++ {
        m[i/4][i%4] = src[i]
    }
    return m
}

func matrixToBytes(m matrix) []byte {
    var out [16]byte
    idx := 0
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            out[idx] = m[i][j]
            idx++
        }
    }
    return out[:]
}

// XOR побайтно
func xorBytes(a, b []byte) []byte {
    n := len(a)
    if len(b) < n {
        n = len(b)
    }
    out := make([]byte, n)
    for i := 0; i < n; i++ {
        out[i] = a[i] ^ b[i]
    }
    return out
}

// Небольшая вспомогательная функция для замены 4 байт через sBox
func switchSBox(arr []byte) []byte {
    for i := 0; i < len(arr); i++ {
        arr[i] = sBox[arr[i]]
    }
    return arr
}

// ---------- Основные операции AES ----------

// SubBytes
func subBytes(state matrix) matrix {
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            state[i][j] = sBox[state[i][j]]
        }
    }
    return state
}

// InvSubBytes
func invSubBytes(state matrix) matrix {
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            state[i][j] = invSBox[state[i][j]]
        }
    }
    return state
}

// ShiftRows
func shiftRows(state matrix) matrix {
    // Перестановки по строкам
    // строка 0 не меняется
    state[0][1], state[1][1], state[2][1], state[3][1] =
        state[1][1], state[2][1], state[3][1], state[0][1]

    state[0][2], state[1][2], state[2][2], state[3][2] =
        state[2][2], state[3][2], state[0][2], state[1][2]

    state[0][3], state[1][3], state[2][3], state[3][3] =
        state[3][3], state[0][3], state[1][3], state[2][3]

    return state
}

// InvShiftRows
func invShiftRows(state matrix) matrix {
    // Обратная перестановка по строкам
    state[0][1], state[1][1], state[2][1], state[3][1] =
        state[3][1], state[0][1], state[1][1], state[2][1]

    state[0][2], state[1][2], state[2][2], state[3][2] =
        state[2][2], state[3][2], state[0][2], state[1][2]

    state[0][3], state[1][3], state[2][3], state[3][3] =
        state[1][3], state[2][3], state[3][3], state[0][3]

    return state
}

// MixColumns
func mixColumns(s matrix) matrix {
    var temp matrix
    for i := 0; i < 4; i++ {
        a0 := s[i][0]
        a1 := s[i][1]
        a2 := s[i][2]
        a3 := s[i][3]

        temp[i][0] = mulBy02(a0) ^ mulBy03(a1) ^ a2 ^ a3
        temp[i][1] = a0 ^ mulBy02(a1) ^ mulBy03(a2) ^ a3
        temp[i][2] = a0 ^ a1 ^ mulBy02(a2) ^ mulBy03(a3)
        temp[i][3] = mulBy03(a0) ^ a1 ^ a2 ^ mulBy02(a3)
    }
    return temp
}

// InvMixColumns
func invMixColumns(s matrix) matrix {
    var temp matrix
    for i := 0; i < 4; i++ {
        a0 := s[i][0]
        a1 := s[i][1]
        a2 := s[i][2]
        a3 := s[i][3]

        temp[i][0] = mulBy0e(a0) ^ mulBy0b(a1) ^ mulBy0d(a2) ^ mulBy09(a3)
        temp[i][1] = mulBy09(a0) ^ mulBy0e(a1) ^ mulBy0b(a2) ^ mulBy0d(a3)
        temp[i][2] = mulBy0d(a0) ^ mulBy09(a1) ^ mulBy0e(a2) ^ mulBy0b(a3)
        temp[i][3] = mulBy0b(a0) ^ mulBy0d(a1) ^ mulBy09(a2) ^ mulBy0e(a3)
    }
    return temp
}

// addRoundKey
func addRoundKey(state, roundKey matrix) matrix {
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            state[i][j] ^= roundKey[i][j]
        }
    }
    return state
}

// ---------- Расширение ключей (KeyExpansion) ----------

// Возвращает массив из (Nr+1) матриц по 4x4, каждая из которых — Round Key
func keyExpansion(key []byte) []matrix {
    length := len(key)
    if _, ok := keySize[length]; !ok {
        log.Fatalln("Ошибка длины ключа в его расширении:", length)
    }
    Nk := length / 4        // кол-во 32-битных слов в ключе
    Nr := keySize[length]   // кол-во раундов
    totalWords := 4 * (Nr + 1)

    // Преобразуем ключ в срез "слов" (каждое слово = 4 байта)
    words := make([][]byte, Nk)
    for i := 0; i < Nk; i++ {
        words[i] = key[4*i : 4*(i+1)]
    }

    // Генерируем остальные слова
    i := Nk
    for i < totalWords {
        temp := make([]byte, 4)
        copy(temp, words[i-1])

        if i%Nk == 0 {
            // Сдвиг байтов (циклический)
            t0 := temp[0]
            temp = temp[1:]
            temp = append(temp, t0)
            // Применяем sBox
            temp = switchSBox(temp)
            // XOR с rCon[i/Nk]
            temp[0] ^= rCon[i/Nk]
        } else if Nk > 6 && (i%Nk == 4) {
            // Для AES-256
            temp = switchSBox(temp)
        }
        // XOR temp c words[i - Nk]
        out := xorBytes(temp, words[i-Nk])
        words = append(words, out)
        i++
    }

    // Преобразуем words в список матриц
    var roundKeys []matrix
    idx := 0
    for j := 0; j < len(words)/4; j++ {
        var m matrix
        for row := 0; row < 4; row++ {
            copy(m[row][:], words[idx+row])
        }
        roundKeys = append(roundKeys, m)
        idx += 4
    }
    return roundKeys
}

// ---------- Шифрование одного блока (16 байт) ----------
func encryptBlock(openText, key []byte) []byte {
    if len(openText) != 16 {
        log.Fatalln("Ошибка на длине блока текста")
    }
    if _, ok := keySize[len(key)]; !ok {
        log.Fatalln("Ошибка на длине ключа при шифровании")
    }

    state := bytesToMatrix(openText)
    roundKeys := keyExpansion(key)
    Nr := keySize[len(key)]

    // Начальный раунд
    state = addRoundKey(state, roundKeys[0])

    // Раунды 1..(Nr-1)
    for r := 1; r < Nr; r++ {
        state = subBytes(state)
        state = shiftRows(state)
        state = mixColumns(state)
        state = addRoundKey(state, roundKeys[r])
    }

    // Финальный раунд
    state = subBytes(state)
    state = shiftRows(state)
    state = addRoundKey(state, roundKeys[Nr])

    return matrixToBytes(state)
}

// Расшифрование одного блока (16 байт)
func decryptBlock(cipherText, key []byte) []byte {
    if len(cipherText) != 16 {
        log.Fatalln("Ошибка на длине блока зашифрованного текста")
    }
    if _, ok := keySize[len(key)]; !ok {
        log.Fatalln("Ошибка на длине ключа при расшифровании")
    }

    state := bytesToMatrix(cipherText)
    roundKeys := keyExpansion(key)
    Nr := keySize[len(key)]

    // Начинаем с последнего roundKey
    state = addRoundKey(state, roundKeys[Nr])
    state = invShiftRows(state)
    state = invSubBytes(state)

    // Проходим раунды (Nr-1) .. 1
    for r := Nr - 1; r > 0; r-- {
        state = addRoundKey(state, roundKeys[r])
        state = invMixColumns(state)
        state = invShiftRows(state)
        state = invSubBytes(state)
    }

    // Финальный раунд
    state = addRoundKey(state, roundKeys[0])

    return matrixToBytes(state)
}

// ---------- Подготовка/запуск шифрования/дешифрования файлика ----------

// Шифрование файла
func preparationAndEncrypt(inputFilePath, keyHex string) {
    file, err := os.Open(inputFilePath)
    if err != nil {
        log.Fatalln("Не удалось открыть файл:", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        log.Fatalln("Ошибка чтения файла:", err)
    }

    // Генерируем имя выходного файла
    outputFileName := inputFilePath + ".crypted"

    // Разбиваем на блоки по 16 байт + паддинг
    blocks := splitIntoBlocks(data, 16)

    // Добавим «псевдо-PKCS7» (упрощённая логика из Python-версии):
    // если последний блок меньше 16, дополняем нулями, 
    // а в последний байт пишем число добавленных байтов (если оно < 15).
    last := blocks[len(blocks)-1]
    var counter byte
    for len(last) < 15 {
        last = append(last, 0x00)
        counter++
    }
    if len(last) == 15 {
        last = append(last, counter)
    }
    blocks[len(blocks)-1] = last

    // Переводим hex-ключ в []byte
    k, err := hex.DecodeString(keyHex)
    if err != nil {
        log.Fatalln("Неверный ключ (hex):", err)
    }

    // Шифруем каждый блок
    var result []byte
    for _, blk := range blocks {
        enc := encryptBlock(blk, k)
        result = append(result, enc...)
    }

    // Записываем результат
    err = os.WriteFile(outputFileName, result, 0644)
    if err != nil {
        log.Fatalln("Ошибка записи:", err)
    }
    fmt.Println("Файл зашифрован:", outputFileName)
}

// Дешифрование файла
func preparationAndDecrypt(inputFilePath, keyHex string) {
    file, err := os.Open(inputFilePath)
    if err != nil {
        log.Fatalln("Не удалось открыть файл:", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        log.Fatalln("Ошибка чтения файла:", err)
    }

    // Условно генерируем имя выходного файла
    outputFileName := "decrypted_" + inputFilePath

    // Разбиваем на блоки по 16 байт
    blocks := splitIntoBlocks(data, 16)

    k, err := hex.DecodeString(keyHex)
    if err != nil {
        log.Fatalln("Неверный ключ (hex):", err)
    }

    var decrypted []byte
    for _, blk := range blocks {
        dec := decryptBlock(blk, k)
        decrypted = append(decrypted, dec...)
    }

    // Снимаем «паддинг»
    // Из последнего блока берём последний байт как счётчик
    if len(decrypted) >= 16 {
        lastByte := decrypted[len(decrypted)-1]
        if lastByte > 0 && lastByte < 15 {
            decrypted = decrypted[:len(decrypted)-1-int(lastByte)]
        }
    }

    err = os.WriteFile(outputFileName, decrypted, 0644)
    if err != nil {
        log.Fatalln("Ошибка записи файла:", err)
    }
    fmt.Println("Файл расшифрован:", outputFileName)
}

// Вспомогательная функция для разделения на блоки
func splitIntoBlocks(data []byte, blockSize int) [][]byte {
    var result [][]byte
    for i := 0; i < len(data); i += blockSize {
        end := i + blockSize
        if end > len(data) {
            end = len(data)
        }
        chunk := make([]byte, end-i)
        copy(chunk, data[i:end])
        result = append(result, chunk)
    }
    if len(result) == 0 {
        // Если вдруг файл пуст
        result = append(result, []byte{})
    }
    return result
}

// ---------- Основная функция ----------

func main() {
    // Пример использования.
    //
    // go run main.go 1 example1.txt 000011112222333344445555666677778888999011101010
    // go run main.go 2 example1.txt.crypted 000011112222333344445555666677778888999011101010

    if len(os.Args) < 3 {
        fmt.Println("Usage: main <mode> <inputFile> <hexKey>")
        fmt.Println(" mode=1 => encrypt, mode=2 => decrypt")
        return
    }

    mode := os.Args[1]
    inputFilePath := os.Args[2]
    hexKey := "000011112222333344445555666677778888999011101010"

    switch mode {
    case "1":
        preparationAndEncrypt(inputFilePath, hexKey)
    case "2":
        preparationAndDecrypt(inputFilePath, hexKey)
    default:
        fmt.Println("Неизвестный режим (1=encrypt, 2=decrypt).")
    }
}
