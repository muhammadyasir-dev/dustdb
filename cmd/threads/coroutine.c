#include <stdio.h>
#include "Coroutine.h"

void coroutine1() {
    printf("Coroutine 1: Step 1\n");
    yield();
    printf("Coroutine 1: Step 2\n");
}

void coroutine2() {
    printf("Coroutine 2: Step 1\n");
    yield();
    printf("Coroutine 2: Step 2\n");
}

int main() {
    Coroutine *c1 = create_coroutine(coroutine1, 1024);
    Coroutine *c2 = create_coroutine(coroutine2, 1024);

    while (!is_finished(c1) || !is_finished(c2)) {
        if (!is_finished(c1)) {
            resume(c1);
        }
        if (!is_finished(c2)) {
            resume(c2);
        }
    }

    destroy_coroutine(c1);
    destroy_coroutine(c2);
    return 0;
}
