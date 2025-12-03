format ELF64

public _start

section '.data' writable

array dw 9, 9, 9, 3, 0, 12, 9, 3, 3, 3
len = 10
result dd 0
outbuf  rb 32
newline db 10

section '.text' executable

uint32_to_str:
    mov rbx, 10
    lea rdi, [outbuf + 31]
    mov byte [rdi], 0
    dec rdi
    mov rcx, 1

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
    

exit:
    mov rax, 60
    mov edi, [result]
    syscall

