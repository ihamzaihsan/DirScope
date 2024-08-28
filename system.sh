#!/bin/bash

# Build the my-ls-1 program
go build -o my-ls-1

# Create test directory and files
TEST_DIR="test_directory"
mkdir -p $TEST_DIR $TEST_DIR/subdir
touch $TEST_DIR/file1.txt $TEST_DIR/file2.txt $TEST_DIR/.hidden_file
touch $TEST_DIR/subdir/subfile.txt
ln -s file1.txt $TEST_DIR/symlink_file
ln -s subdir $TEST_DIR/symlink_dir
mkdir $TEST_DIR/-

# Function to run a test case
run_test_case() {
    echo -e "\nRunning test case: $1"
    echo "Command: ls ${@:2}"
    ./my-ls-1 "${@:2}"
    echo "Press Enter to continue..."
    read
    clear

}


# Run test cases based on audit_questions.md
run_test_case "No arguments"
run_test_case "Single file" "$TEST_DIR/file1.txt"
run_test_case "Directory" "$TEST_DIR"
run_test_case "-l flag" "-l"
run_test_case "-l with file" "-l" "$TEST_DIR/file1.txt"
run_test_case "-l with directory" "-l" "$TEST_DIR"
run_test_case "-l /usr/bin" "-l" "/usr/bin"
run_test_case "-R flag" "-R" "$TEST_DIR"
run_test_case "-a flag" "-a" "$TEST_DIR"
run_test_case "-r flag" "-r" "$TEST_DIR"
run_test_case "-t flag" "-t" "$TEST_DIR"
run_test_case "-la flags" "-la" "$TEST_DIR"
run_test_case "-l -t with directory" "-l" "-t" "$TEST_DIR"
run_test_case "-lRr with directory" "-lRr" "$TEST_DIR"
run_test_case "-l directory and file" "-l" "$TEST_DIR" "-a" "$TEST_DIR/file1.txt"
run_test_case "Multiple slashes" "-lR" "$TEST_DIR//subdir///" "$TEST_DIR/subdir/"
run_test_case "-la /dev" "-la" "/dev"
run_test_case "-alRrt" "-alRrt" "$TEST_DIR"
run_test_case "Directory with - name" "-" "$TEST_DIR/-"
run_test_case "Symlink file with /" "-l" "$TEST_DIR/symlink_file/"
run_test_case "Symlink file without /" "-l" "$TEST_DIR/symlink_file"
run_test_case "Symlink directory with /" "-l" "$TEST_DIR/symlink_dir/"
run_test_case "Symlink directory without /" "-l" "$TEST_DIR/symlink_dir"

# Clean up
rm -rf $TEST_DIR my-ls-1
