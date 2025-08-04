#!/bin/bash

# Script to update all markdown links after repository reorganization
# This script updates internal links to match the new directory structure

echo "🔄 Starting link update process..."

# Create backup
echo "📦 Creating backup of current state..."
cp -r . ../super-backup-$(date +%Y%m%d-%H%M%S)

# Function to update links in a file
update_file_links() {
    local file=$1
    echo "  Updating: $file"
    
    # Planning documents
    sed -i '' \
        -e 's|\(\[\|(\)\./ARCHITECTURE\.md|$1./planning_documents/ARCHITECTURE.md|g' \
        -e 's|\(\[\|(\)\./INTEGRATION\.md|$1./planning_documents/INTEGRATION.md|g' \
        -e 's|\(\[\|(\)\./IMPLEMENTATION\.md|$1./planning_documents/IMPLEMENTATION.md|g' \
        -e 's|\(\[\|(\)\./EDA\.md|$1./planning_documents/EDA.md|g' \
        -e 's|\(\[\|(\)\./EVENTDRIVEN\.md|$1./planning_documents/EVENTDRIVEN.md|g' \
        "$file"
    
    # Guide documents
    sed -i '' \
        -e 's|\(\[\|(\)\./PLAN-USAGE-GUIDE\.md|$1./guides_and_documentation/PLAN-USAGE-GUIDE.md|g' \
        -e 's|\(\[\|(\)\./DECISION-LOG\.md|$1./guides_and_documentation/DECISION-LOG.md|g' \
        -e 's|\(\[\|(\)\./IMPLEMENTATION-ROADMAP\.md|$1./guides_and_documentation/IMPLEMENTATION-ROADMAP.md|g' \
        -e 's|\(\[\|(\)\./MIGRATION-GUIDE\.md|$1./guides_and_documentation/MIGRATION-GUIDE.md|g' \
        -e 's|\(\[\|(\)\./TECHNICAL-GLOSSARY\.md|$1./guides_and_documentation/TECHNICAL-GLOSSARY.md|g' \
        -e 's|\(\[\|(\)\./TECHNOLOGY-STACK\.md|$1./guides_and_documentation/TECHNOLOGY-STACK.md|g' \
        -e 's|\(\[\|(\)\./SECURITY\.md|$1./guides_and_documentation/SECURITY.md|g' \
        -e 's|\(\[\|(\)\./TESTING-STRATEGY\.md|$1./guides_and_documentation/TESTING-STRATEGY.md|g' \
        -e 's|\(\[\|(\)\./MONITORING\.md|$1./guides_and_documentation/MONITORING.md|g' \
        -e 's|\(\[\|(\)\./FAQ\.md|$1./guides_and_documentation/FAQ.md|g' \
        -e 's|\(\[\|(\)\./CONTRIBUTING\.md|$1./guides_and_documentation/CONTRIBUTING.md|g' \
        "$file"
    
    # Main documentation
    sed -i '' \
        -e 's|\(\[\|(\)\./README\.md|$1./main_documentation/README.md|g' \
        -e 's|\(\[\|(\)\./CLAUDE\.md|$1./main_documentation/CLAUDE.md|g' \
        -e 's|\(\[\|(\)\./LICENSE|$1./main_documentation/LICENSE|g' \
        -e 's|\(\[\|(\)\./PROJECT-LINKS\.md|$1./main_documentation/PROJECT-LINKS.md|g' \
        "$file"
    
    # API specifications
    sed -i '' \
        -e 's|\(\[\|(\)\./API-SPECIFICATIONS/|$1./api_specifications/|g' \
        -e 's|\(\[\|(\)\.\.\/API-SPECIFICATIONS/|$1../api_specifications/|g' \
        "$file"
    
    # Examples
    sed -i '' \
        -e 's|\(\[\|(\)\./EXAMPLES/|$1./examples/|g' \
        -e 's|\(\[\|(\)\.\.\/EXAMPLES/|$1../examples/|g' \
        "$file"
}

# Function to calculate relative path
get_relative_path() {
    local from=$1
    local to=$2
    
    # For files in subdirectories, we need to go up first
    if [[ $from == *"planning_documents"* ]] || [[ $from == *"guides_and_documentation"* ]] || [[ $from == *"api_specifications"* ]]; then
        echo "../"
    elif [[ $from == *"examples/"* ]] && [[ $from != *"examples/README.md" ]]; then
        echo "../../"
    elif [[ $from == *".github/"* ]]; then
        if [[ $from == *"ISSUE_TEMPLATE"* ]] || [[ $from == *"workflows"* ]]; then
            echo "../../"
        else
            echo "../"
        fi
    else
        echo "./"
    fi
}

# Update all markdown files
echo "📝 Updating markdown files..."

