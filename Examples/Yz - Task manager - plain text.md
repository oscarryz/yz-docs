
#example 

```js
TASK_FILE: "task.txt"

/*
Adds a new task to the task file
*/
add_task #(task_description String) {
	// 1. Access the reosurce: Open the file in 'apped' mode 
	file.open(TASK_FILE, file.Mode.Append())
	.and_then({
		// 2. Write the new task followed by a newline character
		f File
		r : f.write(task_description ++ "\n")
		println("✅ Task added: `task_description`")
		r
	})
	.or({
		e Error
		print("🚨 Error writing to file: `e`")
	})
}

/*
Reads all task from the task file and prints them.
*/
view_task #() { 
	// 1. Access ther resource: Open the file in 'read' mode
	file.open(TASK_FILE, file.Mode.Read())
	.and_then({ 
		f File
		
		lines: f.read_lines().or({
			e Error
			println("🚨 Error readon file: `e`")
			[String]()
		})
		lines.len() == 0 ? {
			println("📖 Your task list is empty!")
		},  {
			println("\n--- Current Tasks ---")	
			lines.for_each({
				index Int
				task  String
				// The .strip() removes the newline character (\n) from the end
				println("`index`. `task.strip()`")
			})
			println("---------------------\n")
		}
		
	}).or({ 
		e Error;
		println("📖 The task file (`TASK_FILE`) was not found. Your task list is currently empty.")	
	}}
}

// ---- Example Usage 
println("--- Running Task Manager Demo ---")

// A. Add a few task using `add_task` function
add_task("Buy groceries for the week")
add_task("Finish the Yz programming example")
add_task("Schedule dentist appointment")

// B. View all task using the 'view_task' function
view_tasks()

// C. Add one mor task 
add_task("Review the new software design document")

// D. View the updated list 
view_tasks()
```