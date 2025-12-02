format ELF64

public _start

msg db "hello world", 0
len = $-msg

_start:
    mov rax, 4
    mov rbx, 1
    mov rcx, msg
    mov rdx, len
    syscall
    call exit
    
exit:
    mov rax, 60
    mov rdi, 0
    syscall