# Find all markdown files and update them
find . -name "*.md" -type f | while read -r file; do
    # Skip backup directories
    if [[ $file == *"backup"* ]]; then
        continue
    fi
    
    # Get the relative path prefix for this file
    rel_path=$(get_relative_path "$file" ".")
    
    # Create temp file with updated content
    temp_file="${file}.tmp"
    cp "$file" "$temp_file"
    
    # Update links with relative path
    sed -i '' \
        -e "s|\[\([^]]*\)\](\./ARCHITECTURE\.md)|[\1](${rel_path}planning_documents/ARCHITECTURE.md)|g" \
        -e "s|\[\([^]]*\)\](\./INTEGRATION\.md)|[\1](${rel_path}planning_documents/INTEGRATION.md)|g" \
        -e "s|\[\([^]]*\)\](\./IMPLEMENTATION\.md)|[\1](${rel_path}planning_documents/IMPLEMENTATION.md)|g" \
        -e "s|\[\([^]]*\)\](\./EDA\.md)|[\1](${rel_path}planning_documents/EDA.md)|g" \
        -e "s|\[\([^]]*\)\](\./EVENTDRIVEN\.md)|[\1](${rel_path}planning_documents/EVENTDRIVEN.md)|g" \
        -e "s|\[\([^]]*\)\](\./PLAN-USAGE-GUIDE\.md)|[\1](${rel_path}guides_and_documentation/PLAN-USAGE-GUIDE.md)|g" \
        -e "s|\[\([^]]*\)\](\./DECISION-LOG\.md)|[\1](${rel_path}guides_and_documentation/DECISION-LOG.md)|g" \
        -e "s|\[\([^]]*\)\](\./IMPLEMENTATION-ROADMAP\.md)|[\1](${rel_path}guides_and_documentation/IMPLEMENTATION-ROADMAP.md)|g" \
        -e "s|\[\([^]]*\)\](\./MIGRATION-GUIDE\.md)|[\1](${rel_path}guides_and_documentation/MIGRATION-GUIDE.md)|g" \
        -e "s|\[\([^]]*\)\](\./TECHNICAL-GLOSSARY\.md)|[\1](${rel_path}guides_and_documentation/TECHNICAL-GLOSSARY.md)|g" \
        -e "s|\[\([^]]*\)\](\./TECHNOLOGY-STACK\.md)|[\1](${rel_path}guides_and_documentation/TECHNOLOGY-STACK.md)|g" \
        -e "s|\[\([^]]*\)\](\./SECURITY\.md)|[\1](${rel_path}guides_and_documentation/SECURITY.md)|g" \
        -e "s|\[\([^]]*\)\](\./TESTING-STRATEGY\.md)|[\1](${rel_path}guides_and_documentation/TESTING-STRATEGY.md)|g" \
        -e "s|\[\([^]]*\)\](\./MONITORING\.md)|[\1](${rel_path}guides_and_documentation/MONITORING.md)|g" \
        -e "s|\[\([^]]*\)\](\./FAQ\.md)|[\1](${rel_path}guides_and_documentation/FAQ.md)|g" \
        -e "s|\[\([^]]*\)\](\./CONTRIBUTING\.md)|[\1](${rel_path}guides_and_documentation/CONTRIBUTING.md)|g" \
        -e "s|\[\([^]]*\)\](\./API-SPECIFICATIONS/|[\1](${rel_path}api_specifications/|g" \
        -e "s|\[\([^]]*\)\](\./EXAMPLES/|[\1](${rel_path}examples/|g" \
        "$temp_file"
    
    # Move temp file back
    mv "$temp_file" "$file"
    echo "  ✓ Updated: $file"
done

# Special handling for CODEOWNERS file
echo "📋 Updating CODEOWNERS file..."
if [ -f ".github/CODEOWNERS" ]; then
    sed -i '' \
        -e 's|^ARCHITECTURE\.md|planning_documents/ARCHITECTURE.md|g' \
        -e 's|^INTEGRATION\.md|planning_documents/INTEGRATION.md|g' \
        -e 's|^IMPLEMENTATION\.md|planning_documents/IMPLEMENTATION.md|g' \
        -e 's|^EDA\.md|planning_documents/EDA.md|g' \
        -e 's|^EVENTDRIVEN\.md|planning_documents/EVENTDRIVEN.md|g' \
        -e 's|^/API-SPECIFICATIONS/|/api_specifications/|g' \
        -e 's|^/EXAMPLES/|/examples/|g' \
        -e 's|^CONTRIBUTING\.md|guides_and_documentation/CONTRIBUTING.md|g' \
        -e 's|^SECURITY\.md|guides_and_documentation/SECURITY.md|g' \
        -e 's|^CLAUDE\.md|main_documentation/CLAUDE.md|g' \
        -e 's|^README\.md|main_documentation/README.md|g' \
        -e 's|^MIGRATION-GUIDE\.md|guides_and_documentation/MIGRATION-GUIDE.md|g' \
        -e 's|^TESTING-STRATEGY\.md|guides_and_documentation/TESTING-STRATEGY.md|g' \
        .github/CODEOWNERS
    echo "  ✓ Updated: .github/CODEOWNERS"
fi

echo ""
echo "✅ Link update completed!"
echo ""
echo "📋 Next steps:"
echo "1. Review the changes: git diff"
echo "2. Test some links manually to verify they work"
echo "3. Commit the changes: git add -A && git commit -m 'fix: Update internal links for new directory structure'"
echo "4. Push to repository: git push"
echo ""
echo "💡 Tip: The backup was created in the parent directory with timestamp."