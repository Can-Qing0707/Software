$ErrorActionPreference = "Stop"
$BASE_URL = "http://localhost:8080"

function Send-Post($uri, $body) {
    $req = [System.Net.WebRequest]::Create($uri)
    $req.Method = "POST"
    $req.ContentType = "application/json"
    $bytes = [System.Text.Encoding]::UTF8.GetBytes($body)
    $req.ContentLength = $bytes.Length
    $stream = $req.GetRequestStream()
    $stream.Write($bytes, 0, $bytes.Length)
    $stream.Close()
    try {
        $resp = $req.GetResponse()
        $reader = New-Object System.IO.StreamReader($resp.GetResponseStream())
        $text = $reader.ReadToEnd()
        $reader.Close()
        $resp.Close()
        return @{ Success = $true; StatusCode = [int]$resp.StatusCode; Body = ($text | ConvertFrom-Json) }
    } catch [System.Net.WebException] {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $text = $reader.ReadToEnd()
        $reader.Close()
        $resp = $_.Exception.Response
        $resp.Close()
        return @{ Success = $false; StatusCode = [int]$resp.StatusCode; Body = ($text | ConvertFrom-Json) }
    }
}

function Send-Get($uri, $token) {
    $req = [System.Net.WebRequest]::Create($uri)
    $req.Method = "GET"
    $req.ContentType = "application/json"
    if ($token) {
        $req.Headers.Add("Authorization", "Bearer $token")
    }
    try {
        $resp = $req.GetResponse()
        $reader = New-Object System.IO.StreamReader($resp.GetResponseStream())
        $text = $reader.ReadToEnd()
        $reader.Close()
        $resp.Close()
        return @{ Success = $true; StatusCode = [int]$resp.StatusCode; Body = ($text | ConvertFrom-Json) }
    } catch [System.Net.WebException] {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $text = $reader.ReadToEnd()
        $reader.Close()
        $resp = $_.Exception.Response
        $resp.Close()
        return @{ Success = $false; StatusCode = [int]$resp.StatusCode; Body = ($text | ConvertFrom-Json) }
    }
}

$passed = 0
$failed = 0

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Login API Test Suite" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# --- Test 1: Normal login ---
Write-Host "[Test 1] Normal login (admin / 123456)" -ForegroundColor Yellow
$r = Send-Post "$BASE_URL/api/auth/login" '{"username":"admin","password":"123456"}'
if ($r.Success -and $r.Body.code -eq 0 -and $r.Body.data.token) {
    Write-Host "  PASS: code=$($r.Body.code), token=$($r.Body.data.token.SubString(0,20))..." -ForegroundColor Green
    Write-Host "  User: $($r.Body.data.user.username) ($($r.Body.data.user.real_name)), role=$($r.Body.data.user.role)"
    $token = $r.Body.data.token
    $passed++
} else {
    Write-Host "  FAIL: $($r.Body | ConvertTo-Json -Compress)" -ForegroundColor Red
    $failed++
}

# --- Test 2: Wrong password ---
Write-Host "[Test 2] Wrong password" -ForegroundColor Yellow
$r = Send-Post "$BASE_URL/api/auth/login" '{"username":"admin","password":"wrongpass"}'
if (-not $r.Success -and $r.Body.code -eq -1) {
    Write-Host "  PASS: code=$($r.Body.code), message=$($r.Body.message)" -ForegroundColor Green
    $passed++
} else {
    Write-Host "  FAIL: $($r.Body | ConvertTo-Json -Compress)" -ForegroundColor Red
    $failed++
}

# --- Test 3: Non-existent user ---
Write-Host "[Test 3] Non-existent user" -ForegroundColor Yellow
$r = Send-Post "$BASE_URL/api/auth/login" '{"username":"nobody","password":"123456"}'
if (-not $r.Success -and $r.Body.code -eq -1) {
    Write-Host "  PASS: code=$($r.Body.code), message=$($r.Body.message)" -ForegroundColor Green
    $passed++
} else {
    Write-Host "  FAIL: $($r.Body | ConvertTo-Json -Compress)" -ForegroundColor Red
    $failed++
}

