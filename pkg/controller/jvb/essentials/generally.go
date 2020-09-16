package essentials

import (
	"encoding/base64"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//config maps

func ConfigMap(namespace string) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jitsi",
			Namespace: namespace,
		},
		Data: map[string]string{
			"LDAP_VERSION":               "3",
			"XMPP_INTERNAL_MUC_DOMAIN":   "internal-muc.prosody." + namespace,
			"LDAP_AUTH_METHOD":           "bind",
			"JIBRI_BREWERY_MUC":          "jibribrewery",
			"JVB0_PUBLIC_ADDR":           "",
			"JVB1_PUBLIC_ADDR":           "",
			"JVB2_PUBLIC_ADDR":           "",
			"ENABLE_AUTH":                "1",
			"JWT_TOKEN_AUTH_MODULE":      "",
			"JWT_APP_SECRET":             "",
			"TZ":                         "Etc/UTC",
			"JVB_TCP_PORT":               "443",
			"LDAP_BASE":                  "",
			"LOG_LEVEL":                  "",
			"JVB_PORT":                   "10000",
			"LDAP_USE_TLS":               "1",
			"XMPP_INTERNAL_MUC_MODULES":  "",
			"GLOBAL_MODULES":             "",
			"JVB_ENABLE_APIS":            "rest",
			"JVB_TCP_HARVESTER_DISABLED": "true",
			"LDAP_TLS_CACERT_DIR":        "/etc/ssl/certs",
			"LDAP_TLS_CIPHERS":           "SECURE256:SECURE128",
			"XMPP_SERVER":                "prosody." + namespace,
			"XMPP_DOMAIN":                "prosody." + namespace,
			"JWT_ACCEPTED_AUDIENCES":     "",
			"LDAP_FILTER":                "",
			"XMPP_AUTH_DOMAIN":           "auth.prosody." + namespace,
			"JWT_APP_ID":                 "",
			"AUTH_TYPE":                  "internal",
			"XMPP_MUC_DOMAIN":            "muc.prosody." + namespace,
			"ETHERPAD_URL_BASE":          "https://etherpad.fbi.h-da.de",
			"JVB_STUN_SERVERS":           "stun.l.google.com:19302,stun1.l.google.com:19302,stun2.l.google.com:19302",
			"XMPP_BOSH_URL_BASE":         "http://prosody." + namespace + ":5280",
			"XMPP_GUEST_DOMAIN":          "guest.prosody." + namespace,
			"ENABLE_RECORDING":           "1",
			"LDAP_TLS_CHECK_PEER":        "1",
			"LDAP_URL":                   "",
			"ENABLE_GUESTS":              "1",
			"XMPP_MUC_MODULES":           "",
			"JWT_AUTH_TYPE":              "",
			"PUBLIC_URL":                 "",
			"JWT_ACCEPTED_ISSUERS":       "",
			"XMPP_RECORDER_DOMAIN":       "recorder.prosody." + namespace,
			"JWT_ASAP_KEYSERVER":         "",
			"ENABLE_TRANSCRIPTIONS":      "0",
			"JVB_BREWERY_MUC":            "jvbbrewery",
			"XMPP_MODULES":               "",
			"LDAP_TLS_CACERT_FILE":       "/etc/ssl/certs/ca-certificates.crt",
			"GLOBAL_CONFIG":              "",
			"JWT_ALLOW_EMPTY":            "",
		},
	}
	return configMap
}

func ConfigMapWeb(namespace string) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jitsi-web",
			Namespace: namespace,
		},
		Data: map[string]string{
			"default": "" +
				"server {\n      " +
				"listen 8080 default_server;\n      " +
				"include /config/nginx/meet.conf;\n    " +
				"}",
		},
	}
	return configMap
}

func JitsiSecret(namespace string) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jitsi",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"JIBRI_RECORDER_USER":     []byte(base64.StdEncoding.EncodeToString([]byte("recorder"))),
			"JIBRI_RECORDER_PASSWORD": []byte(RandomString(16)),
			"JIBRI_XMPP_USER":         []byte(base64.StdEncoding.EncodeToString([]byte("jibri"))),
			"JIBRI_XMPP_PASSWORD":     []byte(RandomString(16)),
			"JICOFO_AUTH_USER":        []byte(base64.StdEncoding.EncodeToString([]byte("focus"))),
			"JICOFO_AUTH_PASSWORD":    []byte(RandomString(16)),
			"JICOFO_COMPONENT_SECRET": []byte(RandomString(16)),
			"JIGASI_XMPP_USER":        []byte(base64.StdEncoding.EncodeToString([]byte("jigasi"))),
			"JIGASI_XMPP_PASSWORD":    []byte(RandomString(16)),
			"JVB_AUTH_USER":           []byte(base64.StdEncoding.EncodeToString([]byte("jvb"))),
			"JVB_AUTH_PASSWORD":       []byte(RandomString(16)),
			"LDAP_BINDDN":             []byte(base64.StdEncoding.EncodeToString([]byte(""))),
			"LDAP_BINDPW":             []byte(base64.StdEncoding.EncodeToString([]byte(""))),
		},
	}
	return secret
}
