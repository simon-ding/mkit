package dji

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (d *DjiModem) GetSms() (string, error) {
	q := "AT+CMGF?"
	out, err := d.ExecAT(q)
	if err != nil {
		return "", err
	}
	out = extractMsg(out)
	if !strings.Contains(out, "+CMGF: 1") {
		//modem is not in text mode
		d.ExecAT(`AT+CMGF=1`)
	}

	input := `AT+CMGL="ALL"`
	return d.ExecAT(input)
}

func (d *DjiModem) Restart() error {
	cmd := "AT+CFUN=1,1"
	s, err := d.ExecAT(cmd)
	if err != nil {
		return err
	}
	s = strings.TrimSpace(s)
	if s != "OK" {
		return errors.New(s)
	}
	return nil
}

type SimInfo struct {
	Status         string
	PhoneNumber    string
	ICCID          string
	Operator       string
	SignalStrength int
}

/*
at> AT+CNUM

+CNUM: "My Number","+85251426631",145

# OK

at> AT+CPIN?

+CPIN: READY

OK
at> AT+QCCID

+QCCID: 89852312488535969248

OK
at> AT+COPS?

+COPS: 0,0,"CHN-CT",7

OK
at> AT+CSQ

+CSQ: 29,99

OK
at> AT+QNWINFO

+QNWINFO: "FDD LTE","46011","LTE BAND 3",1550

OK
*/
func (d *DjiModem) SimInfo() (*SimInfo, error) {
	info := SimInfo{}
	cmd := "AT+CNUM"
	out, err := d.ExecAT(cmd)
	if err != nil {
		return nil, err
	}
	out = extractMsg(out)
	pp := strings.Split(out, ":")
	if len(pp) > 1 {
		pp1 := strings.Split(pp[1], ",")
		if len(pp1) > 1 {
			info.PhoneNumber = strings.Trim(strings.TrimSpace(pp1[1]), "\"")
		}
	}

	cmd = "AT+CPIN?"
	out, err = d.ExecAT(cmd)
	if err != nil {
		return nil, err
	}
	out = extractMsg(out)
	pp = strings.Split(out, ":")
	if len(pp) > 1 {
		info.Status = strings.TrimSpace(pp[1])
	}

	cmd = "AT+QCCID"
	out, err = d.ExecAT(cmd)
	if err != nil {
		return nil, err
	}
	out = extractMsg(out)
	pp = strings.Split(out, ":")
	if len(pp) > 1 {
		info.ICCID = strings.TrimSpace(pp[1])
	}

	cmd = "AT+COPS?"
	out, err = d.ExecAT(cmd)
	if err != nil {
		return nil, err
	}
	out = extractMsg(out)
	pp = strings.Split(out, ":")
	if len(pp) > 1 {
		pp1 := strings.Split(pp[1], ",")
		if len(pp1) > 2 {
			info.Operator = strings.Trim(pp1[2], "\"")
		}
	}
	cmd = "AT+CSQ"
	out, err = d.ExecAT(cmd)
	if err != nil {
		return nil, err
	}
	out = extractMsg(out)
	pp = strings.Split(out, ":")
	if len(pp) > 1 {
		pp1 := strings.Split(pp[1], ",")
		n, err := strconv.Atoi(strings.TrimSpace(pp1[0]))
		if err != nil {
			return nil, err
		}
		info.SignalStrength = n
	}

	return &info, nil
}

func (d *DjiModem) SetUsbnetMode(i int) error {
	cmd := fmt.Sprintf("AT+QCFG=\"usbnet\",%d", i)
	out, err := d.ExecAT(cmd)
	if err != nil {
		return err
	}
	out = strings.TrimSpace(out)
	if out != "OK" {
		return errors.New(out)
	}
	return nil
}

/*
at> AT+QCFG="usbnet"

+QCFG: "usbnet",1

OK
*/
func (d *DjiModem) GetUsbnetMode() (int, error) {
	cmd := `AT+QCFG="usbnet"`
	out, err := d.ExecAT(cmd)
	if err != nil {
		return -1, err
	}
	out = extractMsg(out)
	pp := strings.Split(out, ":")
	if len(pp) > 1 {
		pp1 := strings.Split(pp[1], ",")
		if len(pp1) > 1 {
			mode, err := strconv.Atoi(strings.TrimSpace(pp1[1]))
			if err != nil {
				return -1, err
			}
			return mode, nil
		}
	}
	return -1, errors.New(out)
}

func extractMsg(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "OK")
	s = strings.TrimSpace(s)
	return s
}
