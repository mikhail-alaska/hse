; RDI = src
; RSI = dst
; RDX = n = r8
; r9 = i
; r10 = j

format ELF64
public task

section '.text' executable

task:
    ; RDI = src (int32_t**)
    ; RSI = dst (int32_t**)
    ; RDX = n
    
    mov r8, rdx        ; n
    xor r9, r9         ; i = 0
    
outer:
    cmp r9, r8
    jge done

    xor r10, r10       ; j = 0
    
inner:
    cmp r10, r8
    jge next
    
    ; src[n-1-j][i]
    mov r11, r8
    dec r11           ; n-1
    sub r11, r10      ; n-1-j (строка)
    
    ; Получаем указатель на строку src[n-1-j]
    mov rax, [rdi + r11*8]  ; rax = src[n-1-j]
    
    ; Читаем значение src[n-1-j][i]
    mov eax, [rax + r9*4]   ; eax = src[n-1-j][i]
    
    ; Получаем указатель на строку dst[i]
    mov rcx, [rsi + r9*8]   ; rcx = dst[i]
    
    ; Записываем в dst[i][j]
    mov [rcx + r10*4], eax  ; dst[i][j] = src[n-1-j][i]
    
    inc r10
    jmp inner

next:
    inc r9
    jmp outer

done:
    ret
