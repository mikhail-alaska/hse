digits = '0123456789'
def convert(number, base, upper=False):
    if base > len(digits):
        return None
    result = ''
    while number > 0:
        result = digits[number % base] + result
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
        poly, flag = [], 0

        for j in a:
            poly.append(int(j))
        for k in range(p):
            sum = 0
            for l in range(len(poly)):
                sum = sum + poly[l] * (k ** (n + 1 - l - 1))

            if (sum % p == 0):
                flag = 1
        if (n%2==0):
            pole_galua = create_galua(p,n)
            for ii in pole_galua[1:]:

                if ii[1:] == [0]*(len(ii)-1):
                    continue
                res = div_poly(poly, ii, p)
                if res == [0]*len(res):
                    flag=1
        if (flag == 0):
            b.append(poly)
    if (len(b) != 0):
        return b
    else:
        ValueError('Нет неприводимых элементов')


def poly_mult(a, b, p, polynom):
    result = [0] * (len(a) + len(b) - 1)
    for i in range(len(a)):
        for j in range(len(b)):
            result[i + j] += a[i] * b[j]
    for i in range(len(result)):
        result[i] %= p
    
    result = div_poly(result, polynom, p)
    result = remove_trailing_zeros(result)
    while len(result) < len(a):
        result.append(0)
    for i in range(len(result)):
        result[i] = int(result[i]) %p
    return result

def find_multip(galua,p, n, neprevod):
    b = []
    for i in range(1, len(galua)):
        if galua[i][1:]== [0]*(len(galua[i])-1) and galua[i][0]==1:
            continue
        flag = 0
        proizved = galua[i]
        st = 1
        step = []
        step.append((proizved, st))
        while flag==0:

            proizved = poly_mult(proizved, galua[i], p, neprevod)
            st += 1
            step.append((proizved, st))
            if proizved[1:]== [0]*(len(proizved)-1) and proizved[0]==1:
                flag = 1

        if st == p**n - 1:
            b.append((galua[i], step))
    return b


def remove_trailing_zeros(arr):
    if arr == [0]*len(arr):
        return []
    last_non_zero_index = len(arr) - next((i for i, x in enumerate(reversed(arr)) if x != 0), 0)
    return arr[:last_non_zero_index]

def frac_module(p, num: float): #деление по модулю дроби
    a, b = float(num).as_integer_ratio()
    a %= p
    for i in range(1,p):
        if b * i % p == 1:
            b = i
            break
    result = float(a * b % p)
    return result

def div_poly(a, b, p):
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



def opposite(field, k, nepriv, p):
    one = [0] * len(field[0])
    one[0] = 1
    for i in range(1, len(field)):
        temp =poly_mult(field[k], field[i], p, nepriv)  
        if temp== one:
            return field[i]


def alpha_in(stepen):
    print('Ввод альфа: Введите', stepen, 'коэфицентов, начиная с свободного члена')
    key_alpha = []
    for i in range(stepen):
        s = -1
        while s<0 or s>=stepen:
            s = int(input())
            if s<0 or s>=stepen:
                print('Попробуйте ввести число еще раз')
        key_alpha.insert(0, s)
    if(key_alpha == [0]*len(key_alpha)):
        return ValueError("error error ошибка ОШИБКА")
    return key_alpha


def beta_in(stepen):
    print('Ввод бета: Введите ', stepen, 'коэфицентов начиная со свободного члена')
    key_beta = []
    for i in range(stepen):
        s = -1
        while s < 0 or s >= stepen:
            s = int(input())
            if s < 0 or s >= stepen:
                print('Попробуйте ввести число еще раз!')
        key_beta.insert(0, s)
    return key_beta




def affine_shifr(message, alphabet, stepen, key_alpha, key_beta, neprevodim):
    shift_text = ''
    for i in range(0, len(message)):
        if message[i] not in (list(alphabet)):
            shift_text +=str(message[i])
        else:
            for j in range(0, len(alphabet)):
                if message[i] == alphabet[j]:

                    bukva = convert(j ,2)
                    bukva_galua =[]
                    for j in bukva.zfill(stepen):

                        bukva_galua.insert(0, float(j))
                    proizved = poly_mult(bukva_galua, key_alpha, 2, neprevodim)
#                    print("index", index,"sum", sum,"proizved", proizved, "bukva_galua", bukva_galua)
                    sum = poly_sum(proizved, key_beta, 2)
                    index = 0
                    for i in range(len(sum)):
                        index = int(index + 2**i  * sum[int(i)])

                    shift_text += alphabet[index%len(alphabet)]
                    break
    return shift_text


