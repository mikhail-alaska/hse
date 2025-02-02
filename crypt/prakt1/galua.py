import math
def convert(number, base, upper=False):
    literals = '123456abcdefghijklmnopqrstuvwxyz'
    if base > len(literals):
        return None
    result = ''
    while number > 0:
        result = literals[number % base] + result
        number //= base
    if upper:
        return result.upper() 
    else:
        return result


def create_galua(p, n):
    field = []
    for i in range(0, p ** n):
        a = convert(i, p, True)
        poly = []
        for j in a.zfill(n):
            poly.append(float(j))
        field.append(poly)
    return field

def poly_sum(a, b, n):
    zip_sum = zip(a, b)
    l = list(map(sum, zip_sum))
    result = []
    for i in range(len(l)):
        x = l[i]
        result.append(x%n)
    return result

def nepr(p, n):
    b = []
    for i in range(p** n, p ** (n + 1)):
        a = convert(i, p, True)
        poly = []
        flag = 0
        for j in a:
            poly.append(int(j))
        for k in range(p):
            sum = 0
            for l in range(len(poly)):
                sum = sum + poly[l] * (k ** (n + 1 - l - 1))
            if (sum % p == 0):
                flag = 1
        if (flag == 0):
            b.append(poly)
    if (len(b) != 0):
        return b
    else:
        return 0


def poly_mult(a, b, p, polynom):
    result = [0] * (len(a) + len(b) - 1)
    for i in range(len(a)):
        for j in range(len(b)):
            result[i + j] += a[i] * b[j]
    for i in range(len(result)):
        result[i] %= p
    result = div_poly(result, polynom)
    result = remove_trailing_zeros(result)
    while len(result) < len(a):
        result.append(0)
    for i in range(len(result)):
        result[i] = int(result[i]) %p
    return result


def remove_trailing_zeros(arr):
    if arr == [0]*len(arr):
        return []
    last_non_zero_index = len(arr) - next((i for i, x in enumerate(reversed(arr)) if x != 0), 0)
    return arr[:last_non_zero_index]

def frac_module(p, num: float):
    a, b = float(num).as_integer_ratio()
    a %= p
    for i in range(1,p):
        if b * i % p == 1:
            b = i
            break
    result = float(a * b % p)
    return result

def div_poly(a, b): #деление полиномов
    c1 = []
    c2 = []
    for i in a:
        c1.append(i)
    for i in b:
        c2.append(i)
    c1 = remove_trailing_zeros(c1)
    c2 = remove_trailing_zeros(c2)
    lc1 = len(c1)
    lc2 = len(c2)

    if b == [0] * len(b):
        raise ValueError("Division by zero")
    elif lc1 < lc2:

        return a
    elif lc1 == lc2 and c1[-1] < c2[-1]:
        c1, c2 = c2,c1
    while len(c1) >= len(c2):
        k = c1[-1]/c2[-1]
        i = len(c1)-1
        j = len(c2) - 1
        while i >=0 and j >= 0:
            c1[i] -= c2[j] * k
            i-=1
            j-=1
        c1 = remove_trailing_zeros(c1)


    while len(c1) != len(c2):
        c1.append(0)


    for i in range(len(c1)):
        while c1[i] < 0:
            c1[i] += p

        if c1[i] % 1 != 0:
            c1[i] = frac_module(p,c1[i])

    return c1



def opposite(field, k, nepriv, p): #поиск обратного элемента
    one = [0] * len(field[0])
    one[0] = 1
    for i in range(1, len(field)):
        if poly_mult(field[k], field[i], p, nepriv) == one:
            return field[i]

