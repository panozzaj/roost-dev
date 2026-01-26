# Code Improvements

Areas identified for potential refactoring and improvement.

## 1. Help Flag Checking Duplication ✅ DONE

**Files**: cert.go, service.go, ports.go, list.go, docs.go, logs.go

**Fix**: Extracted `checkHelpFlag(args []string, usage string) bool` helper in cli.go.

## 2. Config Loading Pattern Repeated ✅ DONE

**Files**: main.go, list.go, logs.go, cert.go, service.go

**Fix**: Extracted `getDefaultConfigDir()` and `getConfigWithDefaults() (*GlobalConfig, string)` helpers in cli.go.

## 3. ANSI Color Codes Scattered ✅ DONE

**Files**: main.go, ports.go, cert.go, service.go, list.go

**Fix**: Consolidated all color codes as constants in cli.go (colorRed, colorGreen, colorYellow, colorCyan, colorGray, colorDim, colorReset).

## 4. Silent Command Failures ✅ DONE

**Files**: ports.go, service.go, main.go

**Fix**: Added clarifying comments explaining why errors are intentionally ignored (e.g., "ignore errors - may not be running"). These are cleanup operations where failure is expected and acceptable.

## 5. Setup Wizard Too Long (189 lines)

**File**: main.go:539-726 `runSetupWizard()`

Contains 3 major sections with repeated status checking and confirmation patterns.

**Fix**: Split into `setupPortsStep()`, `setupCertsStep()`, `setupServiceStep()`.

## 6. API Handler Monolith

**File**: internal/server/api.go (352-line switch with 20+ cases)

Giant switch statement handling all API endpoints.

**Fix**: Use router pattern or handler registration, split into separate handler functions.

## 7. Service Name Resolution Duplicated

**File**: internal/server/api.go (3 locations)

Same parsing logic for service names repeated at lines ~73, ~159, ~246.

**Fix**: Consolidate into single `resolveServiceName()` function.

## 8. Missing Tests

No tests for:

- `runSetupWizard()` / `runTeardownWizard()`
- `runPortsInstall()` / `runPortsUninstall()`
- `runCertInstall()` / `runCertUninstall()`
- `runServiceInstall()` / `runServiceUninstall()`
- Most API endpoints in api.go

## 9. Root Check Duplication ✅ DONE

**Files**: main.go, ports.go, service.go

**Fix**: Extracted `requireNonRoot()` and `isRoot()` helpers in cli.go. Note: Not all usages need the helper (some have specific error messages or different logic).

## Priority Order

1. ~~Extract CLI helper functions (help flag, config loading, colors)~~ ✅
2. ~~Fix silent errors (add logging/comments)~~ ✅
3. Refactor setup/teardown wizards
4. Add tests for critical paths
5. Refactor API handler
