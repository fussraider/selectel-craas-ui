1. **Understand the Vulnerability:**
   - The `JWTSecret` configuration defaults to an empty string `""` if the `JWT_SECRET` environment variable is not set.
   - If `AUTH_ENABLED` is true and `JWT_SECRET` is not provided, the application refuses to start (`backend/cmd/server/main.go`).
   - However, if the `JWT_SECRET` environment variable is somehow empty but not caught by the check, or if we want a safe fallback, an empty string as a secret means tokens are signed with a known (empty) key. This is a critical security vulnerability.
   - Wait, `main.go` actually checks `cfg.JWTSecret == ""` and calls `log.Fatal` if `AuthEnabled` is true. But what if `AuthEnabled` is true and someone sets `JWT_SECRET` to spaces, or what if the application logic relies on `JWTSecret` being non-empty? Or even better, if `JWTSecret` is missing, it's safer to generate a strong random secret by default so that the application works out of the box securely if `JWTSecret` isn't provided, or just fail to load config.
   - Wait, let me re-read the prompt. "Default Empty JWT Secret... The fix is to require a non-empty secret or generate a random one on startup."
   - Generating a random one on startup is a great idea, because it allows the application to be secure without requiring the user to explicitly provide a secret. Although, if it's restarted, sessions will be invalidated. That's a reasonable trade-off for zero-conf security.

2. **Implement the Fix in `config.go`:**
   - In `Load()`, retrieve `JWT_SECRET`.
   - If it's empty, generate a random secure string (e.g., using `crypto/rand`).
   - Assign this string to `JWTSecret`.
   - Let's remove the `log.Fatal` in `main.go` for `JWTSecret == ""` since it will never be empty anymore, or keep it as a defensive check.

3. **Write a helper function to generate a random secret:**
   ```go
   import "crypto/rand"
   import "encoding/hex"

   func generateRandomSecret(length int) string {
       b := make([]byte, length)
       if _, err := rand.Read(b); err != nil {
           log.Fatalf("Failed to generate random JWT secret: %v", err)
       }
       return hex.EncodeToString(b)
   }
   ```

4. **Update `config.go`:**
   ```go
   jwtSecret := getEnv("JWT_SECRET", "")
   if jwtSecret == "" {
       // Generate a 32-byte (256-bit) random secret
       jwtSecret = generateRandomSecret(32)
       log.Println("JWT_SECRET not provided, generated a random secret for this session")
   }
   ```

5. **Update `main.go`:**
   - Remove the `if cfg.JWTSecret == ""` check, or keep it but it will never trigger. Let's keep it or remove it depending on cleanliness. It's safe to remove since we ensure it's set in `config.Load()`. Actually, keeping it as an invariant check isn't bad.

6. **Tests:**
   - Run `go test ./...` in the `backend/` directory to ensure no tests break. Some tests might assume `JWTSecret` can be empty or set it to empty. If they do, they'll just get a random one if they call `config.Load()`.
