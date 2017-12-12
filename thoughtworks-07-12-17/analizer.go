package main

func encontrarValorMinimo(numeros []int) (int, error) {

	if len(numeros) == 0 || numeros == nil {
		return 0, nil
	}

	if len(numeros) == 1 {
		return numeros[0], nil
	}

	menorValor := numeros[0]
	for i := 1; i < len(numeros); i++ {
		if numeros[i] < menorValor {
			menorValor = numeros[i]
		}
	}

	return menorValor, nil
}
