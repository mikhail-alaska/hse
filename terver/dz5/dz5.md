## **Домашнее задание 5**

#### **Задача 1**
$$\int_{-1}^{1} (1-x^2)dx = 1 \implies $$
$$
\int_{-1}^{1} c\cdot(1-x^2)dx = \int_{-1}^{1}1dx - \int_{-1}^{1}x^2dx =
2 - (\frac{x^3}{3}|_{-1}^{1}) = 2 - (\frac{2}{3}) = \frac{4}{3}\implies
$$
$$c\cdot\frac{4}{3} = 1 \Rightarrow c = \frac{3}{4}$$

**Функция распределения F(x):**
- Для $x<-1$:
$$F(x)=0$$
- Для $x\in[-1,1]$
$$F(x) = \int_{-1}^{x}\frac{3}{4}\cdot(1-t^2)dt\implies F(x)=\frac{3}{4}[t-\frac{t^3}{3}]^{x}_{-1}$$
$$F(x) = \frac{3}{4}\cdot((x-\frac{x^3}{3}) -(-1-\frac{-1}{3})) = \frac{3}{4}(x-\frac{x^3}{3} + \frac{2}{3}) = \\
$$
$$
= \frac{-x^3}{4} +\frac{3x}{4}+\frac{1}{2}
$$
- Для $x>1$:
$$F(x) = 1$$
---
#### **Задача 2**

1. $f(x)\ge 0 \ \ \forall \ x$

    При $x\in[0,5/2]$:
    $$2x-x^3=x(2-x^2)\implies \text{При x = 1, $f(x)>0$}\implies \\
    \text{найдем все решения и посмотрим где $f(x)$ переходит ноль} \implies \\
    x(2-x^2) = 0 \implies x = 0 \\
    2-x^2 = 0 \implies x^2 = 2 \implies x = \pm\sqrt{2}\implies 0<\sqrt2 <2.5\\
    $$
    $f(x)$ не может быть плотностью вероятности

---
#### **Задача 3**
$$Ex=\frac{3}{5}$$

$$\int_0^1(a+bx^2)dx = 1\implies a+b\frac{x^3}{3}|_0^1 = a+\frac{b}{3}=1$$
$$Ex = \int_0^1x(a+bx^2)dx = \int_0^1(ax+bx^3)dx=\frac{3}{5}$$
$$\int_0^1axdx = a\frac{x^2}{2}|_0^1 = \frac{a}{2}\quad\int_0^1bx^3dx=b\frac{x^4}{4}|_0^1=\frac{b}{4}\implies$$
$$\frac{a}{2}+\frac{b}{4}=\frac{3}{5}\quad a+\frac{b}{3}=1 \implies$$
$$\frac{1}{2}-\frac{b}{6}+\frac{b}{4} = 0.5+\frac{b}{12}=0.6\implies b =1.2\implies a = 0.6$$

---
#### **Задача 4**

Сначала найти плотность $f(x)$:
$$f(x)=(0.4x^{1.5}+0.6x)' = 0.4\cdot1.5x^{0.5}+0.6 = 0.6x^{0.5}+0.6$$

$$Ex = \int_0^1xf(x)=0.6\int_0^1x^{1.5}+x =0.6\cdot\frac{x^{2.5}}{2.5}|_0^1+0.6\cdot\frac{x^2}{2}|_0^1=0.6(\frac{2}{5}+\frac{1}{2})=0.54$$

$P(X<9/16|X>1/4) = \frac{P(\frac{1}{4}<X<\frac{9}{16})}{P(X>\frac{1}{4})}$

$$P(X>\frac{1}{4}) = 1-F(1/4) = 1 - (0.4\cdot 0.25^{1.5}+0.6\cdot0.25) = 1 - 0.2 = 0.8$$

$$P(\frac{1}{4}<X<\frac{9}{16}) = F(9/16)-F(1/4) = 0.4\cdot(9/16)^{1.5}+0.6\cdot(9/16)-0.2 =\\=0.50625-0.2=0.30625$$

