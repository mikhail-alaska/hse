global task

section .text
task:
    push rbp
    mov rbp, rsp

    sub rsp, 0x40

    mov qword [rbp - 0x40], 0

    mov [rbp - 0x38], rdi
    mov [rbp - 0x30], rsi
    mov [rbp - 0x28], rdx

    mov qword [rbp - 0x20], 1
    
    mov qword [rbp - 0x18], '0'
    mov [rbp - 0x17], 'x'
    mov [rbp - 0x6], 0
    mov [rbp - 0x5], 0

    mov rdi, [rbp - 0x38]
    mov rcx, [rbp - 0x30]

.column_loop:
    push rcx
    mov rcx, [rbp - 0x28]
    mov rsi, [rdi]
.line_loop:
    mov rax, [rbp - 0x20]
    mov [rsi], rax

    mov rax, [rsi]
    mul qword [rbp - 0x20]
    mov [rsi], rax

    mov rax, [rbp - 0x40]
    test rax, rax
    jz .mul_skip
    
    mov rax, [rsi]
    mul qword [rbp - 0x20]
    mov [rsi], rax

.mul_skip:
    inc [rbp - 0x20]
    add rsi, 8
    xor [rbp - 0x40], 1
    
    loop .line_loop

    add rdi, 8
    
    pop rcx
    loop .column_loop
        
.print_result:
    mov rax, 2
    xor rdi, rdi
    xor rsi, rsi
    xor rdx, rdx
    syscall

    mov rdi, [rbp - 0x38]
    mov rcx, [rbp - 0x30]

.print_column_loop:
    push rcx
    mov rcx, [rbp - 0x28]
    mov rsi, [rdi]

.print_line_loop:
    push rdi
    push rsi
    push rcx

    lea rdi, [rbp - 0x18 + 2]
    mov rsi, [rsi]

    call to_hex
    
    pop rcx
    pop rsi
    pop rdi

    cmp rcx, 1
    je .newline
.space:
    mov [rbp - 0x6], 32
    jmp .write_out
.newline:
    mov [rbp - 0x6], 10
.write_out:
    push rdi
    push rsi
    push rcx
    
    mov rax, 1
    xor rdi, rdi
    lea rsi, [rbp - 0x18]
    mov rdx, 19
    syscall
    
    pop rcx
    pop rsi
    pop rdi

    add rsi, 8

    loop .print_line_loop
    
    add rdi, 8

    pop rcx
    loop .print_column_loop

    mov rax, 3
    xor rdi, rdi
    syscall

    add rsp, 0x40

    pop rbp
    ret

to_hex:
    mov rbx, rsi
    mov rcx, 16

.hex_loop:
    mov rax, rbx
    shr rax, 60

    cmp rax, 9
    jbe .digit
    add rax, 55
    jmp .store

.digit:
    add rax, 48

.store:
    mov [rdi], al
    inc rdi

    shl rbx, 4

    loop .hex_loop

    ret

