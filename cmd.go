package main

import (
	"fmt"

	"qrpay/serve"
)

func main() {
	qr := serve.NewMyQRPay()
	res := qr.Build(serve.BuildQROptions{
		BankBin:    "970407",
		BankNumber: "19033868110065",
		Amount:     "100000",
		Remark:     "test",
	})

	fmt.Println(res)
	//qrNew := serve.NewQRPay("00020101021138580010A000000727012800069704070114190338681560110208QRIBFTTA53037045802VN8300840063043C31")
	//fmt.Println(qrNew.IsValid)
}
