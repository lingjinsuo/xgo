package testcase

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	v8 "gitlab.chuangzhen-sh.net/golang/xgo/business/orderbook/v8"
)

// 解析日志文件中的消息
func parseMessage(line string) (msgType int, msg interface{}, err error) {
	// 提取 MsgType
	msgTypeRe := regexp.MustCompile(`MsgType:(\d+)`)
	msgTypeMatch := msgTypeRe.FindStringSubmatch(line)
	if len(msgTypeMatch) < 2 {
		return 0, nil, fmt.Errorf("cannot find MsgType in line")
	}
	msgType, _ = strconv.Atoi(msgTypeMatch[1])

	// 提取 SecurityCode
	securityCodeRe := regexp.MustCompile(`SecurityCode:(\d+)`)
	securityCodeMatch := securityCodeRe.FindStringSubmatch(line)
	if len(securityCodeMatch) < 2 {
		return 0, nil, fmt.Errorf("cannot find SecurityCode in line")
	}
	securityCode, _ := strconv.ParseUint(securityCodeMatch[1], 10, 32)

	// 提取 OrderId
	orderIdRe := regexp.MustCompile(`OrderId:(\d+)`)
	orderIdMatch := orderIdRe.FindStringSubmatch(line)
	if len(orderIdMatch) < 2 {
		return 0, nil, fmt.Errorf("cannot find OrderId in line")
	}
	orderId, _ := strconv.ParseUint(orderIdMatch[1], 10, 64)

	switch msgType {
	case 30: // AddOrder
		// 提取 Price
		priceRe := regexp.MustCompile(`Price:(\d+)`)
		priceMatch := priceRe.FindStringSubmatch(line)
		price, _ := strconv.ParseUint(priceMatch[1], 10, 32)

		// 提取 Quantity
		quantityRe := regexp.MustCompile(`Quantity:(\d+)`)
		quantityMatch := quantityRe.FindStringSubmatch(line)
		quantity, _ := strconv.ParseUint(quantityMatch[1], 10, 32)

		// 提取 Side
		sideRe := regexp.MustCompile(`Side:(\d+)`)
		sideMatch := sideRe.FindStringSubmatch(line)
		side, _ := strconv.ParseUint(sideMatch[1], 10, 16)

		msg = &v8.AddOrder{
			OrderId:      orderId,
			SecurityCode: uint32(securityCode),
			Price:        uint32(price),
			Quantity:     uint32(quantity),
			Side:         uint16(side),
		}

	case 31: // ModifyOrder
		// 提取 Quantity
		quantityRe := regexp.MustCompile(`Quantity:(\d+)`)
		quantityMatch := quantityRe.FindStringSubmatch(line)
		quantity, _ := strconv.ParseUint(quantityMatch[1], 10, 32)

		// 提取 Side
		sideRe := regexp.MustCompile(`Side:(\d+)`)
		sideMatch := sideRe.FindStringSubmatch(line)
		side, _ := strconv.ParseUint(sideMatch[1], 10, 16)

		msg = &v8.ModifyOrder{
			OrderId:  orderId,
			Quantity: uint32(quantity),
			Side:     uint16(side),
		}

	case 32: // DeleteOrder
		// 提取 Side
		sideRe := regexp.MustCompile(`Side:(\d+)`)
		sideMatch := sideRe.FindStringSubmatch(line)
		side, _ := strconv.ParseUint(sideMatch[1], 10, 16)

		msg = &v8.DeleteOrder{
			OrderId: orderId,
			Side:    uint16(side),
		}

	default:
		return msgType, nil, fmt.Errorf("unknown MsgType: %d", msgType)
	}

	return msgType, msg, nil
}

func RunOBTestCase(securityCode uint32, fileName string) {
	// 1. 初始化股票信息
	orderbook := v8.New(securityCode)

	// 2. 打开文件并逐行读取
	// fileName 是相对于 testcase 目录的路径，如 "test_files/1-18-1-A.txt"
	// 测试运行时的工作目录是 testcase，所以直接使用相对路径
	filePath := "./" + fileName

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// 跳过空行
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 解析消息
		msgType, msg, err := parseMessage(line)
		if err != nil {
			fmt.Printf("第 %d 行解析失败: %v\n", lineNum, err)
			continue
		}

		// 3. 根据 MsgType 调用相应方法
		var opErr error
		switch msgType {
		case 30: // AddOrder
			addMsg := msg.(*v8.AddOrder)
			opErr = orderbook.Add(addMsg)
		case 31: // ModifyOrder
			modifyMsg := msg.(*v8.ModifyOrder)
			opErr = orderbook.Modify(modifyMsg)
		case 32: // DeleteOrder
			deleteMsg := msg.(*v8.DeleteOrder)
			opErr = orderbook.Delete(deleteMsg)
		}

		if opErr != nil {
			fmt.Printf("第 %d 行 MsgType:%d 操作失败: %v\n", lineNum, msgType, opErr)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}

	// 4. 打印 OrderBook
	orderbook.PrintOrderBookTable()
	orderbook.PrintOrderBook(10)
}
