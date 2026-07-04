package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/google/gousb"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const (
	VID = 0x2c7c
	PID = 0x0125
)

func NewDjiModem() (*DjiModem, error) {
	d := &DjiModem{}
	if err := d.connect(); err != nil {
		return nil, err
	}
	return d, nil
}

type DjiModem struct {
	out *gousb.OutEndpoint
	in  *gousb.InEndpoint
}

func (d *DjiModem) connect() error {
	ctx := gousb.NewContext()

	dev, err := ctx.OpenDeviceWithVIDPID(
		gousb.ID(VID),
		gousb.ID(PID),
	)
	if err != nil {
		return err
	}

	cfg, err := dev.Config(1)
	if err != nil {
		return err
	}

	intf, err := cfg.Interface(3, 0)
	if err != nil {
		return err
	}

	out, err := intf.OutEndpoint(4)
	if err != nil {
		return err
	}
	in, err := intf.InEndpoint(6)
	if err != nil {
		return err
	}

	d.out = out
	d.in = in
	return nil
}

func (d *DjiModem) ExecAT(at string) (string, error) {
	_, err := d.out.Write([]byte(at + "\r"))
	if err != nil {
		return "", err
	}
	buf := make([]byte, 4096)

	_, err = d.in.Read(buf)
	if err != nil {
		return "", err
	}
	return cleanOutput(string(buf)), nil
}

func (d *DjiModem) AtShell() {
	for {
		fmt.Print("at> ")
		input := ""
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			continue
		}
		if input == "sms" {
			input = `AT+CMGL="ALL"`
		}
		if !strings.HasPrefix(input, "AT") {
			continue
		}
		output, err := d.ExecAT(input)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			continue
		}
		fmt.Println(output)
	}
}

func cleanOutput(s string) string {
	ss := strings.Split(s, "\n")
	r := ""
	sms := false
	for _, v := range ss {
		v = strings.TrimSpace(v)
		if sms {
			v1, err := DecodeUCS2(v)
			if err == nil {
				v = v1
			}
			sms = false
		}

		if strings.HasPrefix(v, "+CMGL") || strings.HasPrefix(v, "+CMGR") {
			sms = true
		}
		r += v

		if v == "OK" || v == "ERROR" {
			return r
		} else {
			r += "\n"
		}
	}
	return r
}

func DecodeUCS2(hexStr string) (string, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}

	reader := transform.NewReader(
		bytes.NewReader(data),
		unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(),
	)

	out, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
