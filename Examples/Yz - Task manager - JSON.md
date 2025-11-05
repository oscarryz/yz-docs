
#example  

```js
TASK_FILE: "tasks.json"

// Define the structure of a Task object
Task: {
   description String
   created_at String
   done Bool
}

//This function encapsulates the crucial logic for 
// **deserializing** (reading JSON to Task Objects).

/*
Tries to read and deserialize all tasks from the JSON file.
Returns a Result<[Task], Error>
*/
get_all_tasks #(Result([Task], Error)) {
	
	// 1. Access the resource: Try to open the file in 'read' mode
	return file.open(TASK_FILE, file.Mode.Read())
	.and_then({
		f File
		
		// 2. Read the entire file content as a single string
		file_content_result : f.read_all()
        
        // Handle successful read
        file_content_result.and_then({
            content String
            
            // 3. Deserialize the JSON string into a list of Task objects
            // Assume json.from_string attempts to parse the content into the [Task] type
            json.from_string(content, [Task])
            .or({
                e Error
                // Handle JSON parse failure (e.g., file corruption) by recovering
                println("🚨 Warning: Could not deserialize JSON (File corrupted?). Starting fresh.")
                Result.Ok([Task]()) // Return an OK Result with an empty Task list
            })
        })
        // Handle file content read failure (which might happen if the file is empty or missing data)
        .or({
             // Return an OK Result with an empty Task list on read error
             Result.Ok([Task]())
        })
	})
    .or({
        e Error
        // This is primarily the File Open failure (file not found).
		println("📖 JSON file (`TASK_FILE`) not found. Creating a new list.")
        // Return an OK Result with an empty Task list on file access error
        Result.Ok([Task]())
	})
}

// This function handles **serialization** (Task Objects to writing JSON).
/*
Adds a new task object to the JSON file.
*/
add_task #(task_description String) {

    // 1. Define and initialize the new Task object using the new syntax
    new_task : Task(
        description: task_description,
        created_at:  date.time.now().to_string(), 
        done:        false
    )
    
    // 2. Read existing tasks safely using the helper function
    tasks_result : get_all_tasks()
    
    tasks_result.and_then({
        tasks [Task] // tasks is now a list of Task objects
        
        // 3. Append the new task to the list (mutation)
        tasks.add(new_task) 
        
        // 4. Serialize the entire list back to a JSON string
        json_string_result : json.to_string(tasks)
        
        // Handle serialization success
        json_string_result.and_then({
            json_string String
            
            // 5. Access resource: Open the file in 'write' mode ('w' - overwrites)
            file.open(TASK_FILE, file.Mode.Write())
            .and_then({
                f File
                // 6. Write the JSON string to the file
                f.write(json_string)
                println("✅ Task added and saved to JSON: `task_description`")
            })
            .or({
                e Error
                print("🚨 Error writing JSON to file: `e`")
            })
        })
        .or({
            e Error
            print("🚨 Error serializing tasks to JSON: `e`")
        })
    })
    // Note: Since `get_all_tasks` is designed to always return an OK Result (even if empty), 
    // the outer `.or` block is not strictly necessary here, leading to cleaner top-level logic.
}

/*
Reads all task objects and prints them.
*/
view_tasks #() {
    tasks_result : get_all_tasks()

    tasks_result.and_then({
        tasks [Task]
        
        // Conditional check using the new syntax
        tasks.len() == 0 ? {
            println("📖 Your task list is empty!")
        },  {
            println("\n--- Current Tasks ---")	
            lines.for_each({
                index Int
                task Task // 'task' is now a Task object
                
                // Ternary operator on the object property:
                status : task.done ? "[DONE]" : "[PENDING]"
                
                // Accessing struct properties using dot notation
                println("`index`. `status` `task.description` (Created: `task.created_at`)")
            })
            println("---------------------\n")
        }
    })
}

// ---- Example Usage
println("--- Running Task Manager Demo (JSON Persistence) ---")

// A. Add a few task using `add_task` function
add_task("Buy groceries for the week")
add_task("Finish the Yz programming example")
add_task("Schedule dentist appointment")

// B. View all task using the 'view_task' function
view_tasks()

// C. Add one more task
add_task("Review the new software design document")

// D. View the updated list
view_tasks()
```

