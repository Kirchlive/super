# 📁 Repository Reorganization Plan

## Overview
This document provides instructions for reorganizing the repository structure to improve clarity and navigation.

## Target Directory Structure

```
/Users/rob/Development/Super/
├── docs/
│   ├── planning/           # 📚 Planungsdokumente (5 Integrationspläne)
│   │   ├── ARCHITECTURE.md
│   │   ├── INTEGRATION.md
│   │   ├── IMPLEMENTATION.md
│   │   ├── EDA.md
│   │   └── EVENTDRIVEN.md
│   │
│   ├── guides/            # 📖 Anleitungen & Dokumentation
│   │   ├── PLAN-USAGE-GUIDE.md
│   │   ├── DECISION-LOG.md
│   │   ├── IMPLEMENTATION-ROADMAP.md
│   │   ├── MIGRATION-GUIDE.md
│   │   ├── TECHNICAL-GLOSSARY.md
│   │   ├── TECHNOLOGY-STACK.md
│   │   ├── SECURITY.md
│   │   ├── TESTING-STRATEGY.md
│   │   ├── MONITORING.md
│   │   ├── FAQ.md
│   │   └── CONTRIBUTING.md
│   │
│   └── api/              # 📁 API-SPECIFICATIONS/
│       ├── mcp-api.md
│       ├── plugin-api.md
│       └── event-api.md
│
├── examples/             # 💻 EXAMPLES/ (rename to lowercase)
│   ├── README.md
│   ├── simple-plugin/
│   ├── mcp-integration/
│   └── event-handler/
│
├── .github/             # 🔧 GitHub configuration (no change)
│   ├── ISSUE_TEMPLATE/
│   ├── workflows/
│   ├── PULL_REQUEST_TEMPLATE.md
│   ├── CODEOWNERS
│   ├── markdown-link-check-config.json
│   └── cspell.json
│
└── 📋 Root files (Hauptdokumentation)
    ├── README.md
    ├── CLAUDE.md
    ├── LICENSE
    ├── PROJECT-LINKS.md
    ├── DEPLOYMENT-INSTRUCTIONS.md
    ├── .gitignore
    └── (other existing files)
```

## Reorganization Commands

Execute these commands in order:

```bash
# 1. Create new directory structure
mkdir -p docs/planning docs/guides docs/api

# 2. Move planning documents
mv ARCHITECTURE.md docs/planning/
mv INTEGRATION.md docs/planning/
mv IMPLEMENTATION.md docs/planning/
mv EDA.md docs/planning/
mv EVENTDRIVEN.md docs/planning/

# 3. Move guide documents
mv PLAN-USAGE-GUIDE.md docs/guides/
mv DECISION-LOG.md docs/guides/
mv IMPLEMENTATION-ROADMAP.md docs/guides/
mv MIGRATION-GUIDE.md docs/guides/
mv TECHNICAL-GLOSSARY.md docs/guides/
mv TECHNOLOGY-STACK.md docs/guides/
mv SECURITY.md docs/guides/
mv TESTING-STRATEGY.md docs/guides/
mv MONITORING.md docs/guides/
mv FAQ.md docs/guides/
mv CONTRIBUTING.md docs/guides/

# 4. Move API specifications
mv API-SPECIFICATIONS/*.md docs/api/
rmdir API-SPECIFICATIONS

# 5. Rename EXAMPLES to examples (lowercase)
mv EXAMPLES examples

# 6. Update file references
# This will need to be done manually or with a script
```

## Files to Update After Reorganization

These files contain references that need to be updated:

1. **README.md** - Update all document links
2. **CLAUDE.md** - Update planning document references
3. **PROJECT-LINKS.md** - Update internal links
4. **docs/guides/PLAN-USAGE-GUIDE.md** - Update plan file references
5. **.github/CODEOWNERS** - Update file paths
6. **All issue templates** - Update documentation links

## Sample Link Updates

### Before:
```markdown
[ARCHITECTURE.md](./ARCHITECTURE.md)
[PLAN-USAGE-GUIDE.md](./PLAN-USAGE-GUIDE.md)
[API-SPECIFICATIONS/](./API-SPECIFICATIONS/)
```

### After:
```markdown
[ARCHITECTURE.md](./docs/planning/ARCHITECTURE.md)
[PLAN-USAGE-GUIDE.md](./docs/guides/PLAN-USAGE-GUIDE.md)
[API Specifications](./docs/api/)
```

## Automated Update Script

Save this as `update-links.sh`:

```bash
#!/bin/bash

# Update links in markdown files
find . -name "*.md" -type f -exec sed -i '' \
  -e 's|\./ARCHITECTURE\.md|./docs/planning/ARCHITECTURE.md|g' \
  -e 's|\./INTEGRATION\.md|./docs/planning/INTEGRATION.md|g' \
  -e 's|\./IMPLEMENTATION\.md|./docs/planning/IMPLEMENTATION.md|g' \
  -e 's|\./EDA\.md|./docs/planning/EDA.md|g' \
  -e 's|\./EVENTDRIVEN\.md|./docs/planning/EVENTDRIVEN.md|g' \
  -e 's|\./PLAN-USAGE-GUIDE\.md|./docs/guides/PLAN-USAGE-GUIDE.md|g' \
  -e 's|\./DECISION-LOG\.md|./docs/guides/DECISION-LOG.md|g' \
  -e 's|\./IMPLEMENTATION-ROADMAP\.md|./docs/guides/IMPLEMENTATION-ROADMAP.md|g' \
  -e 's|\./MIGRATION-GUIDE\.md|./docs/guides/MIGRATION-GUIDE.md|g' \
  -e 's|\./TECHNICAL-GLOSSARY\.md|./docs/guides/TECHNICAL-GLOSSARY.md|g' \
  -e 's|\./TECHNOLOGY-STACK\.md|./docs/guides/TECHNOLOGY-STACK.md|g' \
  -e 's|\./SECURITY\.md|./docs/guides/SECURITY.md|g' \
  -e 's|\./TESTING-STRATEGY\.md|./docs/guides/TESTING-STRATEGY.md|g' \
  -e 's|\./MONITORING\.md|./docs/guides/MONITORING.md|g' \
  -e 's|\./FAQ\.md|./docs/guides/FAQ.md|g' \
  -e 's|\./CONTRIBUTING\.md|./docs/guides/CONTRIBUTING.md|g' \
  -e 's|\./API-SPECIFICATIONS/|./docs/api/|g' \
  -e 's|\./EXAMPLES/|./examples/|g' \
  {} \;
```

## Benefits of New Structure

1. **Better Organization**: Clear separation between planning, guides, and API docs
2. **Scalability**: Easy to add new document categories
3. **Navigation**: More intuitive for newcomers
4. **Standards**: Follows common documentation practices
5. **Maintainability**: Easier to manage as project grows

## Post-Reorganization Tasks

1. Run tests to ensure all links work
2. Update GitHub Pages configuration if used
3. Verify all workflows still function
4. Create redirects for old URLs if needed
5. Update any external references to the documentation

---

**Note**: Make sure to commit all changes before reorganizing, and test thoroughly after completion.