# protolint alignment plugin [![build](https://github.com/w0rng/protolint-alignment/actions/workflows/ci.yml/badge.svg)](https://github.com/w0rng/protolint-alignment/actions/workflows/ci.yml) &nbsp;[![Coverage Status](https://coveralls.io/repos/github/w0rng/protolint-alignment/badge.svg?branch=main)](https://coveralls.io/github/w0rng/protolint-alignment?branch=main)  

The alignment_rule plugin is a custom rule for Protolint that helps ensure consistent formatting in your Protocol Buffer (.proto) files. It focuses on aligning elements like fields, enums, and messages to improve readability and maintainability.  

## Why use alignment_rule?

- Consistent Formatting: Ensures elements in your .proto files are properly aligned, making the code easier to read and understand.  
- Fix Mode Support: Automatically corrects alignment issues, saving you time on manual adjustments.  
- Customizable Severity: Set the rule's strictness level to match your team's coding standards.  
- Seamless Integration: Works as part of Protolint, fitting naturally into your linting workflow.  
- Team Collaboration: Helps maintain a uniform style across your .proto files, reducing inconsistencies in team projects.  

## Examples
### before:  
```protobuf
syntax     = "proto3";

message outer {
    option (my_option).a = true;
    message inner {
        int64 ival = 1;
    }
    repeated inner inner_message = 2;
    EnumAllowingAlias enum_field = 3;
    map<int32, string> my_map = 4;

    enum enumAllowingAlias {
        option allow_alias = true;
        UNKNOWN = 0;
        STARTED   = 1;
        RUNNING  = 2 [(custom_option) = "hello world"];
    }
}
```  
### after:
```protobuf
syntax = "proto3";

message outer {
    option (my_option).a = true;
    message inner {
        int64 ival = 1;
    }
    repeated inner inner_message = 2;
    EnumAllowingAlias enum_field = 3;
    map<int32, string> my_map    = 4;

    enum enumAllowingAlias {
        option allow_alias = true;
        UNKNOWN            = 0;
        STARTED            = 1;
        RUNNING            = 2 [(custom_option) = "hello world"];
    }
}
```