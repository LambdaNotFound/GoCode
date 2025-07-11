package dynamic_programming

/**
 * 53. Maximum Subarray
 *
 * Given an integer array nums, find the subarray with the largest sum, and return its sum.
 *
 * Kadane's algorithm
 *
 * DynamicProgramming, Time: O(n * sum), Space: O(n)
 *     dp[i] stores the maximum subarray ending at i
 *     dp[i] = dp[i - 1] + nums[i] if dp[i - 1] > 0
 *                                 else dp[i] = nums[i]
 */
func maxSubArray(nums []int) int {
    globalMax, curSum := nums[0], 0 // opt space, replace dp[i] w/ a var
    for _, num := range nums {
        if curSum < 0 {
            curSum = 0
        }
        curSum += num
        if curSum > globalMax {
            globalMax = curSum
        }
    }
    return globalMax
}

/**
 * 322. Coin Change
 * Return the fewest number of coins that you need to make up that amount.
 * If that amount of money cannot be made up by any combination of the coins, return -1.
 *
 * You may assume that you have an infinite number of each kind of coin.
 *
 * DynamicProgramming, Time: O(n), Space: O(n)
 *     dp[i] stores the minimum number of coins used for amount i:
 *     coins[j] is the jth coin
 *
 *     dp[i] = min(dp[i], dp[i - coins[j]] + 1) if (i - coins[j] >= 0)
 *
 *     dp[0] = 0
 */
func coinChange(coins []int, amount int) int {
    dp := make([]int, amount+1)
    for i := range dp {
        dp[i] = amount + 1
    }
    dp[0] = 0
    for i := 1; i <= amount; i++ {
        for _, coin := range coins { // reuse coins of same value
            if coin <= i {
                dp[i] = min(dp[i], dp[i-coin]+1)
            }
        }
    }
    if dp[amount] > amount {
        return -1
    }
    return dp[amount]
}

/**
 * 416. Partition Equal Subset Sum
 *
 * Given an integer array nums, return true if you can partition the
 * array into two subsets such that the sum of the elements in both subsets is equal or false otherwise
 *
 * DynamicProgramming, Time: O(n * sum), Space: O(sum)
 *     dp[i] stores if sum i can be partitioned to 2 equal subsets
 *     target = sum of the nums in the array / 2
 *
 *     Sum1 [a, b, c...], Sum2 [x] => Sum1 == Sum2
 *
 *     dp[i] == true if dp[target - i] is true
 */
func canPartition(nums []int) bool {
    totalSum := 0 // if total sum is odd, can't partition into equal subsets
    for _, num := range nums {
        totalSum += num
    }
    if totalSum%2 != 0 {
        return false
    }

    target := totalSum / 2
    dp := make([]bool, target+1)
    dp[0] = true // base case: sum 0 can always be achieved

    for _, num := range nums {
        // iterate backwards to ensure each number is only considered once
        for j := target; j >= num; j-- { // [num, ..., target]
            dp[j] = dp[j] || dp[j-num] // if Sum == j-num then theres Sum == j
        }
    }
    return dp[target]
}
