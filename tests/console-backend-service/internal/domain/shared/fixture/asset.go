package fixture

import "github.com/kyma-project/kyma/tests/console-backend-service/internal/domain/shared"

func Asset(typeArg string) shared.Asset {
	return shared.Asset{
		Type: typeArg,
	}
}
