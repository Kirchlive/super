# 🛠️ Repository Maintenance Scripts

This directory contains utility scripts for maintaining the repository after reorganization.

## Available Scripts

### 1. `update-links.sh`
**Purpose**: Automatically updates all internal markdown links to match the new directory structure.

**Usage**:
```bash
./update-links.sh
```

**What it does**:
- Creates a backup of the current repository state
- Updates all markdown links to point to the new file locations
- Updates the CODEOWNERS file with new paths
- Handles relative paths correctly for files in subdirectories

**Features**:
- Safe: Creates automatic backup before making changes
- Smart: Calculates correct relative paths based on file location
- Complete: Updates all markdown files and GitHub configuration

### 2. `verify-links.sh`
**Purpose**: Verifies that all internal markdown links point to existing files.

**Usage**:
```bash
./verify-links.sh
```

**What it does**:
- Scans all markdown files for internal links
- Checks if each linked file actually exists
- Reports broken links with their locations
- Provides a summary of valid vs broken links

**Output**:
- ✅ Green checkmarks for valid links
- ❌ Red X marks for broken links
- Summary statistics at the end

## Recommended Workflow

After reorganizing the repository structure:

1. **First, verify current state**:
   ```bash
   ./verify-links.sh
   ```
   This will show you which links are broken.

2. **Update all links automatically**:
   ```bash
   ./update-links.sh
   ```
   This will fix all internal links and create a backup.

3. **Verify the fix**:
   ```bash
   ./verify-links.sh
   ```
   This should now show all links as valid.

4. **Review and commit changes**:
   ```bash
   git diff                    # Review changes
   git add -A                  # Stage all changes
   git commit -m "fix: Update internal links for new directory structure"
   git push
   ```

## Manual Link Patterns

If you need to manually update links, here are the new patterns:

### Old Structure → New Structure
```
./ARCHITECTURE.md → ./planning_documents/ARCHITECTURE.md
./PLAN-USAGE-GUIDE.md → ./guides_and_documentation/PLAN-USAGE-GUIDE.md
./API-SPECIFICATIONS/ → ./api_specifications/
./EXAMPLES/ → ./examples/
./README.md → ./main_documentation/README.md
```

### Relative Paths
From files in subdirectories, use `../` to go up:
- From `planning_documents/`: `../guides_and_documentation/FILE.md`
- From `api_specifications/`: `../planning_documents/FILE.md`
- From `.github/ISSUE_TEMPLATE/`: `../../main_documentation/FILE.md`

## Troubleshooting

**Script won't run**: Make sure it's executable:
```bash
chmod +x update-links.sh verify-links.sh
```

**Backup location**: Backups are created in the parent directory with timestamp:
```
../super-backup-YYYYMMDD-HHMMSS/
```

**macOS vs Linux**: The scripts use `sed -i ''` syntax for macOS. For Linux, remove the empty quotes:
```bash
# macOS
sed -i '' 's/old/new/g' file

# Linux
sed -i 's/old/new/g' file
```

---

*These scripts help maintain link integrity after repository reorganization.*