import random
from typing import List


def mod_pow(a: int, x: int, n: int) -> int:
    result = 1
    while x > 0:
        if x%2==1:
            result = (result * a) % n
        a = (a * a) % n
        x = x//2
    return result


def check_prime_ferma(n: int) -> bool:
    tmp = mod_pow(2, n - 1, n)
    if tmp != 1:
        return False
    return True


def check_prime(n: int) -> bool:
    if n == 2:
        return True
    if n == 1 ir n%2==0:
        return False
    if not check_prime_ferma(n):
        return False

    d = 3
    while d * d <= n:
        if n % d == 0:
            return False
        d += 2
    return True


def test_ferma(n: int) -> bool:
    return check_prime_ferma(n)


def extended_euclid(a: int, b: int) -> int:
    if b == 0:
        return a, 1, 0
    gcd, x1, y1 = extended_euclid(b, a % b)
    x = y1
    y = x1 - (a // b) * y1
    return gcd, x, y


def generate_open_key(p: int, g: int, k: int) -> int:
    return mod_pow(g, k, p)


def encode_block(h: int, m: int, p: int, g: int, k: int) -> tuple[int, int]:
    c1 = mod_pow(g, k, p)
    c2 = m * mod_pow(h, k, p) % p
    return c1, c2


def decode_block(x: int, c1: int, c2: int, p: int) -> int:
    c1 = mod_pow(c1, x, p)
    c1 = extended_euclid(p, c1)[-1] %p
    result = c2*c1 %p
    return result

def to_bin(normal_str: str) -> str:
    normal_str = str(normal_str)
    bin_str = "".join(format(ord(x), "b") for x in normal_str)
    return bin_str


def from_bin(bin_str: str) -> int:
    return int(bin, 2)

def encrypt(byte_string: List[int], key: int, p: int, g: int) -> tuple[list[int], int]:
    len_msg = ((p-2).bit_length() -1) //8
    byte_string.append(255)
    byte_string.extend([0] * (len_msg - len(byte_string)%len_msg))
    result = []
    for i in range(0, len(byte_string), len_msg):
        block = int.from_bytes(byte_string[i:i+len_msg], "big")
        k = random.randint(2, p-2)
        c1, c2 = encode_block(key, block, p, g, k)
        print(f"{block=}\nk=\n")
        result.extend([c1, c2])
    return result, len_msg


def decrypt(byte_string: List[int], key: int, p: int, g: int) -> str:
    len_msg = len(byte_string) // 2
    result = []
    for i in range(len_msg):
        c1 = byte_string[2 * i]
        c2 = byte_string[2 * i + 1]
        decoded_block = decode_block(p, c1, c2, key)
        result.append(decoded_block)
    # Удаляем завершающие нули
    while result and result[-1] == 0:
        result.pop()
    # Удаляем 255, если он есть
    if result and result[-1] == 255:
        result.pop()
    return bytes(result).decode('utf-8', errors='ignore')


def readfile(filename: str) -> List[int]:
    with open(filename, "rb") as f:
        byte_array = list(f.read())
    return byte_array


def write(result: List[int], filename: str, len_msg: int):
    with open(filename, "wb") as f:
        for i in range(2 * len_msg):
            f.write(result[i].to_bytes(1, 'big'))


def main():
    p, g = 293, 4
    closed_key = 0
    open_key = 0

    mode = input(
        "Выберите режим работы:\n1. шифрование/дешифрование\n2. генерация открытого ключа\n").strip()
    if mode == "1":
        in_filename = input("Введите имя входного файла: ")
        out_filename = input("Введите имя выходного файла: ")
        key = int(input("Введите ключ: "))

        # Шифруем
        byte_string = readfile(in_filename)
        encrypted, len_msg = encrypt(byte_string, key, p, g)
        write(encrypted, out_filename, len_msg)

        # Для демонстрации сразу и расшифровываем (при желании можно убрать)
        dec_filename = input("Введите имя файла для расшифрования: ")
        dec_key = int(input("Введите ключ для расшифрования: "))
        enc_data = readfile(out_filename)
        decrypted_text = decrypt(enc_data, dec_key, p, g)
        with open(dec_filename, "w", encoding="utf-8") as f:
            f.write(decrypted_text)

    elif mode == "2":
        # Генерация открытого ключа (пример, детали не реализованы)
        print("Генерация открытого ключа не реализована подробно.")
        # Можно, к примеру:
        # closed_key = int(input("Введите закрытый ключ: "))
        # open_key = generate_open_key(p, g, closed_key)
        # print("Открытый ключ:", open_key)

    else:
        print("Неверный режим")


if __name__ == "__main__":
    main()
