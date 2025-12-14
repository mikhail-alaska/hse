# Названия файлов
TARGET = matrix
SRC = main.c
ASM = pr1_2.asm
OBJ = main.o pr1_2.o

# Компиляторы
CC = gcc
ASM_CC = fasm

# Флаги компиляции и линковки
CFLAGS = -c -O2 -Wall -Wextra -std=c11
LDFLAGS = -no-pie

# Цель по умолчанию
all: $(TARGET)

# Линковка
$(TARGET): $(OBJ)
	$(CC) $(LDFLAGS) -o $@ $^

# Компиляция C-файла
main.o: $(SRC)
	$(CC) $(CFLAGS) $< -o $@

# Компиляция ASM-файла
pr1_2.o: $(ASM)
	$(ASM_CC) $<

# Очистка
clean:
	rm -f $(TARGET) *.o

.PHONY: all clean