Ответ:
$$P(X<9/16|X>1/4) = \frac{P(\frac{1}{4}<X<\frac{9}{16})}{P(X>\frac{1}{4})} = \frac{0.30625}{0.8}\approx0.383$$

---

#### **Задача 5**

$$X\sim N(\mu,\sigma^2)\implies \text{Стандартизация }Z=\frac{X-\mu}{\sigma}$$
1. $P\{X>5\}$
    $$P(X>5)=1-P(X\le5)$$
    Стандартизуем:
    $$Z=\frac{5-10}{6} = -5/6\approx-0.8333$$
    $$P(X\le5)=\Phi(-0.833) = 1-\Phi(0.833)=0.2033$$
    $$P(X>5) = 1 - 0.2033 = 0.7967$$
2. $P(4<X<16)$
    $$P(X<16) = \Phi(Z) = \Phi(\frac{16-10}{6})=\Phi(1)\approx0.8413$$
    $$P(X\le4) \implies Z =\frac{4-10}{6}=-1\implies\Phi(-1)=0.1587$$
    $$P(4<X<16) = 0.8413 - 0.1587 = 0.6826$$
3. $P(X<8)$
    $$Z = \frac{8-10}{6} \approx -0.333$$
    $$P(X<8) = \Phi(-0.333) = 1-\Phi(0.333)\approx 1-0.6293=0.3707$$
4. $P(X<20)$
    $$P(X<20)=\Phi(Z)=\Phi(\frac{20-10}{6}) \approx\Phi(1.6666)\approx0.9515$$
5. $P(X>16)$
    $$P(X>16) = 1-P(x\le16)\, \text{ что посчитано ранее}=1-0.8413=0.1587$$

---

#### **Задача 6**

$$Z = \frac{X-\mu}{\sigma}\implies$$
$$P(X>c) = 1- P(X\le c)\implies Z=\frac{c-\mu}{\sigma}\implies$$
$$1-\Phi(Z)=0.1\implies\Phi(Z)=0.9\implies Z\approx1.28\implies$$
$$1.28=\frac{c-12}{2} \implies c = 14.56$$

---

#### **Задача 7**

Плотность распределения X:

$$
f_X(x) = 
\begin{cases}
\frac{1}{2}, & -1 \le x \le 1 \\[6pt]
0, & \text{иначе}
\end{cases}
$$

- Для $Y=2X-2$:
    
    $$Y = 2X-2\implies Y\in[-4, 0]$$
    

$$
f_Y(y) = 
\begin{cases}
\frac{1}{4}, & -4 \le y \le 0 \\[6pt]
0, & \text{иначе}
\end{cases}
$$

- Для $Z = -X$

    Интервал остался тем же $\implies$ плотность осталась та же
$$
f_Z(z) = 
\begin{cases}
\frac{1}{2}, & -1 \le z \le 1 \\[6pt]
0, & \text{иначе}
\end{cases}
$$

---

#### **Задача 8**

1. $P(-0.5\le X \le -0.1)$
    $$P(-0.5\le X \le -0.1)=\Phi(-0.1)-\Phi(-0.5) = 1-0.5398-1+0.6915 = 0.1517$$

2. $P(1\le X \le 2)$
    $$P(1\le X \le 2)=\Phi(2)-\Phi(1) = 0.97725-0.8413=0.13595$$

**Ответ: $P(-0.5\le X\le -0.1)>P(1\le X\le 2)$**

---

#### **Задача 9**

- Мат. ожидание Y = 
$$E(e^X)=\int_0^{\infty}e^x\cdot 3e^{-3x}dx=3\int_0^{\infty}e^{-2x}dx = 3\cdot\frac{1}{2}=\frac{3}{2}=1.5$$

- Дисперсия Y =
$$E[Y^2]=E[e^{2X}]=\int_0^{\infty}e^{2x}\cdot3e^{-3x}dx=3\int_0^{\infty}e^{-x}dx=3\cdot1=3\implies$$
$$\text{Дисперсия $Y$ = }E[Y^2]-(E[Y])^2=3-1.5\cdot1.5=3-2.25=0.75$$

