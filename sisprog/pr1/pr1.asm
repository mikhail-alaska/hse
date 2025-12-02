format ELF64

public _start

_start:
    call exit
    
exit:
    mov rax, 60
    mov rbx, 55
    syscall
