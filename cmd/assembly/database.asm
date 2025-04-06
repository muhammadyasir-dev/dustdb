section .data
    db_host db "localhost", 0
    db_user db "your_username", 0
    db_pass db "your_password", 0
    db_name db "your_database", 0

    mysql_init db 0
    mysql_real_connect db 0
    mysql_close db 0

section .bss
    conn resq 1

section .text
    extern mysql_init
    extern mysql_real_connect
    extern mysql_close
    extern mysql_library_init
    extern mysql_library_end
    extern printf
    extern exit
    extern strlen

    global _start

_start:
    ; Initialize MySQL library
    push 0
    call mysql_library_init
    add esp, 4

    ; Initialize MySQL connection
    push 0
    call mysql_init
    add esp, 4

    ; Store the connection in conn
    mov [conn], eax

    ; Connect to the database
    push db_name
    push db_pass
    push db_user
    push db_host
    push [conn]
    call mysql_real_connect
    add esp, 20

    ; Check if connection was successful
    test eax, eax
    jz .connection_failed

    ; Connection successful
    ; Here you would typically execute queries

    ; Close the connection
    push [conn]
    call mysql_close
    add esp, 4

    ; Clean up MySQL library
    call mysql_library_end

    ; Exit program
    mov eax, 0
    call exit

.connection_failed:
    ; Handle connection failure
    push msg_connection_failed
    call printf
    add esp, 4

    ; Clean up MySQL library
    call mysql_library_end

    ; Exit program
    mov eax, 1
    call exit

section .data
msg_connection_failed db "Connection to database failed!", 10, 0