################################
def affine_shifr():

    message = str(input())
    alphabet = str(input())
    stepen = 0
    while 2 ** stepen < len(alphabet):
        stepen += 1

    print('Vvod alpha. Vpishite', stepen, 'kofficentof')
    key_alpha = []
    for i in range(stepen):
        simbol = -1
        while simbol<0 or simbol>=stepen:
            simbol = float(input())
            if simbol<0 or simbol>=stepen:
                print('Vvedite echo raz')
        key_alpha.insert(0, simbol)
    if(key_alpha == [0]*len(key_alpha)):
        print('ERROR: Key dont work')
        return 0

    alpha = 0
    for i in range(len(key_alpha)):
        alpha = int(alpha + 2 ** i * key_alpha[int(i)])
    if math.gcd(len(alphabet), alpha) !=1:
        print('Error. Bad alpha key')
        return 0


    print('Vvod beta. Vpishite ', stepen, 'kofficentof')
    key_beta = []
    for i in range(stepen):
        simbol = -1
        while simbol < 0 or simbol >= stepen:
            simbol = float(input())
            if simbol < 0 or simbol >= stepen:
                print('Vvedite echo raz')
        key_beta.insert(0, simbol)
    irreducible = nepr(2, stepen)[1]
    print( key_alpha, key_beta)

    alphabet = alphabet.lower()
    if message == '':
        print('Так не получится. Введите что-нибудь.')
    else:
        shift_text = ''
        for i in range(0, len(message)):
            for j in range(0, len(alphabet)):
                if message[i] == alphabet[j]:

                    bukva = convert(j ,2)
                    print(bukva, j)
                    bukva_galua =[]
                    for j in bukva.zfill(stepen):
                        bukva_galua.insert(0, float(j))
                    print(bukva_galua)
                    prozved = poly_mult(bukva_galua, key_alpha, 2, irreducible)
                    print(prozved)
                    sum = poly_sum(prozved, key_beta, 2)
                    print(sum)
                    index = 0
                    for i in range(len(sum)):
                        index = int(index + 2**i  * sum[int(i)])
                    print(index)
                    shift_text += alphabet[index%len(alphabet)]
                    break
        return shift_text
    #else:
    #    print('Так нельзя.')


def affine_shifr_decode(message,alphabet, key_alpha, key_betta):
    message = message.lower()
    alphabet = alphabet.lower()
    stepen = 0
    while 2 ** stepen < len(alphabet):
        stepen += 1

    if message == '':
        print('Так не получится. Введите что-нибудь.')
    else:
        shifr_text = ''

        for i in range(0, len(message)):
            for j in range(0, len(alphabet)):

                if message[i] == alphabet[j]:
                    bukva = convert(j, 2)
                    bukva_galua = []
                    for k in bukva.zfill(stepen):
                        bukva_galua.insert(0, float(k))
                    print(bukva , j)
                    print(bukva_galua)

                    minus_beta = []
                    for i in range(len(key_betta)):
                        minus_beta.insert(0, key_betta[i])
                    print(minus_beta)
                    sum = poly_sum(bukva_galua, minus_beta, 2)
                    print(sum)

                    galua = create_galua(2, stepen)
                    print(galua)
                    number = 0
                    for i in range(len(galua)):
                        if key_alpha == galua[i]:
                            number = i
                            break
                    print(number)
                    irreducible = nepr(2, stepen)[1]
                    print(nepr(2, stepen))
                    obrat_alpha = opposite(galua, number, irreducible, 2)
                    print(obrat_alpha)
                    prozved = poly_mult(sum, obrat_alpha, 2, irreducible)
                    print(prozved)
                    index = 0
                    for i in range(len(prozved)):
                        index = int(index + (2 ** i)* prozved[int(i)])
                        print(index)
                    shifr_text += alphabet[index%len(alphabet)]
                    print(index)
                    break

        return shifr_text

################################



p = 2
# n = 5
# print("Поле Галуа::")
# galois_field = create_gf(p, n)
# print(galois_field)
# print("Все возможные неприводимые элементы:")
# irreducible = nepr(p, n)
# print(irreducible)
# print("Умножение над полем Галуа.")
# print("Используемый неприводимый элемент поля:",irreducible[1])
# for i in range(1, len(galois_field)):
#     for j in range(1, len(galois_field)):
#         print(galois_field[i], "*",galois_field[j],"=",poly_mult(galois_field[i],galois_field[j],p,irreducible[1]))
# print("Все обратные элементы:")
# for k_opposite in range(1, len(galois_field)):
#     print(galois_field[k_opposite], "^-1 = ", opposite(galois_field, k_opposite, irreducible[1], p), sep="")
#
#print(affine_shifr())

print(affine_shifr_decode('as31xy14rsgxpv14r', 'qazwsxedcrfvtgbyhnujmikolp123456', [ 1,0, 1, 0, 1 ], [1, 0,0, 0, 1]))
