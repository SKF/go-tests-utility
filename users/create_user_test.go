package users_test

import (
	"testing"

	"github.com/SKF/go-tests-utility/users"
	"github.com/stretchr/testify/require"
)

func TestGetTemporaryPassword_HappyCase(t *testing.T) {
	testFunc := users.GetTemporaryPassword
	emailMessage := "<a href=\"https://sandbox.digital-services.skf.com/sign-in?user_name=contracts_gherkin-xy-zw@disposable-emails.enlight.skf.com&password=hopelessly-premium-eft\" class=\"button primary-button\">"

	temporaryPassword, err := testFunc(emailMessage)

	require.Nil(t, err)
	require.Equal(t, "hopelessly-premium-eft", temporaryPassword)
}

func TestGetTemporaryPassword_NotEmptyPassword(t *testing.T) {
	testFunc := users.GetTemporaryPassword
	emailMessage := "<a href=\"https://sandbox.digital-services.skf.com/sign-in?user_name=contracts_gherkin-xy-zw@disposable-emails.enlight.skf.com&password=\" class=\"button primary-button\">"

	_, err := testFunc(emailMessage)

	require.NotNil(t, err)
}

func TestGetTemporaryPassword_PasswordStartingWithSpaceIsNotValid(t *testing.T) {
	testFunc := users.GetTemporaryPassword
	emailMessage := "<a href=\"https://sandbox.digital-services.skf.com/sign-in?user_name=contracts_gherkin-xy-zw@disposable-emails.enlight.skf.com&password= hopelessly-premium-eft\" class=\"button primary-button\">"

	_, err := testFunc(emailMessage)

	require.NotNil(t, err)
}

func TestGetTemporaryPassword_HyperLinkMissingAmpersandShouldError(t *testing.T) {
	testFunc := users.GetTemporaryPassword
	emailMessage := "<a href=\"https://sandbox.digital-services.skf.com/sign-in?user_name=contracts_gherkin-xy-zw@disposable-emails.enlight.skf.com%20password=hopelessly-premium-eft\" class=\"button primary-button\">"

	_, err := testFunc(emailMessage)

	require.NotNil(t, err)
}

func TestGetTemporaryPassword_HyperLinkMissingQueryStringShouldError(t *testing.T) {
	testFunc := users.GetTemporaryPassword
	emailMessage := "<a href=\"https://sandbox.digital-services.skf.com/sign-in%20user_name=contracts_gherkin-xy-zw@disposable-emails.enlight.skf.com&password=hopelessly-premium-eft\" class=\"button primary-button\">"

	_, err := testFunc(emailMessage)

	require.NotNil(t, err)
}
