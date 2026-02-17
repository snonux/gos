# Plan

- Add a new flag `--stats` which would only print out the stats for all social networks, but do nothing else.

## Implementation Outline for `--stats` flag: (ALL DONE)

1.  **Modify `internal/main.go`**: (DONE)
    *   **Define the `--stats` flag**: Add `statsOnly := flag.Bool("stats", false, "Print statistics for all social networks and exit")` near other flag definitions.
    *   **Populate `config.Args`**: Assign `*statsOnly` to `args.StatsOnly` after `flag.Parse()`.
    *   **Conditional logic**: After `flag.Parse()` and before calling `run(ctx, args)`, add a check:
        ```go
        if args.StatsOnly {
            // Call the new function to print all stats
            schedule.PrintAllStats(args)
            return // Exit after printing stats
        }
        ```

2.  **Modify `internal/config/args.go`**: (DONE)
    *   Add `StatsOnly bool` to the `Args` struct.

3.  **Modify `internal/schedule/stats.go`**: (DONE)
    *   **Create a new public function `PrintAllStats(args config.Args)`**:
        *   This function will iterate through `args.Platforms`.
        *   For each platform, it will call `newStats` to gather the statistics.
        *   Then, it will call `stats.RenderTable` to display the statistics for that platform.