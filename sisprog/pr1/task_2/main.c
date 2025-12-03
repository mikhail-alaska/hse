#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <inttypes.h>
#include <time.h>
#include <string.h>

#define CELL_WIDTH 6

void task(char ***src, char ***dst, int32_t n);

void print_matrix(char ***arr, int32_t n) {
  for (int32_t i = 0; i < n; ++i) {
    for (int32_t j = 0; j < n; ++j) {
      printf("%s ", arr[i][j]);
    }
    printf("\n");
  }
  printf("\n");
}

int main(void) {
  srand(time(NULL));

  uint32_t n;
  printf("Введите размер матрицы > ");
  if (scanf("%u", &n) != 1 || n == 0) {
    fprintf(stderr, "Invalid matrix size\n");
    return 1;
  }

  int32_t min_val, max_val;
  printf("Введите минимальное значение > ");
  if (scanf("%d", &min_val) != 1) {
    fprintf(stderr, "Invalid min_val\n");
    return 1;
  }

  printf("Введите максимальное значение > ");
  if (scanf("%d", &max_val) != 1) {
    fprintf(stderr, "Invalid max_val\n");
    return 1;
  }

  if (min_val > max_val) {
    int32_t tmp = min_val;
    min_val = max_val;
    max_val = tmp;
  }

  int64_t range = (int64_t)max_val - (int64_t)min_val + 1;
  if (range <= 0) {
    fprintf(stderr, "Invalid range\n");
    return 1;
  }

  char ***array = calloc(n, sizeof(char **));
  for (uint32_t i = 0; i < n; ++i) {
    array[i] = calloc(n, sizeof(char *));
    for (uint32_t j = 0; j < n; ++j) {
      int32_t val = min_val + rand() % range;
      char buffer[20];
      sprintf(buffer, "%d", val);
      array[i][j] = calloc(strlen(buffer), sizeof(char));
      snprintf(array[i][j], strlen(buffer), "%*d", CELL_WIDTH, val);  // форматируем число
    }
  }

  print_matrix(array, n);

  // ⬇️ Результат — тоже как строки
  char ***result = calloc(n, sizeof(char **));
  for (uint32_t i = 0; i < n; ++i) {
    result[i] = calloc(n, sizeof(char *));
    for (uint32_t j = 0; j < n; ++j) {
      result[i][j] = calloc(CELL_WIDTH + 1, sizeof(char));
    }
  }

  task(array, result, n);

  print_matrix(result, n);

  // ⬇️ Освобождение памяти
  for (uint32_t i = 0; i < n; ++i) {
    for (uint32_t j = 0; j < n; ++j) {
      free(array[i][j]);
      free(result[i][j]);
    }
    free(array[i]);
    free(result[i]);
  }
  free(array);
  free(result);

  return 0;
}
