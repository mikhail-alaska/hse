import random
pravilno = 0
vse = 0

temp = 0
counter = 0
for i in range (100000000):
    r = random.randint(1, 4)
    print(r)
    if r == 1:
        temp+=1
    counter +=1
    if counter > 6:
        counter =0
        if temp >2:
            pravilno +=1
        temp = 0 
        vse += 1

print("=======================================")
print(pravilno)
print(vse)
print((pravilno/vse)*100)


