package tokenUtil

import (
	"fmt"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	token := GenerateToken()
	fmt.Println(token)
	fmt.Println([]byte(token))
}
