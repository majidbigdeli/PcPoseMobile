package hello

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/snksoft/crc"
	"github.com/yerden/go-util/bcd"
)

var jc JavaCallback
var connMap = &sync.Map{}

const (
	connHost = "10.0.2.15"
	connPort = "6000"
	connType = "tcp"
)

func startInternal() {
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected.")

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		id := uuid.New().String()
		connMap.Store(id, c)

		go handleConnection(id, c, connMap)
	}
}

func handleConnection(id string, conn net.Conn, connMap *sync.Map) {

	result := new(Person)
	result.Id = id
	recData := make([]byte, 1024)

	n, err := conn.Read(recData)
	if err != nil {
		panic(err)
	}

	extractMsgLen := n - 4 - 2

	msg4Extract := make([]byte, extractMsgLen)

	index := 4

	for i := 0; i < extractMsgLen; i++ {
		msg4Extract[i] = recData[index]
		index++
	}

	index = 0

	for index < len(msg4Extract) {
		var tagID byte = 0
		tagID = msg4Extract[index]
		index++

		switch tagID {
		case 160:
			result.DllType = handleNormalType(msg4Extract, &index)
		case 161:
			result.DllVersion = handleBCDType(msg4Extract, &index)
		case 162:
			result.ProcessCode = handleBCDType(msg4Extract, &index)
		case 164:
			result.PcID = handleNormalType(msg4Extract, &index)
		case 176:
			i, err := strconv.ParseUint(handleBCDType(msg4Extract, &index), 10, 64)
			if err != nil {
				panic(err)
			}
			result.Amount = strconv.FormatUint(uint64(i), 10)
		case 177:
			i, err := strconv.ParseUint(handleBCDType(msg4Extract, &index), 10, 64)
			if err != nil {
				panic(err)
			}
			result.PayerID = strconv.FormatUint(uint64(i), 10)
		case 182:
			result.MerchantMsg = handleNormalType(msg4Extract, &index)
		case 185:
			result.MerchantadditionalData = handleNormalType(msg4Extract, &index)
		}
	}

	fmt.Println(result)

	b, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return
	}

	callbackMethod(string(b))
}

func handleNormalType(msg4Extract []byte, index *int) string {
	var len byte = 0
	len = msg4Extract[*index]
	*index++
	stringValue := ""
	var i byte
	for i = 0; i < len; i++ {
		stringValue += string(msg4Extract[*index])
		*index++
	}
	return stringValue
}

func handleEncodeNormalType(source []byte, value string, number int, index *int) {
	source[*index] = byte(number)
	src := []byte(value)
	*index++
	source[*index] = byte(len(src))
	*index++
	for i := 0; i < len(src); i++ {
		source[*index] = src[i]
		*index++
	}
}

func handleBCDType(msg4Extract []byte, index *int) string {
	var leng byte = 0
	leng = msg4Extract[*index]
	*index++
	bCDValue := make([]byte, leng)
	var i byte
	for i = 0; i < leng; i++ {
		bCDValue[i] += msg4Extract[*index]
		*index++
	}

	dec := bcd.NewDecoder(bcd.Standard)

	dst := make([]byte, bcd.DecodedLen(len(bCDValue)))

	m, err := dec.Decode(dst, bCDValue)
	if err != nil {
		return ""
	}
	return string(dst[:m])
}

func handleEncodeBCDType(source []byte, value string, number int, index *int) {
	src := []byte(value)
	source[*index] = byte(number)
	dst := make([]byte, bcd.EncodedLen(len(src)))
	m, err := byteArray2BCD(&dst, src)
	if err != nil {
		return
	}

	*index++
	source[*index] = byte(m)
	*index++
	for i := 0; i < m; i++ {
		source[*index] = dst[i]
		*index++
	}
}

func getOutBoundIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func callbackMethod(r string) {
	jc.CallFromGo(r)
}

func byteArray2BCD(dst *[]byte, src []byte) (n int, err error) {

	arrayLen := len(src)

	if arrayLen%2 != 0 {

		arrayLen++
	}

	tmpArray := make([]byte, arrayLen)

	if len(src)%2 != 0 {
		i := 0
		tmpArray[i] = 0

		for _, v := range src {
			i++
			tmpArray[i] = v
		}

	} else {
		for i := 0; i < arrayLen; i++ {
			tmpArray[i] = src[i]
		}
	}

	for j := 0; j < len(tmpArray); j++ {
		tmpArray[j] = byte(tmpArray[j] & 0xF)
	}

	BCDArray := make([]byte, arrayLen/2)

	for i := 0; i < len(BCDArray); i++ {

		hi := tmpArray[i*2] << 4
		lo := tmpArray[i*2+1] & 0xF

		if tmpArray[i*2] > 9 || tmpArray[i*2+1] > 9 || tmpArray[i*2] < 0 || tmpArray[i*2+1] < 0 {
			return 0, nil
		}

		BCDArray[i] = byte(hi | lo)

	}

	*dst = BCDArray

	return len(BCDArray), nil
}

func padLeft(str, pad string, lenght int) string {
	for {
		str = pad + str
		if len(str) == lenght {
			return str
		}
	}
}

func returnResponse(r *Response) {

	source := make([]byte, 1024)

	index := 0

	handleEncodeNormalType(source, "00", 192, &index)
	handleEncodeBCDType(source, r.SerialTransaction, 193, &index)
	handleEncodeBCDType(source, r.TraceNumber, 194, &index)
	handleEncodeNormalType(source, r.TransactionDate, 195, &index)
	handleEncodeNormalType(source, r.TransactionTime, 196, &index)
	handleEncodeBCDType(source, r.PAN, 197, &index)
	handleEncodeBCDType(source, r.TerminalNo, 199, &index)
	handleEncodeBCDType(source, r.AccountNo, 200, &index)
	handleEncodeBCDType(source, "01", 162, &index)
	handleEncodeNormalType(source, r.PcID, 196, &index)
	handleEncodeBCDType(source, r.ReqID, 228, &index)
	handleEncodeBCDType(source, r.Amount, 176, &index)

	CRC_ByteArray := make([]byte, index)
	for i := 0; i < index; i++ {
		CRC_ByteArray[i] = source[i]
	}

	ccittCrc := uint16(crc.CalculateCRC(crc.CRC16, []byte(CRC_ByteArray)))

	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, ccittCrc)

	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("Encoded: % x\n", buf.Bytes())

	Crc16_ByteArray := make([]byte, 2)

	Crc16_ByteArray = buf.Bytes()

	for i := 0; i < len(Crc16_ByteArray); i++ {
		source[index] = Crc16_ByteArray[i]
		index++

	}

	sendBytes_Len := index + 4
	sendBytes := make([]byte, sendBytes_Len)
	strSendLen := padLeft(strconv.Itoa(sendBytes_Len), "0", 4)
	value := []byte(strSendLen)
	for index = 0; index < len(value); index++ {
		sendBytes[index] = value[index]
	}

	j := 0

	for i := index; i < sendBytes_Len; i++ {
		sendBytes[i] = source[j]
		j++
	}

	if v, ok := connMap.Load(r.Id); ok {
		if conn, ok := v.(net.Conn); ok {
			if _, err := conn.Write(sendBytes); err != nil {
				fmt.Println("error on writing to connection", err.Error())
			}
			conn.Close()
		}
	}
}
