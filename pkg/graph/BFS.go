package graph

/**
 * 994. Rotting Oranges
 */
func orangesRotting(grid [][]int) int {
    m, n := len(grid), len(grid[0])
    queue := [][2]int{}
    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            if grid[r][c] == 2 {
                queue = append(queue, [2]int{r, c})
            }
        }
    }

    directions := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

    minutes := 0
    for len(queue) > 0 {
        for _, orange := range queue {
            queue = queue[1:]
            for _, dir := range directions {
                r, c := orange[0]+dir[0], orange[1]+dir[1]
                if r < 0 || r >= m || c < 0 || c >= n {
                    continue
                }
                if grid[r][c] != 1 {
                    continue
                }
                grid[r][c] = 2
                queue = append(queue, [2]int{r, c})
            }
        }
        if len(queue) > 0 {
            minutes++
        }
    }

    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            if grid[r][c] == 1 {
                return -1
            }
        }
    }
    return minutes
}
