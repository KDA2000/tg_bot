package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	var choice int
	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1. Вычислить возведение в степень по модулю")
		fmt.Println("2. Проверить число на простоту (тест Ферма)")
		fmt.Println("3. Вычислить НОД и коэффициенты (ввод с клавиатуры)")
		fmt.Println("4. Сгенерировать случайные a и b и вычислить НОД")
		fmt.Println("5. Сгенерировать простые a и b и вычислить НОД")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			var base, exp, mod int64
			fmt.Print("Введите число: ")
			fmt.Scan(&base)
			fmt.Print("Введите степень: ")
			fmt.Scan(&exp)
			fmt.Print("Введите модуль: ")
			fmt.Scan(&mod)
			result := step(base, exp, mod)
			fmt.Printf("%d в степени %d по модулю %d = %d\n", base, exp, mod, result)

		case 2:
			var n int64
			fmt.Print("Введите число для проверки на простоту: ")
			fmt.Scan(&n)
			isPrime := fermatTest(n)
			if isPrime {
				fmt.Printf("%d - простое число\n", n)
			} else {
				fmt.Printf("%d - составное число (видимо не простое...)\n", n)
			}

		case 3:
			var a, b int64
			fmt.Print("Введите a: ")
			fmt.Scan(&a)
			fmt.Print("Введите b: ")
			fmt.Scan(&b)
			gcd, u, v := extendedGCD(a, b)
			fmt.Printf("НОД(%d, %d) = %d, x = %d, y = %d\n", a, b, gcd, u, v)

		case 4:
			a, b := generateRandomNumbers()
			gcd, u, v := extendedGCD(a, b)
			fmt.Printf("Случайно сгенерированные a и b: %d, %d\n", a, b)
			fmt.Printf("НОД(%d, %d) = %d, x = %d, y = %d\n", a, b, gcd, u, v)

		case 5:
			a, b := generateProstoeNumbers()
			gcd, u, v := extendedGCD(a, b)
			fmt.Printf("Случайно сгенерированные простые a и b: %d, %d\n", a, b)
			fmt.Printf("НОД(%d, %d) = %d, x = %d, y = %d\n", a, b, gcd, u, v)

		default:
			fmt.Println("Неверный выбор")
		}
	}
}

func step(base, exp, mod int64) int64 {

	bigBase := big.NewInt(base)
	bigMod := big.NewInt(mod)
	result := big.NewInt(1)

	for exp > 0 {
		if exp%2 == 1 {
			result.Mul(result, bigBase)
			result.Mod(result, bigMod)
		}
		bigBase.Mul(bigBase, bigBase)
		bigBase.Mod(bigBase, bigMod)
		exp = exp / 2
	}

	return result.Int64()
}

func fermatTest(n int64) bool {

	a := rand.Int63()

	result := big.NewInt(step(a, n-1, n))

	if result.Cmp(big.NewInt(1)) != 0 {
		return false
	}
	return true
}

func generateRandomNumbers() (int64, int64) {
	a := rand.Int63n(1000)
	b := rand.Int63n(1000)
	return a, b
}

func generateProstoeNumbers() (int64, int64) {
	var a, b int64
	for {
		a = rand.Int63n(1000)
		if a > 1 && fermatTest(a) {
			break
		}
	}
	for {
		b = rand.Int63n(1000)
		if b > 1 && b != a && fermatTest(b) {
			break
		}
	}
	return a, b
}

func extendedGCD(a, b int64) (int64, int64, int64) {

	bigU1 := big.NewInt(a)
	bigU2 := big.NewInt(1)
	bigU3 := big.NewInt(0)
	bigV1 := big.NewInt(b)
	bigV2 := big.NewInt(0)
	bigV3 := big.NewInt(1)

	for bigV1.Cmp(big.NewInt(0)) != 0 {
		q := new(big.Int).Div(bigU1, bigV1)

		t1 := new(big.Int).Mod(bigU1, bigV1)
		t2 := new(big.Int).Sub(bigU2, new(big.Int).Mul(q, bigV2))
		t3 := new(big.Int).Sub(bigU3, new(big.Int).Mul(q, bigV3))
		bigU1, bigV1 = bigV1, t1
		bigU2, bigV2 = bigV2, t2
		bigU3, bigV3 = bigV3, t3
	}

	return bigU1.Int64(), bigU2.Int64(), bigU3.Int64()
}
