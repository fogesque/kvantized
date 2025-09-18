# C++ Coding Style Guide

This guide outlines the coding standards enforced by our clang-tidy configuration. **All warnings are treated as errors**, so following these conventions is mandatory.

## Naming Conventions

### Type Names
- **Classes, Structs, Enums**: `CamelCase`
  ```cpp
  class FileManager { };
  struct UserData { };
  enum ColorMode { };
  ```

- **Type Aliases & Typedefs**: `CamelCase`
  ```cpp
  using StringMap = std::map<std::string, std::string>;
  typedef std::vector<int> IntVector;
  ```

### Function Names
- **Public Member Functions**: `CamelCase`
  ```cpp
  class DataProcessor {
  public:
	void ProcessData();
	int GetUserCount();
  }
  ```
- **Private Member Functions**: `camelBack`
  ```cpp
  class DataProcessor {
  private:
	void processData();
	int getUserCount();
  }
  ```

- **Global Functions**: `CamelCase`
  ```cpp
  void ProcessFile();
  int CalculateSum();
  ```

### Variable Names
- **Local Variables & Parameters**: `camelBack`
  ```cpp
  int userCount = 0;
  void SetName(const std::string & firstName);
  ```

- **Member Variables**: `camelBack`
  ```cpp
  class User {
      std::string userName;
      int age;
  };
  ```

### Constants
- **All Constants** (global, local, member): `camelBack`
  ```cpp
  const int maxEntryPoint = 3;
  static constexpr double piValue = 3.14159;
  ```

- **Enum Constants**: `camelBack`
  ```cpp
  enum Status {
      success,
      errorNotFound,
      errorTimeout
  };
  ```

## Function Complexity Limits

Keep functions simple and maintainable:

- **Maximum 100 lines** per function
- **Maximum 50 statements** per function
- **Maximum 15 branch points** (if/else, switch cases, loops)
- **Cognitive complexity threshold: 20** (nested conditions and loops increase complexity)

If a function exceeds these limits, refactor it into smaller functions.

## Code Safety Requirements

### Always Use Braces
Control statements must use braces, even for single lines:
```cpp
// ✅ Good
if (condition) {
    DoSomething();
}

// ❌ Bad - will cause build error
if (condition)
    DoSomething();
```

### Magic Numbers
Avoid unnamed numeric constants. Define named constants instead:
```cpp
// ❌ Bad
if (retryCount > 5) { }

// ✅ Good
const int maxRetries = 5;
if (retryCount > maxRetries) { }
```

### Null Pointers
Always use `nullptr`, never `NULL` or `0`:
```cpp
// ✅ Good
int * ptr = nullptr;

// ❌ Bad
int * ptr = NULL;
int * ptr2 = 0;
```

## Modern C++ Practices

### Range-Based For Loops
Prefer range-based loops when iterating over containers:
```cpp
// ✅ Good
for (const auto & item : container) {
    ProcessItem(item);
}

// ❌ Avoid when possible
for (size_t i = 0; i < container.size(); ++i) {
    ProcessItem(container[i]);
}
```

### Auto Keyword
Use `auto` for variables when the type name is longer than 5 characters and the type is clear from context:
```cpp
// ✅ Good
auto userIterator = userMap.find("john");
auto result = CalculateComplexValue();

// Keep explicit for simple types
int count = 0;
bool flag = true;
```

### Special Member Functions
Use `= default` for default implementations:
```cpp
class Widget {
public:
    Widget() = default;
    ~Widget() = default;
    Widget(const Widget &) = default;
	Widget(Widget &&) = default;
    Widget & operator=(const Widget &) = default;
	Widget & operator=(Widget &&) = default;
};
```

### Member Initialization
Initialize member variables in-class when possible:
```cpp
class User {
    std::string name = "Unknown";
    int age = 0;
    bool active = false;
};
```

## Function Parameters

### Easily Swappable Parameters
Be careful with functions that have multiple parameters of the same type:
```cpp
// ⚠️ Risky - easy to swap width and height
void SetDimensions(int width, int height);

// ✅ Better - use strong types or a struct
struct Dimensions {
    int width;
    int height;
};
void SetDimensions(const Dimensions & dims);
```

For functions with 2+ parameters, consider using named parameters or structs for clarity.

## Build Integration

- **All checks run on both source files and headers**
- **Warnings are treated as errors** - code must pass all checks to compile
- The configuration inherits from parent directories if additional rules exist

## Quick Reference

| Element | Style | Example |
|---------|-------|---------|
| Class/Struct | CamelCase | `UserAccount` |
| Enum | CamelCase | `ErrorCode` |
| Enum Constant | camelBack | `errorTimeout` |
| Public Member Function | CamelCase | `GetUserName()` |
| Private Member Function | camelBack | `getUserName()` |
| Global Function | CamelCase | `ProcessData()` |
| Variable | camelBack | `userCount` |
| Constant | camelBack | `maxSize` |
| Type Alias | CamelCase | `StringMap` |

## Running Clang-Tidy

Integrated in VS Code with kvant-cpp.vsprofile

## Code Flow And Practices

### Const Keyword
Use const always where it is possible:
```cpp
// Methods
std::vector<int> GetData() const;

// Variables
const int width = 200;
const int height = 300;
```

### Pointers
Always use smart pointers:
```cpp
// ✅ Good
auto dataPtr = std::make_shared<Data>(data);

// ❌ Avoid
Data * dataPtr = &data;
```

### Error Code Flow
**Don't use exceptions**. Our team will use custom error codeflow:
```cpp
// Error class
class error {
public:  
  operator bool();
}

// Some class
class SomeClass {
public:  
  error doWork();
  std::tuple<std::vector<int>, error> getData();
}
```
Error handling:
```cpp
// Code Flow
SomeClass some;
auto err = some.doWork();
if (err) {
  return err; // or log(err.message()); or exit(err.code());
}

SomeClass some;
auto [data, err] = some.getData();
if (err) {
  return err; // or {nullptr, err} when multiple return values
}
ProcessData(data); // no else needed
```

## References 

[Google Code Style](https://google.github.io/styleguide/cppguide.html)


## C 

For developing C code, use different codestyle that your team leader provides with Clang Tidy rules