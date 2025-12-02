; RDI = src
; RSI = dst
; EDX = n = r8d
; r9d = i
; r10d = j

format ELF64

public task

section '.text' executable

task:
    mov r8d, edx
    xor r9d, r9d

outer_loop:
    cmp r9d, r8d
    jge done

    xor r10d, r10d

inner_loop:
    cmp r10d, r8d
    jge next_row

    mov r11d, r8d
    dec r11d
    sub r11d, r10d

    mov eax, r11d
    imul eax, r8d
    add eax, r9d

    mov edx, eax
    shl rdx, 2
    mov eax, [rdi + rdx]

    mov edx, r9d
    imul edx, r8d
    add edx, r10d

    mov ecx, edx
    imul rcx, 4
    mov [rsi + rcx], eax

    inc r10d
    jmp inner_loop

next_row:
    inc r9d
    jmp outer_loop

done:
    ret
