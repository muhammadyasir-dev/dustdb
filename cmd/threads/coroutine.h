#ifndef COROUTINE_H
#define COROUTINE_H

#include <stdint.h>
#include <stdbool.h>

typedef struct Coroutine {
    uint8_t *stack;        // Pointer to the coroutine stack
    uint8_t *stack_top;    // Top of the stack
    void (*function)(void); // Coroutine function
    bool is_finished;       // Flag to indicate if the coroutine is finished
} Coroutine;

// Function prototypes
Coroutine* create_coroutine(void (*function)(void), size_t stack_size);
void destroy_coroutine(Coroutine *coroutine);
void yield();
void resume(Coroutine *coroutine);
bool is_finished(Coroutine *coroutine);

#endif // COROUTINE_H
