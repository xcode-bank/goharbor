/*
   Copyright (c) 2016 VMware, Inc. All Rights Reserved.
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package email

import (
	"bytes"
	tlspkg "crypto/tls"
	"net"
	"strconv"
	"time"

	"net/smtp"
	"text/template"

	"github.com/vmware/harbor/src/common/models"
	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/ui/config"
)

// Mail holds information about content of Email
type Mail struct {
	From    string
	To      []string
	Subject string
	Message string
}

var mc models.Email

// SendMail sends Email according to the configurations
func (m Mail) SendMail() error {
	mc, err := config.Email()
	if err != nil {
		return err
	}

	mailTemplate, err := template.ParseFiles("views/mail.tpl")
	if err != nil {
		return err
	}
	mailContent := new(bytes.Buffer)
	err = mailTemplate.Execute(mailContent, m)
	if err != nil {
		return err
	}
	content := mailContent.Bytes()

	auth := smtp.PlainAuth(mc.Identity, mc.Username, mc.Password, mc.Host)
	if mc.SSL {
		err = sendMailWithTLS(m, auth, content)
	} else {
		err = sendMail(m, auth, content)
	}

	return err
}

func sendMail(m Mail, auth smtp.Auth, content []byte) error {
	return smtp.SendMail(mc.Host+":"+strconv.Itoa(mc.Port), auth, m.From, m.To, content)
}

func sendMailWithTLS(m Mail, auth smtp.Auth, content []byte) error {
	conn, err := tlspkg.Dial("tcp", mc.Host+":"+strconv.Itoa(mc.Port), nil)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, mc.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if ok, _ := client.Extension("AUTH"); ok {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	if err = client.Mail(m.From); err != nil {
		return err
	}

	for _, to := range m.To {
		if err = client.Rcpt(to); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(content)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}

// Ping tests the connection and authentication with email server
// If tls is true, a secure connection is established, or the
// connection is insecure, and if starttls is true, Ping trys to
// upgrate the insecure connection to a secure one if email server
// supports it.
// Ping doesn't verify the server's certificate and hostname
// if the parameter insecure is ture when the connection is insecure
func Ping(addr, identity, username, password string,
	timeout int, tls, starttls, insecure bool) (err error) {
	log.Debugf("establishing TCP connection with %s ...", addr)
	conn, err := net.DialTimeout("tcp", addr,
		time.Duration(timeout)*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return
	}

	if tls {
		log.Debugf("establishing SSL/TLS connection with %s ...", addr)
		tlsConn := tlspkg.Client(conn, &tlspkg.Config{
			ServerName:         host,
			InsecureSkipVerify: insecure,
		})
		if err = tlsConn.Handshake(); err != nil {
			return
		}
		defer tlsConn.Close()

		conn = tlsConn
	}

	log.Debugf("creating SMTP client for %s ...", host)
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return
	}
	defer client.Close()

	//swith to SSL/TLS
	if !tls && starttls {
		if ok, _ := client.Extension("STARTTLS"); ok {
			log.Debugf("switching the connection with %s to SSL/TLS ...", addr)
			if err = client.StartTLS(&tlspkg.Config{
				ServerName:         host,
				InsecureSkipVerify: insecure,
			}); err != nil {
				return
			}
		} else {
			log.Debugf("the email server %s does not support STARTTLS", addr)
		}
	}

	if ok, _ := client.Extension("AUTH"); ok {
		log.Debug("authenticating the client...")
		// only support plain auth
		if err = client.Auth(smtp.PlainAuth(identity,
			username, password, host)); err != nil {
			return
		}
	} else {
		log.Debugf("the email server %s does not support AUTH, skip",
			addr)
	}

	log.Debug("ping email server successfully")

	return
}
