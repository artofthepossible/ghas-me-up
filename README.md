# GHAS Me Up

A deliberately vulnerable demo that tells a complete **build-time and runtime** GitHub Advanced
Security story. One small site on GitHub Pages is built three ways (Node.js, Ruby, Go), every build
ships stale dependencies, and GitHub finds and fixes them.

> ⚠️ Intentionally insecure. Do not copy this code into anything real. See [SECURITY.md](SECURITY.md).

## What's in the box

| Path | What it is | The stale part |
|------|------------|----------------|
| `node/` | Node build: Markdown + Handlebars → `node/index.html`, publishes jQuery | `package-lock.json`: lodash 4.17.15, handlebars 4.0.11, marked 0.3.5, minimist 1.2.0, jquery 3.4.1 |
| `ruby/` | Ruby build: kramdown + ERB → `ruby/index.html` | `Gemfile.lock` pins **Ruby 2.7.8** plus kramdown 2.3.0, rack 2.0.8, rexml 3.2.4, addressable 2.7.0, rake 12.3.2 |
| `go/` | Go build: `html/template` → `go/index.html` | `go.mod` pins **go 1.16** plus x/crypto (2020), x/text 0.3.0, yaml.v2 2.2.2, jwt-go 3.2.0 |
| `site/` | Static landing page + runtime widget | DOM XSS: `location.hash` → jQuery `.html()` |
| `.github/workflows/` | Build + deploy, CodeQL, dependency review, PR title lint | Pipeline flaws CodeQL finds in the workflow itself |
| `.github/dependabot.yml` | Version updates for npm, bundler, gomod, github-actions | |

## GHAS features it exercises

| Feature | Where it shows up |
|---------|-------------------|
| Dependency graph + **Dependabot alerts** | Three ecosystems, dozens of advisories |
| **Dependabot security updates** | One-click fix PRs for each alert |
| **Dependabot version updates** | Stale Actions and libraries |
| **Code scanning (CodeQL)** | JS (DOM XSS, weak hashing, insecure randomness), Ruby, Go (weak hash, hard-coded JWT key, disabled TLS verification), Actions (script injection, missing permissions) |
| **Dependency review** | Blocks PRs that add vulnerable packages |
| **Secret scanning + push protection** | Live in the demo script |

The run of show is in **[DEMO.md](DEMO.md)**.

## One-time repository setup

1. **Settings → Pages → Source: GitHub Actions**
2. **Settings → Advanced Security** (or *Code security*): enable Dependency graph, Dependabot alerts,
   Dependabot security updates, Secret scanning and Push protection. Code scanning comes from
   `codeql.yml`. (Private repos need GitHub Code Security / Secret Protection.)
3. Push to `main`. The workflow builds all three apps and deploys the site.

## Build locally

```bash
# Node
cd node && npm ci && npm run build            # -> ../dist/node, ../dist/vendor

# Ruby (needs Ruby 2.7.8 to match the lockfile)
cd ruby && bundle install && bundle exec rake build   # -> ../dist/ruby

# Go (first run generates go.sum)
cd go && go mod tidy && go run . -out ../dist  # -> ../dist/go

# Assemble and serve
cp -R site/. dist/ && python3 -m http.server -d dist 8000
```
