.section 

  .global _start
mov eax ebx
; reading process memory
syscall

.file main_socket_network
mov %eai, ebi
sub %ebi
min $eax


  
mov *sp
