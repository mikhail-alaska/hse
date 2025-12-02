    .section .data
value:
    .long 0x1337beef
zero_count:
    .long 0
hex_prefix:
    .ascii "0x"
hex_value:
    .space 16
    .byte 10, 0

    .section .text
    .globl _start

_start:
    movl value(%rip), %eax
    movq $32, %rcx
    xorl %ebx, %ebx

count_loop:
    shll $1, %eax
    jc bit_is_one
    incl %ebx

bit_is_one:
    loop count_loop

    movl %ebx, zero_count(%rip)

to_hex:
    leaq hex_value(%rip), %rdi
    movq $8, %rcx

hex_loop:
    movl %ebx, %eax
    shrl $28, %eax

    cmpl $9, %eax
    jbe digit
    addl $55, %eax
    jmp store

digit:
    addl $48, %eax

store:
    movb %al, (%rdi)
    incq %rdi

    shll $4, %ebx

    loop hex_loop

print_result:
    movq $2, %rax
    xorq %rdi, %rdi
    xorq %rsi, %rsi
    xorq %rdx, %rdx
    syscall

    movq $1, %rax
    xorq %rdi, %rdi
    leaq hex_prefix(%rip), %rsi
    movq $19, %rdx
    syscall

    movq $3, %rax
    xorq %rdi, %rdi
    syscall

exit:
    movq $60, %rax
    xorq %rdi, %rdi
    syscall

