package snmp

import (
	"time"

	"github.com/gosnmp/gosnmp"
)

var DefaultTimeout = time.Duration(30) * time.Second

type Connection struct {
	*gosnmp.GoSNMP
}

func Dial(config Config) (*Connection, error) {
	params := &gosnmp.GoSNMP{
		Target:        config.Address,
		Port:          161,
		Version:       gosnmp.Version3,
		SecurityModel: gosnmp.UserSecurityModel,
		MsgFlags:      gosnmp.AuthPriv,
		Timeout:       DefaultTimeout,
		SecurityParameters: &gosnmp.UsmSecurityParameters{
			UserName:                 config.Username,
			AuthenticationProtocol:   gosnmp.SHA,
			AuthenticationPassphrase: config.Password,
			PrivacyProtocol:          gosnmp.AES,
			PrivacyPassphrase:        config.Key,
		},
	}

	err := params.Connect()
	if err != nil {
		return nil, err
	}

	return &Connection{params}, nil
}

func (c *Connection) Set(oid string, value int) error {
	pdu := gosnmp.SnmpPDU{
		Name:  oid,
		Type:  gosnmp.Integer,
		Value: value,
	}

	_, err := c.GoSNMP.Set([]gosnmp.SnmpPDU{pdu})
	if err != nil {
		return err
	}

	return nil
}
