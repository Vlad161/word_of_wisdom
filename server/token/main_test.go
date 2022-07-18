//go:generate mockgen -source=contract.go -package=token_test -destination=mock_contract_test.go

package token

func NewTestPrivateValue(targetBits uint, isVerified bool) value {
	return value{TargetBits: targetBits, IsVerified: isVerified}
}
