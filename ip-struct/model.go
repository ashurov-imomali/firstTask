package ip_struct

type IpBit struct {
	Dec int
	Bit []int
}

type IP struct {
	IpBits []IpBit
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
