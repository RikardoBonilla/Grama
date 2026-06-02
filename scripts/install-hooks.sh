#!/usr/bin/env bash
# scripts/install-hooks.sh
# Installs git hooks for this repository.
# Run once after cloning: bash scripts/install-hooks.sh
set -euo pipefail

HOOKS_DIR="$(git rev-parse --git-dir)/hooks"
REPO_ROOT="$(git rev-parse --show-toplevel)"

# ─────────────────────────────────────────────
# pre-commit: block secrets from being staged
# ─────────────────────────────────────────────
cat > "$HOOKS_DIR/pre-commit" << 'HOOK'
#!/usr/bin/env bash
# Pre-commit hook: blocks .env files with real values and common secret patterns.
# "Real value" = anything that is not a known placeholder.
set -euo pipefail

PLACEHOLDERS="CHANGE_ME|changeme|YOUR_|REPLACE_|<.*>|example|placeholder|TODO|FIXME"

fail=0

# 1. Block staged .env files that contain non-placeholder values.
for file in $(git diff --cached --name-only | grep -E '^\.env$|\.env\.' | grep -v '\.env\.example'); do
    if git show ":$file" 2>/dev/null | grep -Eiq \
        '(password|passwd|secret|api_key|jwt|token|dsn)\s*=\s*[^$\s]' ; then
        # Allow if the value looks like a placeholder.
        if ! git show ":$file" 2>/dev/null | grep -Eiq \
            "(password|passwd|secret|api_key|jwt|token|dsn)\s*=\s*.*(${PLACEHOLDERS})"; then
            echo "ERROR: Staged file '$file' appears to contain real credentials."
            echo "       Add it to .gitignore or replace values with placeholders."
            fail=1
        fi
    fi
done

# 2. Block any staged file containing obvious hardcoded secret patterns.
DANGEROUS_PATTERNS=(
    'password\s*=\s*["\x27][^"$\x27]{8,}'
    'secret\s*=\s*["\x27][^"$\x27]{8,}'
    'api_key\s*=\s*["\x27][^"$\x27]{8,}'
    'jwt_secret\s*=\s*["\x27][^"$\x27]{8,}'
    'AKIA[0-9A-Z]{16}'
    'ghp_[0-9A-Za-z]{36}'
    'gho_[0-9A-Za-z]{36}'
)

for file in $(git diff --cached --name-only | grep -v '\.example$' | grep -v 'install-hooks\.sh'); do
    for pattern in "${DANGEROUS_PATTERNS[@]}"; do
        if git show ":$file" 2>/dev/null | grep -Eiq "$pattern"; then
            echo "ERROR: Staged file '$file' matches secret pattern: $pattern"
            fail=1
        fi
    done
done

if [ "$fail" -eq 1 ]; then
    echo ""
    echo "Commit blocked. Fix the issues above before committing."
    echo "If this is a false positive, check scripts/install-hooks.sh."
    exit 1
fi
HOOK

chmod +x "$HOOKS_DIR/pre-commit"
echo "✓ pre-commit hook installed at $HOOKS_DIR/pre-commit"
echo ""
echo "The hook blocks:"
echo "  - Staged .env files with real credentials (non-placeholder values)"
echo "  - Hardcoded passwords, secrets, API keys, JWT secrets"
echo "  - AWS access key patterns (AKIA...)"
echo "  - GitHub token patterns (ghp_, gho_)"
