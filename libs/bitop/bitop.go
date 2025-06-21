package bitop

import (
	"errors"
	"fmt"
	"strconv"
)

// ApplyBitOperation は、reversi.Game の Board の特定行に対して演算を適用し、更新後の行を返します。
// @param row [8]int オセロの1行（0と1と7）
// @param value int 演算対象値（2進数で解釈）
// @param operator string "+" または "*" を指定
// @return [8]int 更新後の行
// @return error 不正な演算子や行の長さ不正時
func ApplyBitOperation(row [8]int, value int, operator string) ([8]int, error) {
	// 対象ビットインデックスを取得
	var bitIndices []int
	var bits []rune

	for i, v := range row {
		if v == 0 || v == 1 {
			bitIndices = append(bitIndices, i)
			bits = append(bits, rune('0'+v))
		}
	}

	if len(bitIndices) == 0 {
		return row, nil // 操作対象がない場合そのまま返す
	}

	// 元の2進数値を算出
	binaryStr := string(bits)
	originalVal, err := strconv.ParseInt(binaryStr, 2, 64)
	if err != nil {
		return row, fmt.Errorf("failed to parse binary: %v", err)
	}

	// 演算
	var newVal int64
	switch operator {
	case "+":
		newVal = originalVal + int64(value)
	case "*":
		newVal = originalVal * int64(value)
	default:
		return row, errors.New("unsupported operator: must be '+' or '*'")
	}

	// 演算結果を2進数に戻し、桁が多い場合は下位ビットを切り捨て（上位ビットだけ残す）
	resultBits := []rune(strconv.FormatInt(newVal, 2))
	if len(resultBits) > len(bitIndices) {
		resultBits = resultBits[:len(bitIndices)]
	} else if len(resultBits) < len(bitIndices) {
		// ゼロ埋め（左詰め）
		pad := make([]rune, len(bitIndices)-len(resultBits))
		for i := range pad {
			pad[i] = '0'
		}
		resultBits = append(pad, resultBits...)
	}

	// 元の row に結果を反映
	for i, idx := range bitIndices {
		row[idx] = int(resultBits[i] - '0')
	}

	return row, nil
}
