import random
from typing import List


def mod_pow(a: int, x: int, n: int) -> int:
    result = 1
    while x > 0:
        if x % 2 == 1:
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
    if n == 1 or n % 2 == 0:
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
    c1 = extended_euclid(p, c1)[-1] % p
    result = c2*c1 % p
    return result


def to_bin(normal_str: str) -> str:
    normal_str = str(normal_str)
    bin_str = "".join(format(ord(x), "b") for x in normal_str)
    return bin_str


def from_bin(bin_str: str) -> int:
    return int(bin, 2)


def encrypt(byte_string: List[int], key: int, p: int, g: int) -> tuple[list[int], int]:
    len_msg = ((p-2).bit_length() - 1) // 8
    byte_string.append(255)
    byte_string.extend([0] * (len_msg - len(byte_string) % len_msg))
    result = []
    for i in range(0, len(byte_string), len_msg):
        block = int.from_bytes(byte_string[i:i+len_msg], "big")
        k = random.randint(2, p-2)
        c1, c2 = encode_block(key, block, p, g, k)
        print(f"{block=}\nk=\n")
        result.extend([c1, c2])
    return result, len_msg


def decrypt(string: str, key: int, p: int, g: int) -> str:
    byte_string = string
    len_msg = ((p-2).bit_length() - 1)//8
    result = []
    for i in range(0, len(byte_string), len_msg*2 + 2):
        c1 = int.from_bytes(byte_string[i:i+len_msg+1], "big")
        c2 = int.from_bytes(byte_string[i+len_msg+1:i+len_msg*2+2], "big")
        decoded_block = decode_block(key, c1, c2, p)
        print(f"{decoded_block=}\n{c1=}\n{c2=}\n")
        result.append(decoded_block)
    return result


def read(filename: str) -> list:
    with open(filename, "rb") as f:
        file = f.read()
    return list(file)


def readfile(filename: str):
    byte_array = []
    with open(filename, "rb") as file:
        byte = file.read(1)
        while byte:
            byte_array.append(int.from_bytes(byte, byteorder="big"))
            byte = file.read(1)
    return byte_array


def write(result: List[int], filename: str, len_msg: int):
    with open(filename, "wb") as f:
        for i in result:
            f.write(i.to_bytes(len_msg, "big"))


p, g = 293, 4
closed_key = 5
open_key = 0


def main(mode: int):
    match mode:
        case 1:
            in_str = read("in.txt")
            open_key = generate_open_key(p, g, closed_key)
            result, len_msg = encrypt(in_str, open_key, p, g)
            write(result, "out_en.txt", len_msg+1)
        case 2:
            in_str = read("out_en.txt")
            result, len_msg = decrypt(in_str, closed_key, p, g)
            write(result, "out_dec.txt", 1)
        case 3:


if __name__ == "__main__":
    main()
