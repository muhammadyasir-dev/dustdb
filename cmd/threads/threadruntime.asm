section .data
    thread_count db 5
    thread_ids dq 0, 0, 0, 0, 0
    message db "Hello from thread %d", 10, 0
    mutex_name db "mutex", 0

section .bss
    mutex resb 8

section .text
    extern printf
    extern pthread_create
    extern pthread_join
    extern pthread_exit
    extern pthread_mutex_init
    extern pthread_mutex_lock
    extern pthread_mutex_unlock
    extern pthread_mutex_destroy
    global _start

_start:
    ; Initialize mutex
    mov rdi, mutex
    xor rsi, rsi
    call pthread_mutex_init

    ; Create threads
    mov rdi, thread_count
    call create_threads

    ; Wait for threads to finish
    mov rdi, thread_count
    call join_threads

    ; Destroy mutex
    mov rdi, mutex
    call pthread_mutex_destroy

    ; Exit
    mov rax, 60          ; syscall: exit
    xor rdi, rdi         ; status 0
    syscall

create_threads:
    ; rdi: number of threads
    mov rbx, 0           ; thread index
create_thread_loop:
    cmp rbx, rdi
    jge .done_creating
    ; Create thread
    mov rsi, rbx         ; thread index as argument
    lea rdx, [thread_ids + rbx * 8] ; address to store thread ID
    mov rax, 0          ; thread attributes (NULL)
    call pthread_create
    inc rbx
    jmp create_thread_loop
.done_creating:
    ret

join_threads:
    ; rdi: number of threads
    mov rbx, 0           ; thread index
join_thread_loop:
    cmp rbx, rdi
    jge .done_joining
    lea rsi, [thread_ids + rbx * 8] ; thread ID
    call pthread_join
    inc rbx
    jmp join_thread_loop
.done_joining:
    ret

thread_function:
    ; rdi: thread index
    push rdi
    mov rdi, message
    mov rsi, rdi
    call printf
    pop rdi
    ; Lock mutex
    mov rdi, mutex
    call pthread_mutex_lock
    ; Critical section (simulated work)
    ; Unlock mutex
    mov rdi, mutex
    call pthread_mutex_unlock
    ; Exit thread
    call pthread_exit
