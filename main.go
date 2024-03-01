package main

import (
	"fmt"
	"log"
	"math"
	"net"
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
		resIp = append(resIp, sIp[i]|fIp[i])
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

func GetClass(ip int) string {
	switch {
	case ip > 0 && ip <= 127:
		return "A"
	case ip > 127 && ip <= 191:
		return "B"
	case ip > 191 && ip <= 223:
		return "C"
	case ip > 223 && ip <= 239:
		return "D"
	case ip > 239 && ip <= 247:
		return "E"
	default:
		return "WRONG IP"
	}
}

func GetStrIp(bits []IpBit) string {
	strIp := ``
	for _, bit := range bits {
		strIp += strconv.Itoa(bit.Dec) + "."
	}
	return strIp[:len(strIp)-1]
}

func LocalIP() error {
	adr, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return err
	}
	strIp, strMask := GetIpAndMask(adr[3].String())
	ip, err := GetIp(strIp, strMask)

	if err != nil {
		log.Println(err)
		return err
	}
	///Применить маску подсети к ip
	var startIp []IpBit
	for i, bit := range ip.IpBits {
		bitRes := AndOperation(ip.Msk.Bit[i].Bit, bit.Bit)
		dec := BinToDesc(bitRes)
		startIp = append(startIp, IpBit{dec, bitRes})
	}
	var nanMsk []IpBit

	for _, bit := range ip.Msk.Bit {
		bitRes := NotOperator(bit.Bit)
		dec := BinToDesc(bitRes)
		nanMsk = append(nanMsk, IpBit{dec, bitRes})
	}

	var endIp []IpBit
	for i, bit := range nanMsk {
		bitRes := OrOperation(bit.Bit, startIp[i].Bit)
		dec := BinToDesc(bitRes)
		endIp = append(endIp, IpBit{dec, bitRes})
	}
	///Логируем все результаты
	log.Println("IP:", GetStrIp(ip.IpBits))
	log.Println("Mask:", ip.Msk.DecMsk)
	log.Println("Class:", GetClass(ip.IpBits[0].Dec))
	log.Println("Start IP:", GetStrIp(startIp))
	log.Println("End IP:", GetStrIp(endIp))
	return nil
}

func PrintIP() error {
	s := ``
	log.Println("Print IP/Mask: ")
	_, err := fmt.Scan(&s)
	if err != nil {
		log.Println(err)
		return err
	}
	strIp, strMask := GetIpAndMask(s)
	ip, err := GetIp(strIp, strMask)
	if err != nil {
		log.Println(err)
		return err
	}
	///Применить маску подсети к ip
	var startIp []IpBit
	for i, bit := range ip.IpBits {
		bitRes := AndOperation(ip.Msk.Bit[i].Bit, bit.Bit)
		dec := BinToDesc(bitRes)
		startIp = append(startIp, IpBit{dec, bitRes})
	}
	var nanMsk []IpBit

	for _, bit := range ip.Msk.Bit {
		bitRes := NotOperator(bit.Bit)
		dec := BinToDesc(bitRes)
		nanMsk = append(nanMsk, IpBit{dec, bitRes})
	}

	var endIp []IpBit
	for i, bit := range nanMsk {
		bitRes := OrOperation(bit.Bit, startIp[i].Bit)
		dec := BinToDesc(bitRes)
		endIp = append(endIp, IpBit{dec, bitRes})
	}
	///Логируем все результаты
	log.Println("IP:", GetStrIp(ip.IpBits))
	log.Println("Mask:", ip.Msk.DecMsk)
	log.Println("Class:", GetClass(ip.IpBits[0].Dec))
	log.Println("Start IP:", GetStrIp(startIp))
	log.Println("End IP:", GetStrIp(endIp))
	return nil
}

func main() {
	for true {
		log.Println("Choose operation:\n 1 - Local IP\n 2 - Print IP \n 3 - Exit")
		n := 0
		_, err := fmt.Scan(&n)
		if err != nil {
			log.Println(err)
			return
		}
		switch n {
		case 1:
			err := LocalIP()
			if err != nil {
				return
			}
		case 2:
			err := PrintIP()
			if err != nil {
				return
			}
		case 3:
			log.Println("Good Bye :)")
			return
		default:
			log.Println("Bad request (-_-)")
		}
	}

}
