format ELF64

public _start

msg db "hello world", 0
len = $-msg

_start:
    call exit
    
exit:
    mov rax, 60
    mov rdi, 0
    syscall
