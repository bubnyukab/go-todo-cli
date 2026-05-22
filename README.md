# To Do CLI

I wanted to build a custom todo cli. It will be a repl, have time tracking (when you started the task, when you ended the task), task completion, day-to-day progress tracking, in the future maybe analytics (so I will save every todo day into a db).

# Setup

`go run .`

# Controls

`InputView:
    - "enter" to add todo
    - "tab" to switch to List
    - "ctrl+c" to quit
ListView:
    - "j/↓ move down, k/↑ move up"
    - "tab" to switch to Input
    - "space" to mark todo "done"
    - "ctrl+c" to quit
    - "ctrl+e" to edit current todo
    - "ctrl+d" to delete current todo`

# Todo

[ ] Add timestamp tracking (click start on task and it start a timer, when marked done it stops and save the time to see how long the task took (ideas for further analytics maybe))
[ ] Better UI/UX
[ ] Add Per-Day todo lists. Automatically when new day starts, saves the todo list from yesterday, creates a new one that can be filled for today. Saves todo lists per day, in folders per month in folders per year. Also can be used for analytics (for example percentage of completed tasks, time per task, weekly/monthly kpis)
