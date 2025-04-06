#ifndef DYNAMIC_ARRAY_H
#define DYNAMIC_ARRAY_H

#include <stdbool.h>

typedef struct {
    char *data;      // Pointer to the character array
    size_t size;     // Current size of the array
    size_t capacity; // Current capacity of the array
} DynamicArray;

// Function prototypes
DynamicArray* createDynamicArray(size_t initialCapacity);
void appendCharacter(DynamicArray *array, char c);
void deleteCharacter(DynamicArray *array, size_t index);
void freeDynamicArray(DynamicArray *array);
void printDynamicArray(DynamicArray *array);

#endif // DYNAMIC_ARRAY_H
