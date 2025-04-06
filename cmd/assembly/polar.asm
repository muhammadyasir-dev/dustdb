section .data
    ; Socket options
    sockfd db 0
    addr_in db 0
    addr_len db 16
    buffer db 1024 dup(0)

section .text
    global _start

_start:
    ; Create socket
    mov rax, 41          ; syscall: socket
    mov rdi, 2           ; AF_INET
    mov rsi, 1           ; SOCK_STREAM
    mov rdx, 0           ; protocol
    syscall
    mov [sockfd], rax    ; save socket file descriptor

    ; Bind socket
    mov rax, 49          ; syscall: bind
    mov rdi, [sockfd]    ; socket file descriptor
    lea rsi, [addr_in]   ; pointer to sockaddr_in
    mov rdx, [addr_len]  ; length of sockaddr_in
    syscall

    ; Listen on socket
    mov rax, 50          ; syscall: listen
    mov rdi, [sockfd]    ; socket file descriptor
    mov rsi, 5           ; backlog
    syscall

poll:
    ; Accept incoming connection
    mov rax, 43          ; syscall: accept
    mov rdi, [sockfd]    ; socket file descriptor
    lea rsi, [addr_in]   ; pointer to sockaddr_in
    mov rdx, [addr_len]  ; length of sockaddr_in
    syscall

    ; Now you can read from the accepted socket (not implemented here)
    ; For example, you could use recv syscall to read data

    ; Loop back to poll for new connections
    jmp poll

    ; Exit
    mov rax, 60          ; syscall: exit
    xor rdi, rdi         ; status 0
    syscall
