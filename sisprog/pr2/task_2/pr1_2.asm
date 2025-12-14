; RDI = src
; RSI = dst
; RDX = n = r8
; r9 = i
; r10 = j

format ELF64
public task

section '.text' executable

task:
    mov r8, rdx
    xor r9, r9
    
outer:
    cmp r9, r8
    jge done

    xor r10, r10
    
inner:
    cmp r10, r8
    jge next
    
    mov r11, r8
    dec r11
    sub r11, r10 ; n-1-j
    
    mov rax, [rdi + r11*8]  ; rax = src[n-1-j]
    
    mov eax, [rax + r9*4]   ; eax = src[n-1-j][i]
    
    mov rcx, [rsi + r9*8]   ; rcx = dst[i]
    
    mov [rcx + r10*4], eax  ; dst[i][j] = src[n-1-j][i]
    
    inc r10
    jmp inner

next:
    inc r9
    jmp outer

done:
    ret
