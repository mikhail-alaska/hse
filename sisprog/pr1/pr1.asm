format ELF64

public _start

section '.data' writable

msg db "hello world", 0
len = $-msg

section '.text' executable

_start:
    mov rax, 1
    mov rdi, 1
    mov rsi, msg
    mov rdx, len
    syscall
    call exit
    
exit:
    mov rax, 60
    mov rdi, 0
    syscall
