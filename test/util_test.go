package test

import (
	"github.com/wormholes-org/wormholesclient/tools"
	"fmt"
	"testing"
)

func TestToHex16(t *testing.T) {
	num := "0"
	rs := tools.ToHex16(num)
	fmt.Println(rs)
}

//sig 0x9cc9f3dd8df0f5314984fb6ba344f39ce26cb6d82433b79852e1887e627603dd1157bb50609c80a725b61676e2339915af09fe2117637b36dad232210921a2a41c
//address 0xc4946604072b06730157F0eFe13Eeec9341b4805
//func TestSign(t *testing.T) {
//	msg := "0xde0b6b3a76400000x80000000000000000000000000000000000000010x591813F0D13CE48f0E29081200a96565f466B2120x100000"
//	//sig, _ := tools.Sign([]byte(tx), "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9")
//	//fmt.Println("sig", hexutil.Encode(sig))
//	sig := tools.Signature("http://192.168.4.237:8574", msg)
//	fmt.Println("sig", sig)
//
//	//sig := "0x61cb06acc802f8db5195b820111e5a8caaa8ddf76dac2fdc1a18bab4272fd6db1cb72a470bca9b1f6f8e6d4a49eec6c0acafd2072325b6aa7910cb50556a13d21c"
//	//addr, _ := tools.RecoverAddress(msg, hexutil.Encode(sig))
//	addr, _ := tools.RecoverAddress(msg, sig)
//	fmt.Println("address", addr)
//}

func TestPriKeyToAddress(t *testing.T) {
	priKey := "87ff9ec48300e8df51cdf5cf92197629f02b7b4c3f2c19b8a7882b41d9791a4e"
	accoount, fromKey, _ := tools.PriKeyToAddress(priKey)
	fmt.Println(accoount)
	fmt.Println(fromKey)
}
