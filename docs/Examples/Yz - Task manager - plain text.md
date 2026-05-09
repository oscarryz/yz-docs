
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
		print("Task added: ${task_description}")
		r
	})
	.or({ e Error
		print("Error writing to file: ${e}")
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
			print("Error reading file: ${e}")
			[String]()
		})
		lines.length() == 0 ? {
			print("Your task list is empty!")
		}, {
			print("\n--- Current Tasks ---")
			lines.each({ index Int; task String
				// The .strip() removes the newline character (\n) from the end
				print("${index}. ${task.strip()}")
			})
			print("---------------------\n")
		}

	}).or({ e Error
		print("The task file (${TASK_FILE}) was not found. Your task list is currently empty.")
	})
}

// ---- Example Usage
print("--- Running Task Manager Demo ---")

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
