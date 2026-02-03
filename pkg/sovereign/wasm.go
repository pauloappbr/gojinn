package sovereign

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
)

// Usamos um "Magic Footer" para saber se o arquivo foi assinado por nós
// Gojinn Signature Magic: "GJSIG"
var MagicFooter = []byte{0x47, 0x4A, 0x53, 0x49, 0x47}

// SignWasm anexa a assinatura ao FINAL do arquivo
// Formato Final: [WASM Original] + [Assinatura (64 bytes)] + [MagicFooter (5 bytes)]
func SignWasm(wasmBytes []byte, privKey ed25519.PrivateKey) ([]byte, error) {
	// 1. Assina o conteúdo original
	signature := ed25519.Sign(privKey, wasmBytes)

	// 2. Cria o novo binário
	var buf bytes.Buffer
	buf.Write(wasmBytes)   // Código
	buf.Write(signature)   // Assinatura (64 bytes fixos)
	buf.Write(MagicFooter) // Marcador (5 bytes)

	return buf.Bytes(), nil
}

// VerifyWasm verifica a assinatura e RETORNA APENAS O WASM ORIGINAL (Limpo)
func VerifyWasm(signedBytes []byte, trustedKeys []ed25519.PublicKey) ([]byte, error) {
	totalLen := len(signedBytes)
	footerLen := len(MagicFooter)
	sigLen := ed25519.SignatureSize // 64 bytes

	// Verificação básica de tamanho
	minSize := footerLen + sigLen
	if totalLen < minSize {
		return nil, fmt.Errorf("arquivo muito pequeno para conter assinatura")
	}

	// 1. Verificar o Magic Footer
	footer := signedBytes[totalLen-footerLen:]
	if !bytes.Equal(footer, MagicFooter) {
		// Se não tem nosso footer, assumimos que não está assinado (ou formato inválido)
		return nil, fmt.Errorf("arquivo não possui assinatura Gojinn válida (Magic Footer missing)")
	}

	// 2. Extrair os componentes
	// Onde termina o WASM original?
	wasmEnd := totalLen - footerLen - sigLen

	originalWasm := signedBytes[:wasmEnd]
	signature := signedBytes[wasmEnd : totalLen-footerLen]

	// 3. Verificar contra as chaves confiáveis
	verified := false
	for _, key := range trustedKeys {
		if ed25519.Verify(key, originalWasm, signature) {
			verified = true
			break
		}
	}

	if !verified {
		return nil, fmt.Errorf("assinatura digital inválida ou chave não confiável")
	}

	// 4. RETORNAR O CÓDIGO LIMPO
	// Isso é crucial! O Wazero só deve receber o 'originalWasm', sem o lixo no final.
	return originalWasm, nil
}
