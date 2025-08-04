sorttiere jetzt die Repo nach dieser Vorgabe mit den Namen auf Englisch.
/Users/rob/Development/Super/
├── 📋 Hauptdokumentation
│   ├── README.md (mit Links zu Kirchlive/opencode und Kirchlive/superclaude)
│   ├── CLAUDE.md (AI-Kontext)
│   ├── LICENSE (MIT License)
│   └── PROJECT-LINKS.md (alle wichtigen Links)
│
├── 📚 Planungsdokumente (5 Integrationspläne)
│   ├── ARCHITECTURE.md
│   ├── INTEGRATION.md
│   ├── IMPLEMENTATION.md
│   ├── EDA.md
│   └── EVENTDRIVEN.md
│
├── 📖 Anleitungen & Dokumentation
│   ├── PLAN-USAGE-GUIDE.md
│   ├── DECISION-LOG.md
│   ├── IMPLEMENTATION-ROADMAP.md
│   ├── MIGRATION-GUIDE.md
│   ├── TECHNICAL-GLOSSARY.md
│   ├── TECHNOLOGY-STACK.md
│   ├── SECURITY.md
│   ├── TESTING-STRATEGY.md
│   ├── MONITORING.md
│   ├── FAQ.md
│   └── CONTRIBUTING.md
│
├── 📁 API-SPECIFICATIONS/
│   ├── mcp-api.md
│   ├── plugin-api.md
│   └── event-api.md
│
├── 💻 EXAMPLES/
│   ├── README.md
│   ├── simple-plugin/
│   ├── mcp-integration/
│   └── event-handler/
│
└── 🔧 .github/
├── ISSUE\_TEMPLATE/
│   ├── bug\_report.md
│   ├── feature\_request.md
│   └── documentation.md










  1. Alle internen Links aktualisieren - Da Dateien verschoben wurden, müssen
  Markdown-Links aktualisiert werden
  2. Änderungen committen:
  git add -A
  git commit -m "refactor: Reorganize repository structure

  - Move documents into categorized directories
  - Use English directory names for consistency
  - Maintain all existing content and configurations
  - Improve navigation and scalability"

  git push
  3. Root README.md aktualisieren - Hinweis zur neuen Struktur hinzufügen
  4. GitHub Workflows überprüfen - Sicherstellen, dass sie noch die korrekten
  Dateipfade finden

  Möchtest du, dass ich ein Skript erstelle, um automatisch alle internen
  Markdown-Links zu aktualisieren, damit sie die neuen Dateispeicherorte
  widerspiegeln?

> Gute Idee, bitte.

⏺ Ich erstelle ein Skript, das automatisch alle internen Markdown-Links in der
  neuen Verzeichnisstruktur aktualisiert:

⏺ Write(update-links.sh)
  ⎿  Wrote 157 lines to update-links.sh                                       
     #!/bin/bash
     # Script to update all markdown links after repository reorganization
     # This script updates internal links to match the new directory structure
     echo "🔄 Starting link update process..."
     # Create backup
     echo "📦 Creating backup of current state..."
     cp -r . ../super-backup-$(date +%Y%m%d-%H%M%S)
     … +147 lines (ctrl+r to expand)

⏺ Bash(chmod +x /Users/rob/Development/Super/update-links.sh)
  ⎿  (No content)

⏺ Jetzt erstelle ich auch ein Verifizierungsskript, um sicherzustellen, dass
  alle Links korrekt funktionieren:

⏺ Write(verify-links.sh)
  ⎿  Wrote 104 lines to verify-links.sh                                 
     #!/bin/bash
     # Script to verify all internal markdown links after reorganization
     # This script checks if all linked files actually exist
     echo "🔍 Starting link verification..."
     echo ""
     # Colors for output
     RED='\033[0;31m'
     … +94 lines (ctrl+r to expand)

⏺ Bash(chmod +x /Users/rob/Development/Super/verify-links.sh)
  ⎿  (No content)

✢ Inferring… (53s · ↑ 2.9k tokens · esc to interrupt)

╭──────────────────────────────────────────────────────────────────────────────────╮
│ >                   


























├── workflows/
│   └── docs-validation.yml
├── PULL\_REQUEST\_TEMPLATE.md
├── CODEOWNERS
├── markdown-link-check-config.json
└── cspell.json
