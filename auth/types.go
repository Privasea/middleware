package auth

type Auth interface {
	Check(userToken string, apiServerName string, apiMethod string, apiPath string, address string) (bool, error)
}
