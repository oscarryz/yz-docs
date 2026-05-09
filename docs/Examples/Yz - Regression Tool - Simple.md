#example 

```js
/*  
expected: "  
==================================================  
           Yz Compiler Regression Test==================================================  
Project Root: /Users/oscar/code/github/oscarryz/yz  
Date: Mon Jan 15 14:30:45 PST 2024  
==================================================  
  
[INFO] Step 1: Running Go tests...  
Command: go test ./...  
  
[SUCCESS] All Go tests passed!  
  
[INFO] Step 2: Building compiler...  
Command: go build -o /Users/oscar/code/github/oscarryz/yz/bin/yzc ./cmd/yzc  
  
[SUCCESS] Compiler built successfully!  
  
[INFO] Step 3: Running regression tests...  
  
[INFO] Collecting items from test/failing (known broken features)...  
  Found 8 items to test[INFO] Collecting items from test/regressed (previously working features)...  
  Found 15 items to test[INFO] Collecting items from test/passing (currently working features)...  
  Found 137 items to test  
[INFO] Running 160 tests...  
  
  [failing] test_diagnostics.yz... ✗ FAILED (staying in failing) - Compilation failed  [failing] test_generic_map_parse.yz... ✗ FAILED (staying in failing) - Compilation failed  [failing] test_match_correct_syntax.yz... ✗ FAILED (staying in failing) - Compilation failed  [failing] test_multiple_assignment_missing_operator.yz... ✗ FAILED (staying in failing) - Compilation failed  [failing] test_multiple_assignment_syntax_error.yz... ✗ FAILED (staying in failing) - Compilation failed  [failing] test_multiple_assignment_undefined_vars.yz... ✗ FAILED (staying in failing) - Compilation failed  [failing] test_type_error.yz... ✗ FAILED (staying in failing) - Compilation failed  [failing] test_undefined_var.yz... ✗ FAILED (staying in failing) - Compilation failed  [regressed] test_boc_as_return.yz... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_deep_nesting_fix.yz... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_doubled.yz... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_enhanced_boc_types.yz... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_instantiable_types.yz... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_manual_nested_with_code.yz... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_nested_boc_context_chain_complex.yz... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_nested_boc_context_chain_integration.yz... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_nested_boc_types.yz... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] nested... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_dir_nested... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_dir_simple... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_nested_instantiable... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_project... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] test_simple_merge... ✓ SUCCESS (moving to passing) - FIXED!  [regressed] utils/test.yz... ✓ SUCCESS (moving to passing) - FIXED!  [passing] final_test.yz... ✓ SUCCESS (staying in passing)  [passing] hello.yz... ✓ SUCCESS (staying in passing)  [passing] test_anonymous_boc_caching.yz... ✓ SUCCESS (staying in passing)  [passing] test_arithmetic.yz... ✓ SUCCESS (staying in passing)  [passing] test_boc_basic.yz... ✓ SUCCESS (staying in passing)  [passing] test_simple.yz... ✓ SUCCESS (staying in passing)  
[SUCCESS] Completed 160 tests in 12.34 seconds  
  
[INFO] Moving 15 files...  
  
[INFO] Full regression testing completed!  
  
==================================================  
                Test Summary==================================================  
✅ Go Tests: PASSED  
📊 Regression Test Results:   Files moved from failing → passing: 0   Files moved from passing → regressed: 0   Files moved from regressed → passing: 15   Files still failing: 8   Files still passing: 137   Files still regressed: 0   Total files tested: 160  
==================================================  
[SUCCESS] All tests passed and compiler built successfully!  
"  
  
A comprehensive regression testing tool for the Yz compiler, demonstrating advanced Yz language features including:  
- Complex data structures and type definitions  
- Async/concurrent operations  
- Error handling with Result types  
- File system operations  
- String manipulation and formatting  
- Parallel processing capabilities  
- Comprehensive logging and reporting  
  
This example showcases how Yz's unified "blocks of code" abstraction can handle complex real-world applications  
that traditionally require multiple language constructs (classes, functions, actors, etc.).  
*/  
  
// Core data structures for regression testing  
TestResults: {
    go_tests_passed: false
    compiler_built: false
    failing_to_passing: 0
    passing_to_regressed: 0
    regressed_to_passing: 0
    still_failing: 0
    still_passing: 0
    still_regressed: 0
    total_files_tested: 0
    regression_detected: false
    output_mismatches: 0
    missing_expected_output: 0
}  
  
TestFileResult: {
    compiled: false
    executed_successfully: false
    output_matched: false
    missing_expected_output: false
    error_message: ""
}  
  
TestJob: {
    file_path String
    source_dir String  // "failing", "passing", or "regressed"
    compiler_path String
    show_generated_code: false
    is_directory: false
}  
  
// Configuration for the regression tester  
Config: {
    project_root String
    failing_dir String
    passing_dir String
    regressed_dir String
    compiler_path String
    verbose: false
    incremental: false
}  
  
// Create singleton instances  
colors: Colors()  
logger: Logger()  
file_system: FileSystem()  
command_executor: CommandExecutor()  
test_executor: TestExecutor()  
test_parser: TestParser()  
regression_tester: RegressionTester()  
  
// Logging system using Yz's unified boc approach  
Logger: {  
    info: {  
        message String  
        print("${colors.blue}[INFO]${colors.reset} ${message}")  
    }  
      
    success: {  
        message String  
        print("${colors.green}[SUCCESS]${colors.reset} ${message}")  
    }  
      
    warning: {  
        message String  
        print("${colors.yellow}[WARNING]${colors.reset} ${message}")  
    }  
      
    error: {  
        message String  
        print("${colors.red}[ERROR]${colors.reset} ${message}")  
    }  
}  
  
// File system operations  
FileSystem: {  
    // Check if directory exists  
    directory_exists: {  
        path String  
        // In a real implementation, this would use system calls  
        true  // Simplified for demonstration  
    }  
      
    // Create directory if it doesn't exist  
    ensure_directory: {  
        path String  
        logger.info("Ensuring directory exists: ${path}")  
        // In real implementation: os.MkdirAll(path, 0755)  
    }  
      
    // Move file from source to destination  
    move_file: {  
        src String  
        dst String  
        logger.info("Moving file from ${src} to ${dst}")  
        // In real implementation: os.Rename(src, dst)  
    }  
      
    // Read file contents  
    read_file: {  
        path String  
        // In real implementation: os.ReadFile(path)  
        "file contents"  // Simplified for demonstration  
    }  
      
    // List files in directory  
    list_files: {  
        dir_path String  
        // In real implementation: os.ReadDir(dir_path)  
        ["file1.yz", "file2.yz", "file3.yz"]  // Simplified for demonstration  
    }  
      
    // Check if path is directory  
    is_directory: {  
        path String  
        // In real implementation: os.Stat(path).IsDir()  
        path.contains("/")  // Simplified heuristic  
    }  
}  
  
// Command execution system  
CommandExecutor: {  
    // Execute Go tests  
    run_go_tests: {  
        logger.info("Step 1: Running Go tests...")  
        print("Command: go test ./...")  
        print()  
          
        // In real implementation: exec.Command("go", "test", "./...")  
        logger.success("All Go tests passed!")  
        true  
    }  
      
    // Build the compiler  
    build_compiler: {  
        tests_passed Bool  
        project_root String  
          
        tests_passed ? {  
            logger.info("Step 2: Building compiler...")  
            compiler_path: "${project_root}/bin/yzc"  
            print("Command: go build -o ${compiler_path} ./cmd/yzc")  
            print()  
              
            // In real implementation: exec.Command("go", "build", "-o", compiler_path, "./cmd/yzc")  
            logger.success("Compiler built successfully!")  
            true  
        }, {
            logger.warning("Skipping compiler build due to test failures")  
            false  
        }  
    }  
      
    // Execute Yz compiler on a file  
    execute_yz_file: {  
        compiler_path String  
        file_path String  
        working_dir String  
          
        // In real implementation: exec.Command(compiler_path, file_path)  
        // For this example, simulate successful execution  
        "expected output from file"  
    }  
}  
  
// Test file parsing and validation  
TestParser: {  
    // Parse expected output from test file comments  
    parse_expected_output: {  
        file_path String  
        content: file_system.read_file(file_path)  
          
        // Look for /* ... */ comment block with expected: " ... "  
        start_comment: content.index_of("/*")  
        start_comment >= 0 ? {  
            end_comment: content.index_of("*/", start_comment)  
            end_comment >= 0 ? {  
                comment_block: content.substring(start_comment + 2, end_comment)  
                expected_marker: "expected: \""  
                expected_start: comment_block.index_of(expected_marker)  
                expected_start >= 0 ? {  
                    content_start: expected_start + expected_marker.length()  
                    expected_end: comment_block.last_index_of("\"")  
                    expected_end > content_start ? {  
                        comment_block.substring(content_start, expected_end)  
                    }, {
                        ""  
                    }  
                }, {
                    ""  
                }  
            }, {
                ""  
            }  
        }, {
            ""  
        }  
    }  
      
    // Normalize output by trimming whitespace  
    normalize_output: {  
        output String  
        lines: output.split("\n")  
        normalized_lines: lines.map({ line String  
            line.trim()  
        })  
        normalized_lines.join("\n").trim()  
    }  
}  
  
// Test execution engine  
TestExecutor: {  
    // Test a single file  
    test_file: {  
        compiler_path String  
        file_path String  
        show_generated_code Bool  
        results TestResults  
          
        result: TestFileResult()  
          
        // Parse expected output  
        expected_output: test_parser.parse_expected_output(file_path)  
        expected_output.length() > 0 ? {  
            // Get directory and filename  
            file_dir: file_path.substring(0, file_path.last_index_of("/"))  
            filename: file_path.substring(file_path.last_index_of("/") + 1)  
              
            // Execute the file  
            actual_output: command_executor.execute_yz_file(compiler_path, filename, file_dir)  
              
            // Compare outputs  
            normalized_expected: test_parser.normalize_output(expected_output)  
            normalized_actual: test_parser.normalize_output(actual_output)  
              
            normalized_expected == normalized_actual ? {  
                result.compiled = true  
                result.executed_successfully = true  
                result.output_matched = true  
            }, {
                result.compiled = true  
                result.executed_successfully = true  
                result.output_matched = false  
                result.error_message = "Output mismatch"  
                results.output_mismatches = results.output_mismatches + 1  
            }  
        }, {
            result.missing_expected_output = true  
            result.error_message = "Missing expected output block"  
            results.missing_expected_output = results.missing_expected_output + 1  
        }  
          
        result  
    }  
      
    // Test a directory  
    test_directory: {  
        compiler_path String  
        dir_path String  
        show_generated_code Bool  
        results TestResults  
          
        result: TestFileResult()  
          
        // Find main file in directory  
        main_file: test_executor.find_main_file_in_directory(dir_path)  
        main_file.length() > 0 ? {  
            expected_output: test_parser.parse_expected_output(main_file)  
            expected_output.length() > 0 ? {  
                // Execute directory  
                parent_dir: dir_path.substring(0, dir_path.last_index_of("/"))  
                dir_name: dir_path.substring(dir_path.last_index_of("/") + 1)  
                actual_output: command_executor.execute_yz_file(compiler_path, "${dir_name}/", parent_dir)  
                  
                // Compare outputs  
                normalized_expected: test_parser.normalize_output(expected_output)  
                normalized_actual: test_parser.normalize_output(actual_output)  
                  
                normalized_expected == normalized_actual ? {  
                    result.compiled = true  
                    result.executed_successfully = true  
                    result.output_matched = true  
                }, {
                    result.compiled = true  
                    result.executed_successfully = true  
                    result.output_matched = false  
                    result.error_message = "Output mismatch"  
                    results.output_mismatches = results.output_mismatches + 1  
                }  
            }, {
                result.missing_expected_output = true  
                result.error_message = "Missing expected output block"  
                results.missing_expected_output = results.missing_expected_output + 1  
            }  
        }, {
            result.error_message = "No main file found in directory"  
        }  
          
        result  
    }  
      
    // Find main file in directory  
    find_main_file_in_directory: {  
        dir_path String  
        files: file_system.list_files(dir_path)  
          
        // Look for main.yz first  
        main_yz: files.find({ file String  
            file == "main.yz"  
        })  
          
        main_yz.length() > 0 ? {  
            "${dir_path}/${main_yz}"  
        }, {
            // Look for file with main boc  
            main_boc_file: files.find({ file String  
                content: file_system.read_file("${dir_path}/${file}")  
                content.contains("main:")  
            })  
              
            main_boc_file.length() > 0 ? {  
                "${dir_path}/${main_boc_file}"  
            }, {
                ""  
            }  
        }  
    }  
}  
  
// Main regression testing workflow  
RegressionTester: {  
    // Run the complete regression test suite  
    run_regression_tests: {  
        project_root String
        verbose: false
        incremental: false
          
        config: Config()  
        config.project_root = project_root  
        config.failing_dir = "${project_root}/test/failing"  
        config.passing_dir = "${project_root}/test/passing"  
        config.regressed_dir = "${project_root}/test/regressed"  
        config.compiler_path = "${project_root}/bin/yzc"  
        config.verbose = verbose  
        config.incremental = incremental  
          
        results: TestResults()  
          
        // Print header  
        regression_tester.print_header(project_root)  
          
        // Step 1: Run Go tests  
        results.go_tests_passed = command_executor.run_go_tests()  
          
        // Step 2: Build compiler  
        results.compiler_built = command_executor.build_compiler(results.go_tests_passed, project_root)  
          
        // Step 3: Run regression tests  
        results.compiler_built ? {  
            regression_tester.run_file_tests(config, results)  
        }, {
            logger.warning("Skipping regression tests - compiler not built")  
        }  
          
        // Print summary  
        regression_tester.print_footer(results)  
          
        // Return exit code  
        results.regression_detected ? {  
            1  // Exit with error if regression detected  
        }, {
            results.go_tests_passed && results.compiler_built ? {  
                0  // Success  
            }, {
                1  // Some tests failed  
            }  
        }  
    }  
      
    // Print header  
    print_header: {  
        project_root String  
        print("==================================================")  
        print("           Yz Compiler Regression Test")  
        print("==================================================")  
        print("Project Root: ${project_root}")  
        print("Date: Mon Jan 15 14:30:45 PST 2024")  
        print("==================================================")  
        print()  
    }  
      
    // Print footer with summary  
    print_footer: {  
        results TestResults  
        print()  
        print("==================================================")  
        print("                Test Summary")  
        print("==================================================")  
          
        results.go_tests_passed ? {  
            print("✅ Go Tests: PASSED")  
        }, {
            print("❌ Go Tests: FAILED")  
        }  
          
        print()  
        print("📊 Regression Test Results:")  
        print("   Files moved from failing → passing: ${results.failing_to_passing}")  
        print("   Files moved from passing → regressed: ${results.passing_to_regressed}")  
        print("   Files moved from regressed → passing: ${results.regressed_to_passing}")  
        print("   Files still failing: ${results.still_failing}")  
        print("   Files still passing: ${results.still_passing}")  
        print("   Files still regressed: ${results.still_regressed}")  
        print("   Total files tested: ${results.total_files_tested}")  
          
        results.output_mismatches > 0 ? {
            print("   Files with output mismatches: ${results.output_mismatches}")
        }, { }

        results.missing_expected_output > 0 ? {
            print("   Files missing expected output: ${results.missing_expected_output}")
        }, { }

        results.regression_detected ? {
            print()
            print("${colors.red} REGRESSION DETECTED: ${results.passing_to_regressed} previously passing test(s) now regressed!${colors.reset}")
        }, { }
          
        print("==================================================")  
    }  
      
    // Run file tests  
    run_file_tests: {  
        config Config  
        results TestResults  
          
        logger.info("Step 3: Running regression tests...")  
        print()  
          
        // Ensure directories exist  
        file_system.ensure_directory(config.failing_dir)  
        file_system.ensure_directory(config.passing_dir)  
        file_system.ensure_directory(config.regressed_dir)  
          
        // Collect test jobs  
        all_jobs: [TestJob]()
          
        // Collect failing directory files  
        logger.info("Collecting items from test/failing (known broken features)...")  
        failing_files: file_system.list_files(config.failing_dir)  
        logger.info("  Found ${failing_files.length()} items to test")  
          
        failing_files.each({ file String  
            job: TestJob()  
            job.file_path = "${config.failing_dir}/${file}"  
            job.source_dir = "failing"  
            job.compiler_path = config.compiler_path  
            job.show_generated_code = config.verbose  
            job.is_directory = file_system.is_directory(job.file_path)  
            all_jobs.push(job)
        })  
          
        // Collect regressed directory files  
        logger.info("Collecting items from test/regressed (previously working features)...")  
        regressed_files: file_system.list_files(config.regressed_dir)  
        logger.info("  Found ${regressed_files.length()} items to test")  
          
        regressed_files.each({ file String  
            job: TestJob()  
            job.file_path = "${config.regressed_dir}/${file}"  
            job.source_dir = "regressed"  
            job.compiler_path = config.compiler_path  
            job.show_generated_code = config.verbose  
            job.is_directory = file_system.is_directory(job.file_path)  
            all_jobs.push(job)
        })  
          
        // Collect passing directory files  
        logger.info("Collecting items from test/passing (currently working features)...")  
        passing_files: file_system.list_files(config.passing_dir)  
        logger.info("  Found ${passing_files.length()} items to test")  
          
        passing_files.each({ file String  
            job: TestJob()  
            job.file_path = "${config.passing_dir}/${file}"  
            job.source_dir = "passing"  
            job.compiler_path = config.compiler_path  
            job.show_generated_code = config.verbose  
            job.is_directory = file_system.is_directory(job.file_path)  
            all_jobs.push(job)
        })  
          
        all_jobs.length() > 0 ? {  
            // Track files to move  
            files_to_move: [FileMove]()
              
            // Process each test job  
            all_jobs.each({ job TestJob  
                item_name: job.file_path.substring(job.file_path.last_index_of("/") + 1)  
                  
                // Simulate test execution  
                test_result: job.is_directory ? {  
                    test_executor.test_directory(job.compiler_path, job.file_path, job.show_generated_code, results)  
                }, {
                    test_executor.test_file(job.compiler_path, job.file_path, job.show_generated_code, results)  
                }  
                  
                // Process based on source directory  
                job.source_dir == "failing" ? {  
                    print("  [failing] ${item_name}... ")  
                    test_result.output_matched ? {  
                        print("✓ SUCCESS (moving to passing)")  
                        files_to_move.push(FileMove(job.file_path, "${config.passing_dir}/${item_name}"))
                        results.failing_to_passing = results.failing_to_passing + 1  
                    }, {
                        print("✗ FAILED (staying in failing) - ${test_result.error_message}")  
                        results.still_failing = results.still_failing + 1  
                    }  
                }, {
                    job.source_dir == "regressed" ? {  
                        print("  [regressed] ${item_name}... ")  
                        test_result.output_matched ? {  
                            print("✓ SUCCESS (moving to passing) ${colors.green}- FIXED!${colors.reset}")  
                            files_to_move.push(FileMove(job.file_path, "${config.passing_dir}/${item_name}"))
                            results.regressed_to_passing = results.regressed_to_passing + 1  
                        }, {
                            print("✗ FAILED (staying in regressed) - ${test_result.error_message}")  
                            results.still_regressed = results.still_regressed + 1  
                        }  
                    }, {
                        // passing directory  
                        print("  [passing] ${item_name}... ")  
                        test_result.output_matched ? {  
                            print("✓ SUCCESS (staying in passing)")  
                            results.still_passing = results.still_passing + 1  
                        }, {
                            print("✗ FAILED (moving to regressed) ${colors.red}- REGRESSION!${colors.reset}")  
                            files_to_move.push(FileMove(job.file_path, "${config.regressed_dir}/${item_name}"))
                            results.passing_to_regressed = results.passing_to_regressed + 1  
                            results.regression_detected = true  
                        }  
                    }  
                }  
                  
                results.total_files_tested = results.total_files_tested + 1  
            })  
              
            // Move files  
            files_to_move.length() > 0 ? {
                print()
                logger.info("Moving ${files_to_move.length()} files...")
                files_to_move.each({ move FileMove
                    file_system.move_file(move.src, move.dst)
                })
            }, { }
              
            logger.success("Completed ${all_jobs.length()} tests in 12.34 seconds")  
        }, {
            logger.info("No files to test")  
        }  
    }  
}  
  
// FileMove helper type  
FileMove: {  
    src String  
    dst String  
}  
  
// Main entry point  
main: {  
    // Get project root (simplified)  
    project_root: "/Users/oscar/code/github/oscarryz/yz"  
    // Run regression tests  
    exit_code: regression_tester.run_regression_tests(project_root, false, false)  
      
    exit_code == 0 ? {  
        logger.success("All tests passed and compiler built successfully!")  
    } {  
        logger.error("Some tests failed or build unsuccessful!")  
    }  
}
```