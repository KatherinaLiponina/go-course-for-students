package tests

import (
	"context"
	"homework10/internal/ports"
	"strings"
	"testing"

	"github.com/KatherinaLiponina/validation"
	"github.com/stretchr/testify/assert"
)

type validationStruct struct {
	Title string `validate:"title"`
	Text  string `validate:"text"`
}

func FuzzTestServer(f *testing.F) {
	ctx, cf := context.WithCancel(context.Background())
	endChan := make(chan int)
	hsrv, _ := ports.CreateServer(ctx, endChan)
	httpclient := getTestClient(hsrv.Addr)

	usr, _ := httpclient.createUser("Admin", "mail@mail.com")

	testcases := []string{"Hello, world", " ", "", "!12345", strings.Repeat("a", 101)}
	for _, tc := range testcases {
		f.Add(tc)
	}
	
	f.Fuzz(func(t *testing.T, s string) {
		_, err := httpclient.createAd(usr.Data.ID, s, "Test text")
		validationErr := validation.Validate(validationStruct{Title: s, Text: "Test text"})
		if validationErr == nil {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	})

	cf()
	<-endChan
}