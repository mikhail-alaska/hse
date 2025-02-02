import matplotlib.pyplot as plt

N = [2682, 2230, 2042, 1624, 966, 482, 424]
lambda_values = [623.4, 579, 546, 491.6, 435.8, 407.8, 404.7]

from scipy.interpolate import interp1d
interp_func = interp1d(N, lambda_values, kind='linear', fill_value='extrapolate')

N_check = [2564, 1576, 1332, 1268]
hh = 1239.7
counter = 3
for n in N_check:
    lambda_val = interp_func(n)
    delta_E = hh / lambda_val
    ridberg = 1 / ((0.25 - (1 / (counter ** 2))) * (lambda_val))
    print(f"N = {n}, λ = {lambda_val}, ΔE = {delta_E:.2e}, Ридберг = {ridberg*100}")
    counter+=1
plt.figure(figsize=(8, 6))
plt.plot(N, lambda_values, marker='o', linestyle='-', color='b', label='График \u03bb(N)')
plt.scatter(N_check, interp_func(N_check), color='r', label='Заданные точки', zorder=5)

plt.title('График зависимости λ от N')
plt.xlabel('N')
plt.ylabel('λ')
plt.grid(True)
plt.legend()

plt.show()
