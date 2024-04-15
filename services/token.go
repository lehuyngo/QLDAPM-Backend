package services

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/util"

func EncryptToken(raw string) (string, error) {
	return util.Encrypt(raw, Config.Token.SecretKey)
}

func DecryptToken(token string) (string, error) {
	return util.Decrypt(token, Config.Token.SecretKey)
}
