package main

import "testing"

// func TestDeveSomarDoisNumerosCorretamente(t *testing.T) {
// 	somaDoisNumeros := soma(1, 2)

// 	if somaDoisNumeros != 3 {
// 		t.Errorf("A soma deveria ser 3, veio %d", somaDoisNumeros)
// 	}
// }

func TestDeveRetornarMesmoNumeroSe1(t *testing.T) {
	mesmoNumero, _ := encontrarValorMinimo([]int{1})
	if mesmoNumero != 1 {
		t.Errorf("Valor incorreto, veio %d e deveria vir 1", 1)
	}
}

func TestDeveRetornar1sePassar1e2(t *testing.T) {
	resultado, _ := encontrarValorMinimo([]int{1, 2})
	if resultado != 1 {
		t.Errorf("Valor incorreto, veio %d e deveria vir 1", resultado)
	}
}
func TestDeveRetornar2sePassar2e3(t *testing.T) {
	resultado, _ := encontrarValorMinimo([]int{2, 3})
	if resultado != 2 {
		t.Errorf("Valor incorreto, veio %d e deveria vir 2", resultado)
	}
}

func TestDeveRetornar2sePassar3e2(t *testing.T) {
	listaNumerica := []int{3, 2}
	resultado, _ := encontrarValorMinimo(listaNumerica)
	if resultado != listaNumerica[1] {
		t.Errorf("Valor incorreto, veio %d e deveria vir %d", resultado, listaNumerica[1])
	}
}

func TestDeveRetornarMensagemDeErroSeAListaEstiverVazia(t *testing.T) {
	listaNumerica := []int{}
	_, erro := encontrarValorMinimo(listaNumerica)

	if erro == nil {
		t.Errorf("A lista deve estar vazia")
	}
}
