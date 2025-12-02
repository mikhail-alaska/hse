#include <stdlib.h>
#include <stdint.h>
#include <time.h>
#include <stdio.h>

void task(uint64_t** arr, size_t size1, size_t size2);

int main() {
    srand(time(NULL));

    size_t rand_max;
    printf("Input rand_max for array sizes > ");
    scanf("%lu", &rand_max);

    size_t size1 = (rand() % rand_max) + 1;
    size_t size2 = (rand() % rand_max) + 1;
    
    uint64_t** array = calloc(sizeof(uint64_t*), size1);   
    for(size_t i = 0; i < size1; ++i) {
        array[i] = calloc(sizeof(uint64_t), size2);
    }

    task(array, size1, size2);

    for(size_t i = 0; i < size1; ++i) {
        free(array[i]);
    }
    free(array);
    return 0;
}
