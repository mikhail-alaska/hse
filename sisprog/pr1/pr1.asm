format ELF64

public _start

section '.data' writable

array dw 5, 7, 10, 3, 0, 12, 8, 1, 2, 4
len = $-array
result dd 0

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

    mov rax, 1
    mov rdi, 1
    mov rsi, msg
    mov rdx, len
    syscall
    
exit:

    mov rax, 60
    mov rdi, 0
    syscall


