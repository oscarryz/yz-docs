
#example

```js
TASK_FILE: "task.txt"

/*
Adds a new task to the task file
*/
add_task #(task_description String) {
	// 1. Access the resource: Open the file in 'append' mode
	file.open(TASK_FILE, file.Mode.Append())
	.and_then({ f File
		// 2. Write the new task followed by a newline character
		r: f.write(task_description ++ "\n")
		println("Task added: `task_description`")
		r
	})
	.or({ e Error
		print("Error writing to file: `e`")
	})
}

/*
Reads all tasks from the task file and prints them.
*/
view_tasks #() {
	// 1. Access the resource: Open the file in 'read' mode
	file.open(TASK_FILE, file.Mode.Read())
	.and_then({ f File

		lines: f.read_lines().or({ e Error
			println("Error reading file: `e`")
			[String]()
		})
		lines.length() == 0 ? {
			println("Your task list is empty!")
		}, {
			println("\n--- Current Tasks ---")
			lines.each({ index Int; task String
				// The .strip() removes the newline character (\n) from the end
				println("`index`. `task.strip()`")
			})
			println("---------------------\n")
		}

	}).or({ e Error
		println("The task file (`TASK_FILE`) was not found. Your task list is currently empty.")
	})
}

// ---- Example Usage
println("--- Running Task Manager Demo ---")

// A. Add a few tasks using `add_task` function
add_task("Buy groceries for the week")
add_task("Finish the Yz programming example")
add_task("Schedule dentist appointment")

// B. View all tasks using the 'view_tasks' function
view_tasks()

// C. Add one more task
add_task("Review the new software design document")

// D. View the updated list
view_tasks()
```
