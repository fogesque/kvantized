#!/usr/bin/env bash
# init-cpp-project.sh — Initialize a new C++ project from the kvantized template.
#
# Usage: ./init-cpp-project.sh [project-name]
#   project-name  Optional. Defaults to the current directory name.
#
# What this script creates:
#   vcpkg.json                — vcpkg manifest
#   vcpkg-configuration.json  — vcpkg registry configuration
#   CMakeLists.txt            — top-level CMake file
#   .clang-format             — downloaded from fogesque/kvantized
#   .clang-tidy               — downloaded from fogesque/kvantized
#   CMakePresets.json         — downloaded from fogesque/kvantized
#   agent/                    — agent skills folder, downloaded from fogesque/kvantized

set -euo pipefail

# ── Configuration ─────────────────────────────────────────────────────────────

REPO_RAW="https://raw.githubusercontent.com/fogesque/kvantized/main"

# ── Helpers ───────────────────────────────────────────────────────────────────

die() { echo "ERROR: $*" >&2; exit 1; }

download() {
    local url="$1"
    local dest="$2"
    if command -v curl &>/dev/null; then
        curl -fsSL "$url" -o "$dest" || die "Failed to download $url"
    elif command -v wget &>/dev/null; then
        wget -q "$url" -O "$dest" || die "Failed to download $url"
    else
        die "Neither curl nor wget is available."
    fi
}

# ── Resolve project name ───────────────────────────────────────────────────────

PROJECT_NAME="${1:-$(basename "$(pwd)")}"

# vcpkg names must be lowercase alphanumeric + hyphens
VCPKG_NAME=$(echo "$PROJECT_NAME" | tr '[:upper:]' '[:lower:]' | tr '_' '-' | tr -cd 'a-z0-9-')

if [[ -z "$VCPKG_NAME" ]]; then
    die "Could not derive a valid vcpkg project name from '$PROJECT_NAME'."
fi

echo "Initializing C++ project: $PROJECT_NAME (vcpkg name: $VCPKG_NAME)"
echo

# ── 1. vcpkg.json ─────────────────────────────────────────────────────────────

cat > vcpkg.json <<EOF
{
    "name": "${VCPKG_NAME}",
    "version": "0.1.0",
    "dependencies": []
}
EOF
echo "[+] Created vcpkg.json"

# ── 2. vcpkg-configuration.json ───────────────────────────────────────────────

if [[ -z "${VCPKG_ROOT:-}" ]]; then
    die "VCPKG_ROOT is not set. Please set it to your vcpkg installation directory."
fi

VCPKG_BASELINE=$(git -C "$VCPKG_ROOT" rev-parse HEAD) \
    || die "Failed to get vcpkg baseline from $VCPKG_ROOT"

cat > vcpkg-configuration.json <<EOF
{
    "default-registry": {
        "kind": "git",
        "baseline": "${VCPKG_BASELINE}",
        "repository": "https://github.com/microsoft/vcpkg"
    },
    "registries": []
}
EOF
echo "[+] Created vcpkg-configuration.json (baseline: ${VCPKG_BASELINE})"

# ── 3. CMakeLists.txt ─────────────────────────────────────────────────────────

cat > CMakeLists.txt <<EOF
cmake_minimum_required(VERSION 3.25)
project(${PROJECT_NAME} LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 23)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

add_subdirectory()
EOF
echo "[+] Created CMakeLists.txt"

# ── 4. .clang-format and .clang-tidy ──────────────────────────────────────────

download "${REPO_RAW}/linter/.clang-format" ".clang-format"
echo "[+] Downloaded .clang-format"

download "${REPO_RAW}/linter/.clang-tidy" ".clang-tidy"
echo "[+] Downloaded .clang-tidy"

# ── 5. CMakePresets.json ──────────────────────────────────────────────────────

download "${REPO_RAW}/vscode/CMakePresets.json" "CMakePresets.json"
echo "[+] Downloaded CMakePresets.json"

# ── 6. agent/ folder ──────────────────────────────────────────────────────────

AGENT_FILES=(
    "agent/skills/ARCHITECTURE.md"
    "agent/skills/BUILD.md"
    "agent/skills/CODESTYLE.md"
    "agent/skills/GIT.md"
)

mkdir -p agent/skills

for file in "${AGENT_FILES[@]}"; do
    download "${REPO_RAW}/${file}" "${file}"
    echo "[+] Downloaded ${file}"
done

# ── 7. configure ─────────────────────────────────────────────────────────────

cat > configure <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

rm -rf build
cmake -S . -B build --preset amd64-linux-debug
EOF
chmod +x configure
echo "[+] Created configure"

# ── Done ──────────────────────────────────────────────────────────────────────

echo
echo "Done. Project '${PROJECT_NAME}' is ready."
echo "Next steps:"
echo "  1. Add your source subdirectory to CMakeLists.txt add_subdirectory()"
echo "  2. Configure with: cmake --preset amd64-linux-debug"
