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
    mov eax, [result]

    cmp eax, 10
    jl  one_digit

    cmp eax, 20
    je twenty

    mov byte [outbuf], '1'

    sub al, 10
    add al, '0'
    mov [outbuf+1], al

    mov rax, 1
    mov rdi, 1
    mov rsi, outbuf
    mov rdx, 2
    syscall
    ret

twenty:
    mov byte [outbuf], '2'
    mov byte [outbuf+1], '0'
    mov rax, 1
    mov rdi, 1
    mov rsi, outbuf
    mov rdx, 2
    syscall
    ret


one_digit:
    add al, '0'
    mov [outbuf], al

    mov rax, 1
    mov rdi, 1
    mov rsi, outbuf
    mov rdx, 1
    syscall
    ret                    ; <-- возвращаемся в _start