def affine_shifr_decode(message,alphabet,stepen, key_alpha, key_betta, neprevodim):

    shifr_text = ''

    for i in range(0, len(message)):
        if message[i] not in (list(alphabet)):
            shifr_text += message[i]

        else:
            for j in range(0, len(alphabet)):

                if message[i] == alphabet[j]:
                    bukva = convert(j, 2)
                    bukva_galua = []
                    for k in bukva.zfill(stepen):
                        bukva_galua.insert(0, float(k))

                    minus_beta = key_betta
                    sum = poly_sum(bukva_galua, minus_beta, 2)
                    galua = create_galua(2, stepen)


                    obrat_alpha = opposite(galua, galua.index(key_alpha), neprevodim, 2)
                    prozved = poly_mult(sum, obrat_alpha, 2, neprevodim)

                    index = 0
                    for i in range(len(prozved)):
                        index = int(index + (2 ** i)* prozved[int(i)])
                    print("opposite", obrat_alpha,"index", index,"sum", sum, "bukva_galua", bukva_galua)
                    shifr_text += alphabet[index%len(alphabet)]
                    break

    return shifr_text

################################
print('Хотите поработать с полем - press [y]\n'
      'Хотите поработать с шифрованием - press [n]')
if str(input()).lower() == 'y':
    print('Все многочлены представлены в виде [число, x, x^2, x^3...]')
    print('Задайте размеры Поля Галуа\n'
          'Последовательно введите параметры p и n')
#    p = int(input())
#    n = int(input())
    p = 2
    n = 3

    print("Поле Галуа::")
    pole_gal = create_galua(p, n)
    print(pole_gal)
    print("Все возможные неприводимые элементы:")
    neprevodim_all = nepr(p, n)
    for i in range(len(neprevodim_all)):
        print(i+1, neprevodim_all[i])
    print('Введите номер элемента с которым будем работать')
#    nepriv_vibran = neprevodim_all[int(input())-1]
    nepriv_vibran = neprevodim_all[0]
    print("Используемый неприводимый элемент поля:",nepriv_vibran)

    # print("Умножение над полем Галуа:")
    # for i in range(1, len(pole_gal)):
    #     for j in range(1, len(pole_gal)):
    #         print(pole_gal[i], "*",pole_gal[j],"=",poly_mult(pole_gal[i],pole_gal[j],p,nepriv_vibran))

    # print("\nВсе обратные элементы:")
    # for k_opposite in range(1, len(pole_gal)):
    #     print(pole_gal[k_opposite], "^-1 = ", opposite(pole_gal, k_opposite, nepriv_vibran, p), sep="")

    obraz = find_multip(pole_gal, p , n , nepriv_vibran)
    print('Выберите образующий группы:')
    for i in range(len(obraz)):
        print(i+1, obraz[i][0])
    ind = int(input())-1
    for i in range(len(obraz[ind][1])):
        print('Степень', obraz[ind][1][i][1], obraz[ind][1][i][0])


        # print("\nОбразующий группы:", obraz[0])
        # for k_obr in range(len(obraz[1])):
        #     print('Степень', obraz[1][k_obr][1] , obraz[1][k_obr][0])

    flagock = 0
    while flagock == 0:
        print('Если вы хотите выполнить действия с многочленами - press [y]' '\nЕсли больше задач нет - press [n]')
        if (str(input()).lower()=='y'):
            print('Введите два многочлена 1+ x +2x^2 + 3 x^3 ... в виде 123...' '\nОбратите внимание длина многочлена -', n)
            mnogoch1 = list(map(int, input()))
            mnogoch2 = list(map(int, input()))

            print('Если хотите выполнить сложение - press [y]'
                  '\nУмножение -  press [n]')

            if (str(input()).lower()=='y'):
                print('Результат сложения: ', poly_sum(mnogoch1, mnogoch2, p))
            else:
                print('Результат умножения:', poly_mult(mnogoch1, mnogoch2, p, nepriv_vibran))
        else:
            flagock=1

else:
    print("Работа с шифром")
    message = ''
    while message == '':
        print("Введите сообщение")
        message = str(input())
    message = message.lower()

    alphabet = '123456abcdefghijklmnopqrstuvwxyz'

    stepen = 0
    while 2 ** stepen < len(alphabet):
        stepen += 1
#    key_a = [1.0,1.0,0.0,0.0,0.0]    
#    key_b = [1.0,1.0,0.0,0.0,0.0]    
    key_a = alpha_in(stepen)
    key_b = beta_in(stepen)

    print("Все возможные неприводимые элементы:")
    neprevodim_all = nepr(2, stepen)
    for i in range(len(neprevodim_all)):
        print(i+1, neprevodim_all[i])
    print('Введите номер элемента с которым будем работать')
    nepriv_vibran = neprevodim_all[int(input())-1]
#    nepriv_vibran = neprevodim_all[1]
    print('Что будем делать?'
          '\nШифровать - press [y]'
          '\nДешифровать - press [n]')
    if (str(input()).lower()=='y'):
        shifr = affine_shifr(message, alphabet, stepen, key_a, key_b, nepriv_vibran)
        print('Зашифрованное сообщение:\n', shifr)
    else:
        deshifr = affine_shifr_decode(message, alphabet, stepen, key_a, key_b, nepriv_vibran)
        print('Дешифрованное сообщение:\n' , deshifr)
