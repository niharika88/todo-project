# TODO Manager [Golang]

## Introduction
 - Use golang and gRPC to implement a microservices capable of storing TODO list data in mongodb.  
 - User should be able to create, query, update, and delete todo list items.
Bonus: use bazel to build/deploy your go binary + container

Bonus: use bazel to build/deploy go binary + container

## Functionality
 - Runs a grpc server that allows users to:
   - create a todo item for a user, by default the status for todo is assumed to be NOT DONE/PENDING
   - query a todo based on ID
   - Update any values in the todo with a given ID
   - Delete todo list item
   - Query ALL the todos or todos based on a filter
     - Filters are:
       - Todo status - can be DONE/PENDING/ALL
       - USER ID - only returns todos for that user
 - Stores all todods in local mondodb instance
   - username/passowrd as configured in the config file - dev.env
 - The implementation creates a service layer interface, so that new functionalities can be easily added
 - Also contains bazel docker rules to build and create a docker image with mongodb container running alongside it
 - Please find attached screenshots for manual test runs

### List of available RPCs (from Evans explained below)
```
+-------------+--------+-------------------+--------------------+
|   SERVICE   |  RPC   |   REQUEST TYPE    |   RESPONSE TYPE    |
+-------------+--------+-------------------+--------------------+
| ToDoService | Create | CreateItemRequest | TodoResponse       |
| ToDoService | Get    | GetItemByID       | TodoResponse       |
| ToDoService | Update | UpdateItemRequest | TodoResponse       |
| ToDoService | Delete | DeleteItemRequest | DeleteItemResponse |
| ToDoService | GetAll | GetItemsRequest   | ToDo               |
+-------------+--------+-------------------+--------------------+
```

## Setup instructions
 * Install
   * Go
   * Docker
   * Bazel
   * Mongodb
     * User/password config as present in dev.env
   * Evans 
     * https://github.com/ktr0731/evans
   
 * In project root, run
   * go mod tidy
 * To clean previous state for subsequent runs
   * bazel clean --expunge
 * Generate build files
   * bazel run //:gazelle
 * Update repos/deps
   *  bazel run :gazelle -- update-repos -from_file=go.mod -build_file_proto_mode=disable_global -to_macro=repositories.bzl%go_repositories
 * Run the server
   * bazel run //cmd
 * Open another terminal and run 
   * evans --host localhost --port 8080 -r repl
 * show available rpcs
   * show service
 * Now call different rpcs
    * call Create
    * call Get
    * call Update
    * call Delete
    * call GetAll
 * To run the tests:
   * bazel test --test_output=errors //... --@io_bazel_rules_docker//transitions:enable=no
 * Get mongodb latest image and run a mongodb container
   * docker run -d -p 27017:27017 --name bazel-mongo -v mongo-data:/data/db  mongo:latest
 * Build docker image for our project and run it
   * bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64  //cmd:image --@io_bazel_rules_docker//transitions:enable=no -- -p 8080:8080 --name bazel-docker-app


## Assumptions and future additions:
 * Since we do not have user auth/sessions, we are expecting user_id in create todo requests, ideally it can be taken from current logged in user.
 * No validation for user_id sent in create/update Todos for same reason as above.
 * Task completion is just stored as a boolean value but could be kept as an enum for better handling (since proto removes default value of false for a boolean)
 * Can have sorting/etc based on created_at/update_at
 * Can have support for scheduling todos also

