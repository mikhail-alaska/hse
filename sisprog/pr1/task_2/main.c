#include <inttypes.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

void task(int32_t **arr, int32_t **dst, int32_t n);

static void print_matrix(int32_t **arr, int32_t n, const char *title) {
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

  int32_t n;
  printf("Введите размер матрицы > ");
  if (scanf("%d", &n) != 1 || n == 0) {
    fprintf(stderr, "Invalid rand_max\n");
    return 1;
  }


  int32_t min_val, max_val;
  printf("Введите минимальное значение для случайных чисел в матрице > ");
  if (scanf("%d", &min_val) != 1) {
    fprintf(stderr, "Invalid min_val\n");
    return 1;
  }

  printf("Введите максимальное значение для случайных чисел в матрице > ");
  if (scanf("%d", &max_val) != 1) {
    fprintf(stderr, "Invalid rand_max\n");
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

  print_matrix(result, n, "Rotated matrix (clockwise)");

  for (size_t i = 0; i < n; ++i) {
    free(array[i]);
  }
  free(array);

  for (size_t i = 0; i < n; ++i) {
    free(result[i]);
  }
  free(result);
  return 0;
}
