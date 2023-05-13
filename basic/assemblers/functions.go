package assemblers

import (
	"fmt"
	"strings"
)

func Dispatch(parms map[string]interface{}, instruction string, context map[string]interface{}, labels map[string]string) (machineCode string, err error) {
	switch instruction {
	case "jal":
		machineCode, err = Jal(parms, context)
	case "jalr":
		machineCode, err = Jalr(parms, context)
	case "lui":
		machineCode, err = Lui(parms, context)
	case "auipc":
		machineCode, err = Auipc(parms, context)
	case "ebreak":
		machineCode, err = Ebreak(parms, context)
	case "lb", "lh", "lw", "lbu", "lhu":
		machineCode, err = Loads(parms, context)
	case "sb", "sh", "sw":
		machineCode, err = Stores(parms, context)
	case "add", "sub", "xor", "or", "and", "sll", "srl", "sra", "slt", "sltu":
		machineCode, err = RtypeAlu(parms, context)
	case "addi", "xori", "ori", "andi", "slli", "srli", "srai", "slti", "sltiu":
		machineCode, err = ItypeAlu(parms, context, labels)
	case "beq", "bne", "blt", "bge", "bltu", "bgeu":
		machineCode, err = BtypeBranch(parms, context)
	case "csrrw", "csrrs", "csrrc", "csrrwi", "csrrsi", "csrrci":
		machineCode, err = ItypeCSR(parms, context)
	case "mret":
		machineCode, err = MRet(parms, context)
	}

	if err != nil {
		return "", err
	}

	return machineCode, nil
}

func GetFields(instruction string, ass string) []string {
	switch instruction {
	case "addi":
		return GetAddiFields(ass)
	case "jal":
		return GetJalFields(ass)
	case "jalr":
		return GetJalrFields(ass)
	case "beq", "bne", "blt", "bge", "bltu", "bgeu":
		return GetBranchFields(ass)
	case "lb", "lh", "lw", "lbu", "lhu":
		return GetLoadsFields(ass)
	}

	return nil
}

func GetLabel(instruction string, ass string) (label string, err error) {
	fields := GetFields(instruction, ass)
	if fields == nil {
		return "", fmt.Errorf("no fields found for instruction: " + instruction)
	}

	switch instruction {
	case "jal", "jalr":
		return fields[3], nil
	case "addi", "beq", "bne", "blt", "bge", "bltu", "bgeu":
		return fields[4], nil
	case "lb", "lh", "lw", "lbu", "lhu":
		return fields[3], nil
	}

	return "", fmt.Errorf("unknown instruction: " + instruction)
}

func GetLabelRef(instruction string, ass string, labels map[string]string) (label string, value string, err error) {
	label, err = GetLabel(instruction, ass)
	if err != nil {
		return "", "", err
	}

	label = strings.ReplaceAll(label, "@", "")

	value = labels[label]

	return label, value, nil
}
