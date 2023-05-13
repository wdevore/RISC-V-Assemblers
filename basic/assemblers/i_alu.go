package assemblers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wdevore/RISC-V-Assemblers/utils"
)

func GetAddiExpr() *regexp.Regexp {
	rxpr, _ := regexp.Compile(`([a-z]+)[ ]+([xa0-9]+),[ ]*([xa0-9]+),[ ]*([@\w-]+)`)
	return rxpr
}

func GetAddiFields(ass string) []string {
	rxpr := GetAddiExpr()

	return rxpr.FindStringSubmatch(ass)
}

//	addi rd, rs1, imm
//
// Example: addi x3, x1, x2
// OR
//
//	addi x4, x0, @StringOK
func ItypeAlu(parms map[string]interface{}, json map[string]interface{}, labels map[string]string) (macCode string, err error) {
	VerboseEnabled := parms["VerboseEnabled"] == "Yes"
	if VerboseEnabled {
		fmt.Println("### BEGIN Itype ###")
	}
	ass := fmt.Sprintf("%s", json["Assembly"])

	fields := GetAddiFields(ass)

	instru := fields[1]

	rd := fields[2]
	if VerboseEnabled {
		fmt.Println("Destination register: ", rd)
	}
	rs1 := fields[3]
	if VerboseEnabled {
		fmt.Println("Rs1 register: ", rs1)
	}

	// Determine what form the Immediate is, for example, 42 or @xxx
	// If "@" is present then it is an absolute address specified by
	// a Label
	imm := fields[4]
	if VerboseEnabled {
		fmt.Println("Immediate: ", imm)
	}

	ti := ""
	produced := []byte{}

	immAsLabel := strings.Contains(imm, "@")
	immAsWA := strings.Contains(imm, "WA:")

	if immAsLabel {
		// Convert label to BA address
		// Remove "@" dereference

		_, value, err := GetLabelRef("addi", ass, labels)
		if err != nil {
			return "", err
		}

		immInt, err := utils.StringHexToInt(value)
		if err != nil {
			return "", err
		}

		// Convert from word-addressing to byte-addressing
		immInt *= 4

		produced = utils.IntToBinaryArray(immInt)
	} else if immAsWA {
		immInt, err := utils.StringHexToInt(imm)
		if err != nil {
			return "", err
		}

		// Convert from word-addressing to byte-addressing
		immInt *= 4
		if VerboseEnabled {
			fmt.Println("immInt: ", immInt)
		}
		ti = utils.IntToBinaryString(immInt)
		produced = utils.BinaryStringToArray(ti)
	} else {
		immInt, err := utils.StringHexToInt(imm)
		if err != nil {
			return "", err
		}

		if VerboseEnabled {
			fmt.Println("immInt: ", immInt)
		}
		ti = utils.IntToBinaryString(immInt)
		produced = utils.BinaryStringToArray(ti)
	}

	if VerboseEnabled {
		fmt.Println("produced: ", produced)
	}

	instruction := make([]byte, 32)

	// The LSB is at [31] (i.e. reversed)
	//  0                                                             31   memory order
	// [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	//  31                                                            0    logical order
	//  MSB                                                           LSB

	// Immediate
	instruction[31] = produced[20]
	if instru == "srai" {
		instruction[30] = 1
	} else {
		instruction[30] = produced[21]
	}
	instruction[29] = produced[22]
	instruction[28] = produced[23]
	instruction[27] = produced[24]
	instruction[26] = produced[25]
	instruction[25] = produced[26]
	instruction[24] = produced[27]
	instruction[23] = produced[28]
	instruction[22] = produced[29]
	instruction[21] = produced[30]
	instruction[20] = produced[31]

	// Rs1
	rs1Int, err := utils.StringRegToInt(rs1)
	if err != nil {
		return "", err
	}

	rs1Arr := utils.IntToBinaryArray(rs1Int)
	instruction[19] = rs1Arr[27]
	instruction[18] = rs1Arr[28]
	instruction[17] = rs1Arr[29]
	instruction[16] = rs1Arr[30]
	instruction[15] = rs1Arr[31]

	if VerboseEnabled {
		fmt.Println("___ ", instru, " ___")
	}
	// func3
	switch instru {
	case "addi":
		instruction[14] = 0
		instruction[13] = 0
		instruction[12] = 0
	case "xori":
		instruction[14] = 1
		instruction[13] = 0
		instruction[12] = 0
	case "ori":
		instruction[14] = 1
		instruction[13] = 1
		instruction[12] = 0
	case "andi":
		instruction[14] = 1
		instruction[13] = 1
		instruction[12] = 1
	case "slli":
		instruction[14] = 0
		instruction[13] = 0
		instruction[12] = 1
	case "srli":
		instruction[14] = 1
		instruction[13] = 0
		instruction[12] = 1
	case "srai":
		instruction[14] = 1
		instruction[13] = 0
		instruction[12] = 1
	case "slti":
		instruction[14] = 0
		instruction[13] = 1
		instruction[12] = 0
	case "sltiu":
		instruction[14] = 0
		instruction[13] = 1
		instruction[12] = 1
	}

	// rd
	rdInt, err := utils.StringRegToInt(rd)
	if err != nil {
		return "", err
	}

	rdArr := utils.IntToBinaryArray(rdInt)
	instruction[11] = rdArr[27]
	instruction[10] = rdArr[28]
	instruction[9] = rdArr[29]
	instruction[8] = rdArr[30]
	instruction[7] = rdArr[31]

	//            6     0
	// Set Opcode 0010011
	instruction[6] = 0
	instruction[5] = 0
	instruction[4] = 1
	instruction[3] = 0
	instruction[2] = 0
	instruction[1] = 1
	instruction[0] = 1

	instr := utils.BinaryArrayToString(instruction, true)

	if VerboseEnabled {
		fmt.Println("   imm11:0     |  rs1 | funct3 |   rd  |  opcode")
		fmt.Printf("%v     %v    %v    %v    %v\n", instr[0:12], instr[12:17], instr[17:20], instr[20:25], instr[25:32])
		fmt.Println("Instruction Bin: ", instr)
		fmt.Printf("Nibbles: %v %v %v %v %v %v %v %v\n", instr[0:4], instr[4:8], instr[8:12], instr[12:16], instr[16:20], instr[20:24], instr[24:28], instr[28:32])
		fmt.Println("### END Itype ###")
	}

	return utils.BinaryStringToHexString(instr, false), nil
}
