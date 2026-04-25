package vpn

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

// GenerateUserUUID генерирует новый UUID для пользователя
func GenerateUserUUID() string {
	return uuid.New().String()
}

// GenerateVLESSURI генерирует ссылку VLESS+TLS для клиента HAPP
// Параметры:
//   - uuid: уникальный идентификатор пользователя
//   - host: IP или домен сервера (например, "83.97.78.70")
//   - port: порт сервера (обычно 443)
func GenerateVLESSURI(uuid, host string, port int) string {
	params := url.Values{}
	params.Set("encryption", "none")
	params.Set("security", "tls")
	params.Set("type", "tcp")
	params.Set("sni", host)
	params.Set("allowInsecure", "1")
	params.Set("fp", "chrome")

	return fmt.Sprintf("vless://%s@%s:%d?%s#MyVPN",
		uuid, host, port, params.Encode())
}

// GenerateVLESSRealityURI генерирует ссылку VLESS+Reality+gRPC (для будущего использования)
// Параметры:
//   - uuid: уникальный идентификатор пользователя
//   - host: домен для маскировки (например, "www.microsoft.com")
//   - serverIP: реальный IP сервера
//   - port: порт сервера (443)
//   - publicKey: публичный ключ Reality
//   - shortID: короткий идентификатор Reality
//   - serviceName: имя сервиса gRPC (например, "vpn")
func GenerateVLESSRealityURI(uuid, host, serverIP, publicKey, shortID, serviceName string, port int) string {
	params := url.Values{}
	params.Set("encryption", "none")
	params.Set("security", "reality")
	params.Set("sni", host)
	params.Set("pbk", publicKey)
	params.Set("sid", shortID)
	params.Set("type", "grpc")
	params.Set("serviceName", serviceName)
	params.Set("fp", "chrome")

	return fmt.Sprintf("vless://%s@%s:%d?%s#MyVPN-Reality",
		uuid, serverIP, port, params.Encode())
}
