; RDI = src
; RSI = dst
; EDX = n = r8d
; r9d = i

format ELF64

public rotate_clockwise

section '.text' executable

rotate_clockwise:
    mov     r8d, edx
    xor     r9d, r9d

outer_loop:
    cmp     r9d, r8d        ; if (i >= n) -> done
    jge     .done

    xor     r10d, r10d      ; r10d = j = 0

inner_loop:
    cmp     r10d, r8d       ; if (j >= n) -> next row
    jge     .next_row

    ; src_row = n - 1 - j
    mov     r11d, r8d       ; r11d = n
    dec     r11d            ; r11d = n - 1
    sub     r11d, r10d      ; r11d = n - 1 - j

    ; index_src = src_row * n + i
    mov     eax, r11d       ; eax = src_row
    imul    eax, r8d        ; eax = src_row * n
    add     eax, r9d        ; eax = src_row * n + i

    ; offset_src_bytes = index_src * 4
    mov     edx, eax
    shl     rdx, 2          ; rdx = offset in bytes
    mov     eax, [rdi + rdx] ; eax = src[index_src]

    ; index_dst = i * n + j
    mov     edx, r9d        ; edx = i
    imul    edx, r8d        ; edx = i * n
    add     edx, r10d       ; edx = i * n + j

    mov     ecx, edx
    shl     rcx, 2          ; rcx = offset_dst_bytes
    mov     [rsi + rcx], eax ; dst[index_dst] = value

    inc     r10d            ; j++
    jmp     inner_loop

.next_row:
    inc     r9d             ; i++
    jmp     outer_loop

.done:
    ret
