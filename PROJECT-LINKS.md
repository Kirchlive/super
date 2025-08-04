Sehr gut, Fortsetz.# Project Links and Resources

## 🔗 Main Repositories

### Core Projects
- **Integration Documentation** (This Repo): [github.com/Kirchlive/super](https://github.com/Kirchlive/super)
- **OpenCode**: [github.com/Kirchlive/opencode](https://github.com/Kirchlive/opencode)
- **SuperClaude**: [github.com/Kirchlive/superclaude](https://github.com/Kirchlive/superclaude)

## 📚 External Resources

### Technologies
- **Model Context Protocol (MCP)**: [modelcontextprotocol.io](https://modelcontextprotocol.io)
- **HashiCorp go-plugin**: [github.com/hashicorp/go-plugin](https://github.com/hashicorp/go-plugin)
- **NATS JetStream**: [docs.nats.io/jetstream](https://docs.nats.io/jetstream)
- **gRPC**: [grpc.io](https://grpc.io)

### Tools & Libraries
- **Go**: [golang.org](https://golang.org)
- **TypeScript**: [typescriptlang.org](https://www.typescriptlang.org)
- **Docker**: [docker.com](https://www.docker.com)
- **Kubernetes**: [kubernetes.io](https://kubernetes.io)

### Monitoring & Observability
- **OpenTelemetry**: [opentelemetry.io](https://opentelemetry.io)
- **Prometheus**: [prometheus.io](https://prometheus.io)
- **Grafana**: [grafana.com](https://grafana.com)

## 🤝 Community

### Discussion & Support
- **Issues**: [github.com/Kirchlive/super/issues](https://github.com/Kirchlive/super/issues)
- **Discussions**: [github.com/Kirchlive/super/discussions](https://github.com/Kirchlive/super/discussions)

### Contributing
- See [CONTRIBUTING.md](./guides_and_documentation/CONTRIBUTING.md) for contribution guidelines
- Check [open issues](https://github.com/Kirchlive/super/issues) for ways to help
- Join discussions for feature requests and architectural decisions

## 📖 Documentation Structure

### Planning Documents
- [ARCHITECTURE.md](./planning_documents/ARCHITECTURE.md) - Plugin ecosystem architecture
- [INTEGRATION.md](./planning_documents/INTEGRATION.md) - Quick SuperClaude integration
- [IMPLEMENTATION.md](./planning_documents/IMPLEMENTATION.md) - Full SuperClaude implementation
- [EDA.md](./planning_documents/EDA.md) - Event-driven plugin platform
- [EVENTDRIVEN.md](./planning_documents/EVENTDRIVEN.md) - Enterprise MCP + EDA solution

### Reference Documents
- [PLAN-USAGE-GUIDE.md](./guides_and_documentation/PLAN-USAGE-GUIDE.md) - Comprehensive comparison
- [TECHNICAL-GLOSSARY.md](./guides_and_documentation/TECHNICAL-GLOSSARY.md) - Technical terms
- [API-SPECIFICATIONS/](./api_specifications/) - API documentation
- [EXAMPLES/](./examples/) - Implementation examples

## 🚀 Quick Links for Developers

### Getting Started
1. [README.md](./README.md) - Project overview
2. [PLAN-USAGE-GUIDE.md](./guides_and_documentation/PLAN-USAGE-GUIDE.md) - Choose your approach
3. [EXAMPLES/](./examples/) - Code examples

### Implementation
1. [TECHNOLOGY-STACK.md](./guides_and_documentation/TECHNOLOGY-STACK.md) - Tech stack details
2. [TESTING-STRATEGY.md](./guides_and_documentation/TESTING-STRATEGY.md) - Testing approach
3. [SECURITY.md](./guides_and_documentation/SECURITY.md) - Security guidelines

### Operations
1. [MONITORING.md](./guides_and_documentation/MONITORING.md) - Observability setup
2. [MIGRATION-GUIDE.md](./guides_and_documentation/MIGRATION-GUIDE.md) - Migration strategies
3. [FAQ.md](./guides_and_documentation/FAQ.md) - Common questions







 ✅ Skripte erfolgreich erstellt!

  Ich habe drei nützliche Dateien für dich erstellt:

  1. update-links.sh - Automatische Link-Aktualisierung

  - Erstellt automatisch ein Backup
  - Aktualisiert alle internen Markdown-Links auf die neue Verzeichnisstruktur
  - Behandelt relative Pfade korrekt
  - Aktualisiert auch die CODEOWNERS-Datei

  2. verify-links.sh - Link-Verifizierung

  - Überprüft alle Markdown-Dateien auf interne Links
  - Zeigt defekte Links in Rot an
  - Zeigt funktionierende Links in Grün an
  - Gibt eine Zusammenfassung am Ende aus

  3. SCRIPTS-README.md - Dokumentation

  - Erklärt die Verwendung beider Skripte
  - Zeigt den empfohlenen Workflow
  - Enthält Troubleshooting-Hinweise

  🚀 Empfohlene Verwendung:

  # 1. Zuerst den aktuellen Status prüfen
  ./verify-links.sh

  # 2. Links automatisch aktualisieren
  ./update-links.sh

  # 3. Ergebnis verifizieren
  ./verify-links.sh

  # 4. Änderungen committen
  git add -A
  git commit -m "fix: Update internal links for new directory structure"
  git push

  Die Skripte sind bereits ausführbar gemacht und bereit zur Verwendung. Sie
  helfen dabei, die neue Verzeichnisstruktur problemlos zu implementieren!

> Sehr gut fortsetzen.

⏺ Ich führe jetzt die Skripte aus und committe die Änderungen:

⏺ Bash(./verify-links.sh | head -50)
  ⎿  🔍 Starting link verification...

     📂 Scanning markdown files...
     … +48 lines (ctrl+r to expand)

⏺ Wie erwartet sind viele Links defekt, da die Dateien in neue Verzeichnisse
  verschoben wurden. Ich führe jetzt das Update-Skript aus:

⏺ Bash(./update-links.sh)
  ⎿  🔄 Starting link update process...    
     📦 Creating backup of current state...
     📝 Updating markdown files...