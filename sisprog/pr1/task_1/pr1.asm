format ELF64

public _start

section '.data' writable

array dw 9, 9, 9, 3, 0, 12, 9, 3, 3, 3
len = 10
result dd 0
outbuf  rb 32

section '.text' executable

_start:
    mov rdi, array
    mov ecx, len
    xor r8d, r8d
    mov bx, 3

loop_start:
    mov ax, [rdi]
    xor dx, dx
    div bx

    movzx eax, dx
    add r8d, eax

    add rdi, 2
    loop loop_start

    mov [result], r8d
    
    call print_result

exit:
    mov rax, 60
    xor edi, edi
    syscall


print_result:
    mov eax, [result]      ; <-- читаем число из result в EAX

    cmp eax, 10            ; <-- сравниваем с 10
    jl  .one_digit         ; <-- если меньше 10 — одна цифра

    ; ---- двухзначные числа: 10..20 ----
    mov byte [outbuf], '1' ; <-- первая цифра — '1'

    sub al, 10             ; <-- al = (число - 10) → диапазон 0..10
    add al, '0'            ; <-- превращаем остаток в символ '0'..'9'
    mov [outbuf+1], al     ; <-- записываем вторую цифру в буфер

    mov rax, 1             ; <-- sys_write
    mov rdi, 1             ; <-- дескриптор: stdout
    mov rsi, outbuf        ; <-- адрес буфера
    mov rdx, 2             ; <-- длина: 2 символа
    syscall
    ret                    ; <-- возвращаемся в _start

.one_digit:
    ; ---- однозначные числа: 0..9 ----
    add al, '0'            ; <-- превращаем число 0..9 в символ '0'..'9'
    mov [outbuf], al       ; <-- кладём символ в буфер

    mov rax, 1             ; <-- sys_write
    mov rdi, 1             ; <-- stdout
    mov rsi, outbuf        ; <-- адрес буфера
    mov rdx, 1             ; <-- длина: 1 байт
    syscall
    ret                    ; <-- возвращаемся в _start
