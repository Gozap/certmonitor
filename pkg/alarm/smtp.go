/*
 * Copyright 2018 Gozap, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package alarm

import (
	"net/smtp"

	"crypto/tls"
	"fmt"
	"net"
	"net/mail"

	"github.com/Sirupsen/logrus"
)

type SMTPConfig struct {
	Username string `yml:"username"`
	Password string `yml:"password"`
	From     string `yml:"from"`
	Server   string `yml:"Server"`
}

func SMTPExampleConfig() *SMTPConfig {
	return &SMTPConfig{
		Username: "mritd",
		Password: "password",
		From:     "mritd@mritd.me",
		Server:   "smtp.qq.com:465",
	}
}

func (cfg *SMTPConfig) Send(targets []string, message string) {
	for _, t := range targets {
		err := cfg.sendEmail(t, message)
		if err != nil {
			logrus.Printf("Email alarm send failed [%s]: %s", t, err)
		}
	}
}

// dial using TLS/SSL
func (cfg *SMTPConfig) dial(addr string) (*tls.Conn, error) {
	/*
		// TLS config
		tlsconfig := &tls.Config{
			// InsecureSkipVerify controls whether a client verifies the
			// server's certificate chain and host name.
			// If InsecureSkipVerify is true, TLS accepts any certificate
			// presented by the server and any host name in that certificate.
			// In this mode, TLS is susceptible to man-in-the-middle attacks.
			// This should be used only for testing.
			InsecureSkipVerify: false,
			// ServerName indicates the name of the server requested by the client
			// in order to support virtual hosting. ServerName is only set if the
			// client is using SNI (see
			// http://tools.ietf.org/html/rfc4366#section-3.1).
			ServerName: host,
			// MinVersion contains the minimum SSL/TLS version that is acceptable.
			// If zero, then TLS 1.0 is taken as the minimum.
			MinVersion: tls.VersionSSL30,
			// MaxVersion contains the maximum SSL/TLS version that is acceptable.
			// If zero, then the maximum version supported by this package is used,
			// which is currently TLS 1.2.
			MaxVersion: tls.VersionSSL30,
		}
	*/
	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	return tls.Dial("tcp", addr, nil)
}

// compose message according to "from, to, subject, body"
func (cfg *SMTPConfig) composeMsg(to string, subject string, body string) (message string) {
	// Setup headers
	headers := make(map[string]string)
	headers["From"] = cfg.From
	headers["To"] = to
	headers["Subject"] = subject
	// Setup message
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r", k, v)
	}
	message += "\r" + body
	return message
}

// send email over SSL
func (cfg *SMTPConfig) sendEmail(toAddr string, body string) (err error) {
	host, _, _ := net.SplitHostPort(cfg.Server)
	// get SSL connection
	conn, err := cfg.dial(cfg.Server)
	if err != nil {
		return
	}
	// create new SMTP client
	smtpClient, err := smtp.NewClient(conn, host)
	if err != nil {
		return
	}
	// Set up authentication information.
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, host)
	// auth the smtp client
	err = smtpClient.Auth(auth)
	if err != nil {
		return
	}
	// set To && From address, note that from address must be same as authorization user.
	from := mail.Address{Address: cfg.Username}
	to := mail.Address{Address: toAddr}
	err = smtpClient.Mail(from.Address)
	if err != nil {
		return
	}
	err = smtpClient.Rcpt(to.Address)
	if err != nil {
		return
	}
	// Get the writer from SMTP client
	writer, err := smtpClient.Data()
	if err != nil {
		return
	}
	// compose message body
	message := cfg.composeMsg(to.String(), "Website Certificate Monitor", body)
	// write message to recp
	_, err = writer.Write([]byte(message))
	if err != nil {
		return
	}
	// close the writer
	err = writer.Close()
	if err != nil {
		return
	}
	// Quit sends the QUIT command and closes the connection to the server.
	smtpClient.Quit()
	return nil
}
