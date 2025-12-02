format ELF64

public _start

msg db "hello world"

_start:
    call exit
    
exit:
    mov rax, 60
    mov rdi, 57
    syscall
