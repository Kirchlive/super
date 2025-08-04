#!/bin/bash

# Script to verify all internal markdown links after reorganization
# This script checks if all linked files actually exist

echo "🔍 Starting link verification..."
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
total_links=0
broken_links=0
valid_links=0

# Function to check if file exists from the perspective of the source file
check_link() {
    local source_file=$1
    local link_path=$2
    local source_dir=$(dirname "$source_file")
    
    # Resolve the full path
    if [[ $link_path == /* ]]; then
        # Absolute path from repo root
        full_path=".${link_path}"
    else
        # Relative path
        full_path="${source_dir}/${link_path}"
    fi
    
    # Normalize path (remove ./ and ../)
    full_path=$(cd "$(dirname "$full_path")" 2>/dev/null && pwd)/$(basename "$full_path")
    
    # Check if file exists
    if [ -f "$full_path" ] || [ -d "$full_path" ]; then
        return 0
    else
        return 1
    fi
}

# Find all markdown files
echo "📂 Scanning markdown files..."
echo ""

find . -name "*.md" -type f | while read -r file; do
    # Skip backup directories
    if [[ $file == *"backup"* ]]; then
        continue
    fi
    
    echo "Checking: $file"
    
    # Extract all markdown links from the file
    # Matches [text](path) pattern
    grep -oE '\[[^]]+\]\([^)]+\)' "$file" | while read -r link; do
        # Extract the path from the link
        path=$(echo "$link" | sed -E 's/.*\(([^)]+)\).*/\1/')
        
        # Skip external links
        if [[ $path == http* ]] || [[ $path == mailto* ]] || [[ $path == "#"* ]]; then
            continue
        fi
        
        total_links=$((total_links + 1))
        
        # Check if the linked file exists
        if check_link "$file" "$path"; then
            valid_links=$((valid_links + 1))
            echo -e "  ${GREEN}✓${NC} $path"
        else
            broken_links=$((broken_links + 1))
            echo -e "  ${RED}✗${NC} $path ${RED}(FILE NOT FOUND)${NC}"
        fi
    done
    
    echo ""
done

# Summary
echo "═══════════════════════════════════════"
echo "📊 Link Verification Summary"
echo "═══════════════════════════════════════"
echo -e "Total links checked: ${YELLOW}$total_links${NC}"
echo -e "Valid links: ${GREEN}$valid_links${NC}"
echo -e "Broken links: ${RED}$broken_links${NC}"
echo ""

if [ $broken_links -eq 0 ]; then
    echo -e "${GREEN}✅ All links are valid!${NC}"
else
    echo -e "${RED}⚠️  Found $broken_links broken links that need to be fixed.${NC}"
    echo ""
    echo "To fix broken links:"
    echo "1. Run ./update-links.sh to automatically update links"
    echo "2. Or manually update the broken links in the affected files"
fi

echo ""
echo "💡 Tip: Run this script after any file reorganization to ensure all links are valid."