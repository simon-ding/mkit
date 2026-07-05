package dji

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/google/gousb"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const (
	VID          = 0x2c7c
	PID          = 0x0125
	InterfaceNum = 2
)

func NewDjiModem() (*DjiModem, error) {
	d := &DjiModem{ctx: gousb.NewContext()}
	if err := d.connect(); err != nil {
		return nil, err
	}
	return d, nil
}

type DjiModem struct {
	out *gousb.OutEndpoint
	in  *gousb.InEndpoint
	ctx *gousb.Context
}

func (d *DjiModem) connect() error {

	dev, err := d.ctx.OpenDeviceWithVIDPID(
		gousb.ID(VID),
		gousb.ID(PID),
	)
	if err != nil {
		return err
	}
	desc := dev.Desc

	// for _, cfg := range dev.Desc.Configs {
	// 	fmt.Printf("Config %d\n", cfg.Number)

	// 	for _, intf := range cfg.Interfaces {
	// 		fmt.Printf("  Interface %d\n", intf.Number)

	// 		for _, alt := range intf.AltSettings {
	// 			fmt.Printf("    Alt %d, Class=%d SubClass=%d Protocol=%d\n", alt.Number, alt.Class, alt.SubClass, alt.Protocol)

	// 			for _, ep := range alt.Endpoints {
	// 				fmt.Printf(
	// 					"      EP %d Dir=%v Type=%v MaxPacket=%d\n",
	// 					ep.Number,
	// 					ep.Direction,
	// 					ep.TransferType,
	// 					ep.MaxPacketSize,

	// 				)
	// 			}
	// 		}
	// 	}
	// }

	cfgDesc := desc.Configs[1] // Config 1

	intfDesc := cfgDesc.Interfaces[InterfaceNum] // Interface 3

	alt := intfDesc.AltSettings[0] // AltSetting 0

	var (
		inNum  int
		outNum int
	)

	for _, ep := range alt.Endpoints {
		switch ep.Direction {
		case gousb.EndpointDirectionIn:
			if ep.TransferType == gousb.TransferTypeBulk {
				inNum = ep.Number
			}
		case gousb.EndpointDirectionOut:
			if ep.TransferType == gousb.TransferTypeBulk {
				outNum = ep.Number
			}
		}
	}

	if inNum == 0 || outNum == 0 {
		return fmt.Errorf("bulk endpoints not found")
	}

	cfg, err := dev.Config(1)
	if err != nil {
		return err
	}

	intf, err := cfg.Interface(InterfaceNum, 0)
	if err != nil {
		return err
	}

	in, err := intf.InEndpoint(inNum)
	if err != nil {
		return err
	}

	out, err := intf.OutEndpoint(outNum)
	if err != nil {
		return err
	}

	d.out = out
	d.in = in
	return nil
}

func (d *DjiModem) ExecAT(at string) (string, error) {
	d.drainInput()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := d.out.WriteContext(ctx, []byte(at+"\r"))
	if err != nil {
		return "", err
	}

	buf := make([]byte, 4096)
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = d.in.ReadContext(ctx, buf)
	if err != nil {
		return "", err
	}
	return cleanOutput(string(buf)), nil
}

func (d *DjiModem) drainInput() {
	if d.in == nil {
		return
	}

	buf := make([]byte, 64)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, err := d.in.ReadContext(ctx, buf)
		cancel()
		if err != nil {
			return
		}
	}
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

func (d *DjiModem) Close() error {
	return d.ctx.Close()
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
