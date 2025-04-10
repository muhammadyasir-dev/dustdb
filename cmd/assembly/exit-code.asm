section .data
    message db 'Hello, World!', 0xA  ; The message to print, followed by a newline
    message_length equ $ - message     ; Calculate the length of the message

section .text
    global _start                      ; Entry point for the program

_start:
    ; Prepare for the write system call
    mov eax, 4                        ; System call number for sys_write (4)
    mov ebx, 1                        ; File descriptor 1 (stdout)
    mov ecx, message                  ; Pointer to the message
    mov edx, message_length           ; Length of the message

    ; Make the system call
    int 0x80                          ; Call the kernel

    ; Exit the program
    mov eax, 1                        ; System call number for sys_exit (1)
    xor ebx, ebx                      ; Exit code 0
    int 0x80                          ; Call the kernel
