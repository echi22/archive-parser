# Archive Parser

A Go application that parses custom `.env` archive formats and extracts embedded files.

## 🚀 Try it Online

**[Open in GitHub Codespaces](https://codespaces.new/yourusername/archive-parser)**

Once the codespace opens:
\`\`\`bash
# Create test file
echo '**%%DOCUTEST
FILENAME/hello.txt
EXT/.txt
TYPE/TEXT
GUID/123
SHA1/AAF4C61DDCC5E8A2DABEDE0F3B482CD9AEA9434D
_SIG/D.C.hello world**%%' > test.env

# Run parser
./archive-parser test.env

# Check output
cat extracted/hello.txt
\`\`\`

## 🏠 Test Locally

\`\`\`bash
# Clone and test
git clone <repo-url>
cd archive-parser

# Build and test
make test
make build

# Create test file
echo '**%%DOCUTEST
FILENAME/hello.txt
EXT/.txt
TYPE/TEXT
GUID/123
SHA1/AAF4C61DDCC5E8A2DABEDE0F3B482CD9AEA9434D
_SIG/D.C.hello world**%%' > test.env

# Run parser
./archive-parser test.env

# Check output
cat extracted/hello.txt
\`\`\`

## Usage

\`\`\`bash
# Basic usage
./archive-parser archive.env

# Specify output directory
./archive-parser archive.env ./output

# Show help
./archive-parser --help
