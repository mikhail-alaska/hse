section .data 
    value      dd 0x1337beef
    zero_count dd 0
    hex_prefix db '0x'
    hex_value  db 16 dup(0), 10, 0

section .text
    global _start

_start:
    mov eax, [rel value]
    mov rcx, 32
    xor ebx, ebx

count_loop:
    shl eax, 1
    jc bit_is_one
    inc ebx

bit_is_one:
    loop count_loop

    mov [rel zero_count], ebx

to_hex:
    lea rdi, [rel hex_value]
    mov rcx, 8

hex_loop:
    mov eax, ebx
    shr eax, 28

    cmp eax, 9
    jbe digit
    add eax, 55
    jmp store

digit:
    add eax, 48

store:
    mov [rdi], rax
    inc rdi

    shl ebx, 4

    loop hex_loop
    
print_result:
    mov rax, 2
    xor rdi, rdi
    xor rsi, rsi
    xor rdx, rdx
    syscall 

    mov rax, 1
    xor rdi, rdi
    lea rsi, [rel hex_prefix]
    mov rdx, 19
    syscall

    mov rax, 3
    xor rdi, rdi
    syscall

exit:
    mov rax, 60
    xor rdi, rdi
    syscall
