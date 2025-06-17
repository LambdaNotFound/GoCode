package stack

import "strconv"

/**
 * 150. Evaluate Reverse Polish Notation
 */

func evalRPN(tokens []string) int {
    if len(tokens) == 1 {
        num, _ := strconv.Atoi(tokens[0])
        return num
    }

    stack := []int{}
    res := 0
    for _, token := range tokens {
        if num, ok := strconv.Atoi(token); ok == nil {
            stack = append(stack, num)
        } else {
            l := len(stack)
            a, b := stack[l-2], stack[l-1]
            stack = stack[:l-2]

            switch token {
            case "+":
                res = a + b
            case "-":
                res = a - b
            case "*":
                res = a * b
            case "/":
                res = a / b
            }

            stack = append(stack, res)
        }
    }
    return res
}

/**
 * 224. Basic Calculator
 *
 * Given a string s representing a valid expression, implement a basic calculator to evaluate it,
 * and return the result of the evaluation.
 *
 * s consists of digits, '+', '-', '(', ')', and ' '.
 */
func calculate(s string) int {
    if len(s) == 0 {
        return 0
    }

    result, sign, num := 0, 1, 0

    var stk []int
    stk = append(stk, sign)

    for i := range s {
        if s[i] >= '0' && s[i] <= '9' {
            num = num*10 + int(s[i]-'0')
        } else if s[i] == '+' || s[i] == '-' {
            result += sign * num
            sign = stk[len(stk)-1]
            if s[i] != '+' {
                sign *= -1
            }
            num = 0
        } else if s[i] == '(' {
            stk = append(stk, sign)
        } else if s[i] == ')' {
            stk = stk[:len(stk)-1]
        }
    }

    result += sign * num
    return result
}
