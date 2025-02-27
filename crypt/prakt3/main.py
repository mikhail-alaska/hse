import binascii

key_size = {16: 10, 24: 12, 32: 14}

s_box = (
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
    0x8C, 0xA1, 0x89, 0x0D, 0xBF, 0xE6, 0x42, 0x68, 0x41, 0x99, 0x2D, 0x0F, 0xB0, 0x54, 0xBB, 0x16
)

r_con = (
    0x00, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40,
    0x80, 0x1B, 0x36, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00
)

inv_sbox = (
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
    0x17, 0x2B, 0x04, 0x7E, 0xBA, 0x77, 0xD6, 0x26, 0xE1, 0x69, 0x14, 0x63, 0x55, 0x21, 0x0C, 0x7D
)


def mul_by_02(number):
    if number < 0x80:
        res = number << 1
    else:
        res = (number << 1) ^ 0x11B
    return res


def mul_by_03(number):
    return mul_by_02(number) ^ number


def bytes_to_matrix(text):
    return [list(text[i: i + 4]) for i in range(0, len(text), 4)]


def mul_by_09(number):
    return mul_by_02(mul_by_02(mul_by_02(number))) ^ number


def mul_by_0b(number):
    return mul_by_02(mul_by_02(mul_by_02(number))) ^ mul_by_02(number) ^ number


def mul_by_0d(number):
    return (
        mul_by_02(mul_by_02(mul_by_02(number))) ^
        mul_by_02(mul_by_02(number)) ^
        number
    )


def mul_by_0e(number):
    return (
        mul_by_02(mul_by_02(mul_by_02(number))) ^
        mul_by_02(mul_by_02(number)) ^
        mul_by_02(number)
    )


def matrix2bytes(matrix):
    return bytes(sum(matrix, []))


def switch_s_box(matrix):
    for i in range(4):
        matrix[i] = s_box[matrix[i]]
    return matrix


def xor_bytes(a, b):
    return bytes(a_i ^ b_i for a_i, b_i in zip(a, b))


def sub_bytes(matrix):
    for i in range(4):
        for j in range(4):
            matrix[i][j] = s_box[matrix[i][j]]
    return matrix


def inv_sub_bytes(matrix):
    for i in range(4):
        for j in range(4):
            matrix[i][j] = inv_sbox[matrix[i][j]]
    return matrix


def shift_rows(matrix):
    matrix[0][1], matrix[1][1], matrix[2][1], matrix[3][1] = (
        matrix[1][1], matrix[2][1], matrix[3][1], matrix[0][1]
    )
    matrix[0][2], matrix[1][2], matrix[2][2], matrix[3][2] = (
        matrix[2][2], matrix[3][2], matrix[0][2], matrix[1][2]
    )
    matrix[0][3], matrix[1][3], matrix[2][3], matrix[3][3] = (
        matrix[3][3], matrix[0][3], matrix[1][3], matrix[2][3]
    )
    return matrix


def inv_shift_rows(matrix):
    matrix[0][1], matrix[1][1], matrix[2][1], matrix[3][1] = (
        matrix[3][1], matrix[0][1], matrix[1][1], matrix[2][1]
    )
    matrix[0][2], matrix[1][2], matrix[2][2], matrix[3][2] = (
        matrix[2][2], matrix[3][2], matrix[0][2], matrix[1][2]
    )
    matrix[0][3], matrix[1][3], matrix[2][3], matrix[3][3] = (
        matrix[1][3], matrix[2][3], matrix[3][3], matrix[0][3]
    )
    return matrix


def mix_columns(state):
    state2 = [x[:] for x in state]
    for i in range(4):
        s0 = mul_by_02(state[i][0]) ^ mul_by_03(
            state[i][1]) ^ state[i][2] ^ state[i][3]
        s1 = state[i][0] ^ mul_by_02(state[i][1]) ^ mul_by_03(
            state[i][2]) ^ state[i][3]
        s2 = state[i][0] ^ state[i][1] ^ mul_by_02(
            state[i][2]) ^ mul_by_03(state[i][3])
        s3 = mul_by_03(state[i][0]) ^ state[i][1] ^ state[i][2] ^ mul_by_02(
            state[i][3])
        state2[i][0] = s0
        state2[i][1] = s1
        state2[i][2] = s2
        state2[i][3] = s3
    return state2


def inv_mix_columns(state):
    state2 = [x[:] for x in state]
    for i in range(4):
        s0 = mul_by_0e(state[i][0]) ^ mul_by_0b(state[i][1]) ^ mul_by_0d(
            state[i][2]) ^ mul_by_09(state[i][3])
        s1 = mul_by_09(state[i][0]) ^ mul_by_0e(state[i][1]) ^ mul_by_0b(
            state[i][2]) ^ mul_by_0d(state[i][3])
        s2 = mul_by_0d(state[i][0]) ^ mul_by_09(state[i][1]) ^ mul_by_0e(
            state[i][2]) ^ mul_by_0b(state[i][3])
        s3 = mul_by_0b(state[i][0]) ^ mul_by_0d(state[i][1]) ^ mul_by_09(
            state[i][2]) ^ mul_by_0e(state[i][3])
        state2[i][0] = s0
        state2[i][1] = s1
        state2[i][2] = s2
        state2[i][3] = s3
    return state2