# --- Test 4: Empty username ---
Write-Host "[Test 4] Empty username" -ForegroundColor Yellow
$r = Send-Post "$BASE_URL/api/auth/login" '{"username":"","password":"123456"}'
if (-not $r.Success -and $r.Body.code -eq -1) {
    Write-Host "  PASS: code=$($r.Body.code), message=$($r.Body.message)" -ForegroundColor Green
    $passed++
} else {
    Write-Host "  FAIL: $($r.Body | ConvertTo-Json -Compress)" -ForegroundColor Red
    $failed++
}

# --- Test 5: Empty password ---
Write-Host "[Test 5] Empty password" -ForegroundColor Yellow
$r = Send-Post "$BASE_URL/api/auth/login" '{"username":"admin","password":""}'
if (-not $r.Success -and $r.Body.code -eq -1) {
    Write-Host "  PASS: code=$($r.Body.code), message=$($r.Body.message)" -ForegroundColor Green
    $passed++
} else {
    Write-Host "  FAIL: $($r.Body | ConvertTo-Json -Compress)" -ForegroundColor Red
    $failed++
}

# --- Test 6: Missing username field ---
Write-Host "[Test 6] Missing username field" -ForegroundColor Yellow
$r = Send-Post "$BASE_URL/api/auth/login" '{"password":"123456"}'
if (-not $r.Success -and $r.Body.code -eq -1) {
    Write-Host "  PASS: code=$($r.Body.code), message=$($r.Body.message)" -ForegroundColor Green
    $passed++
} else {
    Write-Host "  FAIL: $($r.Body | ConvertTo-Json -Compress)" -ForegroundColor Red
    $failed++
}

# --- Test 7: Token validation (use token to access protected route) ---
if ($token) {
    Write-Host "[Test 7] Token validation (GET /api/users/profile)" -ForegroundColor Yellow
    $r = Send-Get "$BASE_URL/api/users/profile" $token
    if ($r.Success -and $r.Body.code -eq 0) {
        Write-Host "  PASS: code=$($r.Body.code), user=$($r.Body.data.username)" -ForegroundColor Green
        $passed++
    } else {
        Write-Host "  FAIL: $($r.Body | ConvertTo-Json -Compress)" -ForegroundColor Red
        $failed++
    }

    # --- Test 8: Invalid/expired token ---
    Write-Host "[Test 8] Invalid token" -ForegroundColor Yellow
    $r = Send-Get "$BASE_URL/api/users/profile" "invalid_token_xxx"
    if (-not $r.Success -and $r.Body.code -eq -1) {
        Write-Host "  PASS: code=$($r.Body.code), message=$($r.Body.message)" -ForegroundColor Green
        $passed++
    } else {
        Write-Host "  FAIL: $($r.Body | ConvertTo-Json -Compress)" -ForegroundColor Red
        $failed++
    }

    # --- Test 9: Teacher login ---
    Write-Host "[Test 9] Teacher login (teacher01 / 123456)" -ForegroundColor Yellow
    $r = Send-Post "$BASE_URL/api/auth/login" '{"username":"teacher01","password":"123456"}'
    if ($r.Success -and $r.Body.code -eq 0 -and $r.Body.data.user.role -eq "teacher") {
        Write-Host "  PASS: code=$($r.Body.code), role=$($r.Body.data.user.role)" -ForegroundColor Green
        $passed++
    } else {
        Write-Host "  FAIL: $($r.Body | ConvertTo-Json -Compress)" -ForegroundColor Red
        $failed++
    }

    # --- Test 10: Student login ---
    Write-Host "[Test 10] Student login (student01 / 123456)" -ForegroundColor Yellow
    $r = Send-Post "$BASE_URL/api/auth/login" '{"username":"student01","password":"123456"}'
    if ($r.Success -and $r.Body.code -eq 0 -and $r.Body.data.user.role -eq "student") {
        Write-Host "  PASS: code=$($r.Body.code), role=$($r.Body.data.user.role)" -ForegroundColor Green
        $passed++
    } else {
        Write-Host "  FAIL: $($r.Body | ConvertTo-Json -Compress)" -ForegroundColor Red
        $failed++
    }
} else {
    Write-Host "[Test 7-10] SKIPPED - no token from Test 1" -ForegroundColor DarkYellow
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Results: $passed passed, $failed failed, $(10 - $passed - $failed) skipped" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
