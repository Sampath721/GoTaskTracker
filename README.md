# GoTaskTracker
Task Tracker Using Go Language

# commands
1) To Get All Tasks List 
cmd : go run main.go list

2) To Add Task
cmd : go run main.go add --title "your title" --description "your description" --priority low
--priority may be u can define as (low , medium, high)

3) To Mark Task As Completed
cmd : go run main.go complete --index 1 (here one is the key like JSON Element Position Simple terms assume like array[position 1])

4) To Remove/Delete Task
cmd : go run main.go remove --index 0 (here one is the key like JSON Element Position Simple terms assume like array[position 1])
 
5) To Save Task To file
cmd : go run main.go save

6) To Load Tasks From File
cmd : go run main.go load 

7) To Get Tasks Based On Priority
cdm : go run main.go list --priority medium - (medium or high or low)


@reach out kpaavansampath@gmail.com for further clarifications.