def add_round_key(s, k):
    for i in range(4):
        for j in range(4):
            s[i][j] ^= k[i][j]
    return s


def key_expansion(key):
    if len(key) not in key_size:
        print("Ошибка длины ключа в его расширении")
        exit()

    rounds = key_size[len(key)]
    i = len(key) // 4
    Nk = len(key) // 4
    key_schedule = list(bytes_to_matrix(key))

    while i < 4 * (rounds + 1):
        word = list(key_schedule[-1])
        if i % Nk == 0:
            number = word.pop(0)
            word.append(number)
            word = switch_s_box(word)
            word[0] ^= r_con[i // Nk]
        elif Nk > 6 and i % Nk == 4:
            word = switch_s_box(word)

        word = xor_bytes(word, key_schedule[-Nk])
        key_schedule.append(word)
        i += 1

    return [key_schedule[4 * i: 4 * (i + 1)] for i in range(len(key_schedule) // 4)]


def encrypt(open_text, key):
    if len(open_text) != 16:
        print("Ошибка на длине блока текста")
        exit()

    open_text_state = bytes_to_matrix(open_text)
    key_matrix = key_expansion(key)
    if len(key) not in key_size:
        print("Ошибка на длине ключа с шифрованием")
        exit()

    Nr = key_size[len(key)]
    open_text_state = add_round_key(open_text_state, key_matrix[0])

    for i in range(1, Nr):
        open_text_state = sub_bytes(open_text_state)
        open_text_state = shift_rows(open_text_state)
        open_text_state = mix_columns(open_text_state)
        open_text_state = add_round_key(open_text_state, key_matrix[i])

    open_text_state = sub_bytes(open_text_state)
    open_text_state = shift_rows(open_text_state)
    open_text_state = add_round_key(open_text_state, key_matrix[-1])

    return matrix2bytes(open_text_state)


def preparation_and_encrypt(input_file_path, key):
    file = open(input_file_path, "rb")
    text = file.read()
    input_file_path = input_file_path.replace(".", "_")
    output_file_name = input_file_path + ".crypted"
    file.close()

    # Разбиваем на блоки по 16 байт
    t = [text[i: i + 16] for i in range(0, len(text), 16)]

    counter = 0
    while len(t[-1]) < 15:
        t[-1] += b"\x00"
        counter += 1
    if len(t[-1]) == 15:
        t[-1] += counter.to_bytes(1, byteorder="big")

    result = []
    k = binascii.unhexlify(key)
    for letter in t:
        result.append(encrypt(letter, k))

    with open(output_file_name, "wb") as file_output:
        for i in range(len(result)):
            file_output.write(result[i])


def decrypt(cipher_text, key):
    if len(cipher_text) != 16:
        print("Ошибка на длине блока зашифрованного текста")
        exit()

    cipher_state = bytes_to_matrix(cipher_text)
    key_matrix = key_expansion(key)
    if len(key) not in key_size:
        print("Ошибка на длине блока ключа в расшифровании")
        exit()

    Nr = key_size[len(key)]
    cipher_state = add_round_key(cipher_state, key_matrix[-1])
    cipher_state = inv_shift_rows(cipher_state)
    cipher_state = inv_sub_bytes(cipher_state)

    for i in range(Nr - 1, 0, -1):
        cipher_state = add_round_key(cipher_state, key_matrix[i])
        cipher_state = inv_mix_columns(cipher_state)
        cipher_state = inv_shift_rows(cipher_state)
        cipher_state = inv_sub_bytes(cipher_state)

    cipher_state = add_round_key(cipher_state, key_matrix[0])
    return matrix2bytes(cipher_state)


def preparation_and_decrypt(input_file_path, key):
    file = open(input_file_path, "rb")
    text = file.read()
    output_file_name = "decrypted-" + input_file_path[:-7]
    output_file_name = output_file_name[:-1]
    output_file_name = output_file_name.replace("_", ".")
    output_file_name = output_file_name[:-1]
    file.close()

    t = [text[i: i + 16] for i in range(0, len(text), 16)]
    result = []
    k = binascii.unhexlify(key)

    for i in t:
        result.append(decrypt(i, k))

    counter = int.from_bytes(result[-1][15:16], "big")
    if 1 <= counter <= 14:
        result[-1] = result[-1][:15 - counter]

    with open(output_file_name, "wb") as file_output:
        for i in range(len(result)):
            file_output.write(result[i])


def main(mode: int, input_file_path, key):
    match mode:
        case 1:
            preparation_and_encrypt(input_file_path, key)
        case 2:
            preparation_and_decrypt(input_file_path, key)


if __name__ == "__main__":
    key = "000011112222333344445555666677778888999011101010"
    # input_file_path = input()
    # key = input()
    input_file_path = "example1.txt"
    main(1, input_file_path, key)
    input_file_path = "example1_txt.crypted"
    main(2, input_file_path, key)
    # input_file_path = input()
    # key = input()
    # main(2, input_file_path, key)
