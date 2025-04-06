#include <stdio.h>
#include "dynamicarray.h"

int main() {
    DynamicArray *array = createDynamicArray(2); // Start with a capacity of 2

    // Append characters
    appendCharacter(array, 'H');
    appendCharacter(array, 'e');
    appendCharacter(array, 'l');
    appendCharacter(array, 'l');
    appendCharacter(array, 'o');

    printf("Dynamic Array after appending characters: ");
    printDynamicArray(array);

    // Delete a character at index 1 (removing 'e')
    deleteCharacter(array, 1);
    printf("Dynamic Array after deleting character at index 1: ");
    printDynamicArray(array);

    // Free the dynamic array
    freeDynamicArray(array);
    return 0;
}
