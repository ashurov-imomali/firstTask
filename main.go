package main

import (
	"encoding/json"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
)

type IpBit struct {
	Dec int
	Bit []int
}

type Mask struct {
	DecMsk int
	Bit    []IpBit
}

type IP struct {
	IpBits []IpBit
	Msk    Mask
}

func GetIp(strIp, strMask string) (*IP, error) {
	strIpBits := strings.Split(strIp, ".")
	newIp := new(IP)
	for _, ipBit := range strIpBits {
		dec, err := strconv.Atoi(ipBit)
		if err != nil {
			return nil, err
		}
		bit := ToBinary(dec)
		newIp.IpBits = append(newIp.IpBits, IpBit{Dec: dec, Bit: bit})
	}
	mask, err := strconv.Atoi(strMask)
	if err != nil {
		return nil, err
	}
	newIp.Msk = GetMask(mask)
	return newIp, nil
}
func ToBinary(num int) []int {
	bits := []int{0, 0, 0, 0, 0, 0, 0, 0}
	i := 1
	for num != 0 {
		bits[len(bits)-i] = num % 2
		num /= 2
		i++
	}
	return bits
}

func GetIpAndMask(ipAndMask string) (string, string) {
	split := strings.Split(ipAndMask, "/")
	return split[0], split[1]
}

func GetMask(mask int) Mask {
	return Mask{mask, MskToBin(mask)}
}

func MskToBin(msk int) []IpBit {
	var resIpBit []IpBit
	var bits []int
	count := 0
	for msk != 0 {
		bits = append(bits, 1)
		count++
		if count%8 == 0 {
			resIpBit = append(resIpBit, IpBit{BinToDesc(bits), bits})
			bits = []int{}
		}
		msk--
	}
	for count != 33 {
		if count%8 == 0 && len(bits) > 0 {
			resIpBit = append(resIpBit, IpBit{BinToDesc(bits), bits})
			bits = []int{}
		}
		bits = append(bits, 0)
		count++
	}
	return resIpBit
}

func BinToDesc(bits []int) int {
	var result int
	for i, bit := range bits {
		result += bit * int(math.Pow(float64(2), float64(len(bits)-i-1)))
	}
	return result
}

func AndOperation(fIp, sIp []int) []int {
	var resIp []int
	for i := 0; i < len(fIp); i++ {
		resIp = append(resIp, fIp[i]*sIp[i])
	}
	return resIp
}

func OrOperation(fIp, sIp []int) []int {
	var resIp []int
	for i := 0; i < len(sIp); i++ {
		resIp = append(resIp, sIp[i]+fIp[i])
	}
	return resIp
}

func NotOperator(ip []int) []int {
	var resIp []int
	for i := 0; i < len(ip); i++ {
		resIp = append(resIp, (ip[i]-1)*(ip[i]-1))
	}
	return resIp
}

func main() {
	adr, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return
	}
	for _, addr := range adr {
		log.Println(addr.Network())
		log.Println(addr.String())
	}
	return
	strIp, strMask := GetIpAndMask(adr[1].String())
	ip, err := GetIp(strIp, strMask)
	if err != nil {
		log.Println(err)
		return
	}
	bytes, err := json.MarshalIndent(ip, "", "	")
	if err != nil {
		log.Println(err)
		return
	}

	err = os.WriteFile("object.json", bytes, 0777)
	if err != nil {
		log.Println(err)
		return
	}
}
