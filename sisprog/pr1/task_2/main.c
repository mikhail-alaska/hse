#include <inttypes.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

void task(int32_t **arr, int32_t **dst, size_t n);

static void print_matrix(int32_t **arr, size_t n, const char *title) {
  printf("%s (n = %zu):\n", title, n);
  for (size_t i = 0; i < n; ++i) {
    for (size_t j = 0; j < n; ++j) {
      printf("%6" PRId32 " ", arr[i][j]);
    }
    printf("\n");
  }
  printf("\n");
}

int main(void) {
  srand(time(NULL));

  size_t rand_max;
  printf("Введите размер матрицы > ");
  if (scanf("%zu", &rand_max) != 1 || rand_max == 0) {
    fprintf(stderr, "Invalid rand_max\n");
    return 1;
  }

  size_t n = rand_max;

  int32_t min_val, max_val;
  printf("Input min and max values for matrix elements (int32) > ");
  if (scanf("%" SCNd32 " %" SCNd32, &min_val, &max_val) != 2) {
    fprintf(stderr, "Invalid range input\n");
    return 1;
  }
  if (min_val > max_val) {
    int32_t tmp = min_val;
    min_val = max_val;
    max_val = tmp;
  }

  int64_t range = (int64_t)max_val - (int64_t)min_val + 1;
  if (range <= 0) {
    fprintf(stderr, "Invalid range (overflow or empty)\n");
    return 1;
  }

  int32_t **array = calloc(n, sizeof(int32_t *));
  if (!array) {
    fprintf(stderr, "Allocation failed (rows)\n");
    return 1;
  }

  for (size_t i = 0; i < n; ++i) {
    array[i] = calloc(n, sizeof(int32_t));
    if (!array[i]) {
      fprintf(stderr, "Allocation failed (row %zu)\n", i);
      for (size_t k = 0; k < i; ++k) {
        free(array[k]);
      }
      free(array);
      return 1;
    }
  }

  for (size_t i = 0; i < n; ++i) {
    for (size_t j = 0; j < n; ++j) {
      int32_t r = (int32_t)(rand() % range);
      array[i][j] = min_val + r;
    }
  }

  print_matrix(array, n, "Original matrix");
  int32_t **result = calloc(n, sizeof(int32_t *));
  for (size_t i = 0; i < n; ++i) {
    result[i] = calloc(n, sizeof(int32_t));
  }

  task(array, result, n);

  print_matrix(array, n, "Rotated matrix (clockwise)");

  for (size_t i = 0; i < n; ++i) {
    free(array[i]);
  }
  free(array);

  return 0;
}
