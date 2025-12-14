TARGET := matrix
SRC := main.c
ASM := pr1_2.asm
OBJ := main.o pr1_2.o

CC := gcc
ASM_CC := fasm

CFLAGS := -c -O2 -Wall -Wextra -std=c11
LDFLAGS := -no-pie

all: $(TARGET)

$(TARGET): $(OBJ)
	$(CC) $(LDFLAGS) -o $@ $^

main.o: $(SRC)
	$(CC) $(CFLAGS) $< -o $@

pr1_2.o: $(ASM)
	$(ASM_CC) $<

# Очистка
clean:
	rm -f $(TARGET) *.o

.PHONY: all clean
