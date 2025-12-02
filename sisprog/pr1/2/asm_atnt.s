    .globl task

    .text
task:
    pushq   %rbp
    movq    %rsp, %rbp

    subq    $0x40, %rsp

    movq    $0, -0x40(%rbp)

    movq    %rdi, -0x38(%rbp)
    movq    %rsi, -0x30(%rbp)
    movq    %rdx, -0x28(%rbp)

    movq    $1, -0x20(%rbp)
    
    movq    $'0', -0x18(%rbp)
    movb    $'x', -0x17(%rbp)
    movb    $0, -0x6(%rbp)
    movb    $0, -0x5(%rbp)

    movq    -0x38(%rbp), %rdi
    movq    -0x30(%rbp), %rcx

.column_loop:
    pushq   %rcx
    movq    -0x28(%rbp), %rcx
    movq    (%rdi), %rsi

.line_loop:
    movq    -0x20(%rbp), %rax
    movq    %rax, (%rsi)

    movq    (%rsi), %rax
    mulq    -0x20(%rbp)
    movq    %rax, (%rsi)

    movq    -0x40(%rbp), %rax
    testq   %rax, %rax
    jz      .mul_skip
    
    movq    (%rsi), %rax
    mulq    -0x20(%rbp)
    movq    %rax, (%rsi)

.mul_skip:
    incq    -0x20(%rbp)
    addq    $8, %rsi
    xorq    $1, -0x40(%rbp)
    
    loop    .line_loop

    addq    $8, %rdi
    
    popq    %rcx
    loop    .column_loop
        
.print_result:
    movq    $2, %rax
    xorq    %rdi, %rdi
    xorq    %rsi, %rsi
    xorq    %rdx, %rdx
    syscall

    movq    -0x38(%rbp), %rdi
    movq    -0x30(%rbp), %rcx

.print_column_loop:
    pushq   %rcx
    movq    -0x28(%rbp), %rcx
    movq    (%rdi), %rsi

.print_line_loop:
    pushq   %rdi
    pushq   %rsi
    pushq   %rcx

    leaq    -0x16(%rbp), %rdi       # [rbp - 0x18 + 2]
    movq    (%rsi), %rsi

    call    to_hex
    
    popq    %rcx
    popq    %rsi
    popq    %rdi

    cmpq    $1, %rcx
    je      .newline
.space:
    movb    $32, -0x6(%rbp)
    jmp     .write_out
.newline:
    movb    $10, -0x6(%rbp)
.write_out:
    pushq   %rdi
    pushq   %rsi
    pushq   %rcx
    
    movq    $1, %rax
    xorq    %rdi, %rdi
    leaq    -0x18(%rbp), %rsi
    movq    $19, %rdx
    syscall
    
    popq    %rcx
    popq    %rsi
    popq    %rdi

    addq    $8, %rsi

    loop    .print_line_loop
    
    addq    $8, %rdi

    popq    %rcx
    loop    .print_column_loop

    movq    $3, %rax
    xorq    %rdi, %rdi
    syscall

    addq    $0x40, %rsp

    popq    %rbp
    ret

to_hex:
    movq    %rsi, %rbx
    movq    $16, %rcx

.hex_loop:
    movq    %rbx, %rax
    shrq    $60, %rax

    cmpq    $9, %rax
    jbe     .digit
    addq    $55, %rax
    jmp     .store

.digit:
    addq    $48, %rax

.store:
    movb    %al, (%rdi)
    incq    %rdi

    shlq    $4, %rbx

    loop    .hex_loop

    ret
