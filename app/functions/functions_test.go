package functions

import (
	"reflect"
	"testing"
)

func TestNewUUID(t *testing.T) {
	uuid := NewUUID()
	if !(len(uuid) > 0 && uuid != "") {
		t.Error("UUID not generated")
	}
}

func TestRemoveDuplicate(t *testing.T) {
	initMap := []string{"a", "b", "c", "d"}
	resultMap := []string{"a", "b", "c", "d", "c", "d", "a"}
	RemoveDuplicate(&resultMap)
	if !reflect.DeepEqual(initMap, resultMap) {
		t.Error("the slices are different")
	}
}

func TestRound(t *testing.T) {
	if Round(3.14159265359, 0, 4) != 3.1416 {
		t.Error("rounding problems for : roundOn 0 and places 4 ")
	}

	if Round(3.14159265359, 0.1, 6) != 3.141593 {
		t.Error("rounding problems for : roundOn 0.1 and places 6 ")
	}

}

func TestEncrypt_Decrypt(t *testing.T) {
	passphrase := "Flashcards"
	phrase := "ABCD"
	b := []byte(phrase)
	encrypted, _ := Encrypt(b, passphrase)
	decrypted, err := Decrypt(encrypted, passphrase)
	if err != nil {
		t.Error("Decrypt error: " + err.Error())
	}

	if phrase != string(decrypted) {
		t.Error("Encrypt or Decrypt Error. phrase : " + phrase + " decrypted : " + string(decrypted))
	}

}
