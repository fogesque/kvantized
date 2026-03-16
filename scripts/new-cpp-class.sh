#!/usr/bin/env bash
# new-cpp-class.sh — Generate a C++ header and source file pair from codestyle template.
#
# Usage: ./new-cpp-class.sh <ClassName> [namespace]
#   ClassName   Required. PascalCase name of the class (e.g. FightClub).
#   namespace   Optional. Namespace name (e.g. fight). Defaults to empty.
#
# What this script creates (in the current directory):
#   ClassName.hpp  — header file
#   ClassName.cpp  — source file

set -euo pipefail

# ── Helpers ───────────────────────────────────────────────────────────────────

die() { echo "ERROR: $*" >&2; exit 1; }

# ── Arguments ─────────────────────────────────────────────────────────────────

CLASS_NAME="${1:-}"
NAMESPACE="${2:-}"

[[ -z "$CLASS_NAME" ]] && die "ClassName is required. Usage: $0 <ClassName> [namespace]"

# ── Derived paths ──────────────────────────────────────────────────────────────

HEADER_PATH="${CLASS_NAME}.hpp"
SOURCE_PATH="${CLASS_NAME}.cpp"

if [[ -n "$NAMESPACE" ]]; then
    INCLUDE_GUARD="#include <${NAMESPACE}/${CLASS_NAME}.hpp>"
else
    INCLUDE_GUARD="#include <${CLASS_NAME}.hpp>"
fi

for path in "$HEADER_PATH" "$SOURCE_PATH"; do
    if [[ -f "$path" ]]; then
        die "File already exists: $path"
    fi
done

# ── Generate header file ───────────────────────────────────────────────────────

if [[ -n "$NAMESPACE" ]]; then
    cat > "$HEADER_PATH" <<EOF
// ${CLASS_NAME}.hpp
#pragma once

#include <errors/errors.hpp>

namespace ${NAMESPACE}
{

// Forward declarations
class ${CLASS_NAME};

// Type aliases
using ${CLASS_NAME}Ptr = std::shared_ptr<${CLASS_NAME}>;

/// @brief Errors that can occur in ${CLASS_NAME}
namespace ErrorTypes
{
}  // namespace ErrorTypes

///
/// @brief
/// ${CLASS_NAME}
///
class ${CLASS_NAME}
{
public:
    /// [Fabric Methods]

    /// @brief Creates ${CLASS_NAME} instance with given configuration
    static std::tuple<${CLASS_NAME}Ptr, error> Create(const Config & config);

    /// [Construction & Destruction]

#pragma region ${CLASS_NAME}::Construct

    /// @brief Default constructor is deleted
    ${CLASS_NAME}() = delete;

    /// @brief Copy constructor is deleted
    ${CLASS_NAME}(const ${CLASS_NAME} &) = delete;
    
    /// @brief Copy operator is deleted
    ${CLASS_NAME} & operator=(const ${CLASS_NAME} &) = delete;

    /// @brief Move constructor is deleted
    ${CLASS_NAME}(${CLASS_NAME} &&) = delete;

    /// @brief Move operator is deleted
    ${CLASS_NAME} & operator=(${CLASS_NAME} &&) = delete;

    /// @brief Config struct for object construction
    struct Config {
    };

    /// @brief Constructor
    /// @warning Avoid using this constructor since class has static fabric methods
    explicit ${CLASS_NAME}(const Config & config);
    /// @brief Destructor
    ~${CLASS_NAME}();

#pragma endregion

private:
#pragma region ${CLASS_NAME}::PrivateMethods

#pragma endregion

    /// [Properties]

};

}  // namespace ${NAMESPACE}
EOF
else
    cat > "$HEADER_PATH" <<EOF
// ${CLASS_NAME}.hpp
#pragma once

#include <errors/errors.hpp>

// Forward declarations
class ${CLASS_NAME};

// Type aliases
using ${CLASS_NAME}Ptr = std::shared_ptr<${CLASS_NAME}>;

/// @brief Errors that can occur in ${CLASS_NAME}
namespace ErrorTypes
{
}  // namespace ErrorTypes

///
/// @brief
/// ${CLASS_NAME}
///
class ${CLASS_NAME}
{
public:
    /// [Fabric Methods]

    /// @brief Creates ${CLASS_NAME} instance with given configuration
    static std::tuple<${CLASS_NAME}Ptr, error> Create(const Config & config);

    /// [Construction & Destruction]

#pragma region ${CLASS_NAME}::Construct

    /// @brief Default constructor is deleted
    ${CLASS_NAME}() = delete;

    /// @brief Copy constructor is deleted
    ${CLASS_NAME}(const ${CLASS_NAME} &) = delete;

    /// @brief Copy operator is deleted
    ${CLASS_NAME} & operator=(const ${CLASS_NAME} &) = delete;

    /// @brief Move constructor is deleted
    ${CLASS_NAME}(${CLASS_NAME} &&) = delete;

    /// @brief Move operator is deleted
    ${CLASS_NAME} & operator=(${CLASS_NAME} &&) = delete;

    /// @brief Config struct for object construction
    struct Config {
    };

    /// @brief Constructor
    /// @warning Avoid using this constructor since class has static fabric methods
    explicit ${CLASS_NAME}(const Config & config);

    /// @brief Destructor
    ~${CLASS_NAME}();

#pragma endregion

private:
#pragma region ${CLASS_NAME}::PrivateMethods

#pragma endregion

    /// [Properties]

};
EOF
fi

echo "[+] Created ${HEADER_PATH}"

# ── Generate source file ───────────────────────────────────────────────────────

if [[ -n "$NAMESPACE" ]]; then
    cat > "$SOURCE_PATH" <<EOF
// ${CLASS_NAME}.cpp
${INCLUDE_GUARD}

namespace ${NAMESPACE}
{

std::tuple<${CLASS_NAME}Ptr, error> ${CLASS_NAME}::Create(const Config & config)
{
    const auto instance = std::make_shared<${CLASS_NAME}>(config);
    return {instance, nullptr};
}

${CLASS_NAME}::${CLASS_NAME}(const Config & config)
{
}

${CLASS_NAME}::~${CLASS_NAME}() = default;

}  // namespace ${NAMESPACE}
EOF
else
    cat > "$SOURCE_PATH" <<EOF
// ${CLASS_NAME}.cpp
${INCLUDE_GUARD}

std::tuple<${CLASS_NAME}Ptr, error> ${CLASS_NAME}::Create(const Config & config)
{
    const auto instance = std::make_shared<${CLASS_NAME}>(config);
    return {instance, nullptr};
}

${CLASS_NAME}::${CLASS_NAME}(const Config & config)
{
}

${CLASS_NAME}::~${CLASS_NAME}() = default;
EOF
fi

echo "[+] Created ${SOURCE_PATH}"

# ── Done ──────────────────────────────────────────────────────────────────────

echo
echo "Done. Class '${CLASS_NAME}' is ready."
if [[ -n "$NAMESPACE" ]]; then
    echo "  Namespace: ${NAMESPACE}"
fi
echo "  Header: ${HEADER_PATH}"
echo "  Source: ${SOURCE_PATH}"
