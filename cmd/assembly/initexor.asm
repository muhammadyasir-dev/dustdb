section .data
    mem_size equ 1024 * 1024      ; Total memory size (1 MB)
    mem_start resb mem_size       ; Reserve 1 MB of memory for allocation
    free_list db 0                 ; Pointer to the head of the free list

section .bss
    ; Structure for a free block
    struct FreeBlock:
        size resd 1               ; Size of the block
        next resd 1               ; Pointer to the next free block

section .text
    global _start

_start:
    ; Initialize the free list
    mov dword [mem_start], mem_size ; Set the size of the first free block
    mov dword [mem_start + 4], 0    ; Next pointer is NULL
    mov dword [free_list], mem_start ; Set the free list to point to the start

    ; Example allocations
    mov eax, 128                    ; Request 128 bytes
    call allocate_memory
    ; Store the pointer to the allocated memory (for demonstration)
    mov ebx, eax                    ; Store pointer in ebx

    mov eax, 256                    ; Request 256 bytes
    call allocate_memory
    ; Store the pointer to the allocated memory (for demonstration)
    mov ecx, eax                    ; Store pointer in ecx

    ; Free the first allocated block
    call free_memory
    mov eax, ebx                    ; Pointer to the first block
    call free_memory

    ; Exit the program
    mov eax, 1                      ; sys_exit
    xor ebx, ebx                    ; Exit code 0
    int 0x80

; Function to allocate memory
; Input: EAX = size of memory to allocate
; Output: EAX = pointer to allocated memory or 0 if failed
allocate_memory:
    push ebp
    mov ebp, esp
    push eax                        ; Save requested size

    ; Align size to 4 bytes
    add eax, 3
    shr eax, 2
    shl eax, 2

    ; Traverse the free list to find a suitable block
    mov esi, [free_list]           ; Start from the head of the free list
    mov edi, 0                     ; Previous block pointer

find_block:
    cmp esi, 0                     ; Check if we reached the end of the list
    je no_block_found              ; No suitable block found
    mov ebx, [esi]                 ; Get the size of the current block
    cmp ebx, eax                   ; Is the block big enough?
    jl next_block                  ; No, check the next block

    ; Found a suitable block
    ; If the block is larger than needed, split it
    sub ebx, eax                   ; Remaining size after allocation
    cmp ebx, 8                     ; Minimum size for a free block
    jl allocate_full_block         ; Not enough space to split

    ; Split the block
    mov [esi], ebx                 ; Update the size of the current block
    add esi, ebx                   ; Move to the new block
    mov [esi], eax                 ; Set the size of the allocated block
    mov dword [esi + 4], 0         ; Next pointer is NULL
    jmp done

allocate_full_block:
    ; Allocate the entire block
    mov dword [esi + 4], 0         ; Set next pointer to NULL
    jmp done

next_block:
    mov edi, esi                   ; Update previous block pointer
    mov esi, [esi + 4]            ; Move to the next block
    jmp find_block

no_block_found:
    ; No suitable block found
    xor eax, eax                   ; Return 0 (NULL)

done:
    pop eax                        ; Restore requested size
    pop ebp
    ret

; Function to free memory
; Input: EAX = pointer to memory to free
free_memory:
    push ebp
    mov ebp, esp
    push eax                        ; Save pointer to free

    ; Get the size of the block to free
    mov ebx, [eax]                 ; Get the size of the block
    mov dword [eax + 4], [free_list] ;
